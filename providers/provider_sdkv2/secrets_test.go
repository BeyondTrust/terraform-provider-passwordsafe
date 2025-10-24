package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceFileSecretCreate(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name":  "folder_test",
		"name":         "File Secret name",
		"description":  "File Secret Description",
		"title":        "File Secret title",
		"file_name":    "file.txt",
		"file_content": "P@ssw0rd123!$",
		"owner_id":     1,
		"owner_type":   "User",
	}

	rawData["owners"] = []interface{}{
		map[string]interface{}{
			"owner_id": 1,
			"owner":    "User",
			"email":    "test@test.com",
		},
		map[string]interface{}{
			"owner_id": 2,
			"owner":    "Admin",
			"email":    "test@test.com",
		},
	}

	rawData["urls"] = []interface{}{
		map[string]interface{}{
			"id":            1,
			"credential_id": "User",
			"url":           "test@test.com",
		},
		map[string]interface{}{
			"id":            2,
			"credential_id": "Admin",
			"url":           "test@test.com",
		},
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"file_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"file_content": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestresourceFileSecretCreate",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets/file" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "File Secret Title", "Description": "Title Description", "FileName": "textfile.txt"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFileSecretCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestResourceFileSecretCreateError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name":  "not_found_folder_test",
		"name":         "Secret",
		"description":  "Description",
		"title":        "Secret Title",
		"file_name":    "file.txt",
		"file_content": "P@ssw0rd123!$",
		"owner_id":     1,
		"owner_type":   "User",
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"file_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"file_content": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestresourceFileSecretCreateError",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets/file" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "File Secret Title", "Description": "Title Description", "FileName": "textfile.txt"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceFileSecretCreate(data, authenticate)

	if err.Error() != "folder not_found_folder_test was not found in folder list" {
		t.Errorf("Test case Failed %v, %v", err.Error(), "folder not_found_folder_test was not found in folder list")
	}

}

func TestResourceCredentialSecretCreate(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"name":        "Secret Name",
		"description": "Secret Description",
		"title":       "Secret Title",
		"username":    "username",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}
	rawData["owners"] = []interface{}{
		map[string]interface{}{
			"owner_id": 1,
			"owner":    "User",
			"email":    "test@test.com",
		},
		map[string]interface{}{
			"owner_id": 2,
			"owner":    "Admin",
			"email":    "test@test.com",
		},
	}

	rawData["urls"] = []interface{}{
		map[string]interface{}{
			"id":            1,
			"credential_id": "User",
			"url":           "test@test.com",
		},
		map[string]interface{}{
			"id":            2,
			"credential_id": "Admin",
			"url":           "test@test.com",
		},
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceCredentialSecretCreate",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "Secret Title", "Description": "Title Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceCredentialSecretCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestResourceCredentialSecretCreateError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"name":        "Secret",
		"description": "Credential Description",
		"title":       "Credential Title",
		"username":    "",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}
	rawData["owners"] = []interface{}{
		map[string]interface{}{
			"owner_id": 1,
			"owner":    "User",
			"email":    "test@test.com",
		},
		map[string]interface{}{
			"owner_id": 2,
			"owner":    "Admin",
			"email":    "test@test.com",
		},
	}

	rawData["urls"] = []interface{}{
		map[string]interface{}{
			"id":            1,
			"credential_id": "User",
			"url":           "test@test.com",
		},
		map[string]interface{}{
			"id":            2,
			"credential_id": "Admin",
			"url":           "test@test.com",
		},
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceCredentialSecretCreateError",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "Secret Title", "Description": "Title Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceCredentialSecretCreate(data, authenticate)

	if err.Error() != "The field 'Username' is required." {
		t.Errorf("Test case Failed %v, %v", err.Error(), "The field 'Username' is required.")
	}

}

func TestGetSecretByPathReadContext(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"path":      "path/path2",
		"title":     "credential_in_sub_3",
		"separator": "/",
		"value":     "",
	}

	var resourceSchema = map[string]*schema.Schema{
		"path": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"separator": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/",
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
		name: "TestGetSecretByPathReadContext",
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

			if r.URL.Path == "/secrets-safe/secrets" {
				_, err := w.Write([]byte(`[{"SecretType": "FILE", "Password": "credential_in_sub_3_password","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download" {
				_, err := w.Write([]byte(`fake_password`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := getSecretByPathReadContext(context.Background(), data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestGetSecretByPathReadContextError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"path":      "path/path2",
		"title":     "credential_in_sub_3",
		"separator": "/",
		"value":     "",
	}

	var resourceSchema = map[string]*schema.Schema{
		"path": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"separator": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/",
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
		name: "TestGetSecretByPathReadContextError",
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

			if r.URL.Path == "/secrets-safe/secrets" {
				_, err := w.Write([]byte(`[]`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download" {
				_, err := w.Write([]byte(`fake_password`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	diags := getSecretByPathReadContext(context.Background(), data, authenticate)

	if diags[0].Summary != "error SecretGetSecretByPath, Secret was not found: StatusCode: 404 " {
		t.Errorf("Test case Failed %v, %v", diags[0].Summary, "error SecretGetSecretByPath, Secret was not found: StatusCode: 404 ")
	}

}

func TestResourceTextSecretCreate(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"name":        "Text secret name",
		"description": "Text secret description",
		"title":       "Text secret title",
		"text":        "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}
	rawData["owners"] = []interface{}{
		map[string]interface{}{
			"owner_id": 1,
			"owner":    "User",
			"email":    "test@test.com",
		},
		map[string]interface{}{
			"owner_id": 2,
			"owner":    "Admin",
			"email":    "test@test.com",
		},
	}

	rawData["urls"] = []interface{}{
		map[string]interface{}{
			"id":            1,
			"credential_id": "User",
			"url":           "test@test.com",
		},
		map[string]interface{}{
			"id":            2,
			"credential_id": "Admin",
			"url":           "test@test.com",
		},
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"text": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestSecretSafeFlow",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets/text" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "Secret Title", "Description": "Title Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceTextSecretCreate(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

}

func TestResourceTextSecretCreateError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"name":        "Text secret name",
		"description": "Text secret Description",
		"title":       "Text secret Secret Title",
		"text":        "",
		"owner_id":    1,
		"owner_type":  "User",
	}
	rawData["owners"] = []interface{}{
		map[string]interface{}{
			"owner_id": 1,
			"owner":    "User",
			"email":    "test@test.com",
		},
		map[string]interface{}{
			"owner_id": 2,
			"owner":    "Admin",
			"email":    "test@test.com",
		},
	}

	rawData["urls"] = []interface{}{
		map[string]interface{}{
			"id":            1,
			"credential_id": "User",
			"url":           "test@test.com",
		},
		map[string]interface{}{
			"id":            2,
			"credential_id": "Admin",
			"url":           "test@test.com",
		},
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"text": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestSecretSafeFlow",
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

			if r.URL.Path == "/secrets-safe/folders/cb871861-8b40-4556-820c-1ca6d522adfa/secrets/text" {
				_, err := w.Write([]byte(`{"Id": "01ca9cf3-0751-4a90-4856-08dcf22d7472","Title": "Secret Title", "Description": "Title Description"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceTextSecretCreate(data, authenticate)

	if err.Error() != "The field 'Text' is required." {
		t.Errorf("Test case Failed %v, %v", err.Error(), "The field 'Text' is required.")
	}

}

func TestResourceSecretDelete(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"title":       "Secret to Delete",
		"description": "Secret Description",
		"username":    "testuser",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test secret ID
	data.SetId("01ca9cf3-0751-4a90-4856-08dcf22d7472")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceSecretDelete",
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

			if r.URL.Path == "/Auth/SignAppOut" {
				_, err := w.Write([]byte(`{"Message": "SignOut successful"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

			if r.URL.Path == "/secrets-safe/secrets/01ca9cf3-0751-4a90-4856-08dcf22d7472" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"Message": "Secret deleted successfully"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceSecretDelete(data, authenticate)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}

	// Verify that the ID was cleared
	if data.Id() != "" {
		t.Errorf("Expected ID to be cleared after deletion, but got: %v", data.Id())
	}
}

func TestResourceSecretDeleteEmptyID(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"title":       "Secret to Delete",
		"description": "Secret Description",
		"username":    "testuser",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Don't set ID to test empty ID case

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceSecretDeleteEmptyID",
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

	err := resourceSecretDelete(data, authenticate)

	if err == nil || err.Error() != "secret ID is empty" {
		t.Errorf("Expected 'secret ID is empty' error, but got: %v", err)
	}
}

func TestResourceSecretDeleteError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"title":       "Secret to Delete",
		"description": "Secret Description",
		"username":    "testuser",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	// Set a test secret ID that will return 404
	data.SetId("non-existent-secret-id")

	var authenticate, _ = authentication.Authenticate(*authParams)

	// mock config
	testConfig := SecretTestConfigStringResponse{
		name: "TestResourceSecretDeleteError",
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

			if r.URL.Path == "/secrets-safe/secrets/non-existent-secret-id" && r.Method == "DELETE" {
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte(`{"error": "Secret not found"}`))
				if err != nil {
					t.Error("Test case Failed")
				}
			}

		})),
	}

	apiUrl, _ := url.Parse(testConfig.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	err := resourceSecretDelete(data, authenticate)

	if err == nil {
		t.Errorf("Expected error when deleting non-existent secret, but got nil")
	}
}

func TestResourceSecretDeleteAuthenticationError(t *testing.T) {

	InitializeGlobalConfig()

	rawData := map[string]interface{}{
		"folder_name": "folder_test",
		"title":       "Secret to Delete",
		"description": "Secret Description",
		"username":    "testuser",
		"password":    "P@ssw0rd123!$",
		"owner_id":    1,
		"owner_type":  "User",
	}

	var resourceSchema = map[string]*schema.Schema{
		"folder_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"title": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"owner_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, resourceSchema, rawData)
	data.SetId("01ca9cf3-0751-4a90-4856-08dcf22d7472")

	// Pass nil authentication object to simulate authentication error
	err := resourceSecretDelete(data, nil)

	if err == nil {
		t.Errorf("Expected authentication error when passing nil authentication object, but got nil")
	}
}
