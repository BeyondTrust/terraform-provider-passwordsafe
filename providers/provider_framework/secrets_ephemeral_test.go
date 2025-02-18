package providerv2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestEphemeralSecret(t *testing.T) {

	// mocking Password Safe API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case "/Auth/connect/token":
			_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/Auth/SignAppIn":
			_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

			if err != nil {
				t.Error("Test case Failed")
			}

		case "/secrets-safe/secrets":
			_, err := w.Write([]byte(`[{"SecretType": "SECRET", "Password": "fake_password_a#$%!","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error("Test case Failed")
			}
		}

	}))

	serverURL, _ := url.Parse(server.URL + "/")

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
				Config: testSecretEphemeralResourceUsingApiKey(serverURL.String()),
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
				Config: testSecretEphemeralResourceUsingOauth(serverURL.String()),
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

func testSecretEphemeralResourceUsingApiKey(serverURL string) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = "fake_api_key_82a8a8e48b488d"
			client_id = ""
			client_secret = ""
			url =  %[1]q
			api_account_name = "apikey_user"
			client_certificates_folder_path = ""
			client_certificate_name = ""
			client_certificate_password = ""
			api_version = "3.1"
		}

		ephemeral "passwordsafe_secret_ephemeral" "test" {
		title = "secret_title"
		path = "secret_path"
		separator = "/"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_secret_ephemeral.test
		}

		resource "echo" "test" {}

		`, serverURL)
}

func testSecretEphemeralResourceUsingOauth(serverURL string) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = ""
			client_id = "fake_client_id_35e7dd5093ae"
			client_secret = "fake_cliente_secret_JPMhmYTRSpeHJY"
			url =  %[1]q
			api_account_name = ""
			client_certificates_folder_path = ""
			client_certificate_name = ""
			client_certificate_password = ""
			api_version = "3.1"
		}

		ephemeral "passwordsafe_secret_ephemeral" "test" {
		title = "secret_title"
		path = "secret_path"
		separator = "/"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_secret_ephemeral.test
		}

		resource "echo" "test" {}

	`, serverURL)
}
