// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils provides common utilities for Terraform provider operations.
package utils

import (
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	managed_systems "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_systems"
)

// DeleteManagedSystemByID is a helper function to delete a managed system by ID
// This follows Terraform provider patterns by extracting only the API call logic
func DeleteManagedSystemByID(authenticationObj authentication.AuthenticationObj, managedSystemID int, zapLogger logging.Logger) error {
	// instantiating managed system obj
	managedSystemObj, err := managed_systems.NewManagedSystem(authenticationObj, zapLogger)
	if err != nil {
		return err
	}

	err = managedSystemObj.DeleteManagedSystemById(managedSystemID)
	if err != nil {
		return err
	}

	return nil
}
