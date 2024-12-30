// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"

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
	return &schema.Resource{
		Create: resourceCredentialSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"folder_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"owner_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owners": getOwnersSchema(),
			"password_rule_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"urls": getUrlsSchema(),
		},
	}

}

// resourceTextSecret Resource.
func resourceTextSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTextSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"folder_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"text": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"owner_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owners": getOwnersSchema(),
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"urls": getUrlsSchema(),
		},
	}

}

// resourceFileSecret Resource.
func resourceFileSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"folder_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"owner_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"owners": getOwnersSchema(),
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"urls": getUrlsSchema(),
		},
	}

}

// Create context for resourceCredentialSecret Resource.
func resourceCredentialSecretCreate(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	folderName := d.Get("folder_name").(string)

	signApinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	title := d.Get("title").(string)
	description := d.Get("description").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	ownerId := d.Get("owner_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	var owners []entities.OwnerDetails

	if ownerType == "User" {
		mainOwner := entities.OwnerDetails{
			OwnerId: signApinResponse.UserId,
			Owner:   signApinResponse.UserName,
			Email:   signApinResponse.EmailAddress,
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

	secret := entities.SecretCredentialDetails{
		Title:       title,
		Description: description,
		Username:    username,
		Password:    password,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      owners,
		Notes:       notes,
		Urls:        urls,
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

	signApinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	title := d.Get("title").(string)
	description := d.Get("description").(string)
	text := d.Get("text").(string)
	ownerId := d.Get("owner_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	var owners []entities.OwnerDetails

	if ownerType == "User" {
		mainOwner := entities.OwnerDetails{
			OwnerId: signApinResponse.UserId,
			Owner:   signApinResponse.UserName,
			Email:   signApinResponse.EmailAddress,
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

	secret := entities.SecretTextDetails{
		Title:       title,
		Description: description,
		Text:        text,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      owners,
		Notes:       notes,
		Urls:        urls,
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

	signApinResponse, err := autenticate(d, m)
	if err != nil {
		return err
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)

	title := d.Get("title").(string)
	description := d.Get("description").(string)
	fileContent := d.Get("file_content").(string)
	ownerId := d.Get("owner_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)
	file_name := d.Get("file_name").(string)

	var owners []entities.OwnerDetails

	if ownerType == "User" {
		mainOwner := entities.OwnerDetails{
			OwnerId: signApinResponse.UserId,
			Owner:   signApinResponse.UserName,
			Email:   signApinResponse.EmailAddress,
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

	secret := entities.SecretFileDetails{
		Title:       title,
		Description: description,
		OwnerId:     ownerId,
		OwnerType:   ownerType,
		Owners:      owners,
		Notes:       notes,
		Urls:        urls,
		FileContent: fileContent,
		FileName:    file_name,
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
