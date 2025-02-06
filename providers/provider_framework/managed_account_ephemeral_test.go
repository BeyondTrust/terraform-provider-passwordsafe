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

func TestEphemeralManagedAcount(t *testing.T) {

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

		case "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/ManagedAccounts":
			_, err := w.Write([]byte(`{"SystemId":1,"AccountId":10}`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/Requests":
			_, err := w.Write([]byte(`124`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/Credentials/124":
			_, err := w.Write([]byte(`"fake_credential"`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case "/Requests/124/checkin":
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
				// test usging aki key authentication
				Config: testManagedAccountEphemeralResourceUsingApiKey(serverURL.String()),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("fake_credential"),
					),
				},
			},

			{
				// test using oauth authentication
				Config: testManagedAccountEphemeralResourceUsingOauth(serverURL.String()),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("fake_credential"),
					),
				},
			},
		},
	})
}

func testManagedAccountEphemeralResourceUsingApiKey(serverURL string) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = "fake_api_key_82a8a8e48b488d"
			client_id = "fake_client_id_35e7dd5093ae"
			client_secret = "fake_cliente_secret_JPMhmYTRSpeHJY"
			url =  %[1]q
			api_account_name = "apikey_user"
			client_certificates_folder_path = ""
			client_certificate_name = ""
			client_certificate_password = ""
			api_version = "3.1"
		}

		ephemeral "passwordsafe_managed_acccount_ephemeral" "test" {
		system_name = "server01"
		account_name = "managed_account_01"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_managed_acccount_ephemeral.test
		}

		resource "echo" "test" {}

		`, serverURL)
}

func testManagedAccountEphemeralResourceUsingOauth(serverURL string) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = ""
			client_id = "6138d050-e266-4b05-9ced-35e7dd5093ae"
			client_secret = "W8dx3BMkkxe4OpdsJPMhmYTRSpeHJYA/NVmcnmPZv5s="
			url =  %[1]q
			api_account_name = "apikey_user"
			client_certificates_folder_path = ""
			client_certificate_name = ""
			client_certificate_password = ""
			api_version = "3.1"
		}

		ephemeral "passwordsafe_managed_acccount_ephemeral" "test" {
		system_name = "server01"
		account_name = "managed_account_01"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_managed_acccount_ephemeral.test
		}

		resource "echo" "test" {}

		`, serverURL)
}

func testManagedAccountEphemeralResourceError(serverURL string) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = ""
			client_id = ""
			client_secret = ""
			url =  %[1]q
			api_account_name = "apikey_user"
			client_certificates_folder_path = ""
			client_certificate_name = ""
			client_certificate_password = ""
			api_version = "3.1"
		}

		ephemeral "passwordsafe_managed_acccount_ephemeral" "test" {
		system_name = "server01"
		account_name = "managed_account_01"
		}

		provider "echo" {
		data = ephemeral.passwordsafe_managed_acccount_ephemeral.test
		}

		resource "echo" "test" {}

		`, serverURL)
}
