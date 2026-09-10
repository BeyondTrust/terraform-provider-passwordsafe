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

var ManganedAccountEphemeralOauthConfig entities.PasswordSafeTestConfig = entities.PasswordSafeTestConfig{
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	URL:                          "",
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
	ephemeral "passwordsafe_managed_acccount_ephemeral" "test" {
	system_name = "server01"
	account_name = "managed_account_01"
	}

	provider "echo" {
	data = ephemeral.passwordsafe_managed_acccount_ephemeral.test
	}

	resource "echo" "test" {}`,
}

func TestEphemeralManagedAcount(t *testing.T) {

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

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/ManagedAccounts":
			_, err := w.Write([]byte(`{"SystemId":1,"AccountId":10}`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Requests":
			_, err := w.Write([]byte(`124`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Credentials/124":
			_, err := w.Write([]byte(`"fake_credential"`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/Requests/124/checkin":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	server.URL = server.URL + constants.APIPath
	ManganedAccountEphemeralOauthConfig.URL = server.URL

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
				// test using oauth authentication
				Config: utils.TestResourceConfig(ManganedAccountEphemeralOauthConfig),
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

// ManganedAccountEphemeralSeparatorConfig uses a system name containing the "/"
// character, which used to be split as a path separator and rejected with
// "empty managed account list" (BIPS-37662).
var ManganedAccountEphemeralSeparatorConfig entities.PasswordSafeTestConfig = entities.PasswordSafeTestConfig{
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	URL:                          "",
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
	ephemeral "passwordsafe_managed_acccount_ephemeral" "test" {
	system_name = "Accounts - AD/EntraID"
	account_name = "managed_account_01"
	}

	provider "echo" {
	data = ephemeral.passwordsafe_managed_acccount_ephemeral.test
	}

	resource "echo" "test" {}`,
}

func TestEphemeralManagedAcountSystemNameWithSeparator(t *testing.T) {

	var gotSystemName string

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err = w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))

		case constants.APIPath + "/Auth/SignAppIn":
			_, err = w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

		case constants.APIPath + "/Auth/Signout":
			_, err = w.Write([]byte(``))

		case constants.APIPath + "/ManagedAccounts":
			gotSystemName = r.URL.Query().Get("systemName")
			_, err = w.Write([]byte(`{"SystemId":1,"AccountId":10}`))

		case constants.APIPath + "/Requests":
			_, err = w.Write([]byte(`124`))

		case constants.APIPath + "/Credentials/124":
			_, err = w.Write([]byte(`"fake_credential"`))

		case constants.APIPath + "/Requests/124/checkin":
			_, err = w.Write([]byte(``))
		}

		if err != nil {
			t.Error(err.Error())
		}

	}))

	server.URL = server.URL + constants.APIPath
	ManganedAccountEphemeralSeparatorConfig.URL = server.URL

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
				Config: utils.TestResourceConfig(ManganedAccountEphemeralSeparatorConfig),
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

	if gotSystemName != "Accounts - AD/EntraID" {
		t.Errorf("expected systemName query %q, got %q", "Accounts - AD/EntraID", gotSystemName)
	}
}

func TestEphemeralManagedAcountNotFound(t *testing.T) {

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

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/ManagedAccounts":
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`"Managed Account not found"`))
			if err != nil {
				t.Error(err.Error())
			}
		}
	}))

	server.URL = server.URL + constants.APIPath
	ManganedAccountEphemeralOauthConfig.URL = server.URL

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
				Config:      utils.TestResourceConfig(ManganedAccountEphemeralOauthConfig),
				ExpectError: regexp.MustCompile("Managed Account not found"),
			},
		},
	})
}
