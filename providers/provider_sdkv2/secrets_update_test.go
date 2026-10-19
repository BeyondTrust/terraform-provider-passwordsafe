package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecretId = "01ca9cf3-0751-4a90-4856-08dcf22d7472"
	testFolderId = "cb871861-8b40-4556-820c-1ca6d522adfa"
	testTotpSeed = "JBSWY3DPEHPK3PXP"
)

// recordedRequest is a secret create or update request received by the mock server.
type recordedRequest struct {
	method      string
	path        string
	contentType string
	body        []byte
}

// secretMockServer mocks the folder list and the secret create/update endpoints, and records
// every create and update request.
type secretMockServer struct {
	server   *httptest.Server
	mu       sync.Mutex
	requests []recordedRequest
}

func newSecretMockServer(t *testing.T) *secretMockServer {
	t.Helper()
	mock := &secretMockServer{}
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/secrets-safe/folders/" && r.Method == "GET" {
			_, _ = w.Write([]byte(`[{"Id": "` + testFolderId + `","Name": "folder_test"}, {"Id": "a4af73dc-4e89-41ec-eb9a-08dcf22d3aba","Name": "folder2"}]`))
			return
		}

		body, _ := io.ReadAll(r.Body)
		mock.mu.Lock()
		mock.requests = append(mock.requests, recordedRequest{method: r.Method, path: r.URL.Path, contentType: r.Header.Get("Content-Type"), body: body})
		mock.mu.Unlock()

		if r.Method == "PUT" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_, _ = w.Write([]byte(`{"Id": "` + testSecretId + `","Title": "Secret Title"}`))
	}))
	t.Cleanup(mock.server.Close)
	return mock
}

// recorded returns the create and update requests received so far.
func (mock *secretMockServer) recorded() []recordedRequest {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	return append([]recordedRequest(nil), mock.requests...)
}

// newSecretTestMeta returns provider meta for apiVersion, pointing at the mock server.
func newSecretTestMeta(t *testing.T, mock *secretMockServer, version string) *providerMeta {
	t.Helper()
	InitializeGlobalConfig()

	params := *authParams
	params.APIVersion = version
	authenticate, err := authentication.Authenticate(params)
	require.NoError(t, err)

	apiUrl, _ := url.Parse(mock.server.URL + "/")
	authenticate.ApiUrl = *apiUrl

	return &providerMeta{
		authObj:   authenticate,
		signAppin: entities.SignAppinResponse{UserId: 1, UserName: "test", Name: "Test User", EmailAddress: "test@beyondtrust.com"},
	}
}

// credentialSecretAttributes returns the attributes of a credential secret, merged with overrides.
func credentialSecretAttributes(overrides map[string]interface{}) map[string]interface{} {
	attributes := map[string]interface{}{
		"folder_name": "folder_test",
		"title":       "Secret Title",
		"description": "Secret Description",
		"username":    "username",
		"password":    "SafeText",
		"owner_id":    1,
		"owner_type":  "User",
		"group_id":    1,
	}
	for key, value := range overrides {
		attributes[key] = value
	}
	return attributes
}

// applySecretUpdate plans config against a prior state built from prior and applies the plan,
// as terraform apply does. It returns the resulting state and the apply error.
func applySecretUpdate(t *testing.T, r *schema.Resource, prior, config map[string]interface{}, meta *providerMeta) (*terraform.InstanceState, error) {
	t.Helper()

	priorData := r.Data(nil)
	priorData.SetId(testSecretId)
	for key, value := range prior {
		require.NoError(t, priorData.Set(key, value))
	}
	state := priorData.State()

	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), meta)
	require.NoError(t, err)
	require.NotNil(t, diff, "expected a planned change")

	newState, diags := r.Apply(context.Background(), state, diff, meta)
	return newState, errorFromDiags(diags)
}

// errorFromDiags joins the error summaries of diags, or returns nil when there is no error.
func errorFromDiags(diags diag.Diagnostics) error {
	if !diags.HasError() {
		return nil
	}
	var messages []string
	for _, d := range diags {
		messages = append(messages, d.Summary)
	}
	return errors.New(strings.Join(messages, "; "))
}

// decodeBody decodes a JSON request body into a map.
func decodeBody(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &payload), "body: %s", body)
	return payload
}

func TestResourceCredentialSecretCreateWithTotp(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	data := schema.TestResourceDataRaw(t, resourceCredentialSecret().Schema, credentialSecretAttributes(map[string]interface{}{
		"totp_enabled":    true,
		"totp_secret_key": testTotpSeed,
	}))

	require.NoError(t, resourceCredentialSecretCreate(data, meta))
	assert.Equal(t, testSecretId, data.Id())

	requests := mock.recorded()
	require.Len(t, requests, 1)
	assert.Equal(t, "POST", requests[0].method)
	assert.Equal(t, "/secrets-safe/folders/"+testFolderId+"/secrets", requests[0].path)

	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, true, payload["TotpEnabled"])
	assert.Equal(t, testTotpSeed, payload["TotpSecretKey"])
}

func TestResourceCredentialSecretCreateTotpRequiresApiVersion33(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.2")

	data := schema.TestResourceDataRaw(t, resourceCredentialSecret().Schema, credentialSecretAttributes(map[string]interface{}{
		"totp_enabled":    true,
		"totp_secret_key": testTotpSeed,
	}))

	err := resourceCredentialSecretCreate(data, meta)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "require API version 3.3")
	assert.Empty(t, mock.recorded(), "no request should be sent")
}

func TestResourceCredentialSecretCreateTotpDisabledOnOlderApiVersion(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.1")

	// totp_enabled = false is the default and must not be sent, so it works below API v3.3.
	data := schema.TestResourceDataRaw(t, resourceCredentialSecret().Schema, credentialSecretAttributes(map[string]interface{}{
		"totp_enabled": false,
	}))

	require.NoError(t, resourceCredentialSecretCreate(data, meta))

	requests := mock.recorded()
	require.Len(t, requests, 1)
	assert.NotContains(t, string(requests[0].body), "Totp")
}

func TestResourceCredentialSecretCreateTotpSeedWithoutEnabled(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	data := schema.TestResourceDataRaw(t, resourceCredentialSecret().Schema, credentialSecretAttributes(map[string]interface{}{
		"totp_secret_key": testTotpSeed,
	}))

	err := resourceCredentialSecretCreate(data, meta)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "'TotpSecretKey' can only be provided when 'TotpEnabled' is true")
	assert.Empty(t, mock.recorded(), "no request should be sent")
}

func TestResourceCredentialSecretUpdate(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")
	totp := map[string]interface{}{"totp_enabled": true, "totp_secret_key": testTotpSeed}

	prior := credentialSecretAttributes(totp)
	config := credentialSecretAttributes(totp)
	config["title"] = "New Title"
	config["password"] = "NewSafeText"

	newState, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)
	assert.Equal(t, testSecretId, newState.ID)
	assert.Equal(t, "New Title", newState.Attributes["title"])

	requests := mock.recorded()
	require.Len(t, requests, 1)
	assert.Equal(t, "PUT", requests[0].method)
	assert.Equal(t, "/secrets-safe/secrets/"+testSecretId, requests[0].path)

	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, testFolderId, payload["FolderId"])
	assert.Equal(t, "New Title", payload["Title"])
	assert.Equal(t, "Secret Description", payload["Description"])
	assert.Equal(t, "username", payload["Username"])
	assert.Equal(t, "NewSafeText", payload["Password"])
	// Unchanged TOTP settings are left untouched: neither TOTP field is sent.
	assert.NotContains(t, payload, "TotpEnabled")
	assert.NotContains(t, payload, "TotpSecretKey")

	owners, ok := payload["Owners"].([]interface{})
	require.True(t, ok, "Owners: %v", payload["Owners"])
	require.Len(t, owners, 1)
	assert.Equal(t, float64(1), owners[0].(map[string]interface{})["UserId"])
	assert.Equal(t, float64(1), owners[0].(map[string]interface{})["GroupId"])
}

func TestResourceCredentialSecretUpdateEnablesTotp(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	prior := credentialSecretAttributes(nil)
	config := credentialSecretAttributes(map[string]interface{}{"totp_enabled": true, "totp_secret_key": testTotpSeed})

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, true, payload["TotpEnabled"])
	assert.Equal(t, testTotpSeed, payload["TotpSecretKey"])
}

func TestResourceCredentialSecretUpdateReplacesTotpSeed(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	prior := credentialSecretAttributes(map[string]interface{}{"totp_enabled": true, "totp_secret_key": testTotpSeed})
	config := credentialSecretAttributes(map[string]interface{}{"totp_enabled": true, "totp_secret_key": "KRSXG5CTMVRXEZLU"})

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, true, payload["TotpEnabled"])
	assert.Equal(t, "KRSXG5CTMVRXEZLU", payload["TotpSecretKey"])
}

func TestResourceCredentialSecretUpdateDisablesTotp(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	// Removing both TOTP attributes from the configuration disables TOTP.
	prior := credentialSecretAttributes(map[string]interface{}{"totp_enabled": true, "totp_secret_key": testTotpSeed})
	config := credentialSecretAttributes(nil)

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, false, payload["TotpEnabled"])
	assert.NotContains(t, payload, "TotpSecretKey")
}

func TestResourceCredentialSecretUpdateTotpRequiresApiVersion33(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.2")

	prior := credentialSecretAttributes(nil)
	config := credentialSecretAttributes(map[string]interface{}{"title": "New Title", "totp_enabled": true, "totp_secret_key": testTotpSeed})

	newState, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "require API version 3.3")
	assert.Empty(t, mock.recorded(), "no update request should be sent")

	// The failed update keeps the prior state, so the change is planned again.
	require.NotNil(t, newState)
	assert.Equal(t, "Secret Title", newState.Attributes["title"])
	assert.NotEqual(t, "true", newState.Attributes["totp_enabled"])
	assert.Empty(t, newState.Attributes["totp_secret_key"])
}

func TestResourceCredentialSecretUpdateWithoutTotpOnOlderApiVersion(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.1")

	prior := credentialSecretAttributes(nil)
	config := credentialSecretAttributes(map[string]interface{}{"notes": "New notes"})

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, "New notes", payload["Notes"])
	assert.NotContains(t, payload, "TotpEnabled")
}

func TestResourceCredentialSecretUpdateApiVersion30Owners(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.0")

	prior := credentialSecretAttributes(nil)
	config := credentialSecretAttributes(map[string]interface{}{"title": "New Title"})

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, "User", payload["OwnerType"])

	owners, ok := payload["Owners"].([]interface{})
	require.True(t, ok, "Owners: %v", payload["Owners"])
	require.Len(t, owners, 1)
	assert.Equal(t, float64(1), owners[0].(map[string]interface{})["OwnerId"])
	assert.NotContains(t, owners[0], "UserId")
}

func TestResourceTextSecretUpdate(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	prior := credentialSecretAttributes(nil)
	delete(prior, "username")
	delete(prior, "password")
	prior["text"] = "old text"
	config := map[string]interface{}{}
	for key, value := range prior {
		config[key] = value
	}
	config["text"] = "new text"

	_, err := applySecretUpdate(t, resourceTextSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	assert.Equal(t, "PUT", requests[0].method)
	assert.Equal(t, "/secrets-safe/secrets/"+testSecretId+"/text", requests[0].path)

	payload := decodeBody(t, requests[0].body)
	assert.Equal(t, testFolderId, payload["FolderId"])
	assert.Equal(t, "new text", payload["Text"])
	assert.Equal(t, "Secret Title", payload["Title"])
}

func TestResourceFileSecretUpdate(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	prior := credentialSecretAttributes(nil)
	delete(prior, "username")
	delete(prior, "password")
	prior["file_name"] = "file.txt"
	prior["file_content"] = "old content"
	config := map[string]interface{}{}
	for key, value := range prior {
		config[key] = value
	}
	config["file_content"] = "new content"

	_, err := applySecretUpdate(t, resourceFileSecret(), prior, config, meta)
	require.NoError(t, err)

	requests := mock.recorded()
	require.Len(t, requests, 1)
	assert.Equal(t, "PUT", requests[0].method)
	assert.Equal(t, "/secrets-safe/secrets/"+testSecretId+"/file", requests[0].path)

	mediaType, params, err := mime.ParseMediaType(requests[0].contentType)
	require.NoError(t, err)
	require.Equal(t, "multipart/form-data", mediaType)

	parts := map[string]string{}
	reader := multipart.NewReader(strings.NewReader(string(requests[0].body)), params["boundary"])
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		content, _ := io.ReadAll(part)
		parts[part.FormName()] = string(content)
	}

	assert.Equal(t, "new content", parts["file"])
	metadata := decodeBody(t, []byte(parts["secretmetadata"]))
	assert.Equal(t, testFolderId, metadata["FolderId"])
	assert.Equal(t, "file.txt", metadata["FileName"])
}

func TestResourceSecretUpdateFolderNotFound(t *testing.T) {
	mock := newSecretMockServer(t)
	meta := newSecretTestMeta(t, mock, "3.3")

	// folder_name is ForceNew, so a missing folder only happens when it was deleted or renamed.
	prior := credentialSecretAttributes(map[string]interface{}{"folder_name": "deleted_folder"})
	config := credentialSecretAttributes(map[string]interface{}{"folder_name": "deleted_folder", "title": "New Title"})

	_, err := applySecretUpdate(t, resourceCredentialSecret(), prior, config, meta)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "folder deleted_folder was not found in folder list")
	assert.Empty(t, mock.recorded(), "no update request should be sent")
}

func TestResourceSecretFolderNameForceNew(t *testing.T) {
	for name, r := range map[string]*schema.Resource{
		"credential": resourceCredentialSecret(),
		"text":       resourceTextSecret(),
		"file":       resourceFileSecret(),
	} {
		t.Run(name, func(t *testing.T) {
			assert.True(t, r.Schema["folder_name"].ForceNew)
		})
	}
}
