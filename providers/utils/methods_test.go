package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/entities"
	"testing"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"go.uber.org/zap"
)

var authParams *authentication.AuthenticationParametersObj
var apiVersion string = "3.1"
var zapLogger logging.Logger

func TestTestResourceConfig(t *testing.T) {
	config := entities.PasswordSafeTestConfig{
		APIKey:                       "test-api-key",
		ClientID:                     "test-client-id",
		ClientSecret:                 "test-client-secret",
		URL:                          "https://example.com",
		APIAccountName:               "test-account",
		ClientCertificatesFolderPath: "/path/to/certs",
		ClientCertificateName:        "cert.pem",
		ClientCertificatePassword:    "password123",
		APIVersion:                   "v1",
		Resource:                     "resource-content",
	}

	result := TestResourceConfig(config)

	tests := []struct {
		name     string
		expected string
	}{
		{"APIKey", `api_key = "test-api-key"`},
		{"ClientID", `client_id = "test-client-id"`},
		{"ClientSecret", `client_secret = "test-client-secret"`},
		{"URL", `url = "https://example.com"`},
		{"APIAccountName", `api_account_name = "test-account"`},
		{"ClientCertificatesFolderPath", `client_certificates_folder_path = "/path/to/certs"`},
		{"ClientCertificateName", `client_certificate_name = "cert.pem"`},
		{"ClientCertificatePassword", `client_certificate_password = "password123"`},
		{"APIVersion", `api_version = "v1"`},
		{"Resource", `resource-content`},
	}

	for _, tt := range tests {
		if !strings.Contains(result, tt.expected) {
			t.Errorf("Expected output to contain %s, but it was missing", tt.expected)
		}
	}
}

func InitializeGlobalConfig() {

	logger, _ := zap.NewDevelopment()

	zapLogger = logging.NewZapLogger(logger)

	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second

	authParams = &authentication.AuthenticationParametersObj{
		HTTPClient:                 *httpClientObj,
		BackoffDefinition:          backoffDefinition,
		EndpointURL:                constants.FakeApiUrl,
		APIVersion:                 apiVersion,
		ClientID:                   "fakeone_a654+9sdf7+8we4f",
		ClientSecret:               "fakeone_a654+9sdf7+8we4f",
		ApiKey:                     "",
		Logger:                     zapLogger,
		RetryMaxElapsedTimeSeconds: 300,
	}
}

func TestAuthenticate(t *testing.T) {

	InitializeGlobalConfig()

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

		}

	}))

	var authenticateObj, _ = authentication.Authenticate(*authParams)
	server.URL = server.URL + constants.APIPath
	apiUrl, _ := url.Parse(server.URL)

	authenticateObj.ApiUrl = *apiUrl

	var signInCount uint64
	var mu sync.Mutex

	_, err := Authenticate(*authenticateObj, &mu, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}

	// Increment counter
	_, err = Authenticate(*authenticateObj, &mu, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthenticateErrorGettingToken(t *testing.T) {

	InitializeGlobalConfig()

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{"error": "invalid_client"}`))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	var authenticateObj, _ = authentication.Authenticate(*authParams)
	server.URL = server.URL + constants.APIPath
	apiUrl, _ := url.Parse(server.URL)

	authenticateObj.ApiUrl = *apiUrl

	var signInCount uint64
	var mu sync.Mutex

	expectedError := `error - status code: 400 - {"error": "invalid_client"}`

	_, err := Authenticate(*authenticateObj, &mu, &signInCount, zapLogger)
	if err.Error() != expectedError {
		t.Errorf("Test case Failed %v, %v", err.Error(), expectedError)
	}

}

func TestSignOut(t *testing.T) {

	InitializeGlobalConfig()

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {

		case constants.APIPath + "/Auth/Signout":
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	var authenticateObj, _ = authentication.Authenticate(*authParams)
	server.URL = server.URL + constants.APIPath
	apiUrl, _ := url.Parse(server.URL)

	authenticateObj.ApiUrl = *apiUrl

	var signInCount uint64
	var muOut sync.Mutex

	err := SignOut(*authenticateObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}

	// decrement counter, don't signout case
	err = SignOut(*authenticateObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}
}

func TestSignOutError(t *testing.T) {

	InitializeGlobalConfig()

	// mocking Password Safe API
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// mocking Response according to the endpoint path
		switch r.URL.Path {

		case constants.APIPath + "/Auth/Signout":
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(``))
			if err != nil {
				t.Error(err.Error())
			}
		}

	}))

	var authenticateObj, _ = authentication.Authenticate(*authParams)
	server.URL = server.URL + constants.APIPath
	apiUrl, _ := url.Parse(server.URL)

	authenticateObj.ApiUrl = *apiUrl

	var signInCount uint64
	var muOut sync.Mutex

	signInCount = 1

	expectedError := `error - status code: 400 - `

	err := SignOut(*authenticateObj, &muOut, &signInCount, zapLogger)
	if err.Error() != expectedError {
		t.Errorf("Test case Failed %v, %v", err.Error(), expectedError)
	}

}

// Test ValidateChangeFrequencyDays function
func TestValidateChangeFrequencyDays(t *testing.T) {
	tests := []struct {
		name                string
		changeFrequencyType string
		changeFrequencyDays int
		expectError         bool
		expectedError       string
	}{
		{
			name:                "Valid xdays with valid days",
			changeFrequencyType: "xdays",
			changeFrequencyDays: 30,
			expectError:         false,
		},
		{
			name:                "Valid xdays with minimum days",
			changeFrequencyType: "xdays",
			changeFrequencyDays: 1,
			expectError:         false,
		},
		{
			name:                "Valid xdays with maximum days",
			changeFrequencyType: "xdays",
			changeFrequencyDays: 999,
			expectError:         false,
		},
		{
			name:                "Invalid xdays with days too low",
			changeFrequencyType: "xdays",
			changeFrequencyDays: 0,
			expectError:         true,
			expectedError:       "error in change Frequency field, (min=1, max=999)",
		},
		{
			name:                "Invalid xdays with days too high",
			changeFrequencyType: "xdays",
			changeFrequencyDays: 1000,
			expectError:         true,
			expectedError:       "error in change Frequency field, (min=1, max=999)",
		},
		{
			name:                "Non-xdays type with any days value",
			changeFrequencyType: "first",
			changeFrequencyDays: 0,
			expectError:         false,
		},
		{
			name:                "Non-xdays type with high days value",
			changeFrequencyType: "last",
			changeFrequencyDays: 2000,
			expectError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChangeFrequencyDays(tt.changeFrequencyType, tt.changeFrequencyDays)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Error())
				}
			}
		})
	}
}

// Test DeleteManagedSystemByID function
func TestDeleteManagedSystemByID(t *testing.T) {
	InitializeGlobalConfig()

	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			case "/BeyondTrust/api/public/v3/ManagedSystems/123":
				if r.Method == "DELETE" {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(`{}`))
					if err != nil {
						t.Error(err.Error())
					}
				}
			}
		}))
		defer server.Close()

		authenticateObj, _ := authentication.Authenticate(*authParams)
		server.URL = server.URL + constants.APIPath
		apiUrl, _ := url.Parse(server.URL)
		authenticateObj.ApiUrl = *apiUrl

		err := DeleteManagedSystemByID(*authenticateObj, 123, zapLogger)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
	})

	// Test case 2: Error when deleting managed system
	t.Run("Delete error", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			case "/BeyondTrust/api/public/v3/ManagedSystems/123":
				if r.Method == "DELETE" {
					w.WriteHeader(http.StatusBadRequest)
					_, err := w.Write([]byte(`{"error": "not found"}`))
					if err != nil {
						t.Error(err.Error())
					}
				}
			}
		}))
		defer server.Close()

		authenticateObj, _ := authentication.Authenticate(*authParams)
		server.URL = server.URL + constants.APIPath
		apiUrl, _ := url.Parse(server.URL)
		authenticateObj.ApiUrl = *apiUrl

		err := DeleteManagedSystemByID(*authenticateObj, 123, zapLogger)
		if err == nil {
			t.Error("Expected error when deleting managed system, but got none")
		}
	})
}

// Test GetCreateManagedSystemCommonAttributes function
func TestGetCreateManagedSystemCommonAttributes(t *testing.T) {
	attributes := GetCreateManagedSystemCommonAttributes()

	// Test that we get a map with expected attributes
	expectedAttributes := []string{
		"managed_system_id",
		"managed_system_name",
		"contact_email",
		"description",
		"timeout",
		"password_rule_id",
		"release_duration",
		"max_release_duration",
		"isa_release_duration",
		"auto_management_flag",
		"functional_account_id",
		"check_password_flag",
		"change_password_after_any_release_flag",
		"reset_password_on_mismatch_flag",
		"change_frequency_type",
		"change_frequency_days",
		"change_time",
	}

	for _, attrName := range expectedAttributes {
		if _, exists := attributes[attrName]; !exists {
			t.Errorf("Expected attribute '%s' not found", attrName)
		}
	}

	// Test specific attribute properties
	if managedSystemId, ok := attributes["managed_system_id"]; ok {
		if attr, ok := managedSystemId.(schema.Int32Attribute); ok {
			if attr.Required != false || attr.Optional != false || attr.Computed != true {
				t.Error("managed_system_id attribute has incorrect properties")
			}
		} else {
			t.Error("managed_system_id is not Int32Attribute type")
		}
	}

	if contactEmail, ok := attributes["contact_email"]; ok {
		if attr, ok := contactEmail.(schema.StringAttribute); ok {
			if attr.Optional != true {
				t.Error("contact_email attribute should be optional")
			}
		} else {
			t.Error("contact_email is not StringAttribute type")
		}
	}

	if autoMgmtFlag, ok := attributes["auto_management_flag"]; ok {
		if attr, ok := autoMgmtFlag.(schema.BoolAttribute); ok {
			if attr.Optional != true {
				t.Error("auto_management_flag attribute should be optional")
			}
		} else {
			t.Error("auto_management_flag is not BoolAttribute type")
		}
	}
}

// Test GetInt32Attribute function
func TestGetInt32Attribute(t *testing.T) {
	tests := []struct {
		name        string
		description string
		required    bool
		optional    bool
		computed    bool
	}{
		{
			name:        "Required attribute",
			description: "Test required attribute",
			required:    true,
			optional:    false,
			computed:    false,
		},
		{
			name:        "Optional attribute",
			description: "Test optional attribute",
			required:    false,
			optional:    true,
			computed:    false,
		},
		{
			name:        "Computed attribute",
			description: "Test computed attribute",
			required:    false,
			optional:    false,
			computed:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := GetInt32Attribute(tt.description, tt.required, tt.optional, tt.computed)

			if int64Attr, ok := attr.(schema.Int64Attribute); ok {
				if int64Attr.MarkdownDescription != tt.description {
					t.Errorf("Expected description '%s', got '%s'", tt.description, int64Attr.MarkdownDescription)
				}
				if int64Attr.Required != tt.required {
					t.Errorf("Expected required %v, got %v", tt.required, int64Attr.Required)
				}
				if int64Attr.Optional != tt.optional {
					t.Errorf("Expected optional %v, got %v", tt.optional, int64Attr.Optional)
				}
				if int64Attr.Computed != tt.computed {
					t.Errorf("Expected computed %v, got %v", tt.computed, int64Attr.Computed)
				}
			} else {
				t.Error("Expected Int64Attribute type")
			}
		})
	}
}

// Test GetStringAttribute function
func TestGetStringAttribute(t *testing.T) {
	tests := []struct {
		name        string
		description string
		required    bool
		optional    bool
		computed    bool
	}{
		{
			name:        "Required string attribute",
			description: "Test required string attribute",
			required:    true,
			optional:    false,
			computed:    false,
		},
		{
			name:        "Optional string attribute",
			description: "Test optional string attribute",
			required:    false,
			optional:    true,
			computed:    false,
		},
		{
			name:        "Computed string attribute",
			description: "Test computed string attribute",
			required:    false,
			optional:    false,
			computed:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := GetStringAttribute(tt.description, tt.required, tt.optional, tt.computed)

			if stringAttr, ok := attr.(schema.StringAttribute); ok {
				if stringAttr.MarkdownDescription != tt.description {
					t.Errorf("Expected description '%s', got '%s'", tt.description, stringAttr.MarkdownDescription)
				}
				if stringAttr.Required != tt.required {
					t.Errorf("Expected required %v, got %v", tt.required, stringAttr.Required)
				}
				if stringAttr.Optional != tt.optional {
					t.Errorf("Expected optional %v, got %v", tt.optional, stringAttr.Optional)
				}
				if stringAttr.Computed != tt.computed {
					t.Errorf("Expected computed %v, got %v", tt.computed, stringAttr.Computed)
				}
			} else {
				t.Error("Expected StringAttribute type")
			}
		})
	}
}

// Test GetBoolAttribute function
func TestGetBoolAttribute(t *testing.T) {
	tests := []struct {
		name        string
		description string
		required    bool
		optional    bool
		computed    bool
	}{
		{
			name:        "Required bool attribute",
			description: "Test required bool attribute",
			required:    true,
			optional:    false,
			computed:    false,
		},
		{
			name:        "Optional bool attribute",
			description: "Test optional bool attribute",
			required:    false,
			optional:    true,
			computed:    false,
		},
		{
			name:        "Computed bool attribute",
			description: "Test computed bool attribute",
			required:    false,
			optional:    false,
			computed:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := GetBoolAttribute(tt.description, tt.required, tt.optional, tt.computed)

			if boolAttr, ok := attr.(schema.BoolAttribute); ok {
				if boolAttr.MarkdownDescription != tt.description {
					t.Errorf("Expected description '%s', got '%s'", tt.description, boolAttr.MarkdownDescription)
				}
				if boolAttr.Required != tt.required {
					t.Errorf("Expected required %v, got %v", tt.required, boolAttr.Required)
				}
				if boolAttr.Optional != tt.optional {
					t.Errorf("Expected optional %v, got %v", tt.optional, boolAttr.Optional)
				}
				if boolAttr.Computed != tt.computed {
					t.Errorf("Expected computed %v, got %v", tt.computed, boolAttr.Computed)
				}
			} else {
				t.Error("Expected BoolAttribute type")
			}
		})
	}
}

// Test DeleteAssetByID function
func TestDeleteAssetByID(t *testing.T) {
	InitializeGlobalConfig()

	// Test case 1: Successful deletion
	t.Run("Successful deletion", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			case "/BeyondTrust/api/public/v3/Assets/123":
				if r.Method == "DELETE" {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(`{}`))
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
		defer server.Close()

		authenticateObj, _ := authentication.Authenticate(*authParams)
		server.URL = server.URL + constants.APIPath
		apiUrl, _ := url.Parse(server.URL)
		authenticateObj.ApiUrl = *apiUrl

		var signInCount uint64
		var mu sync.Mutex
		var muOut sync.Mutex

		err := DeleteAssetByID(*authenticateObj, 123, &mu, &muOut, &signInCount, zapLogger)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
	})

	// Test case 2: Error when deleting asset
	t.Run("Delete error", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			case "/BeyondTrust/api/public/v3/Assets/123":
				if r.Method == "DELETE" {
					w.WriteHeader(http.StatusBadRequest)
					_, err := w.Write([]byte(`{"error": "not found"}`))
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
		defer server.Close()

		authenticateObj, _ := authentication.Authenticate(*authParams)
		server.URL = server.URL + constants.APIPath
		apiUrl, _ := url.Parse(server.URL)
		authenticateObj.ApiUrl = *apiUrl

		var signInCount uint64
		var mu sync.Mutex
		var muOut sync.Mutex

		err := DeleteAssetByID(*authenticateObj, 123, &mu, &muOut, &signInCount, zapLogger)
		if err == nil {
			t.Error("Expected error when deleting asset, but got none")
		}
	})

	// Test case 3: Authentication error
	t.Run("Authentication error", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case constants.APIPath + "/Auth/connect/token":
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte(`{"error": "invalid_client"}`))
				if err != nil {
					t.Error(err.Error())
				}
			}
		}))
		defer server.Close()

		authenticateObj, _ := authentication.Authenticate(*authParams)
		server.URL = server.URL + constants.APIPath
		apiUrl, _ := url.Parse(server.URL)
		authenticateObj.ApiUrl = *apiUrl

		var signInCount uint64
		var mu sync.Mutex
		var muOut sync.Mutex

		err := DeleteAssetByID(*authenticateObj, 123, &mu, &muOut, &signInCount, zapLogger)
		if err == nil {
			t.Error("Expected error due to authentication failure, but got none")
		}
	})
}
