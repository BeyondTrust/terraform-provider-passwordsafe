package provider

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecretTestConfigStringResponse struct {
	name     string
	server   *httptest.Server
	response string
}

func TestSecretSafeFlow(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":        "MySafe",
		"description": "A secure safe for testing",
	}
	var resourceSchema = map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestSecretSafeFlow",
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

			case "/secrets-safe/safes/":
				_, err := w.Write([]byte(`{"Id": "5b6fc3fb-fa78-48f9-9796-08dd18b16b5b","Name": "Safe Title", "Description": "Safe Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceSafeCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestSecretFolderFlow(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A secure folder for testing",
		"parent_folder_name": "folder_test",
		"user_group_id":      1,
	}
	var resourceSchema = map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_group_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestSecretFolderFlow",
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

			if r.URL.Path == "/secrets-safe/folders/" && r.Method == "GET" {
				_, err := w.Write([]byte(`[{"Id": "cb871861-8b40-4556-820c-1ca6d522adfa","Name": "folder_test"}, {"Id": "a4af73dc-4e89-41ec-eb9a-08dcf22d3aba","Name": "folder2"}]`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}
			if r.URL.Path == "/secrets-safe/folders/" && r.Method == "POST" {
				_, err := w.Write([]byte(`{"Id": "cb871861-8b40-4556-820c-1ca6d522adfa","Name": "Folder Title", "Description": "Folder Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFolderCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestSecretFolderFlowError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A secure folder for testing",
		"parent_folder_name": "",
		"user_group_id":      1,
	}
	var resourceSchema = map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_group_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestSecretFolderFlowError",
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

			if r.URL.Path == "/secrets-safe/folders/" && r.Method == "GET" {
				_, err := w.Write([]byte(`[{"Id": "cb871861-8b40-4556-820c-1ca6d522adfa","Name": "folder_test"}, {"Id": "a4af73dc-4e89-41ec-eb9a-08dcf22d3aba","Name": "folder2"}]`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}
			if r.URL.Path == "/secrets-safe/folders/" && r.Method == "POST" {
				_, err := w.Write([]byte(`{"Id": "cb871861-8b40-4556-820c-1ca6d522adfa","Name": "Folder Title", "Description": "Folder Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFolderCreate(data, authenticate)

	if err.Error() != "parent folder name must not be empty" {
		t.Errorf("Test case Failed %v, %v", err.Error(), "parent folder name must not be empty")
	}

}

func TestResourceFolderDelete(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A folder to delete",
		"parent_folder_name": "folder_test",
		"user_group_id":      1,
	}

	var resourceSchema = map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user_group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test folder ID
	data.SetId("cb871861-8b40-4556-820c-1ca6d522adfa")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceFolderDelete",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"Message": "Folder deleted successfully"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFolderDelete(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

	// Verify that the ID was cleared
	if data.Id() != "" {
		t.Errorf("Expected ID to be cleared after deletion, but got: %v", data.Id())
	}
}

func TestResourceFolderDeleteEmptyID(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A folder to delete",
		"parent_folder_name": "folder_test",
		"user_group_id":      1,
	}

	var resourceSchema = map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user_group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Don't set ID to test empty ID case

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceFolderDeleteEmptyID",
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

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFolderDelete(data, authenticate)

	if err == nil || err.Error() != "folder ID is empty" {
		t.Errorf("Expected 'folder ID is empty' error, but got: %v", err)
	}
}

func TestResourceFolderDeleteError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A folder to delete",
		"parent_folder_name": "folder_test",
		"user_group_id":      1,
	}

	var resourceSchema = map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user_group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test folder ID that will return 404
	data.SetId("non-existent-folder-id")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceFolderDeleteError",
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

			if r.URL.Path == "/secrets-safe/folders/non-existent-folder-id" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"error": "Folder not found"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFolderDelete(data, authenticate)

	if err == nil {
		t.Errorf("Expected error when deleting non-existent folder, but got nil")
	}
}

func TestResourceFolderDeleteAuthenticationError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"name":               "MyFolder",
		"description":        "A folder to delete",
		"parent_folder_name": "folder_test",
		"user_group_id":      1,
	}

	var resourceSchema = map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"user_group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	data.SetId("cb871861-8b40-4556-820c-1ca6d522adfa")

	// Pass nil authentication object to simulate authentication error
	err := resourceFolderDelete(data, nil)

	if err == nil {
		t.Errorf("Expected authentication error when passing nil authentication object, but got nil")
	}
}
