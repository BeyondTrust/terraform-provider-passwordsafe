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

func TestCreateManagedSystemByWorkGroup(t *testing.T) {

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

		case constants.APIPath + "/Workgroups/5/ManagedSystems":
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
		resource "passwordsafe_managed_system_by_workgroup" "system_by_workgroup" {
 			workgroup_id                           = "5"
			entity_type_id                         = 1
			host_name                              = "example-host"
			ip_address                             = "192.168.1.1"
			dns_name                               = "example.local"
			instance_name                          = "example-instance"
			is_default_instance                    = true
			template                               = "example-template"
			forest_name                            = "example-forest"
			use_ssl                                = false
			platform_id                            = 2
			netbios_name                           = "EXAMPLE"
			contact_email                          = "admin@example.com"
			description                            = "Primary database system from terraform 2025"
			port                                   = 5432
			timeout                                = 30
			ssh_key_enforcement_mode               = 0
			password_rule_id                       = 0
			dss_key_rule_id                        = 0
			login_account_id                       = 0
			account_name_format                    = 1
			oracle_internet_directory_id           = "example-dir-id"
			oracle_internet_directory_service_name = "example-service"
			release_duration                       = 60
			max_release_duration                   = 120
			isa_release_duration                   = 30
			auto_management_flag                   = false
			functional_account_id                  = 0
			elevation_command                      = "sudo su -"
			check_password_flag                    = true
			change_password_after_any_release_flag = false
			reset_password_on_mismatch_flag        = true
			change_frequency_type                  = "last"
			change_frequency_days                  = 7
			change_time                            = "02:00"
			access_url                             = "https://example.com"
			remote_client_type                     = "ssh"
			application_host_id                    = 5001
			is_application_host                    = false
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
						"passwordsafe_managed_system_by_workgroup.system_by_workgroup",
						tfjsonpath.New("managed_system_id"),
						knownvalue.Int32Exact(13),
					),
				},
			},
		},
	})
}

// The argument "platform_id" is required, but no definition was found.
func TestCreateManagedByWorkGroupAccountBadData(t *testing.T) {

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

		case constants.APIPath + "/Workgroups/5/ManagedSystems":
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
		resource "passwordsafe_managed_system_by_workgroup" "system_by_workgroup" {
 			workgroup_id                           = "55"
			entity_type_id                         = 1
			host_name                              = "example-host"
			ip_address                             = "192.168.1.1"
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
				ExpectError: regexp.MustCompile("The argument \"platform_id\" is required, but no definition was found."),
			},
		},
	})
}

func TestDeleteManagedSystemByWorkGroup(t *testing.T) {
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

		case constants.APIPath + "/Workgroups/5/ManagedSystems":
			if r.Method == http.MethodPost {
				_, err := w.Write([]byte(`{"ManagedSystemID": 13, "EntityTypeID": 1, "AssetID": 1, "SystemName": "test-workgroup-system"}`))
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
		resource "passwordsafe_managed_system_by_workgroup" "test_system_by_workgroup" {
			workgroup_id                           = "5"
			entity_type_id                         = 1
			host_name                              = "test-host"
			platform_id                            = 2
			contact_email                          = "admin@example.com"
			description                            = "Test managed system for workgroup deletion"
			port                                   = 5432
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
						"passwordsafe_managed_system_by_workgroup.test_system_by_workgroup",
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
