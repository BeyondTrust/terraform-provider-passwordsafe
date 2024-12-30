package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	p := Provider()
	assert.NotNil(t, p, "Provider should not be nil")

	err := p.InternalValidate()
	assert.NoError(t, err, "Provider validation should not return an error")
}

func TestProviderConfigureWithApiKey(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             "https://example.com",
		"api_account_name":                "test-account",
		"client_id":                       "test-client-id",
		"client_secret":                   "test-client-secret",
		"api_key":                         "test-api-key",
		"verify_ca":                       true,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	// Call the function and check results
	_, diags := providerConfigure(context.Background(), resourceData)

	assert.Empty(t, diags, "Diagnostics should be empty if no errors")
}

func TestProviderConfigureWithCredentials(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             "https://example.com",
		"api_account_name":                "test-account",
		"client_id":                       "test-client-id",
		"client_secret":                   "test-client-secret",
		"api_key":                         "",
		"verify_ca":                       true,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	// Call the function and check results
	_, diags := providerConfigure(context.Background(), resourceData)

	assert.Empty(t, diags, "Diagnostics should be empty if no errors")
}

func TestProviderConfigureEmtyUrl(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             "",
		"api_account_name":                "test-account",
		"client_id":                       "test-client-id",
		"client_secret":                   "test-client-secret",
		"api_key":                         "test-api-key",
		"verify_ca":                       true,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	autenticate, diags := providerConfigure(context.Background(), resourceData)

	if autenticate != nil {
		t.Errorf("Error %v", diags)
	}

	if diags[0].Detail != "Please add a proper URL" {
		t.Errorf("Test case Failed %v, %v", diags[0].Detail, "Please add a proper URL")
	}

}

func TestProviderConfigureEmptyCredentials(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             "https://example.com",
		"api_account_name":                "",
		"client_id":                       "",
		"client_secret":                   "",
		"api_key":                         "",
		"verify_ca":                       true,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	autenticate, diags := providerConfigure(context.Background(), resourceData)

	if autenticate != nil {
		t.Errorf("Error %v", diags)
	}

	if diags[0].Detail != "Please add a valid credential (API Key / Client Credentials)" {
		t.Errorf("Test case Failed %v, %v", diags[0].Detail, "Please add a valid credential (API Key / Client Credentials)")
	}

}

func TestProviderConfigureEmptyAccountName(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             "https://example.com",
		"api_account_name":                "",
		"client_id":                       "",
		"client_secret":                   "",
		"api_key":                         "fake_api_key",
		"verify_ca":                       true,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	autenticate, diags := providerConfigure(context.Background(), resourceData)

	if autenticate != nil {
		t.Errorf("Error %v", diags)
	}

	if diags[0].Detail != "Please add a proper Account Name" {
		t.Errorf("Test case Failed %v, %v", diags[0].Detail, "Please add a proper Account Name")
	}

}
