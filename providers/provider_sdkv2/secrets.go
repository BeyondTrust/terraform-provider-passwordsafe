// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"maps"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getSecretByPath DataSource.
func getSecretByPath() *schema.Resource {
	return &schema.Resource{
		Description: "Secret Datasource, get secret.",
		ReadContext: getSecretByPathReadContext,
		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"separator": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceCredentialSecret Resource.
func resourceCredentialSecret() *schema.Resource {

	commonAttributes := getCreateSecretSchema()
	credentialSecretAttributes := map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}

	maps.Copy(credentialSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "Credential secret Resource, creates credential secret.",
		Create:      resourceCredentialSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceSecretUpdate,
		Delete:      resourceSecretDelete,

		Schema: credentialSecretAttributes,
	}

}

// resourceTextSecret Resource.
func resourceTextSecret() *schema.Resource {

	commonAttributes := getCreateSecretSchema()
	textSecretAttributes := map[string]*schema.Schema{
		"text": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}

	maps.Copy(textSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "Text secret Resource, creates text secret.",
		Create:      resourceTextSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceSecretUpdate,
		Delete:      resourceSecretDelete,
		Schema:      textSecretAttributes,
	}

}

// resourceFileSecret Resource.
func resourceFileSecret() *schema.Resource {

	commonAttributes := getCreateSecretSchema()
	fileSecretAttributes := map[string]*schema.Schema{
		"file_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"file_content": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}

	maps.Copy(fileSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "File secret Resource, creates file secret.",
		Create:      resourceFileSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceSecretUpdate,
		Delete:      resourceSecretDelete,
		Schema:      fileSecretAttributes,
	}

}

// Create context for resourceCredentialSecret Resource.
func resourceCredentialSecretCreate(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	folderName := d.Get("folder_name").(string)

	signAppinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	secret := entities.SecretCredentialDetails{
		Title:       title,
		Description: description,
		Username:    username,
		Password:    password,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      getOwnerDetails(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
		Urls:        getUrlDetails(d, ownerType, groupId, signAppinResponse),
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, secret)

	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId(createdSecret.Id)
	return nil
}

// Create context for resourceTextSecret Resource.
func resourceTextSecretCreate(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	folderName := d.Get("folder_name").(string)

	signAppinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	text := d.Get("text").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	secret := entities.SecretTextDetails{
		Title:       title,
		Description: description,
		Text:        text,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      getOwnerDetails(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
		Urls:        getUrlDetails(d, ownerType, groupId, signAppinResponse),
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, secret)

	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId(createdSecret.Id)
	return nil
}

// Create context for resourceFileSecret Resource.
func resourceFileSecretCreate(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	folderName := d.Get("folder_name").(string)

	signAppinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	fileName := d.Get("file_name").(string)
	fileContent := d.Get("file_content").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	secret := entities.SecretFileDetails{
		Title:       title,
		Description: description,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      getOwnerDetails(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
		Urls:        getUrlDetails(d, ownerType, groupId, signAppinResponse),
		FileContent: fileContent,
		FileName:    fileName,
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, secret)

	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId(createdSecret.Id)
	return nil
}

// Read context for resourceSecret Resource.
func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Update context for resourceSecret Resource.
func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Delete context for resourceSecret Resource.
func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Read context for getSecretByPath Datasource.
func getSecretByPathReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	authenticationObj := m.(*auth.AuthenticationObj)

	secretPath := d.Get("path").(string)
	secretTitle := d.Get("title").(string)
	separator := d.Get("separator").(string)

	_, err := autenticate(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)
	secret, err := secretObj.GetSecret(secretPath+separator+secretTitle, separator)

	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("value", secret)

	if err != nil {
		return diag.FromErr(err)
	}

	err = signOut(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash(secret))

	return diags
}

func getOwnerDetails(d *schema.ResourceData, ownerType string, groupId int, signAppinResponse entities.SignAppinResponse) []entities.OwnerDetails {
	var owners []entities.OwnerDetails

	if ownerType == "User" {
		mainOwner := entities.OwnerDetails{
			GroupId: groupId,
			OwnerId: signAppinResponse.UserId,
			Owner:   signAppinResponse.UserName,
			Email:   signAppinResponse.EmailAddress,
		}
		owners = append(owners, mainOwner)
	}

	ownersRaw, _ := d.GetOk("owners")
	if ownersRaw != nil {
		for _, ownerRaw := range ownersRaw.([]interface{}) {
			ownerMap := ownerRaw.(map[string]interface{})
			owner := entities.OwnerDetails{
				OwnerId: ownerMap["owner_id"].(int),
				Owner:   ownerMap["owner"].(string),
				Email:   ownerMap["email"].(string),
			}
			owners = append(owners, owner)
		}
	}

	urlsRaw, _ := d.GetOk("urls")
	var urls []entities.UrlDetails
	if urlsRaw != nil {
		for _, urlRaw := range urlsRaw.([]interface{}) {
			urlMap := urlRaw.(map[string]interface{})

			id, _ := uuid.Parse(urlMap["id"].(string))
			credentialId, _ := uuid.Parse(urlMap["credential_id"].(string))

			url := entities.UrlDetails{
				Id:           id,
				CredentialId: credentialId,
				Url:          urlMap["url"].(string),
			}
			urls = append(urls, url)
		}
	}

	return owners
}

func getUrlDetails(d *schema.ResourceData, ownerType string, groupId int, signAppinResponse entities.SignAppinResponse) []entities.UrlDetails {

	urlsRaw, _ := d.GetOk("urls")
	var urls []entities.UrlDetails
	if urlsRaw != nil {
		for _, urlRaw := range urlsRaw.([]interface{}) {
			urlMap := urlRaw.(map[string]interface{})

			id, _ := uuid.Parse(urlMap["id"].(string))
			credentialId, _ := uuid.Parse(urlMap["credential_id"].(string))

			url := entities.UrlDetails{
				Id:           id,
				CredentialId: credentialId,
				Url:          urlMap["url"].(string),
			}
			urls = append(urls, url)
		}
	}

	return urls
}

func getCreateSecretSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"folder_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"title": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"owner_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"group_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"owner_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"owners": getOwnersSchema(),
		"notes": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"urls": getUrlsSchema(),
	}
}
