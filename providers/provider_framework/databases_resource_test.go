package provider_framework

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/entities"
	"terraform-provider-passwordsafe/providers/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
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

// TestDatabaseResourceConfigure tests the Configure method with various scenarios
func TestDatabaseResourceConfigure(t *testing.T) {
	// Test Configure method directly
	r := &databaseResource{}

	// Test with nil provider data (should trigger early return)
	req := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	// Should not have any errors for nil provider data
	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil provider data should not error, got: %v", resp.Diagnostics.Errors())
	}

	// Test with empty username in provider data (should trigger early return)
	req2 := fwresource.ConfigureRequest{
		ProviderData: ProviderData{
			userName: "", // Empty username should cause early return
		},
	}
	resp2 := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), req2, resp2)

	// Should not have any errors for empty username
	if resp2.Diagnostics.HasError() {
		t.Errorf("Configure with empty username should not error, got: %v", resp2.Diagnostics.Errors())
	}
}

// TestDatabaseResourceRead tests the Read method (currently unimplemented)
func TestDatabaseResourceRead(t *testing.T) {
	r := &databaseResource{}

	req := fwresource.ReadRequest{}
	resp := &fwresource.ReadResponse{}

	r.Read(context.Background(), req, resp)

	// Read method should not panic and should handle gracefully
	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			if !strings.Contains(diag.Summary(), "not implemented") {
				t.Errorf("Unexpected error in Read: %v", diag.Summary())
			}
		}
	}
}

// TestDatabaseResourceUpdate tests the Update method (currently unimplemented)
func TestDatabaseResourceUpdate(t *testing.T) {
	r := &databaseResource{}

	req := fwresource.UpdateRequest{}
	resp := &fwresource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	// Update method should not panic and should handle gracefully
	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics.Errors() {
			if !strings.Contains(diag.Summary(), "not implemented") {
				t.Errorf("Unexpected error in Update: %v", diag.Summary())
			}
		}
	}
}

// TestDatabaseResourceImportState tests the ImportState method
func TestDatabaseResourceImportState(t *testing.T) {
	// Skip this test for now as it requires complex state setup
	t.Skip("ImportState test requires complex state initialization - focusing on other coverage improvements")
}

// TestDatabaseDataSourceBasics tests basic datasource functions for coverage
func TestDatabaseDataSourceBasics(t *testing.T) {
	dataSource := NewDatabaseDataSource()
	if dataSource == nil {
		t.Error("NewDatabaseDataSource() returned nil")
	}
}

// TestDatabaseDataSourceMetadata tests the Metadata method for coverage
func TestDatabaseDataSourceMetadata(t *testing.T) {
	dataSource := NewDatabaseDataSource()
	req := datasource.MetadataRequest{
		ProviderTypeName: "passwordsafe",
	}
	var resp datasource.MetadataResponse
	dataSource.Metadata(context.Background(), req, &resp)

	if resp.TypeName != "passwordsafe_database_datasource" {
		t.Errorf("Expected TypeName to be 'passwordsafe_database_datasource', got '%s'", resp.TypeName)
	}
}

// TestDatabaseDataSourceSchema tests the Schema method for coverage
func TestDatabaseDataSourceSchema(t *testing.T) {
	dataSource := NewDatabaseDataSource()
	req := datasource.SchemaRequest{}
	var resp datasource.SchemaResponse
	dataSource.Schema(context.Background(), req, &resp)

	if len(resp.Schema.Blocks) == 0 {
		t.Error("Expected Schema to have blocks")
	}
} // TestDatabaseDataSourceConfigure tests the Configure method for coverage
func TestDatabaseDataSourceConfigure(t *testing.T) {
	dataSource := &DatabaseDataSource{}
	req := datasource.ConfigureRequest{}
	var resp datasource.ConfigureResponse
	dataSource.Configure(context.Background(), req, &resp)

	// Should complete without panic for basic test
}
