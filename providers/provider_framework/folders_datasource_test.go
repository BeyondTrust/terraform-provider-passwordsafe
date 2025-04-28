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

var foldersListConfig = entities.PasswordSafeTestConfig{
	APIKey:                       "",
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
		data "passwordsafe_folder_datasource" "folders_list" {
		}`,
}

func TestSecretGetFoldersList(t *testing.T) {

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

		case constants.APIPath + "/secrets-safe/folders/":
			_, err := w.Write([]byte(`[{"Id": "cb871861-8b40-4556-820c-1ca6d522adfa","Name": "folder1"}, {"Id": "a4af73dc-4e89-41ec-eb9a-08dcf22d3aba","Name": "folder2"}]`))
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

	foldersListConfig.URL = server.URL

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
				// test using oauth authentication, get folders list
				Config: utils.TestResourceConfig(foldersListConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.passwordsafe_folder_datasource.folders_list",
						tfjsonpath.New("folders"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}

func TestSecretGetFoldersListNotFound(t *testing.T) {

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

		case constants.APIPath + "/secrets-safe/folders/":
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

	foldersListConfig.URL = server.URL

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
				Config:      utils.TestResourceConfig(foldersListConfig),
				ExpectError: regexp.MustCompile("Error getting folders list"),
			},
		},
	})
}
