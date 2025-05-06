// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceFolder Resource.
func resourceFolder() *schema.Resource {
	return &schema.Resource{
		Description: "Folder Resource, creates folder",
		Create:      resourceFolderCreate,
		Read:        resourceFolderRead,
		Update:      resourceFolderUpdate,
		Delete:      resourceFolderDelete,

		Schema: map[string]*schema.Schema{
			"parent_folder_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}

}

// Create context for resourceFolder Resource.
func resourceFolderCreate(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	parent_folder_name := d.Get("parent_folder_name").(string)

	_, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	userGroupId := d.Get("user_group_id").(int)

	folder := entities.FolderDetails{
		Name:        name,
		Description: description,
		UserGroupId: userGroupId,
		FolderType:  "FOLDER",
	}

	createdFolder, err := secretObj.CreateFolderFlow(parent_folder_name, folder)

	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId(createdFolder.Id.String())
	return nil
}

// Read context for resourceFolder Resource.
func resourceFolderRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Update context for resourceFolder Resource.
func resourceFolderUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Delete context for resourceFolder Resource.
func resourceFolderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
