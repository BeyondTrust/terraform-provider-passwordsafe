package provider_framework

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/entities"
	"terraform-provider-passwordsafe/providers/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const totpCodeSecretId = "9152f5b6-07d6-4955-175a-08db047219ce"

// newTotpCodeMockServer mocks the Password Safe API; the TOTP code endpoint answers with
// totpStatus and totpBody.
func newTotpCodeMockServer(t *testing.T, totpStatus int, totpBody string) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, _ = w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))

		case constants.APIPath + "/Auth/SignAppIn":
			_, _ = w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

		case constants.APIPath + "/secrets-safe/secrets/" + totpCodeSecretId + "/totp/code":
			if r.Method != "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(totpStatus)
			_, _ = w.Write([]byte(totpBody))

		case constants.APIPath + "/Auth/Signout":
			_, _ = w.Write([]byte(``))
		}
	}))
	t.Cleanup(server.Close)
	server.URL = server.URL + constants.APIPath
	return server
}

// totpCodeTestConfig returns the test configuration for the TOTP code ephemeral resource.
func totpCodeTestConfig(serverUrl, apiVersion, secretId string) string {
	return utils.TestResourceConfig(entities.PasswordSafeTestConfig{
		ClientID:     constants.FakeClientId,
		ClientSecret: constants.FakeClientSecret,
		URL:          serverUrl,
		APIVersion:   apiVersion,
		Resource: `
		ephemeral "passwordsafe_secret_totp_code_ephemeral" "test" {
		secret_id = "` + secretId + `"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_secret_totp_code_ephemeral.test
		}

		resource "echo" "test" {}`,
	})
}

// runTotpCodeTest runs a single-step test of the TOTP code ephemeral resource.
func runTotpCodeTest(t *testing.T, step resource.TestStep) {
	t.Helper()
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
			"echo":         echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{step},
	})
}

func TestEphemeralSecretTotpCode(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusOK, `{"Code": "123456", "ValidFromUtc": "2026-10-19T10:00:00Z", "ValidToUtc": "2026-10-19T10:00:30Z"}`)

	runTotpCodeTest(t, resource.TestStep{
		Config: totpCodeTestConfig(server.URL, "3.3", totpCodeSecretId),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("code"), knownvalue.StringExact("123456")),
			statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("valid_from_utc"), knownvalue.StringExact("2026-10-19T10:00:00Z")),
			statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("valid_to_utc"), knownvalue.StringExact("2026-10-19T10:00:30Z")),
		},
	})
}

func TestEphemeralSecretTotpCodeNotAvailable(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusNotFound, `"TOTP is not enabled"`)

	runTotpCodeTest(t, resource.TestStep{
		Config:      totpCodeTestConfig(server.URL, "3.3", totpCodeSecretId),
		ExpectError: regexp.MustCompile("secret not found or TOTP is not enabled for this secret"),
	})
}

func TestEphemeralSecretTotpCodeNotCredentialSecret(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusBadRequest, `"Secret is not a credential"`)

	runTotpCodeTest(t, resource.TestStep{
		Config:      totpCodeTestConfig(server.URL, "3.3", totpCodeSecretId),
		ExpectError: regexp.MustCompile("secret is not a credential secret"),
	})
}

func TestEphemeralSecretTotpCodeInvalidSecretId(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusOK, `{"Code": "123456"}`)

	runTotpCodeTest(t, resource.TestStep{
		Config:      totpCodeTestConfig(server.URL, "3.3", "not-a-uuid"),
		ExpectError: regexp.MustCompile("must be a valid UUID"),
	})
}
