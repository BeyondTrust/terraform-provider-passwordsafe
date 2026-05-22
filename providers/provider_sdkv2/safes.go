// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"fmt"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSafe Resource.
func resourceSafe() *schema.Resource {
	return &schema.Resource{
		Description: "Safes Resource, creates safe",

		Create: resourceSafeCreate,
		Read:   resourceSafeRead,
		Update: resourceSafeUpdate,
		Delete: resourceSafeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

}

// Create context for resourceSafe Resource.
func resourceSafeCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)

	secretObj, _ := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	safe := entities.FolderDetails{
		Name:        name,
		Description: description,
		FolderType:  "SAFE",
	}

	createdSafe, err := secretObj.CreateFolderFlow("", safe)
	if err != nil {
		return err
	}

	d.SetId(createdSafe.Id.String())
	return nil
}

// Read context for resourceSafe Resource.
func resourceSafeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Update context for resourceSafe Resource.
func resourceSafeUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Delete context for resourceSafe Resource.
func resourceSafeDelete(d *schema.ResourceData, m interface{}) error {
	if m == nil {
		return fmt.Errorf("authentication object is nil")
	}

	meta := m.(*providerMeta)

	secretObj, err := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)
	if err != nil {
		return err
	}

	// Get the safe ID from the resource data
	safeID := d.Id()
	if safeID == "" {
		return fmt.Errorf("safe ID is empty")
	}

	// Delete the safe using the DeleteSafeById method
	err = secretObj.DeleteSafeById(safeID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
