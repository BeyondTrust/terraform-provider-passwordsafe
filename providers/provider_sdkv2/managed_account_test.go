package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	"terraform-provider-passwordsafe/providers/constants"
	"testing"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/zap"
)

// the recommended version is 3.1. If no version is specified,
// the default API version 3.0 will be used
var apiVersion string = "3.1"

var authParams *authentication.AuthenticationParametersObj

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

func TestResourceManagedAccountCreate(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system01",
		"account_name": "account_name",
		"password":     "password",
	}
	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestresourceManagedAccountCreate",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response according to the endpoint path
			switch r.URL.Path {

			case "/Auth/connect/token":
				_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/SignAppIn":
				_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/Signout":
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/ManagedSystems":
				_, err := w.Write([]byte(`[{"ManagedSystemID":5, "SystemName":"system01", "EntityTypeID": 4}]`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/ManagedSystems/5/ManagedAccounts":
				_, err := w.Write([]byte(`{"ManagedSystemID":5, "ManagedAccountID":10, "AccountName": "Managed_account_name"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			default:
				http.NotFound(w, r)
			}
		})),
		response: "fake_credential",
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceManagedAccountCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestResourceManagedAccountCreateError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system0101",
		"account_name": "account_name",
		"password":     "password",
	}

	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestresourceManagedAccountCreateError",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response according to the endpoint path
			switch r.URL.Path {

			case "/Auth/connect/token":
				_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/SignAppIn":
				_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/Signout":
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/ManagedSystems":
				_, err := w.Write([]byte(`[{"ManagedSystemID":5, "SystemName":"system01", "EntityTypeID": 4}]`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/ManagedSystems/5/ManagedAccounts":
				_, err := w.Write([]byte(`{"ManagedSystemID":5, "ManagedAccountID":10, "AccountName": "Managed_account_name"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			default:
				http.NotFound(w, r)
			}
		})),
		response: "fake_credential",
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceManagedAccountCreate(data, authenticate)
	if err.Error() != "managed system system0101 was not found in managed system list" {
		t.Errorf("Test case Failed %v, %v", err.Error(), " managed system system0101 was not found in managed system list")
	}

}

func TestGetManagedAccountReadContext(t *testing.T) {

	rawData := map[string]interface{}{
		"system_name":  "system_name",
		"account_name": "account_name",
		"value":        "/",
	}

	var resourceSchema = map[string]*schema.Schema{
		"system_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"account_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestGetManagedAccountReadContext",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			switch r.URL.Path {

			case "/Auth/connect/token":
				_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/SignAppIn":
				_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Auth/Signout":
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/ManagedAccounts":
				_, err := w.Write([]byte(`{"SystemId":1,"AccountId":10}`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Requests":
				_, err := w.Write([]byte(`124`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Credentials/124":
				_, err := w.Write([]byte(`"fake_credential"`))
				if err != nil {
					t.Error("Test case Failed")
				}

			case "/Requests/124/checkin":
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error("Test case Failed")
				}

			default:
				http.NotFound(w, r)
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := getManagedAccountReadContext(context.Background(), data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestResourceManagedAccountDelete(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system01",
		"account_name": "account_name",
		"password":     "password",
	}
	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test managed account ID
	data.SetId("123")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceManagedAccountDelete",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Mocking Response according to the endpoint path
			if r.URL.Path == "/Auth/connect/token" {
				_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/Auth/SignAppIn" {
				_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/Auth/Signout" {
				_, err := w.Write([]byte(`{"Message": "SignOut successful"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/ManagedAccounts/123" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"Message": "Managed account deleted successfully"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceManagedAccountDelete(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

	// Verify that the ID was cleared
	if data.Id() != "" {
		t.Errorf("Expected ID to be cleared after deletion, but got: %v", data.Id())
	}
}

func TestResourceManagedAccountDeleteInvalidID(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system01",
		"account_name": "account_name",
		"password":     "password",
	}
	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set an invalid ID that can't be converted to int
	data.SetId("invalid_id")

	var authenticate, _ = authentication.Authenticate(*authParams)

	err := resourceManagedAccountDelete(data, authenticate)

	if err == nil {
		t.Errorf("Expected error for invalid ID, but got nil")
	}
}

func TestResourceManagedAccountDeleteError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system01",
		"account_name": "account_name",
		"password":     "password",
	}
	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test managed account ID that will return 404
	data.SetId("999")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceManagedAccountDeleteError",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Mocking Response according to the endpoint path
			if r.URL.Path == "/Auth/connect/token" {
				_, err := w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/Auth/SignAppIn" {
				_, err := w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/ManagedAccounts/999" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"error": "Managed account not found"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceManagedAccountDelete(data, authenticate)

	if err == nil {
		t.Errorf("Expected error when deleting non-existent managed account, but got nil")
	}
}

func TestResourceManagedAccountDeleteAuthenticationError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"system_name":  "system01",
		"account_name": "account_name",
		"password":     "password",
	}
	var resourceSchema = getManagedAccountSchema()

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	data.SetId("123")

	// Pass nil authentication object to simulate authentication error
	err := resourceManagedAccountDelete(data, nil)

	if err == nil {
		t.Errorf("Expected authentication error when passing nil authentication object, but got nil")
	}
}
