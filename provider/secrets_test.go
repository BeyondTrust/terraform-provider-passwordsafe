package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/zap"
)

func TestResourceFileSecretCreate(t *testing.T) {

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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

	logger, _ := zap.NewDevelopment()

	// authenticate config
	zapLogger := logging.NewZapLogger(logger)
	httpClientObj, _ := utils.GetHttpClient(5, false, "", "", zapLogger)
	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.MaxElapsedTime = time.Second
	var authenticate, _ = authentication.Authenticate(*httpClientObj, backoffDefinition, "https://fake.api.com:443/BeyondTrust/api/public/v3/", "fakeone_a654+9sdf7+8we4f", "fakeone_aasd156465sfdef", zapLogger, 300)

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
