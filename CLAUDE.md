# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Terraform provider for BeyondTrust Password Safe / BeyondInsight, distributed via the public Terraform Registry as `BeyondTrust/passwordsafe`. The provider wraps `github.com/BeyondTrust/go-client-library-passwordsafe` — the Go SDK is the source of truth for API behavior; this repo is the Terraform-facing adapter.

## Build, test, lint

```bash
# Build the provider binary (versioned name is what Terraform expects in the plugin dir)
go build -o terraform-provider-passwordsafe_<major>_<minor>_<build>

# Unit tests — TF_ACC=1 is required (the acceptance-test framework is used even
# for the unit-style tests in this repo). Either invocation works:
cd providers && TF_ACC=1 go test ./...
# or, equivalent, from repo root (this is what sonarqube.yml CI uses):
TF_ACC=1 go test ./providers/...

# Single test
cd providers && TF_ACC=1 go test ./provider_framework -run TestName

# Coverage (matches CI)
cd providers && TF_ACC=1 go test -cover -coverprofile=coverage.txt ./...

# Lint (CI runs both)
golangci-lint run
gocognit -ignore "_test|testdata" -over 10 .   # CI fails any function with cognitive complexity >10
```

CI is in `.github/workflows/` — `golint.yml` (lint + gocognit), `coverage-report.yml` (tests), `automation-tests.yaml` (release-time integration), `release.yml`, `wiz.yml`, `sonarqube.yml`.

Go version is pinned to **1.26.4** (see `go.mod`; CI reads it via `go-version-file` in the `.github/actions/setup-go` composite action). `.goreleaser.yml` builds release artifacts (CGO disabled, multi-arch).

### Local end-to-end testing against a real Password Safe instance

`config/config.go` reads `PS_API_KEY`, `PS_ACCOUNT_NAME`, `PS_URL`, `CERTIFICATE_PATH`, `CERTIFICATE_NAME`, `CERTIFICATE_PASSWORD` from env (or `.env` via godotenv). The `terraform/` directory holds an end-to-end manifest (`main.tf` + `terraform.tfvars`) used to exercise every resource/data-source/ephemeral against a live instance — `terraform/` is gitignored. `script.bat` shows the Windows install flow (build → copy to `%APPDATA%\terraform.d\plugins\providers\beyondtrust\passwordsafe\<version>\windows_amd64\`).

## Architecture: two providers behind one binary (mux)

`main.go` uses `tf5muxserver` to serve **two** provider implementations under the single `registry.terraform.io/providers/BeyondTrust/passwordsafe` address:

1. **`providers/provider_sdkv2/`** — the original provider, built on `terraform-plugin-sdk/v2`. Owns: `passwordsafe_managed_account`, `passwordsafe_credential_secret`, `passwordsafe_text_secret`, `passwordsafe_file_secret`, `passwordsafe_folder`, `passwordsafe_safe`, and the data sources `passwordsafe_secret` and `passwordsafe_managed_account`.
2. **`providers/provider_framework/`** — the newer provider, built on `terraform-plugin-framework`. Required because **ephemeral resources** are framework-only. Owns: workgroup/asset/database/managed-system/functional-account resources and most data sources, plus `passwordsafe_secret_ephemeral` and `passwordsafe_managed_acccount_ephemeral`.

When adding a new resource decide which provider it belongs in: ephemeral → framework; everything else → either, but new work generally lives in `provider_framework`. Both providers expose the **same** schema attributes (`api_key`, `client_id`, `client_secret`, `url`, `api_account_name`, `verify_ca`, certificate fields, etc.) — they must stay in sync, and a single `provider "passwordsafe"` block in user HCL configures both halves.

### Shared session state (CRITICAL)

Both providers run in the **same process** and must share a single Password Safe session. The shared state lives in `providers/utils/methods.go`:

- `SignInCount uint64` — package-level reference count of active session holders.
- `signAppinResponse` — the cached session.
- `AuthMu sync.Mutex` — **single** mutex guarding both signin and signout.

Rules — violating these has caused a real race-condition bug (see commit `ea692e1`):

- Every Open/Read/Create/Update/Delete that hits the API must call `utils.Authenticate(..., &utils.AuthMu, &utils.SignInCount, ...)` and pair it with `utils.SignOut(...)` using the **same** `AuthMu` and `SignInCount`. Do not introduce a second mutex for signin vs. signout — they must be serialized against each other because the API's signout is user-global.
- `SignOut` must guard against underflow (no decrement when count is 0) and must decrement even on signout error so a transient failure does not pin the counter open forever.
- Concurrent Terraform workers will call into these functions in parallel; assume parallelism and never assume "I'm the only caller."

`utils/methods.go` also contains shared helpers like `DeleteAssetByID` and `ValidateChangeFrequencyDays`. `utils/managed_system_operations.go` and `utils/schema_attributes.go` hold cross-resource helpers that the framework provider's managed-system resources reuse.

### Provider package layout

- `providers/constants/constants.go` — API path constants and fake values used by tests.
- `providers/entities/entities.go` — `PasswordSafeTestConfig` used to render HCL into acceptance tests via `utils.TestResourceConfig`.
- `providers/utils/` — shared auth, schema, helpers (above).
- `providers/provider_sdkv2/` — sdkv2 resources/data sources, one file per resource plus `provider.go` and `common.go`.
- `providers/provider_framework/` — framework resources/data sources/ephemerals, one file per concept plus `provider.go`. The `*_test.go` files alongside use `terraform-plugin-testing` to drive HCL through the provider.

Each provider writes its own debug log: `providerSdkv2.log` and `providerFramework.log` (configured via zap in each `provider.go`). These are gitignored by virtue of being build artifacts.

### Authentication paths

Both providers accept either an **API Key** (`api_key` + `api_account_name`, formatted as `<key>;runas=<account>;`) or **OAuth client credentials** (`client_id` + `client_secret`). Optional client-certificate auth loads a PFX via `utils.GetPFXContent`. Validation lives in `ValidateCredentialsAndConfig` in each provider; keep the rules identical between the two.

## Conventions worth knowing

- **Cognitive complexity**: CI hard-fails any non-test function over 10. Prefer extracting helpers over nested branching.
- **Examples and docs**: every resource/data-source/ephemeral has a folder under `examples/` and a generated page under `docs/`. New resources should add both.
- **Resource naming**: `passwordsafe_<entity>` for resources/data sources, `passwordsafe_<entity>_ephemeral` for ephemerals. There is one historical typo — `passwordsafe_managed_acccount_ephemeral` (three c's) — leave it; renaming is a breaking change for users.
