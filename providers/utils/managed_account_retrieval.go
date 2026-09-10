// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
)

const (
	// MaxSystemNameLength is the maximum length accepted by the API for a system name.
	MaxSystemNameLength = 128
	// MaxAccountNameLength is the maximum length accepted by the API for an account name.
	MaxAccountNameLength = 245
)

// ValidateManagedAccountCoordinates validates the system name and account name
// used to retrieve a managed account credential. It replaces the length checks
// that utils.ValidatePaths performs in the SDK, which GetManagedAccountSecret
// deliberately bypasses.
func ValidateManagedAccountCoordinates(systemName string, accountName string) error {
	if strings.TrimSpace(systemName) == "" {
		return errors.New("empty system name, please provide a valid system name")
	}

	if strings.TrimSpace(accountName) == "" {
		return errors.New("empty account name, please provide a valid account name")
	}

	if len(systemName) > MaxSystemNameLength {
		return fmt.Errorf("invalid system name length, the maximum length is %v characters", MaxSystemNameLength)
	}

	if len(accountName) > MaxAccountNameLength {
		return fmt.Errorf("invalid account name length, the maximum length is %v characters", MaxAccountNameLength)
	}

	return nil
}

// GetManagedAccountSecret retrieves the credential of a single managed account,
// treating system name and account name as the separate fields they are.
//
// The SDK's GetSecret/ManageAccountFlow take one "<system_name>/<account_name>"
// path and split it on the separator, so any system name containing a "/" (for
// example "Accounts - AD/EntraID") is rejected by ValidatePaths and surfaces as
// "empty managed account list". This helper runs the same four API calls as
// ManageAccountFlow without ever concatenating or splitting a path, so the
// separator is just another character in the value.
//
// The SDK gained an equivalent GetManagedAccountSecret method after v1.3.2;
// once the provider pins a release that includes it, this helper can delegate
// to it instead of repeating the flow.
func GetManagedAccountSecret(authenticationObj *auth.AuthenticationObj, managedAccountObj *managed_accounts.ManagedAccountstObj, systemName string, accountName string) (string, error) {
	if err := ValidateManagedAccountCoordinates(systemName, accountName); err != nil {
		return "", err
	}

	query := url.Values{}
	query.Add("systemName", systemName)
	query.Add("accountName", accountName)

	managedAccountGetUrl := authenticationObj.ApiUrl.JoinPath("ManagedAccounts").String() + "?" + query.Encode()
	managedAccount, err := managedAccountObj.ManagedAccountGet(systemName, accountName, managedAccountGetUrl)
	if err != nil {
		return "", err
	}

	createRequestUrl := authenticationObj.ApiUrl.JoinPath("Requests").String()
	requestId, err := managedAccountObj.ManagedAccountCreateRequest(managedAccount.SystemId, managedAccount.AccountId, createRequestUrl)
	if err != nil {
		return "", err
	}

	credentialUrl := authenticationObj.ApiUrl.JoinPath("Credentials", requestId).String()
	credential, err := managedAccountObj.CredentialByRequestId(requestId, credentialUrl)
	if err != nil {
		return "", err
	}

	checkInUrl := authenticationObj.ApiUrl.JoinPath("Requests", requestId, "checkin").String()
	if _, err := managedAccountObj.ManagedAccountRequestCheckIn(requestId, checkInUrl); err != nil {
		return "", err
	}

	// The API returns the credential as a JSON string, so it arrives quoted.
	// Fall back to the raw value when it is not, instead of dropping it.
	secretValue, unquoteErr := strconv.Unquote(credential)
	if unquoteErr != nil {
		return credential, nil
	}

	return secretValue, nil
}
