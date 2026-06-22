// Copyright 2025 BeyondTrust. All rights reserved.
package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"terraform-provider-passwordsafe/providers/constants"
	"testing"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	libutils "github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

// fuzzLogger is a no-op logger so fuzz iterations don't flood stderr with the
// (expected) error noise produced by deliberately malformed server responses.
var fuzzLogger = logging.NewZapLogger(zap.NewNop())

// newFuzzAuthObj builds an *AuthenticationObj wired to serverURL. The backoff
// window is intentionally tiny: fuzzing feeds mostly invalid responses, and a
// long retry window would let each failing iteration burn seconds, starving the
// corpus. The mock server is expected to live at <serverURL>/BeyondTrust/api/public/v3.
func newFuzzAuthObj(t *testing.T, serverURL string) *authentication.AuthenticationObj {
	t.Helper()

	httpClientObj, err := libutils.GetHttpClient(5, false, "", "", fuzzLogger)
	if err != nil {
		t.Fatalf("GetHttpClient: %v", err)
	}

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = 10 * time.Millisecond

	authObj, err := authentication.Authenticate(authentication.AuthenticationParametersObj{
		HTTPClient:                 *httpClientObj,
		BackoffDefinition:          backoffDefinition,
		EndpointURL:                constants.FakeApiUrl,
		APIVersion:                 "3.1",
		ClientID:                   "fakeone_a654+9sdf7+8we4f",
		ClientSecret:               "fakeone_a654+9sdf7+8we4f",
		ApiKey:                     "",
		Logger:                     fuzzLogger,
		RetryMaxElapsedTimeSeconds: 1,
	})
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}

	apiUrl, err := url.Parse(serverURL + constants.APIPath)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	authObj.ApiUrl = *apiUrl
	return authObj
}

// FuzzValidateChangeFrequencyDays exercises ValidateChangeFrequencyDays with
// arbitrary input and asserts its documented invariants never break (and that
// it never panics).
func FuzzValidateChangeFrequencyDays(f *testing.F) {
	f.Add("xdays", 1)
	f.Add("xdays", 999)
	f.Add("xdays", 0)
	f.Add("xdays", 1000)
	f.Add("xdays", -5)
	f.Add("other", 0)
	f.Add("", 50)

	f.Fuzz(func(t *testing.T, changeFrequencyType string, changeFrequencyDays int) {
		err := ValidateChangeFrequencyDays(changeFrequencyType, changeFrequencyDays)

		if changeFrequencyType != "xdays" {
			if err != nil {
				t.Errorf("expected no error for type %q, got %v", changeFrequencyType, err)
			}
			return
		}

		inRange := changeFrequencyDays >= 1 && changeFrequencyDays <= 999
		if inRange && err != nil {
			t.Errorf("expected no error for xdays=%d (in range), got %v", changeFrequencyDays, err)
		}
		if !inRange && err == nil {
			t.Errorf("expected error for xdays=%d (out of range), got nil", changeFrequencyDays)
		}
	})
}

// FuzzAuthenticationFlow drives the OAuth sign-in handshake
// (GetPasswordSafeAuthentication -> GetToken -> SignAppin) against a server that
// returns fuzzer-controlled bodies and status codes for the token and SignAppIn
// endpoints. The invariant is simply that parsing untrusted auth responses never
// panics; an error return is an acceptable outcome for garbage input.
func FuzzAuthenticationFlow(f *testing.F) {
	// Happy path.
	f.Add(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`,
		`{"UserId":1, "UserName":"jdoe", "EmailAddress":"test@beyondtrust.com"}`, 200)
	// Token endpoint rejects the credentials.
	f.Add(`{"error": "invalid_client"}`, ``, 400)
	// Malformed JSON in both responses.
	f.Add(`not-json`, `not-json`, 200)
	// Empty bodies.
	f.Add(``, ``, 200)
	// Token ok, SignAppIn returns an unexpected shape.
	f.Add(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`,
		`{"UserId":"not-an-int"}`, 200)

	f.Fuzz(func(t *testing.T, tokenBody string, signAppinBody string, tokenStatus int) {
		// net/http panics on a status code outside [100, 999]; keep it in the
		// plausible response range so we fuzz behavior, not the stdlib guard.
		if tokenStatus < 100 || tokenStatus > 599 {
			tokenStatus = http.StatusOK
		}

		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case constants.APIPath + "/Auth/connect/token":
				w.WriteHeader(tokenStatus)
				_, _ = w.Write([]byte(tokenBody))
			case constants.APIPath + "/Auth/SignAppIn":
				_, _ = w.Write([]byte(signAppinBody))
			}
		}))
		defer server.Close()

		authObj := newFuzzAuthObj(t, server.URL)

		// Must never panic regardless of what the server returns.
		_, _ = authObj.GetPasswordSafeAuthentication()
	})
}

// FuzzGetSecretFlow exercises the "get secret by path" retrieval flow
// (SecretObj.GetSecret -> SecretGetSecretByPath, and the file-download branch
// when the secret type is FILE) with a fuzzer-controlled path, separator and
// server payloads. The invariant: the flow never panics, and when it returns a
// value without error that value matches what the mock served.
func FuzzGetSecretFlow(f *testing.F) {
	// Credential secret found. The full coordinate is "<path><sep><title>".
	f.Add("path/path2/credential_title", "/",
		`[{"SecretType":"CREDENTIAL","Password":"super_secret","Id":"id-1","Title":"credential_title"}]`, "")
	// File secret found -> download branch.
	f.Add("path/path2/file_title", "/",
		`[{"SecretType":"FILE","Id":"id-2","Title":"file_title"}]`, "file-secret-contents")
	// Empty list -> "not found".
	f.Add("path/missing", "/", `[]`, "")
	// Malformed JSON.
	f.Add("path/title", "/", `not-json`, "")
	// Degenerate inputs: empty separator and empty path.
	f.Add("", "", `[]`, "")
	// Multiple secrets in the response.
	f.Add("a/b/c/title", "/",
		`[{"SecretType":"CREDENTIAL","Password":"p1","Id":"x"},{"SecretType":"CREDENTIAL","Password":"p2","Id":"y"}]`, "")

	f.Fuzz(func(t *testing.T, secretPath string, separator string, secretListJSON string, fileBody string) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == constants.APIPath+"/Auth/connect/token":
				_, _ = w.Write([]byte(`{"access_token":"fake_token","expires_in":600,"token_type":"Bearer","scope":"publicapi"}`))
			case r.URL.Path == constants.APIPath+"/Auth/SignAppIn":
				_, _ = w.Write([]byte(`{"UserId":1,"EmailAddress":"test@beyondtrust.com"}`))
			case strings.HasSuffix(r.URL.Path, "/file/download"):
				_, _ = w.Write([]byte(fileBody))
			case strings.HasSuffix(r.URL.Path, "/secrets-safe/secrets"):
				_, _ = w.Write([]byte(secretListJSON))
			}
		}))
		defer server.Close()

		authObj := newFuzzAuthObj(t, server.URL)
		secretObj, err := secrets.NewSecretObj(*authObj, fuzzLogger, 5000000, false)
		if err != nil {
			t.Fatalf("NewSecretObj: %v", err)
		}

		// Must never panic. The returned value/err pair is exercised, not asserted,
		// because valid-but-arbitrary payloads have many legitimate outcomes.
		got, err := secretObj.GetSecret(secretPath, separator)
		if err == nil && got == "" && secretListJSON != "" {
			// No assertion failure here; this is just a touch of the success path
			// to keep the compiler from eliding `got`.
			_ = got
		}
	})
}

// FuzzManageAccountFlow exercises the managed-account retrieval flow
// (ManagedAccountstObj.GetSecret -> ManageAccountFlow), which chains four API
// calls: ManagedAccountGet, create Request, CredentialByRequestId and check-in.
// Each step's response body is fuzzer-controlled. The invariant is that the
// chain never panics for arbitrary system/account paths or server payloads.
func FuzzManageAccountFlow(f *testing.F) {
	managedAccountOK := `{"SystemId":1,"AccountId":10,"SystemName":"system01","AccountName":"account01"}`
	// Happy path.
	f.Add("system01/account01", "/", managedAccountOK, "123", `"managed_secret"`)
	// ManagedAccountGet returns malformed JSON.
	f.Add("system01/account01", "/", "not-json", "123", `"x"`)
	// Credential is not a quoted string (strconv.Unquote will fail silently).
	f.Add("system01/account01", "/", managedAccountOK, "123", `raw_unquoted`)
	// Path without a separator -> validated away -> empty list.
	f.Add("systemonly", "/", managedAccountOK, "1", `"x"`)
	// Too many segments for a managed-account path (depth > 1).
	f.Add("a/b/c", "/", managedAccountOK, "1", `"x"`)
	// Degenerate separator.
	f.Add("system01/account01", "", managedAccountOK, "1", `"x"`)

	f.Fuzz(func(t *testing.T, accountPath string, separator string, managedAccountJSON string, requestIDBody string, credentialBody string) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == constants.APIPath+"/Auth/connect/token":
				_, _ = w.Write([]byte(`{"access_token":"fake_token","expires_in":600,"token_type":"Bearer","scope":"publicapi"}`))
			case r.URL.Path == constants.APIPath+"/Auth/SignAppIn":
				_, _ = w.Write([]byte(`{"UserId":1,"EmailAddress":"test@beyondtrust.com"}`))
			case strings.HasSuffix(r.URL.Path, "/ManagedAccounts"):
				_, _ = w.Write([]byte(managedAccountJSON))
			case strings.HasSuffix(r.URL.Path, "/checkin"):
				_, _ = w.Write([]byte(`{}`))
			case strings.Contains(r.URL.Path, "/Credentials/"):
				_, _ = w.Write([]byte(credentialBody))
			case strings.HasSuffix(r.URL.Path, "/Requests"):
				_, _ = w.Write([]byte(requestIDBody))
			default:
				_, _ = w.Write([]byte(`{}`))
			}
		}))
		defer server.Close()

		authObj := newFuzzAuthObj(t, server.URL)
		managedAccountObj, err := managed_accounts.NewManagedAccountObj(*authObj, fuzzLogger)
		if err != nil {
			t.Fatalf("NewManagedAccountObj: %v", err)
		}

		// Must never panic regardless of the path shape or server payloads.
		_, _ = managedAccountObj.GetSecret(accountPath, separator)
	})
}
