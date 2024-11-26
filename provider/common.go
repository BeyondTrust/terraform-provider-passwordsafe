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
