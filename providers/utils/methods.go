// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"errors"
	"fmt"
	"sync"
	"terraform-provider-passwordsafe/providers/entities"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/assets"
	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	libraryEntitites "github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
)

func TestResourceConfig(config entities.PasswordSafeTestConfig) string {
	return fmt.Sprintf(`
		provider "passwordsafe" {
			api_key = "%s"
			client_id = "%s"
			client_secret = "%s"
			url = "%s"
			verify_ca=false
			api_account_name = "%s"
			client_certificates_folder_path = "%s"
			client_certificate_name = "%s"
			client_certificate_password = "%s"
			api_version = "%s"
		}

		%s`,

		config.APIKey,
		config.ClientID,
		config.ClientSecret,
		config.URL,
		config.APIAccountName,
		config.ClientCertificatesFolderPath,
		config.ClientCertificateName,
		config.ClientCertificatePassword,
		config.APIVersion,
		config.Resource,
	)
}

var signAppinResponse libraryEntitites.SignAppinResponse

// Authenticate gets Password Safe authentication, sharing a session across
// concurrent callers via the reference count guarded by mu.
func Authenticate(authenticationObj auth.AuthenticationObj, mu *sync.Mutex, signInCount *uint64, zapLogger logging.Logger) (libraryEntitites.SignAppinResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	if *signInCount > 0 {
		*signInCount++
		zapLogger.Debug(fmt.Sprintf("%v %v", "Already signed in", *signInCount))
		return signAppinResponse, nil
	}

	resp, err := authenticationObj.GetPasswordSafeAuthentication()
	if err != nil {
		zapLogger.Error(err.Error())
		return libraryEntitites.SignAppinResponse{}, err
	}
	signAppinResponse = resp
	*signInCount++
	zapLogger.Debug(fmt.Sprintf("%v %v", "signin", *signInCount))
	return signAppinResponse, nil
}

// SignOut releases this caller's reference to the shared session. The same
// mutex used by Authenticate must be passed in so signin and signout can never
// run concurrently — the API's signout is user-global and would otherwise tear
// down a session another worker is racing to use.
func SignOut(authenticationObj auth.AuthenticationObj, mu *sync.Mutex, signInCount *uint64, zapLogger logging.Logger) error {
	mu.Lock()
	defer mu.Unlock()

	if *signInCount > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", *signInCount))
		*signInCount--
		return nil
	}

	err := authenticationObj.SignOut()
	// Decrement regardless of error: this caller is leaving, and on signout
	// failure the cached session state is unreliable, so the next Authenticate
	// should establish a fresh one.
	*signInCount--
	if err != nil {
		return err
	}
	zapLogger.Debug(fmt.Sprintf("%v %v", "signout user", *signInCount))
	return nil
}

// ValidateChangeFrequencyDays validate Change Frequency Days field
func ValidateChangeFrequencyDays(changeFrequencyType string, changeFrequencyDays int) error {
	if changeFrequencyType == "xdays" {
		if changeFrequencyDays < 1 || changeFrequencyDays > 999 {
			return errors.New("error in change Frequency field, (min=1, max=999)")
		}
	}
	return nil
}

// DeleteAssetByID deletes an asset by its ID using the provided authentication object
func DeleteAssetByID(authenticationObj auth.AuthenticationObj, assetID int32, mu *sync.Mutex, signInCount *uint64, zapLogger logging.Logger) error {
	_, err := Authenticate(authenticationObj, mu, signInCount, zapLogger)
	if err != nil {
		return fmt.Errorf("error getting Authentication: %w", err)
	}

	// instantiating asset obj
	assetObj, err := assets.NewAssetObj(authenticationObj, zapLogger)
	if err != nil {
		return fmt.Errorf("error creating asset object: %w", err)
	}

	// deleting the asset by ID
	err = assetObj.DeleteAssetById(int(assetID))
	if err != nil {
		return fmt.Errorf("error deleting asset: %w", err)
	}

	err = SignOut(authenticationObj, mu, signInCount, zapLogger)
	if err != nil {
		return fmt.Errorf("error signing out: %w", err)
	}

	return nil
}
