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

var functionalAccountsListConfig = entities.PasswordSafeTestConfig{
	APIKey:                       "",
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
		data "passwordsafe_functional_account_datasource" "functional_accounts_list" {
		}`,
}

func TestGetFunctionalAccountsList(t *testing.T) {

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
			_, err := w.Write([]byte(`[ { "FunctionalAccountID": 1, "PlatformID": 4 }, { "FunctionalAccountID": 2, "PlatformID": 3 }, { "FunctionalAccountID": 3, "PlatformID": 47 }, { "FunctionalAccountID": 4, "PlatformID": 1 }, { "FunctionalAccountID": 5, "PlatformID": 1 }, { "FunctionalAccountID": 6, "PlatformID": 1 } ]`))
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

	functionalAccountsListConfig.URL = server.URL

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
				// test using oauth authentication, get functional accounts list
				Config: utils.TestResourceConfig(functionalAccountsListConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.passwordsafe_functional_account_datasource.functional_accounts_list",
						tfjsonpath.New("accounts"),
						knownvalue.ListSizeExact(6),
					),
				},
			},
		},
	})
}

func TestGetFunctionalAccountsListNotFound(t *testing.T) {

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
			// not found mock
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(``))
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

	functionalAccountsListConfig.URL = server.URL

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
				Config:      utils.TestResourceConfig(functionalAccountsListConfig),
				ExpectError: regexp.MustCompile("Error getting functional accounts list"),
			},
		},
	})
}
