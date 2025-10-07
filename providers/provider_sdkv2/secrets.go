// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"fmt"
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
				Sensitive: true,
			},
		},
	}
}

// resourceCredentialSecret Resource.
func resourceCredentialSecret() *schema.Resource {

	commonAttributes := getCreateSecretCommonSchema()
	credentialSecretAttributes := map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Sensitive: true,
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

	commonAttributes := getCreateSecretCommonSchema()
	textSecretAttributes := map[string]*schema.Schema{
		"text": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Sensitive: true,
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

	commonAttributes := getCreateSecretCommonSchema()
	fileSecretAttributes := map[string]*schema.Schema{
		"file_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"file_content": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			Sensitive: true,
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

	secretDetailsConfig := entities.SecretDetailsBaseConfig{
		Title:       title,
		Description: description,
		Urls:        getUrlsDetailsList(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
	}

	secretCredentialDetailsConfig30 := entities.SecretCredentialDetailsConfig30{
		SecretDetailsBaseConfig: secretDetailsConfig,
		Username:                username,
		Password:                password,
		OwnerId:                 ownerId,
		OwnerType:               ownerType,
		Owners:                  getOwnerDetailsOwnerIdList(d, ownerType, groupId, signAppinResponse),
	}

	secretCredentialDetailsConfig31 := entities.SecretCredentialDetailsConfig31{
		SecretDetailsBaseConfig: secretDetailsConfig,
		Username:                username,
		Password:                password,
		Owners:                  getOwnerDetailsGroupIdList(d, ownerType, groupId, signAppinResponse),
	}

	// Configure input object according to API version.
	configMap := map[string]interface{}{
		"3.0": secretCredentialDetailsConfig30,
		"3.1": secretCredentialDetailsConfig31,
	}

	credentialSecretDetails, exists := configMap[authenticationObj.ApiVersion]

	if !exists {
		return fmt.Errorf("unsupported API version: %v", authenticationObj.ApiVersion)
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, credentialSecretDetails)

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

	secretDetailsConfig := entities.SecretDetailsBaseConfig{
		Title:       title,
		Description: description,
		Urls:        getUrlsDetailsList(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
	}

	secretTextDetailsConfig30 := entities.SecretTextDetailsConfig30{
		SecretDetailsBaseConfig: secretDetailsConfig,
		Text:                    text,
		OwnerId:                 ownerId,
		OwnerType:               ownerType,
		Owners:                  getOwnerDetailsOwnerIdList(d, ownerType, groupId, signAppinResponse),
	}

	secretTextDetailsConfig31 := entities.SecretTextDetailsConfig31{
		SecretDetailsBaseConfig: secretDetailsConfig,
		Text:                    text,
		Owners:                  getOwnerDetailsGroupIdList(d, ownerType, groupId, signAppinResponse),
	}

	// Configure input object according to API version.
	configMap := map[string]interface{}{
		"3.0": secretTextDetailsConfig30,
		"3.1": secretTextDetailsConfig31,
	}

	textSecretDetails, exists := configMap[authenticationObj.ApiVersion]

	if !exists {
		return fmt.Errorf("unsupported API version: %v", authenticationObj.ApiVersion)
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, textSecretDetails)

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

	secretDetailsConfig := entities.SecretDetailsBaseConfig{
		Title:       title,
		Description: description,
		Urls:        getUrlsDetailsList(d, ownerType, groupId, signAppinResponse),
		Notes:       notes,
	}

	secretFileDetailsConfig30 := entities.SecretFileDetailsConfig30{
		SecretDetailsBaseConfig: secretDetailsConfig,
		FileContent:             fileContent,
		FileName:                fileName,
		OwnerId:                 ownerId,
		OwnerType:               ownerType,
		Owners:                  getOwnerDetailsOwnerIdList(d, ownerType, groupId, signAppinResponse),
	}

	secretFileDetailsConfig31 := entities.SecretFileDetailsConfig31{
		SecretDetailsBaseConfig: secretDetailsConfig,
		FileContent:             fileContent,
		FileName:                fileName,
		Owners:                  getOwnerDetailsGroupIdList(d, ownerType, groupId, signAppinResponse),
	}

	// Configure input object according to API version.
	configMap := map[string]interface{}{
		"3.0": secretFileDetailsConfig30,
		"3.1": secretFileDetailsConfig31,
	}

	fileSecretDetails, exists := configMap[authenticationObj.ApiVersion]

	if !exists {
		return fmt.Errorf("unsupported API version: %v", authenticationObj.ApiVersion)
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, fileSecretDetails)

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

// getOwnerDetailsOwnerIdList get Owners details list.
func getOwnerDetailsOwnerIdList(d *schema.ResourceData, ownerType string, groupId int, signAppinResponse entities.SignAppinResponse) []entities.OwnerDetailsOwnerId {
	var owners []entities.OwnerDetailsOwnerId

	mainOwner := entities.OwnerDetailsOwnerId{
		OwnerId: signAppinResponse.UserId,
		Owner:   signAppinResponse.UserName,
		Email:   signAppinResponse.EmailAddress,
	}
	owners = append(owners, mainOwner)

	ownersRaw, _ := d.GetOk("owners")
	if ownersRaw != nil {
		for _, ownerRaw := range ownersRaw.([]interface{}) {
			ownerMap := ownerRaw.(map[string]interface{})
			owner := entities.OwnerDetailsOwnerId{
				OwnerId: ownerMap["owner_id"].(int),
				Owner:   ownerMap["owner"].(string),
				Email:   ownerMap["email"].(string),
			}
			owners = append(owners, owner)
		}
	}

	return owners
}

// getOwnerDetailsGroupIdList get Owners details list.
func getOwnerDetailsGroupIdList(d *schema.ResourceData, ownerType string, groupId int, signAppinResponse entities.SignAppinResponse) []entities.OwnerDetailsGroupId {
	var owners []entities.OwnerDetailsGroupId

	mainOwner := entities.OwnerDetailsGroupId{
		GroupId: groupId,
		UserId:  signAppinResponse.UserId,
		Name:    signAppinResponse.Name,
		Email:   signAppinResponse.EmailAddress,
	}
	owners = append(owners, mainOwner)

	ownersRaw, _ := d.GetOk("owners")
	if ownersRaw != nil {
		for _, ownerRaw := range ownersRaw.([]interface{}) {
			ownerMap := ownerRaw.(map[string]interface{})
			owner := entities.OwnerDetailsGroupId{
				GroupId: groupId,
				UserId:  ownerMap["user_id"].(int),
				Name:    ownerMap["name"].(string),
				Email:   ownerMap["email"].(string),
			}
			owners = append(owners, owner)
		}
	}

	return owners
}

// getUrlsDetailsList get urls details list.
func getUrlsDetailsList(d *schema.ResourceData, ownerType string, groupId int, signAppinResponse entities.SignAppinResponse) []entities.UrlDetails {

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

// getCreateSecretCommonSchema get common attributes to create credential, file and text secrets.
func getCreateSecretCommonSchema() map[string]*schema.Schema {
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
