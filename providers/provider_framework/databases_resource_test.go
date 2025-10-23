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

func TestCreateDatabase(t *testing.T) {

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

		case constants.APIPath + "/Assets/25/Databases":
			_, err := w.Write([]byte(`{ "DatabaseID": 1001, "AssetID": 25, "PlatformID": 10, "InstanceName": "SQLInstance10ss", "IsDefaultInstance": false, "Port": 1433, "Version": "15.0", "Template": "StandardTemplate" }`))
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

	configDatabase := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_database" "database" {
			asset_id            = "25"
			platform_id         = 1
			instance_name       = "primary-db-instance"
			is_default_instance = true
			port               = 5432
			version            = "13.3"
			template           = "standard-template"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	configDatabase.URL = server.URL

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
				// test using oauth authentication, creating database.
				Config: utils.TestResourceConfig(configDatabase),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_database.database",
						tfjsonpath.New("database_id"),
						knownvalue.Int32Exact(1001),
					),
				},
			},
		},
	})
}

func TestCreateDatabaseBadRequest(t *testing.T) {

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

		case constants.APIPath + "/Assets/30/Databases":
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

	configDatabase := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_database" "database" {
			asset_id            = "30"
			platform_id         = 10
			instance_name       = "primary-db-instance-error"
			is_default_instance = true
			port               = 5432
			version            = "13.3"
			template           = "standard-template"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	configDatabase.URL = server.URL

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
				Config:      utils.TestResourceConfig(configDatabase),
				ExpectError: regexp.MustCompile("error - status code: 400"),
			},
		},
	})
}

func TestDeleteDatabase(t *testing.T) {

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

		case constants.APIPath + "/Assets/25/Databases":
			if r.Method == http.MethodPost {
				// Create database response
				_, err := w.Write([]byte(`{ "DatabaseID": 1001, "AssetID": 25, "PlatformID": 10, "InstanceName": "SQLInstance10ss", "IsDefaultInstance": false, "Port": 1433, "Version": "15.0", "Template": "StandardTemplate" }`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Databases/1001":
			if r.Method == http.MethodDelete {
				// Delete database response - success
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	configDatabase := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_database" "database" {
			asset_id            = "25"
			platform_id         = 10
			instance_name       = "primary-db-instance"
			is_default_instance = false
			port               = 1433
			version            = "15.0"
			template           = "StandardTemplate"
		}`,
	}

	server.URL = server.URL + constants.APIPath
	configDatabase.URL = server.URL

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
				// Create database
				Config: utils.TestResourceConfig(configDatabase),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_database.database",
						tfjsonpath.New("database_id"),
						knownvalue.Int32Exact(1001),
					),
				},
			},
			{
				// Delete database by removing from config
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configDatabase.APIKey,
					ClientID:                     configDatabase.ClientID,
					ClientSecret:                 configDatabase.ClientSecret,
					APIAccountName:               configDatabase.APIAccountName,
					ClientCertificatesFolderPath: configDatabase.ClientCertificatesFolderPath,
					ClientCertificateName:        configDatabase.ClientCertificateName,
					ClientCertificatePassword:    configDatabase.ClientCertificatePassword,
					APIVersion:                   configDatabase.APIVersion,
					URL:                          configDatabase.URL,
					Resource:                     ``, // Empty resource to trigger delete
				}),
			},
		},
	})
}

func TestDeleteDatabaseNotFound(t *testing.T) {

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

		case constants.APIPath + "/Assets/25/Databases":
			if r.Method == http.MethodPost {
				// Create database response
				_, err := w.Write([]byte(`{ "DatabaseID": 9999, "AssetID": 25, "PlatformID": 10, "InstanceName": "SQLInstance10ss", "IsDefaultInstance": false, "Port": 1433, "Version": "15.0", "Template": "StandardTemplate" }`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Databases/9999":
			if r.Method == http.MethodDelete {
				// Delete database response - not found
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"error": "Database not found"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	configDatabase := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_database" "database" {
			asset_id            = "25"
			platform_id         = 10
			instance_name       = "primary-db-instance"
			is_default_instance = false
			port               = 1433
			version            = "15.0"
			template           = "StandardTemplate"
		}`,
	}

	server.URL = server.URL + constants.APIPath
	configDatabase.URL = server.URL

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
				// Create database
				Config: utils.TestResourceConfig(configDatabase),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_database.database",
						tfjsonpath.New("database_id"),
						knownvalue.Int32Exact(9999),
					),
				},
			},
			{
				// Delete database by removing from config - expect error due to 404
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configDatabase.APIKey,
					ClientID:                     configDatabase.ClientID,
					ClientSecret:                 configDatabase.ClientSecret,
					APIAccountName:               configDatabase.APIAccountName,
					ClientCertificatesFolderPath: configDatabase.ClientCertificatesFolderPath,
					ClientCertificateName:        configDatabase.ClientCertificateName,
					ClientCertificatePassword:    configDatabase.ClientCertificatePassword,
					APIVersion:                   configDatabase.APIVersion,
					URL:                          configDatabase.URL,
					Resource:                     ``, // Empty resource to trigger delete
				}),
				ExpectError: regexp.MustCompile("error - status code: 404"),
			},
		},
	})
}

func TestDeleteDatabaseServerError(t *testing.T) {

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

		case constants.APIPath + "/Assets/25/Databases":
			if r.Method == http.MethodPost {
				// Create database response
				_, err := w.Write([]byte(`{ "DatabaseID": 5000, "AssetID": 25, "PlatformID": 10, "InstanceName": "SQLInstance10ss", "IsDefaultInstance": false, "Port": 1433, "Version": "15.0", "Template": "StandardTemplate" }`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Databases/5000":
			if r.Method == http.MethodDelete {
				// Delete database response - server error
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error": "Internal server error"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	configDatabase := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_database" "database" {
			asset_id            = "25"
			platform_id         = 10
			instance_name       = "primary-db-instance"
			is_default_instance = false
			port               = 1433
			version            = "15.0"
			template           = "StandardTemplate"
		}`,
	}

	server.URL = server.URL + constants.APIPath
	configDatabase.URL = server.URL

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
				// Create database
				Config: utils.TestResourceConfig(configDatabase),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_database.database",
						tfjsonpath.New("database_id"),
						knownvalue.Int32Exact(5000),
					),
				},
			},
			{
				// Delete database by removing from config - expect error due to server error
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configDatabase.APIKey,
					ClientID:                     configDatabase.ClientID,
					ClientSecret:                 configDatabase.ClientSecret,
					APIAccountName:               configDatabase.APIAccountName,
					ClientCertificatesFolderPath: configDatabase.ClientCertificatesFolderPath,
					ClientCertificateName:        configDatabase.ClientCertificateName,
					ClientCertificatePassword:    configDatabase.ClientCertificatePassword,
					APIVersion:                   configDatabase.APIVersion,
					URL:                          configDatabase.URL,
					Resource:                     ``, // Empty resource to trigger delete
				}),
				ExpectError: regexp.MustCompile("error - status code: 500"),
			},
		},
	})
}
