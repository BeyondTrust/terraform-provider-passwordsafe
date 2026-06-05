package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	p := Provider()
	assert.NotNil(t, p, "Provider should not be nil")

	err := p.InternalValidate()
	assert.NoError(t, err, "Provider validation should not return an error")
}

// newSignInMockServer returns a mock Password Safe server that handles the
// OAuth token exchange and SignAppIn. Tests that drive providerConfigure end
// to end need this because the new InitSharedAuth performs the handshake
// during Configure (the old providerConfigure only built the authObj).
func newSignInMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, _ = w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
		case constants.APIPath + "/Auth/SignAppIn":
			_, _ = w.Write([]byte(`{"UserId":1, "UserName":"test", "EmailAddress":"test@beyondtrust.com"}`))
		}
	}))
}

func TestProviderConfigureWithApiKey(t *testing.T) {
	server := newSignInMockServer(t)
	defer server.Close()
	utils.ResetSharedAuthForTest()

	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             server.URL + constants.APIPath,
		"api_account_name":                "test-account",
		"client_id":                       "",
		"client_secret":                   "",
		"api_key":                         "test-api-key",
		"verify_ca":                       false,
		"client_certificate_name":         "",
		"client_certificates_folder_path": "",
		"client_certificate_password":     "",
	})

	// Call the function and check results
	_, diags := providerConfigure(context.Background(), resourceData)

	assert.Empty(t, diags, "Diagnostics should be empty if no errors")
}

func TestProviderConfigureWithCredentials(t *testing.T) {
	server := newSignInMockServer(t)
	defer server.Close()
	utils.ResetSharedAuthForTest()

	resourceData := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"url":                             server.URL + constants.APIPath,
		"api_account_name":                "test-account",
		"client_id":                       "00000000-0000-0000-0000-000000000001",
		"client_secret":                   "00000000-0000-0000-0000-000000000002",
		"api_key":                         "",
		"verify_ca":                       false,
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

	authenticate, diags := providerConfigure(context.Background(), resourceData)

	if authenticate != nil {
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

	authenticate, diags := providerConfigure(context.Background(), resourceData)

	if authenticate != nil {
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

	authenticate, diags := providerConfigure(context.Background(), resourceData)

	if authenticate != nil {
		t.Errorf("Error %v", diags)
	}

	if diags[0].Detail != "Please add a proper Account Name" {
		t.Errorf("Test case Failed %v, %v", diags[0].Detail, "Please add a proper Account Name")
	}

}
