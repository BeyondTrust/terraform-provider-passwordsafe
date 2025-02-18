// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"fmt"
	"sync/atomic"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// autenticate get Password Safe authentication.
func autenticate(d *schema.ResourceData, m interface{}) (entities.SignApinResponse, error) {
	authenticationObj := m.(*auth.AuthenticationObj)
	var err error

	mu.Lock()
	if atomic.LoadUint64(&signInCount) > 0 {
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "Already signed in", atomic.LoadUint64(&signInCount)))
		mu.Unlock()

	} else {
		signApinResponse, err = authenticationObj.GetPasswordSafeAuthentication()
		if err != nil {
			mu.Unlock()
			zapLogger.Error(err.Error())
			return entities.SignApinResponse{}, err
		}
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "signin", atomic.LoadUint64(&signInCount)))
		mu.Unlock()
	}

	return signApinResponse, nil
}

// signOut sign Password Safe out
func signOut(d *schema.ResourceData, m interface{}) error {
	authenticationObj := m.(*auth.AuthenticationObj)
	var err error

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

	return nil
}

// getOwnersSchema get Owners schema.
func getOwnersSchema() *schema.Schema {

	schema := schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"owner_id": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
				"owner": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"email": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return &schema
}

// getUrlsSchema get Urls schema.
func getUrlsSchema() *schema.Schema {

	schema := schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"credential_id": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"url": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	return &schema
}

func getManagedAccountSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
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
	}

	return schema
}
