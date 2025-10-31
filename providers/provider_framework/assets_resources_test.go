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

func TestCreateAsset(t *testing.T) {

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

		case constants.APIPath + "/workgroups/work_group_name/assets":
			_, err := w.Write([]byte(`{ "WorkgroupID": 1, "AssetID": 36, "AssetName": "Asset created by Workgroup Name", "AssetType": "Server", "DnsName": "server01.local", "DomainName": "test.com", "IPAddress": "192.168.1.1", "OperatingSystem": "Ubuntu 22.04", "CreateDate": "2025-02-27T22:57:27.127Z", "LastUpdateDate": "2025-02-27T22:57:27.127Z", "Description": "Primary application server" }`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/workgroups/20/assets":
			_, err := w.Write([]byte(`{ "WorkgroupID": 1, "AssetID": 36, "AssetName": "Asset created by Workgroup Id", "AssetType": "Server", "DnsName": "server01.local", "DomainName": "test.com", "IPAddress": "192.168.1.1", "OperatingSystem": "Ubuntu 22.04", "CreateDate": "2025-02-27T22:57:27.127Z", "LastUpdateDate": "2025-02-27T22:57:27.127Z", "Description": "Primary application server" }`))
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

	configAssetByWorkGroupName := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_asset_by_workgroup_name" "asset" {
 			work_group_name= "work_group_name"
			ip_address = "192.168.1.1"
			asset_name = "Asset created by Workgroup Name"
			dns_name = "server01.local"
			domain_name = "test.com"
			asset_type = "Server"
			description = "Primary application server"
			operating_system = "Ubuntu 22.04"
		}`,
	}

	configAssetByWorkGroupId := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_asset_by_workgroup_id" "asset" {
 			work_group_id= "20"
			ip_address = "192.168.1.1"
			asset_name = "Asset created by Workgroup Id"
			dns_name = "server01.local"
			domain_name = "test.com"
			asset_type = "Server"
			description = "Primary application server"
			operating_system = "Ubuntu 22.04"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	configAssetByWorkGroupName.URL = server.URL
	configAssetByWorkGroupId.URL = server.URL

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
				// test using oauth authentication, creating asset by workgroup name.
				Config: utils.TestResourceConfig(configAssetByWorkGroupName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_asset_by_workgroup_name.asset",
						tfjsonpath.New("asset_type"),
						knownvalue.StringExact("Server"),
					),
				},
			},

			{
				// test using oauth authentication, creating asset by workgroup Id.
				Config: utils.TestResourceConfig(configAssetByWorkGroupId),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_asset_by_workgroup_id.asset",
						tfjsonpath.New("asset_type"),
						knownvalue.StringExact("Server"),
					),
				},
			},
		},
	})
}

func TestCreateAssetBadRequest(t *testing.T) {

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

		case constants.APIPath + "/workgroups/work_group_name/assets":
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"Bad request"}`))
			if err != nil {
				t.Error(err.Error())
			}

		case constants.APIPath + "/workgroups/20/assets":
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

	configAssetByWorkGroupIdError := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_asset_by_workgroup_id" "asset" {
			work_group_id= "20"
			ip_address = "192.168.1.1"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	configAssetByWorkGroupIdError.URL = server.URL

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
				Config:      utils.TestResourceConfig(configAssetByWorkGroupIdError),
				ExpectError: regexp.MustCompile("error - status code: 400"),
			},
		},
	})
}

func TestCreateAssetInvalidCredentials(t *testing.T) {

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"Invalid Credentials"}`))
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
		}

	}))

	configAssetByWorkGroupNameError := entities.PasswordSafeTestConfig{
		APIKey:                       "",
		ClientID:                     constants.FakeClientId,
		ClientSecret:                 constants.FakeClientSecret,
		APIAccountName:               "",
		ClientCertificatesFolderPath: "",
		ClientCertificateName:        "",
		ClientCertificatePassword:    "",
		APIVersion:                   "3.1",
		Resource: `
		resource "passwordsafe_asset_by_workgroup_name" "asset" {
			work_group_name= "work_group_name"
			ip_address = "192.168.1.1"
		}`,
	}

	server.URL = server.URL + constants.APIPath

	configAssetByWorkGroupNameError.URL = server.URL

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
				Config:      utils.TestResourceConfig(configAssetByWorkGroupNameError),
				ExpectError: regexp.MustCompile("error - status code: 400"),
			},
		},
	})
}

func TestDeleteAssetByWorkGroupId(t *testing.T) {
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

		case constants.APIPath + "/workgroups/20/assets":
			if r.Method == http.MethodPost {
				_, err := w.Write([]byte(`{"WorkgroupID": 20, "AssetID": 36, "AssetName": "Test Asset for Deletion", "AssetType": "Server", "DnsName": "test-server.local", "DomainName": "test.com", "IPAddress": "192.168.1.100", "OperatingSystem": "Ubuntu 22.04", "CreateDate": "2025-02-27T22:57:27.127Z", "LastUpdateDate": "2025-02-27T22:57:27.127Z", "Description": "Test asset for deletion by workgroup ID"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Assets/36":
			if r.Method == http.MethodDelete {
				// DELETE endpoint for asset
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

	configAsset := entities.PasswordSafeTestConfig{
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
		resource "passwordsafe_asset_by_workgroup_id" "test_asset_by_workgroup_id" {
			work_group_id    = "20"
			ip_address       = "192.168.1.100"
			asset_name       = "Test Asset for Deletion"
			dns_name         = "test-server.local"
			domain_name      = "test.com"
			asset_type       = "Server"
			description      = "Test asset for deletion by workgroup ID"
			operating_system = "Ubuntu 22.04"
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
				// Create asset
				Config: utils.TestResourceConfig(configAsset),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_asset_by_workgroup_id.test_asset_by_workgroup_id",
						tfjsonpath.New("asset_id"),
						knownvalue.Int32Exact(36),
					),
				},
			},
			{
				// Delete asset by removing from config
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configAsset.APIKey,
					ClientID:                     configAsset.ClientID,
					ClientSecret:                 configAsset.ClientSecret,
					APIAccountName:               configAsset.APIAccountName,
					ClientCertificatesFolderPath: configAsset.ClientCertificatesFolderPath,
					ClientCertificateName:        configAsset.ClientCertificateName,
					ClientCertificatePassword:    configAsset.ClientCertificatePassword,
					APIVersion:                   configAsset.APIVersion,
					URL:                          configAsset.URL,
					Resource:                     "", // Empty resource to trigger deletion
				}),
			},
		},
	})
}

func TestDeleteAssetByWorkGroupName(t *testing.T) {
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

		case constants.APIPath + "/workgroups/test_workgroup/assets":
			if r.Method == http.MethodPost {
				_, err := w.Write([]byte(`{"WorkgroupID": 21, "AssetID": 37, "AssetName": "Test Asset for Name Deletion", "AssetType": "Server", "DnsName": "test-name-server.local", "DomainName": "test.com", "IPAddress": "192.168.1.101", "OperatingSystem": "Ubuntu 22.04", "CreateDate": "2025-02-27T22:57:27.127Z", "LastUpdateDate": "2025-02-27T22:57:27.127Z", "Description": "Test asset for deletion by workgroup name"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}

		case constants.APIPath + "/Assets/37":
			if r.Method == http.MethodDelete {
				// DELETE endpoint for asset
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

	configAsset := entities.PasswordSafeTestConfig{
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
		resource "passwordsafe_asset_by_workgroup_name" "test_asset_by_workgroup_name" {
			work_group_name  = "test_workgroup"
			ip_address       = "192.168.1.101"
			asset_name       = "Test Asset for Name Deletion"
			dns_name         = "test-name-server.local"
			domain_name      = "test.com"
			asset_type       = "Server"
			description      = "Test asset for deletion by workgroup name"
			operating_system = "Ubuntu 22.04"
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
				// Create asset
				Config: utils.TestResourceConfig(configAsset),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"passwordsafe_asset_by_workgroup_name.test_asset_by_workgroup_name",
						tfjsonpath.New("asset_id"),
						knownvalue.Int32Exact(37),
					),
				},
			},
			{
				// Delete asset by removing from config
				Config: utils.TestResourceConfig(entities.PasswordSafeTestConfig{
					APIKey:                       configAsset.APIKey,
					ClientID:                     configAsset.ClientID,
					ClientSecret:                 configAsset.ClientSecret,
					APIAccountName:               configAsset.APIAccountName,
					ClientCertificatesFolderPath: configAsset.ClientCertificatesFolderPath,
					ClientCertificateName:        configAsset.ClientCertificateName,
					ClientCertificatePassword:    configAsset.ClientCertificatePassword,
					APIVersion:                   configAsset.APIVersion,
					URL:                          configAsset.URL,
					Resource:                     "", // Empty resource to trigger deletion
				}),
			},
		},
	})
}
