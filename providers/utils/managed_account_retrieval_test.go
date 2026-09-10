package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"terraform-provider-passwordsafe/providers/constants"
	"testing"

	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
)

func TestValidateManagedAccountCoordinates(t *testing.T) {
	testCases := []struct {
		name        string
		systemName  string
		accountName string
		expectedErr string
	}{
		{
			name:        "valid coordinates",
			systemName:  "system01",
			accountName: "managed_account_01",
			expectedErr: "",
		},
		{
			name:        "system name containing the path separator is valid",
			systemName:  "Accounts - AD/EntraID",
			accountName: "managed_account_01",
			expectedErr: "",
		},
		{
			name:        "empty system name",
			systemName:  "   ",
			accountName: "managed_account_01",
			expectedErr: "empty system name",
		},
		{
			name:        "empty account name",
			systemName:  "system01",
			accountName: "",
			expectedErr: "empty account name",
		},
		{
			name:        "system name too long",
			systemName:  strings.Repeat("a", MaxSystemNameLength+1),
			accountName: "managed_account_01",
			expectedErr: "invalid system name length",
		},
		{
			name:        "account name too long",
			systemName:  "system01",
			accountName: strings.Repeat("a", MaxAccountNameLength+1),
			expectedErr: "invalid account name length",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateManagedAccountCoordinates(testCase.systemName, testCase.accountName)

			if testCase.expectedErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q, got nil", testCase.expectedErr)
			}

			if !strings.Contains(err.Error(), testCase.expectedErr) {
				t.Errorf("expected error containing %q, got %q", testCase.expectedErr, err.Error())
			}
		})
	}
}

// managedAccountCredentialServer mocks the four endpoints the managed account
// retrieval flow calls, recording the query values the ManagedAccounts request
// was made with.
func managedAccountCredentialServer(t *testing.T, gotSystemName *string, gotAccountName *string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err = w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))

		case constants.APIPath + "/Auth/SignAppIn":
			_, err = w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

		case constants.APIPath + "/Auth/Signout":
			_, err = w.Write([]byte(``))

		case constants.APIPath + "/ManagedAccounts":
			*gotSystemName = r.URL.Query().Get("systemName")
			*gotAccountName = r.URL.Query().Get("accountName")
			_, err = w.Write([]byte(`{"SystemId":1,"AccountId":10}`))

		case constants.APIPath + "/Requests":
			_, err = w.Write([]byte(`124`))

		case constants.APIPath + "/Credentials/124":
			_, err = w.Write([]byte(`"fake_credential"`))

		case constants.APIPath + "/Requests/124/checkin":
			_, err = w.Write([]byte(``))

		default:
			http.NotFound(w, r)
		}

		if err != nil {
			t.Error(err.Error())
		}
	}))
}

func TestGetManagedAccountSecretSystemNameWithSeparator(t *testing.T) {
	InitializeGlobalConfig()

	var gotSystemName, gotAccountName string
	server := managedAccountCredentialServer(t, &gotSystemName, &gotAccountName)
	defer server.Close()

	authObj := newAuthObjAtServer(t, server)

	managedAccountObj, err := managed_accounts.NewManagedAccountObj(*authObj, zapLogger)
	if err != nil {
		t.Fatalf("NewManagedAccountObj: %v", err)
	}

	// A system name containing the "/" separator used to be rejected by the
	// SDK's path validation, surfacing as "empty managed account list".
	secret, err := GetManagedAccountSecret(authObj, managedAccountObj, "Accounts - AD/EntraID", "managed_account_01")
	if err != nil {
		t.Fatalf("GetManagedAccountSecret: %v", err)
	}

	if secret != "fake_credential" {
		t.Errorf("expected secret %q, got %q", "fake_credential", secret)
	}

	if gotSystemName != "Accounts - AD/EntraID" {
		t.Errorf("expected systemName query %q, got %q", "Accounts - AD/EntraID", gotSystemName)
	}

	if gotAccountName != "managed_account_01" {
		t.Errorf("expected accountName query %q, got %q", "managed_account_01", gotAccountName)
	}
}

func TestGetManagedAccountSecretInvalidCoordinates(t *testing.T) {
	// Validation runs before any API call, so no authentication object or
	// managed account object is needed.
	_, err := GetManagedAccountSecret(nil, nil, "", "managed_account_01")

	if err == nil {
		t.Fatal("expected an error for an empty system name, got nil")
	}

	if !strings.Contains(err.Error(), "empty system name") {
		t.Errorf("expected an empty system name error, got %q", err.Error())
	}
}

func TestGetManagedAccountSecretNotFound(t *testing.T) {
	InitializeGlobalConfig()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.URL.Path {
		case constants.APIPath + "/Auth/connect/token":
			_, err = w.Write([]byte(`{"access_token": "fake_token", "expires_in": 600, "token_type": "Bearer", "scope": "publicapi"}`))

		case constants.APIPath + "/Auth/SignAppIn":
			_, err = w.Write([]byte(`{"UserId":1, "EmailAddress":"test@beyondtrust.com"}`))

		case constants.APIPath + "/ManagedAccounts":
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte(`"Managed Account not found"`))

		default:
			http.NotFound(w, r)
		}

		if err != nil {
			t.Error(err.Error())
		}
	}))
	defer server.Close()

	authObj := newAuthObjAtServer(t, server)

	managedAccountObj, err := managed_accounts.NewManagedAccountObj(*authObj, zapLogger)
	if err != nil {
		t.Fatalf("NewManagedAccountObj: %v", err)
	}

	_, err = GetManagedAccountSecret(authObj, managedAccountObj, "system01", "managed_account_01")

	if err == nil {
		t.Fatal("expected an error when the managed account is not found, got nil")
	}
}
