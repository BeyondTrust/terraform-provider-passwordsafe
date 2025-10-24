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

func TestCreateManagedSystemByDatabase(t *testing.T) {

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

		case constants.APIPath + "/Databases/2/ManagedSystems":
			_, err := w.Write([]byte(`{"ManagedSystemID": 13, "EntityTypeID": 1, "AssetID": 1}`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	config := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.0",
		Resource: `
		resource "passwordsafe_managed_system_by_database" "system_by_database" {
 			database_id                            = "2"
			contact_email                          = "admin@example.com"
			description                            = "Managed system for example DB"
			timeout                                = 30
			password_rule_id                       = 101
			release_duration                       = 60
			max_release_duration                   = 120
			isa_release_duration                   = 45
			auto_management_flag                   = true
			functional_account_id                  = 1234
			check_password_flag                    = true
			change_password_after_any_release_flag = false
			reset_password_on_mismatch_flag        = true
			change_frequency_type                  = "xdays"
			change_frequency_days                  = 15
			change_time                            = "03:00"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	config.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		// load providers
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
		},
		Steps: []resource.TestStep{
			{
				// test using oauth authentication
				Config: utils.TestResourceConfig(config),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_managed_system_by_database.system_by_database",
						tfjsonpath.New("managed_system_id"),
						knownvalue.Int32Exact(13),
					),
				},
			},
		},
	})
}

// The argument "database_id" is required, but no definition was found.
func TestCreateManagedSystemByDatabaseBadData(t *testing.T) {

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

		case constants.APIPath + "/Databases/5/ManagedSystems":
			_, err := w.Write([]byte(`{"ManagedSystemID": 13, "EntityTypeID": 1, "AssetID": 1}`))
			if err != nil {
				t.Error("Test case Failed")
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	config := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_managed_system_by_database" "system_by_database" {
 			contact_email = "admin@example.com"
  			description   = "Managed system for example DB"
  			timeout       = 30
		}`,
	}

	server.URL = server.URL + constants.APIPath

	config.URL = server.URL

	resource.Test(t, resource.TestCase{

		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		// load providers
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
		},

		Steps: []resource.TestStep{
			{
				// test using oauth authentication
				Config:      utils.TestResourceConfig(config),
				ExpectError: regexp.MustCompile("The argument \"database_id\" is required, but no definition was found."),
			},
		},
	})
}

func TestDeleteManagedSystemByDatabase(t *testing.T) {
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

		case constants.APIPath + "/Databases/2/ManagedSystems":
			if r.Method == http.MethodPost {
				_, err := w.Write([]byte(`{"ManagedSystemID": 13, "EntityTypeID": 1, "AssetID": 1, "SystemName": "test-db-system"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/ManagedSystems/13":
			if r.Method == http.MethodDelete {
				// DELETE endpoint for managed system
				w.WriteHeader(http.StatusOK)
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}
	}))

	server.URL = server.URL + constants.APIPath

	configManagedSystem := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.0",
		URL:                          server.URL,
		Resource: `
		resource "passwordsafe_managed_system_by_database" "test_system_by_database" {
			database_id                            = "2"
			contact_email                          = "admin@example.com"
			description                            = "Test managed system for database deletion"
			timeout                                = 30
			password_rule_id                       = 101
			release_duration                       = 60
			max_release_duration                   = 120
			isa_release_duration                   = 90
			auto_management_flag                   = false
			change_frequency_type                  = "xdays"
			change_frequency_days                  = 30
			change_time                            = "02:00"
		}`,
	}

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
				// Create managed system
				Config: utils.TestResourceConfig(configManagedSystem),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_managed_system_by_database.test_system_by_database",
						tfjsonpath.New("managed_system_id"),
						knownvalue.Int32Exact(13),
					),
				},
			},
			{
				// Delete managed system by removing from config
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configManagedSystem.APIKey,
					ClientID:                     configManagedSystem.ClientID,
					ClientSecret:                 configManagedSystem.ClientSecret,
					APIAccountName:               configManagedSystem.APIAccountName,
					ClientCertificatesFolderPath: configManagedSystem.ClientCertificatesFolderPath,
					ClientCertificateName:        configManagedSystem.ClientCertificateName,
					ClientCertificatePassword:    configManagedSystem.ClientCertificatePassword,
					APIVersion:                   configManagedSystem.APIVersion,
					URL:                          configManagedSystem.URL,
					Resource:                     "", // Empty resource to trigger deletion
				}),
			},
		},
	})
}
