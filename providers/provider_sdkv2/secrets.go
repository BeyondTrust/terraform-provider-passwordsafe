// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"fmt"
	"maps"
	"strings"

	libconstants "github.com/BeyondTrust/go-client-library-passwordsafe/api/constants"
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
			"decrypt": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to decrypt the secret value when retrieving it.",
			},
			"value": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
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
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"totp_enabled": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether TOTP is enabled for the secret. Requires API version 3.3 and `totp_secret_key`.",
		},
		"totp_secret_key": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "TOTP seed: a Base32 secret or an otpauth:// URI (max 1024 characters). Requires API version 3.3 and `totp_enabled = true`. Changing it replaces the seed.",
		},
	}

	maps.Copy(credentialSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "Credential secret Resource, creates credential secret.",
		Create:      resourceCredentialSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceCredentialSecretUpdate,
		Delete:      resourceSecretDelete,

		Schema: credentialSecretAttributes,
	}

}

// resourceTextSecret Resource.
func resourceTextSecret() *schema.Resource {

	commonAttributes := getCreateSecretCommonSchema()
	textSecretAttributes := map[string]*schema.Schema{
		"text": &schema.Schema{
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}

	maps.Copy(textSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "Text secret Resource, creates text secret.",
		Create:      resourceTextSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceTextSecretUpdate,
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
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}

	maps.Copy(fileSecretAttributes, commonAttributes)

	return &schema.Resource{
		Description: "File secret Resource, creates file secret.",
		Create:      resourceFileSecretCreate,
		Read:        resourceSecretRead,
		Update:      resourceFileSecretUpdate,
		Delete:      resourceSecretDelete,
		Schema:      fileSecretAttributes,
	}

}

// Create context for resourceCredentialSecret Resource.
func resourceCredentialSecretCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)
	folderName := d.Get("folder_name").(string)

	secretObj, _ := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)
	totpSecretKey, _ := d.Get("totp_secret_key").(string)

	credentialSecretInput := entities.SecretCredentialInput{
		SecretDetailsBaseConfig: entities.SecretDetailsBaseConfig{
			Title:       title,
			Description: description,
			Urls:        getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
			Notes:       notes,
		},
		Username:        username,
		Password:        password,
		TotpEnabled:     getTotpEnabledOnCreate(d),
		TotpSecretKey:   totpSecretKey,
		OwnerId:         ownerId,
		OwnerType:       ownerType,
		OwnersByOwnerId: getOwnerDetailsOwnerIdList(d, ownerType, groupId, meta.signAppin),
		OwnersByGroupId: getOwnerDetailsGroupIdList(d, ownerType, groupId, meta.signAppin),
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, credentialSecretInput)
	if err != nil {
		return err
	}

	d.SetId(createdSecret.Id)
	return nil
}

// Create context for resourceTextSecret Resource.
func resourceTextSecretCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)
	folderName := d.Get("folder_name").(string)

	secretObj, _ := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)

	text := d.Get("text").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	textSecretInput := entities.SecretTextInput{
		SecretDetailsBaseConfig: entities.SecretDetailsBaseConfig{
			Title:       title,
			Description: description,
			Urls:        getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
			Notes:       notes,
		},
		Text:            text,
		OwnerId:         ownerId,
		OwnerType:       ownerType,
		OwnersByOwnerId: getOwnerDetailsOwnerIdList(d, ownerType, groupId, meta.signAppin),
		OwnersByGroupId: getOwnerDetailsGroupIdList(d, ownerType, groupId, meta.signAppin),
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, textSecretInput)
	if err != nil {
		return err
	}

	d.SetId(createdSecret.Id)
	return nil
}

// Create context for resourceFileSecret Resource.
func resourceFileSecretCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)
	folderName := d.Get("folder_name").(string)

	secretObj, _ := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)

	fileName := d.Get("file_name").(string)
	fileContent := d.Get("file_content").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	ownerId := d.Get("owner_id").(int)
	groupId := d.Get("group_id").(int)
	ownerType := d.Get("owner_type").(string)
	notes := d.Get("notes").(string)

	fileSecretInput := entities.SecretFileInput{
		SecretDetailsBaseConfig: entities.SecretDetailsBaseConfig{
			Title:       title,
			Description: description,
			Urls:        getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
			Notes:       notes,
		},
		FileName:        fileName,
		FileContent:     fileContent,
		OwnerId:         ownerId,
		OwnerType:       ownerType,
		OwnersByOwnerId: getOwnerDetailsOwnerIdList(d, ownerType, groupId, meta.signAppin),
		OwnersByGroupId: getOwnerDetailsGroupIdList(d, ownerType, groupId, meta.signAppin),
	}

	createdSecret, err := secretObj.CreateSecretFlow(folderName, fileSecretInput)
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

// getTotpEnabledOnCreate returns the TotpEnabled value to create a credential secret with.
// TOTP is off by default, so false is not sent: this keeps credential secrets creatable on
// API versions below 3.3, which reject any TOTP field.
func getTotpEnabledOnCreate(d *schema.ResourceData) *bool {
	if enabled, _ := d.Get("totp_enabled").(bool); !enabled {
		return nil
	}
	enabled := true
	return &enabled
}

// getTotpUpdate returns the TOTP fields of a credential secret update. When neither TOTP
// attribute changed, nothing is sent (nil) and the server leaves the TOTP setting untouched.
// Otherwise totp_enabled is sent as is (unset means false, which disables TOTP). The seed is
// sent only when enabling: the API rejects a seed alongside TotpEnabled=false, so a retained
// totp_secret_key must not block disabling TOTP.
func getTotpUpdate(d *schema.ResourceData) (*bool, string) {
	if !d.HasChanges("totp_enabled", "totp_secret_key") {
		return nil, ""
	}
	enabled := d.Get("totp_enabled").(bool)
	if !enabled {
		return &enabled, ""
	}
	return &enabled, d.Get("totp_secret_key").(string)
}

// newSecretUpdateObj builds the secret client and resolves the ID of the secret's folder,
// which the update PUT requires. folder_name is ForceNew, so it is the current folder.
func newSecretUpdateObj(d *schema.ResourceData, meta *providerMeta) (*secrets.SecretObj, uuid.UUID, error) {
	secretObj, err := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)
	if err != nil {
		return nil, uuid.Nil, err
	}

	folderName := d.Get("folder_name").(string)
	folders, err := secretObj.SecretGetFoldersAndSafes("secrets-safe/folders/", libconstants.SecretGetFolders)
	if err != nil {
		return nil, uuid.Nil, err
	}

	for _, folder := range folders {
		if folder.Name == strings.TrimSpace(folderName) {
			folderId, err := uuid.Parse(folder.Id)
			return secretObj, folderId, err
		}
	}

	return nil, uuid.Nil, fmt.Errorf("folder %v was not found in folder list", folderName)
}

// getSecretUpdateOwners builds the owners of a secret update the same way create does; the
// library sends the owners model of the session API version.
func getSecretUpdateOwners(d *schema.ResourceData, signAppin entities.SignAppinResponse) entities.SecretUpdateOwners {
	ownerType := d.Get("owner_type").(string)
	groupId := d.Get("group_id").(int)

	return entities.SecretUpdateOwners{
		OwnerId:         d.Get("owner_id").(int),
		OwnerType:       ownerType,
		OwnersByOwnerId: getOwnerDetailsOwnerIdList(d, ownerType, groupId, signAppin),
		OwnersByGroupId: getOwnerDetailsGroupIdList(d, ownerType, groupId, signAppin),
	}
}

// failUpdate keeps the prior state when an update fails, so attributes that were not applied
// are planned again on the next run.
func failUpdate(d *schema.ResourceData, err error) error {
	d.Partial(true)
	return err
}

// Update context for resourceCredentialSecret Resource.
// The PUT replaces the secret, so every attribute is sent.
func resourceCredentialSecretUpdate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)

	secretObj, folderId, err := newSecretUpdateObj(d, meta)
	if err != nil {
		return failUpdate(d, err)
	}

	totpEnabled, totpSecretKey := getTotpUpdate(d)
	ownerType := d.Get("owner_type").(string)
	groupId := d.Get("group_id").(int)

	details := entities.SecretCredentialUpdateConfig33{
		FolderId:           folderId,
		Title:              d.Get("title").(string),
		Description:        d.Get("description").(string),
		Notes:              d.Get("notes").(string),
		Urls:               getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		SecretUpdateOwners: getSecretUpdateOwners(d, meta.signAppin),
		TotpEnabled:        totpEnabled,
		TotpSecretKey:      totpSecretKey,
	}

	if _, err = secretObj.UpdateCredentialSecret(d.Id(), details); err != nil {
		return failUpdate(d, err)
	}
	return nil
}

// Update context for resourceTextSecret Resource.
// The PUT replaces the secret, so every attribute is sent.
func resourceTextSecretUpdate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)

	secretObj, folderId, err := newSecretUpdateObj(d, meta)
	if err != nil {
		return failUpdate(d, err)
	}

	ownerType := d.Get("owner_type").(string)
	groupId := d.Get("group_id").(int)

	details := entities.SecretTextUpdateConfig{
		FolderId:           folderId,
		Title:              d.Get("title").(string),
		Description:        d.Get("description").(string),
		Notes:              d.Get("notes").(string),
		Urls:               getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
		Text:               d.Get("text").(string),
		SecretUpdateOwners: getSecretUpdateOwners(d, meta.signAppin),
	}

	if _, err = secretObj.UpdateTextSecret(d.Id(), details); err != nil {
		return failUpdate(d, err)
	}
	return nil
}

// Update context for resourceFileSecret Resource.
// The multipart PUT replaces the secret, so every attribute and the file content are sent.
func resourceFileSecretUpdate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*providerMeta)

	secretObj, folderId, err := newSecretUpdateObj(d, meta)
	if err != nil {
		return failUpdate(d, err)
	}

	ownerType := d.Get("owner_type").(string)
	groupId := d.Get("group_id").(int)

	details := entities.SecretFileUpdateConfig{
		FolderId:           folderId,
		Title:              d.Get("title").(string),
		Description:        d.Get("description").(string),
		Notes:              d.Get("notes").(string),
		Urls:               getUrlsDetailsList(d, ownerType, groupId, meta.signAppin),
		SecretUpdateOwners: getSecretUpdateOwners(d, meta.signAppin),
		FileName:           d.Get("file_name").(string),
		FileContent:        d.Get("file_content").(string),
	}

	if err = secretObj.UpdateFileSecret(d.Id(), details); err != nil {
		return failUpdate(d, err)
	}
	return nil
}

// Delete context for resourceSecret Resource.
func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	if m == nil {
		return fmt.Errorf("authentication object is nil")
	}

	meta := m.(*providerMeta)

	secretObj, err := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, false)
	if err != nil {
		return err
	}

	// Get the secret ID from the resource data
	secretID := d.Id()
	if secretID == "" {
		return fmt.Errorf("secret ID is empty")
	}

	// Delete the secret using the proper secret deletion method
	err = secretObj.DeleteSecretById(secretID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// Read context for getSecretByPath Datasource.
func getSecretByPathReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	meta := m.(*providerMeta)

	secretPath := d.Get("path").(string)
	secretTitle := d.Get("title").(string)
	separator := d.Get("separator").(string)
	decrypt := d.Get("decrypt").(bool)
	coordinate := secretPath + separator + secretTitle

	secretObj, _ := secrets.NewSecretObj(*meta.authObj, zapLogger, 5000000, decrypt)
	secret, err := secretObj.GetSecret(coordinate, separator)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("value", secret); err != nil {
		return diag.FromErr(err)
	}

	// Derive the resource ID from the non-secret secret coordinates
	// (path + separator + title) rather than from the secret value. Terraform
	// resource IDs are never redacted and are printed verbatim to console/CI
	// logs, so a secret-derived ID would leak crackable material about the
	// secret value.
	d.SetId(coordinate)

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
			// The API does not move a secret to another folder on update.
			ForceNew: true,
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
