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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestCreateFunctionaAccount(t *testing.T) {

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

		case constants.APIPath + "/FunctionalAccounts":
			_, err := w.Write([]byte(`{ "PlatformID": 1, "DomainName": "corp.example.com", "AccountName": "svc-backup" }`))
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

	configFunctioalAccount := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_functional_account" "functional_account" {
			platform_id =           1
			domain_name =           "corp.example.com"
			account_name =          "svc-monitoring"
			display_name =          "FUNCTIONAL_ACCOUNT"
			password =              "P@ssw0rd123!"
			private_key =           "private key value"
			passphrase =            "my-passphrase"
			description =           "Used for monitoring agents to access the platform"
			elevation_command =     "sudo"
			tenant_id =             "123e4567-e89b-12d3-a456-426614174000"
			object_id =             "abc12345-def6-7890-gh12-ijklmnopqrst"
			secret =                "super-secret-value"
			service_account_email = "monitoring@project.iam.gserviceaccount.com"
			azure_instance =        "AzurePublic"
		}`,
	}

	configFunctioalAccount.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
		},
		Steps: []resource.TestStep{
			{
				// test using oauth authentication, creating functional account.
				Config: utils.TestResourceConfig(configFunctioalAccount),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_functional_account.functional_account",
						tfjsonpath.New("domain_name"),
						knownvalue.StringExact("corp.example.com"),
					),
				},
			},
		},
	})
}

func TestCreateFunctionaAccountBadRequest(t *testing.T) {

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

		case constants.APIPath + "/FunctionalAccounts":
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"Bad request"}`))
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

	configFunctioalAccount := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_functional_account" "functional_account" {
			platform_id =           1
			domain_name =           "corp.example.com"
			account_name =          "svc-monitoring"
			display_name =          "FUNCTIONAL_ACCOUNT"
			password =              "P@ssw0rd123!"
		}`,
	}

	configFunctioalAccount.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
		},
		Steps: []resource.TestStep{
			{
				// test using oauth authentication
				Config:      utils.TestResourceConfig(configFunctioalAccount),
				ExpectError: regexp.MustCompile("error - status code: 400"),
			},
		},
	})
}
