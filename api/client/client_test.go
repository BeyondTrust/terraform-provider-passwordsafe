// Copyright 2023 BeyondTrust. All rights reserved.
// Package client implements functions to call Beyondtrust Secret Safe API.
// Unit tests for Client package.
package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"terraform-provider-passwordsafe/api/client/entities"
	"terraform-provider-passwordsafe/config"
	"testing"
)

type UserTestConfig struct {
	name     string
	server   *httptest.Server
	response *entities.User
}

type SecretTestConfig struct {
	name     string
	server   *httptest.Server
	response *entities.Secret
}

type ManagedAccountTestConfig struct {
	name     string
	server   *httptest.Server
	response *entities.ManagedAccount
}

type ManagedAccountCreateRequestConfig struct {
	name     string
	server   *httptest.Server
	response string
}

var apiKey string = config.PS_API_KEY
var apiAccountName = config.PS_ACCOUNT_NAME
var url = config.PS_URL
var clientCertificatePath = config.CERTIFICATE_PATH
var clientCertificateName = config.CERTIFICATE_NAME
var clientCertificatePassword = config.CERTIFICATE_PASSWORD

var apiClient, _ = NewClient(url, apiKey, apiAccountName, false, clientCertificatePath, clientCertificateName, clientCertificatePassword)

func TestSignAppin(t *testing.T) {

	testConfig := UserTestConfig{
		name: "TestSignAppin",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))
		})),
		response: &entities.User{
			UserId:       1,
			EmailAddress: "Felipe",
		},
	}

	response, err := apiClient.SignAppin(testConfig.server.URL + "/" + "TestSignAppin")

	if !reflect.DeepEqual(response, *testConfig.response) {
		t.Errorf("Test case Failed %v, %v", response, *testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestSignOut(t *testing.T) {

	testConfig := UserTestConfig{
		name: "TestSignOut",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(``))
		})),
		response: nil,
	}

	err := apiClient.SignOut(testConfig.server.URL)

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestSecretGetFileSecret(t *testing.T) {

	testConfig := UserTestConfig{
		name: "TestSecretGetFileSecret",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`fake_password`))
		})),
	}

	response, err := apiClient.SecretGetFileSecret("1", testConfig.server.URL)

	if response != "fake_password" {
		t.Errorf("Test case Failed %v, %v", response, *testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestSecretGetSecretByPath(t *testing.T) {

	testConfig := SecretTestConfig{
		name: "TestSecretGetSecretByPath",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response
			w.Write([]byte(`[{"Password": "credential_in_sub_3_password","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))
		})),
		response: &entities.Secret{
			Id:       "9152f5b6-07d6-4955-175a-08db047219ce",
			Title:    "credential_in_sub_3",
			Password: "credential_in_sub_3_password",
		},
	}

	response, err := apiClient.SecretGetSecretByPath("path1/path2", "fake_title", "/", testConfig.server.URL)

	if response != *testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, *testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestManagedAccountGet(t *testing.T) {

	testConfig := ManagedAccountTestConfig{
		name: "TestManagedAccountGet",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response
			w.Write([]byte(`{"SystemId": 1,"AccountId": 10}`))
		})),
		response: &entities.ManagedAccount{
			SystemId:  1,
			AccountId: 10,
		},
	}

	response, err := apiClient.ManagedAccountGet("fake_system_name", "fake_account_name", testConfig.server.URL)

	if response != *testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, *testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestManagedAccountCreateRequest(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestManagedAccountCreateRequest",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response
			w.Write([]byte(`124`))
		})),
		response: "124",
	}

	response, err := apiClient.ManagedAccountCreateRequest(1, 10, testConfig.server.URL)

	if response != testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestCredentialByRequestId(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestCredentialByRequestId",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response
			w.Write([]byte(`fake_credential`))
		})),
		response: "fake_credential",
	}

	response, err := apiClient.CredentialByRequestId("124", testConfig.server.URL)

	if response != testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestManagedAccountRequestCheckIn(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestManagedAccountRequestCheckIn",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response
			w.Write([]byte(``))
		})),
		response: "",
	}

	response, err := apiClient.ManagedAccountRequestCheckIn("124", testConfig.server.URL)

	if response != testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestManageAccountFlow(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestManageAccountFlow",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case fmt.Sprintf("/ManagedAccounts"):
				w.Write([]byte(`{"SystemId":1,"AccountId":10}`))

			case "/Requests":
				w.Write([]byte(`124`))

			case "/Credentials/124":
				w.Write([]byte(`"fake_credential"`))

			case "/Requests/124/checkin":
				w.Write([]byte(``))

			default:
				http.NotFound(w, r)
			}
		})),
		response: "fake_credential",
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["ManagedAccountGetPath"] = fmt.Sprintf("ManagedAccounts?systemName=%v&accountName=%v", "system_name_test", "account_name_test")
	paths["ManagedAccountCreateRequestPath"] = "Requests"
	paths["CredentialByRequestIdPath"] = "Credentials/%v"
	paths["ManagedAccountRequestCheckInPath"] = "Requests/%v/checkin"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	response, err := apiClient.ManageAccountFlow("system_name_test", "account_name_test", paths)

	if response != testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestManageAccountFlowFailedSignAppin(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestManageAccountFlowFailedSignAppin",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`Unauthorized`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case fmt.Sprintf("/ManagedAccounts"):
				w.Write([]byte(`{"SystemId":1,"AccountId":10}`))

			case "/Requests":
				w.Write([]byte(`124`))

			case "/Credentials/124":
				w.Write([]byte(`"fake_credential"`))

			case "/Requests/124/checkin":
				w.Write([]byte(``))

			default:
				http.NotFound(w, r)
			}
		})),
		response: "got a non 200 status code: 401 - Unauthorized",
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["ManagedAccountGetPath"] = fmt.Sprintf("ManagedAccounts?systemName=%v&accountName=%v", "system_name_test", "account_name_test")
	paths["ManagedAccountCreateRequestPath"] = "Requests"
	paths["CredentialByRequestIdPath"] = "Credentials/%v"
	paths["ManagedAccountRequestCheckInPath"] = "Requests/%v/checkin"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	_, err := apiClient.ManageAccountFlow("system_name_test", "account_name_test", paths)

	if err.Error() != testConfig.response {
		t.Errorf("Test case Failed %v, %v", err.Error(), testConfig.response)
	}

}

func TestManageAccountFlowFailedManagedAccounts(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestManageAccountFlowFailedManagedAccounts",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case fmt.Sprintf("/ManagedAccounts"):
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`"Managed Account not found"`))

			case "/Requests":
				w.Write([]byte(`124`))

			case "/Credentials/124":
				w.Write([]byte(`"fake_credential"`))

			case "/Requests/124/checkin":
				w.Write([]byte(``))

			default:
				http.NotFound(w, r)
			}
		})),
		response: `got a non 200 status code: 404 - "Managed Account not found"`,
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["ManagedAccountGetPath"] = fmt.Sprintf("ManagedAccounts?systemName=%v&accountName=%v", "system_name_test", "account_name_test")
	paths["ManagedAccountCreateRequestPath"] = "Requests"
	paths["CredentialByRequestIdPath"] = "Credentials/%v"
	paths["ManagedAccountRequestCheckInPath"] = "Requests/%v/checkin"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	_, err := apiClient.ManageAccountFlow("system_name_test", "account_name_test", paths)

	if err.Error() != testConfig.response {
		t.Errorf("Test case Failed %v, %v", err.Error(), testConfig.response)
	}

}

func TestSecretFlow(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestSecretFlow",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case "/secrets-safe/secrets":
				w.Write([]byte(`[{"SecretType": "FILE", "Password": "credential_in_sub_3_password","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))

			case "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download":
				w.Write([]byte(`fake_password`))

			default:
				http.NotFound(w, r)
			}
		})),
		response: "fake_password",
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["SecretGetSecretByPathPath"] = "secrets-safe/secrets?title=%v&path=%v&separator=%v"
	paths["SecretGetFileSecretPath"] = "secrets-safe/secrets/%v/file/download"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	response, err := apiClient.SecretFlow("path1\\path2", "credential_in_sub_3", "\\", paths)

	if response != testConfig.response {
		t.Errorf("Test case Failed %v, %v", response, testConfig.response)
	}

	if err != nil {
		t.Errorf("Test case Failed: %v", err)
	}
}

func TestSecretFlowFailedSignAppin(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestSecretFlowFailedSignAppin",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":

				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`Unauthorized`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case "/secrets-safe/secrets":
				w.Write([]byte(`[{"SecretType": "FILE", "Password": "credential_in_sub_3_password","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))

			case "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download":
				w.Write([]byte(`fake_password`))

			default:
				http.NotFound(w, r)
			}
		})),
		response: "got a non 200 status code: 401 - Unauthorized",
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["SecretGetSecretByPathPath"] = "secrets-safe/secrets?title=%v&path=%v&separator=%v"
	paths["SecretGetFileSecretPath"] = "secrets-safe/secrets/%v/file/download"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	_, err := apiClient.SecretFlow("path1\\path2", "credential_in_sub_3", "\\", paths)

	if err.Error() != testConfig.response {
		t.Errorf("Test case Failed %v, %v", err.Error(), testConfig.response)
	}
}

func TestSecretFlowFailedSecretNotFound(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestSecretFlowFailedSecretNotFound",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case "/secrets-safe/secrets":
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"Secret does not exist."}`))

			case "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download":
				w.Write([]byte(`fake_password`))

			default:
				http.NotFound(w, r)
			}
		})),
		response: `got a non 200 status code: 404 - {"error":"Secret does not exist."}`,
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["SecretGetSecretByPathPath"] = "secrets-safe/secrets?title=%v&path=%v&separator=%v"
	paths["SecretGetFileSecretPath"] = "secrets-safe/secrets/%v/file/download"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	_, err := apiClient.SecretFlow("path1\\path2", "credential_in_sub_3", "\\", paths)

	if err.Error() != testConfig.response {
		t.Errorf("Test case Failed %v, %v", err.Error(), testConfig.response)
	}
}

func TestSecretFlowFailedSecretNotFoundDowload(t *testing.T) {

	testConfig := ManagedAccountCreateRequestConfig{
		name: "TestSecretFlowFailedSecretNotFoundDowload",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mocking Response accorging to the endpoint path
			switch r.URL.Path {

			case "/Auth/SignAppin":
				w.Write([]byte(`{"UserId":1, "EmailAddress":"Felipe"}`))

			case "/Auth/Signout":
				w.Write([]byte(``))

			case "/secrets-safe/secrets":

				w.Write([]byte(`[{"SecretType": "FILE", "Password": "credential_in_sub_3_password","Id": "9152f5b6-07d6-4955-175a-08db047219ce","Title": "credential_in_sub_3"}]`))

			case "/secrets-safe/secrets/9152f5b6-07d6-4955-175a-08db047219ce/file/download":
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"Secret does not exist.}`))

			default:
				http.NotFound(w, r)
			}
		})),
		response: `got a non 200 status code: 404 - {"error":"Secret does not exist.}`,
	}

	paths := make(map[string]string)

	paths["SignAppinOutPath"] = "Auth/Signout"
	paths["SignAppinPath"] = "Auth/SignAppin"
	paths["SecretGetSecretByPathPath"] = "secrets-safe/secrets?title=%v&path=%v&separator=%v"
	paths["SecretGetFileSecretPath"] = "secrets-safe/secrets/%v/file/download"

	// Changing actual endpoint by fake server endpoint
	apiClient.url = testConfig.server.URL

	_, err := apiClient.SecretFlow("path1\\path2", "credential_in_sub_3", "\\", paths)

	if err.Error() != testConfig.response {
		t.Errorf("Test case Failed %v, %v", err.Error(), testConfig.response)
	}
}
