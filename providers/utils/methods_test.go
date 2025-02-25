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
		EndpointURL:                "https://fake.api.com:443/BeyondTrust/api/public/v3/",
		APIVersion:                 apiVersion,
		ClientID:                   "fakeone_a654+9sdf7+8we4f",
		ClientSecret:               "fakeone_a654+9sdf7+8we4f",
		ApiKey:                     "",
		Logger:                     zapLogger,
		RetryMaxElapsedTimeSeconds: 300,
	}
}

func TestAutenticate(t *testing.T) {

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

	_, err := Autenticate(*authenticateObj, &mu, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}

	// Increment counter
	_, err = Autenticate(*authenticateObj, &mu, &signInCount, zapLogger)
	if err != nil {
		t.Error(err)
	}
}

func TestAutenticateErrorGettingToken(t *testing.T) {

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

	_, err := Autenticate(*authenticateObj, &mu, &signInCount, zapLogger)
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
