package provider_framework

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
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

var SecretEphemeralOauthConfig entities.PasswordSafeTestConfig = entities.PasswordSafeTestConfig{
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	URL:                          "",
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
	ephemeral "passwordsafe_secret_ephemeral" "test" {
	title = "secret_title"
	path = "secret_path"
	separator = "/"
	}

	provider "echo" {
	data = ephemeral.passwordsafe_secret_ephemeral.test
	}

	resource "echo" "test" {}`,
}

func TestEphemeralSecret(t *testing.T) {

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Auth/SignAppIn":
			_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/secrets-safe/secrets":
			_, err := w.Write([]byte(`[{"SecretType": "SECRET", "Password": "fake_password_a#$%!","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	server.URL = server.URL + constants.APIPath

	APIKeyConfig := entities.PasswordSafeTestConfig{
		APIKey:                       "fake_api_key_82a8a8e48b488d",
		ClientID:                     "",
		ClientSecret:                 "",
		URL:                          "",
		APIAccountName:               "api_key_username",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		ephemeral "passwordsafe_secret_ephemeral" "test" {
		title = "secret_title"
		path = "secret_path"
		separator = "/"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_secret_ephemeral.test
		}

		resource "echo" "test" {}`,
	}

	APIKeyConfig.URL = server.URL

	SecretEphemeralOauthConfig.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		// load providers, echo is just for test purposes
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
			"echo":         echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				// test using aki key authentication
				Config: utils.TestResourceConfig(APIKeyConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("fake_password_a#$%!"),
					),
				},
			},
			{
				// test using oauth authentication
				Config: utils.TestResourceConfig(SecretEphemeralOauthConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("fake_password_a#$%!"),
					),
				},
			},
		},
	})
}

// TestEphemeralSecretCustomSeparator verifies that when the user provides a
// non-default `separator` (here "|"), the ephemeral secret resource forwards
// the same separator to the Password Safe API. The mock server asserts the
// expected query parameters (path, title, separator) instead of using the
// default "/" and returns a known value the test checks back through echo.
func TestEphemeralSecretCustomSeparator(t *testing.T) {

	customSeparator := "|"
	expectedPath := "secret_path"
	expectedTitle := "secret_title"

	// handlerErr captures assertion failures from the HTTP handler goroutine.
	// We cannot call t.Errorf/t.Error from a non-test goroutine because that
	// can panic with "testing: t.Errorf called after test finished" or have
	// the failure swallowed. Instead, store the error under handlerMu and
	// assert it from the main test goroutine after resource.Test returns.
	var (
		handlerMu  sync.Mutex
		handlerErr error
	)
	recordHandlerErr := func(err error) {
		handlerMu.Lock()
		defer handlerMu.Unlock()
		// Only keep the first error to avoid losing the original cause.
		if handlerErr == nil {
			handlerErr = err
		}
	}

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
			if err != nil {
				recordHandlerErr(err)
			}

		case constants.APIPath + "/Auth/SignAppIn":
			_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

			if err != nil {
				recordHandlerErr(err)
			}

		case constants.APIPath + "/secrets-safe/secrets":
			// Assert the custom separator was forwarded to the API. The
			// underlying client library passes path, title and separator
			// as query parameters; if the provider had hard-coded "/" the
			// separator value would not match.
			query := r.URL.Query()
			if got := query.Get("separator"); got != customSeparator {
				recordHandlerErr(fmt.Errorf("expected separator=%q in request, got %q", customSeparator, got))
			}
			if got := query.Get("path"); got != expectedPath {
				recordHandlerErr(fmt.Errorf("expected path=%q in request, got %q", expectedPath, got))
			}
			if got := query.Get("title"); got != expectedTitle {
				recordHandlerErr(fmt.Errorf("expected title=%q in request, got %q", expectedTitle, got))
			}
			_, err := w.Write([]byte(`[{"SecretType": "SECRET", "Password": "fake_password_custom_sep","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))
			if err != nil {
				recordHandlerErr(err)
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				recordHandlerErr(err)
			}
		}

	}))

	server.URL = server.URL + constants.APIPath

	customSepConfig := entities.PasswordSafeTestConfig{
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		URL:                          server.URL,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		ephemeral "passwordsafe_secret_ephemeral" "test" {
		title = "secret_title"
		path = "secret_path"
		separator = "|"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_secret_ephemeral.test
		}

		resource "echo" "test" {}`,
	}

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
			"echo":         echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: utils.TestResourceConfig(customSepConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("fake_password_custom_sep"),
					),
				},
			},
		},
	})

	// Surface any assertion error captured by the HTTP handler goroutine
	// from the main test goroutine. This avoids the "t.Errorf called after
	// test finished" panic and ensures the failure is not swallowed.
	handlerMu.Lock()
	defer handlerMu.Unlock()
	if handlerErr != nil {
		t.Errorf("HTTP handler assertion failed: %v", handlerErr)
	}
}

func TestEphemeralSecretNotFound(t *testing.T) {

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Auth/SignAppIn":
			_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/secrets-safe/secrets":
			_, err := w.Write([]byte(`[]`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	server.URL = server.URL + constants.APIPath

	SecretEphemeralOauthConfig.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
			"echo":         echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				// test using oauth authentication
				Config:      utils.TestResourceConfig(SecretEphemeralOauthConfig),
				ExpectError: regexp.MustCompile("error SecretGetSecretByPath, Secret was not found: StatusCode: 404"),
			},
		},
	})
}
