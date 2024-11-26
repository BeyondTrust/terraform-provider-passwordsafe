// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync/atomic"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getManagedAccount DataSource.
func getManagedAccount() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceManagedAccount Resource.
func resourceManagedAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedAccountCreate,
		Read:   resourceManagedAccountRead,
		Update: resourceManagedAccountUpdate,
		Delete: resourceManagedAccountDelete,

		Schema: map[string]*schema.Schema{
			"system_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_principal_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sam_account_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"distinguished_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"passphrase": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_fallback_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"login_account_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_rule_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"api_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"release_notification_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_services_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"restart_services_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_tasks_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"release_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_release_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"isa_release_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrent_requests": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_management_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dss_auto_management_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"check_password_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_password_after_any_release_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reset_password_on_mismatch_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_frequency_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_frequency_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"change_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_change_date": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_own_credentials": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"workgroup_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"change_windows_auto_logon_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_com_plus_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_dcom_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"change_scom_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// Create context for resourceManagedAccount Resource.
func resourceManagedAccountCreate(d *schema.ResourceData, m interface{}) error {

	authenticationObj := m.(*auth.AuthenticationObj)
	system_name := d.Get("system_name").(string)

	_, err := autenticate(d, m)
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

	_, err = manageAccountObj.ManageAccountCreateFlow(system_name, accountDetailsObj)

	if err != nil {
		return err
	}

	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
	} else {
		err = authenticationObj.SignOut()
		if err != nil {
			return err
		}
		zapLogger.Debug(fmt.Sprintf("%v %v", "signout user", atomic.LoadUint64(&signInCount)))
		// decrement counter
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()

	}

	d.SetId(accountDetailsObj.AccountName)
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
	return nil
}

// Read context for getManagedAccount Datasource.
func getManagedAccountReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	authenticationObj := m.(*auth.AuthenticationObj)

	system_name := d.Get("system_name").(string)
	account_name := d.Get("account_name").(string)

	_, err := autenticate(d, m)
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

	d.SetId(hash(gotManagedAccount))

	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
	} else {
		err = authenticationObj.SignOut()
		if err != nil {
			return diag.FromErr(err)
		}
		zapLogger.Debug(fmt.Sprintf("%v %v", "signout user", atomic.LoadUint64(&signInCount)))
		// decrement counter
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()

	}

	return diags
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

	d.SetId(hash(secret))

	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
	} else {
		err = authenticationObj.SignOut()
		if err != nil {
			return diag.FromErr(err)
		}
		zapLogger.Debug(fmt.Sprintf("%v %v", "signout user", atomic.LoadUint64(&signInCount)))
		// decrement counter
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()

	}

	return diags
}

// hash function.
func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
