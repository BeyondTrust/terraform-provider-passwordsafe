// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"errors"
	"fmt"
	"terraform-provider-passwordsafe/providers/entities"
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

// ValidateChangeFrequencyDays validate Change Frequency Days field
func ValidateChangeFrequencyDays(changeFrequencyType string, changeFrequencyDays int) error {
	if changeFrequencyType == "xdays" {
		if changeFrequencyDays < 1 || changeFrequencyDays > 999 {
			return errors.New("error in change Frequency field, (min=1, max=999)")
		}
	}
	return nil
}
