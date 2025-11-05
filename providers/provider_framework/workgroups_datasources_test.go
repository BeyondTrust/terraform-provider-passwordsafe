package provider_framework

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/entities"
	"terraform-provider-passwordsafe/providers/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var workgroupsListConfig = entities.PasswordSafeTestConfig{
	APIKey:                       "",
	ClientID:                     constants.FakeClientId,
	ClientSecret:                 constants.FakeClientSecret,
	APIAccountName:               "",
	ClientCertificatesFolderPath: "",
	ClientCertificateName:        "",
	ClientCertificatePassword:    "",
	APIVersion:                   "3.1",
	Resource: `
		data "passwordsafe_workgroup_datasource" "workgroups_list" {
		}`,
}

func TestGetWorkgroupsList(t *testing.T) {

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

		case constants.APIPath + "/Workgroups":
			_, err := w.Write([]byte(`[{ "OrganizationID": "abcd27cf-791a-4c65-abe9-a6a250b8e4f6", "ID": 1, "Name": "Default Workgroup" }, { "OrganizationID": "abcd27cf-791a-4c65-abe9-a6a250b8e4f6", "ID": 2, "Name": "Name 1" }]`))
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

	workgroupsListConfig.URL = server.URL

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
				// test using oauth authentication, get workgroups list
				Config: utils.TestResourceConfig(workgroupsListConfig),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.passwordsafe_workgroup_datasource.workgroups_list",
						tfjsonpath.New("workgroups"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}

func TestGetWorkgroupsListNotFound(t *testing.T) {

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

		case constants.APIPath + "/Workgroups":
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

	workgroupsListConfig.URL = server.URL

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
				Config:      utils.TestResourceConfig(workgroupsListConfig),
				ExpectError: regexp.MustCompile("Error getting workgroups list"),
			},
		},
	})
}

// Unit tests for workgroup datasource methods
func TestWorkgroupDataSourceBasics(t *testing.T) {
	ds := NewWorkgroupDataSource()

	if ds == nil {
		t.Error("NewWorkgroupDataSource returned nil")
	}
}

func TestWorkgroupDataSourceMetadata(t *testing.T) {
	ds := NewWorkgroupDataSource()
	req := datasource.MetadataRequest{
		ProviderTypeName: "passwordsafe",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "passwordsafe_workgroup_datasource" {
		t.Errorf("Expected TypeName 'passwordsafe_workgroup_datasource', got '%s'", resp.TypeName)
	}
}

func TestWorkgroupDataSourceSchema(t *testing.T) {
	ds := NewWorkgroupDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	if resp.Schema.Description != "Workgroup Datasource, gets workgroups list" {
		t.Errorf("Unexpected schema description: %s", resp.Schema.Description)
	}
}
