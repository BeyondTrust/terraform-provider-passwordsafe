// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"terraform-provider-passwordsafe/providers/entities"

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

// autenticate get Password Safe authentication.
func Autenticate(authenticationObj auth.AuthenticationObj, mu *sync.Mutex, signInCount *uint64, zapLogger logging.Logger) (libraryEntitites.SignAppinResponse, error) {
	var err error

	mu.Lock()
	if atomic.LoadUint64(signInCount) > 0 {
		atomic.AddUint64(signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "Already signed in", atomic.LoadUint64(signInCount)))
		mu.Unlock()

	} else {
		signAppinResponse, err = authenticationObj.GetPasswordSafeAuthentication()
		if err != nil {
			mu.Unlock()
			zapLogger.Error(err.Error())
			return libraryEntitites.SignAppinResponse{}, err
		}
		atomic.AddUint64(signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "signin", atomic.LoadUint64(signInCount)))
		mu.Unlock()
	}

	return signAppinResponse, nil
}

// signOut sign Password Safe out
func SignOut(authenticationObj auth.AuthenticationObj, muOut *sync.Mutex, signInCount *uint64, zapLogger logging.Logger) error {
	var err error

	muOut.Lock()
	if atomic.LoadUint64(signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(signInCount, ^uint64(0))
		muOut.Unlock()
	} else {
		err = authenticationObj.SignOut()
		if err != nil {
			return err
		}
		zapLogger.Debug(fmt.Sprintf("%v %v", "signout user", atomic.LoadUint64(signInCount)))
		// decrement counter
		atomic.AddUint64(signInCount, ^uint64(0))
		muOut.Unlock()

	}

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
