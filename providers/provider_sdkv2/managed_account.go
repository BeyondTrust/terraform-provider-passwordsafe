// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getManagedAccount DataSource.
func getManagedAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Managed Account Datasource, gets managed account.",
		ReadContext: getManagedAccountReadContext,
		Schema: map[string]*schema.Schema{
			"system_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// resourceManagedAccount Resource.
func resourceManagedAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Managed Account Resource, creates managed account.",
		Create:      resourceManagedAccountCreate,
		Read:        resourceManagedAccountRead,
		Update:      resourceManagedAccountUpdate,
		Delete:      resourceManagedAccountDelete,

		Schema: getManagedAccountSchema(),
	}
}

// Create context for resourceManagedAccount Resource.
func resourceManagedAccountCreate(d *schema.ResourceData, m interface{}) error {

	authenticationObj := m.(*auth.AuthenticationObj)
	system_name := d.Get("system_name").(string)

	_, err := authenticate(d, m)
	if err != nil {
		return err
	}

	manageAccountObj, _ := managed_accounts.NewManagedAccountObj(*authenticationObj, zapLogger)

	accountDetailsObj := entities.AccountDetails{
		AccountName:                       d.Get("account_name").(string),
		Password:                          d.Get("password").(string),
		DomainName:                        d.Get("domain_name").(string),
		UserPrincipalName:                 d.Get("user_principal_name").(string),
		SAMAccountName:                    d.Get("sam_account_name").(string),
		DistinguishedName:                 d.Get("distinguished_name").(string),
		PrivateKey:                        d.Get("private_key").(string),
		Passphrase:                        d.Get("passphrase").(string),
		PasswordFallbackFlag:              d.Get("password_fallback_flag").(bool),
		LoginAccountFlag:                  d.Get("login_account_flag").(bool),
		Description:                       d.Get("description").(string),
		ApiEnabled:                        d.Get("api_enabled").(bool),
		ReleaseNotificationEmail:          d.Get("release_notification_email").(string),
		ChangeServicesFlag:                d.Get("change_services_flag").(bool),
		RestartServicesFlag:               d.Get("restart_services_flag").(bool),
		ChangeTasksFlag:                   d.Get("change_tasks_flag").(bool),
		ReleaseDuration:                   d.Get("release_duration").(int),
		MaxReleaseDuration:                d.Get("max_release_duration").(int),
		ISAReleaseDuration:                d.Get("isa_release_duration").(int),
		MaxConcurrentRequests:             d.Get("max_concurrent_requests").(int),
		AutoManagementFlag:                d.Get("auto_management_flag").(bool),
		DSSAutoManagementFlag:             d.Get("dss_auto_management_flag").(bool),
		CheckPasswordFlag:                 d.Get("check_password_flag").(bool),
		ResetPasswordOnMismatchFlag:       d.Get("reset_password_on_mismatch_flag").(bool),
		ChangePasswordAfterAnyReleaseFlag: d.Get("change_password_after_any_release_flag").(bool),
		ChangeFrequencyType:               d.Get("change_frequency_type").(string),
		ChangeFrequencyDays:               d.Get("change_frequency_days").(int),
		ChangeTime:                        d.Get("change_time").(string),
		NextChangeDate:                    d.Get("next_change_date").(string),
		UseOwnCredentials:                 d.Get("use_own_credentials").(bool),
		ChangeWindowsAutoLogonFlag:        d.Get("change_windows_auto_logon_flag").(bool),
		ChangeComPlusFlag:                 d.Get("change_com_plus_flag").(bool),
		ChangeDComFlag:                    d.Get("change_dcom_flag").(bool),
		ChangeSComFlag:                    d.Get("change_scom_flag").(bool),
		ObjectID:                          d.Get("object_id").(string),
	}

	var createResponse entities.CreateManagedAccountsResponse
	createResponse, err = manageAccountObj.ManageAccountCreateFlow(system_name, accountDetailsObj)

	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", createResponse.ManagedAccountID))
	return nil
}

// Read context for resourceManagedAccount Resource.
func resourceManagedAccountRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Update context for resourceManagedAccount Resource.
func resourceManagedAccountUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

// Delete context for resourceManagedAccount Resource.
func resourceManagedAccountDelete(d *schema.ResourceData, m interface{}) error {
	if m == nil {
		return fmt.Errorf("authentication object is nil")
	}

	authenticationObj := m.(*auth.AuthenticationObj)

	_, err := authenticate(d, m)
	if err != nil {
		return err
	}

	manageAccountObj, err := managed_accounts.NewManagedAccountObj(*authenticationObj, zapLogger)
	if err != nil {
		return err
	}

	// Get the managed account ID from the resource data
	managedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = manageAccountObj.DeleteManagedAccountById(managedAccountID)
	if err != nil {
		return err
	}

	err = signOut(d, m)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// Read context for getManagedAccount Datasource.
func getManagedAccountReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	authenticationObj := m.(*auth.AuthenticationObj)

	system_name := d.Get("system_name").(string)
	account_name := d.Get("account_name").(string)

	_, err := authenticate(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	manageAccountObj, _ := managed_accounts.NewManagedAccountObj(*authenticationObj, zapLogger)
	gotManagedAccount, err := manageAccountObj.GetSecret(system_name+"/"+account_name, "/")

	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("value", gotManagedAccount)

	if err != nil {
		return diag.FromErr(err)
	}

	err = signOut(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash(gotManagedAccount))
	return diags
}

// hash function.
func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
