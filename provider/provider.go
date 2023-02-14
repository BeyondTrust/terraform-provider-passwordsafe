// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"terraform-provider-passwordsafe/api/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider Definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"passwordsafe_secret":          getSecretByPath(),
			"passwordsafe_managed_account": getManagedAccount(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The api key for making requests to the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL for the Password Safe instance from which to request a secret.",
			},
			"api_account_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name for the api request to the Password Safe instance. For use when authenticating with an api key.",
			},
			"verify_ca": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether to verify the certificate authority on the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"client_certificates_folder_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path to the Client Certificate associated with the Password Safe instance for use when authenticating with an api key using a Client Certificate.",
			},
			"client_certificate_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the Client Certificate for use when authenticating with an api key using a Client Certificate.",
			},
			"client_certificate_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The password associated with the Client Certificate. For use when authenticating with an api key using a Client Certificate",
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Provider Init Config.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	apikey := d.Get("api_key").(string)
	url := d.Get("url").(string)
	accountname := d.Get("api_account_name").(string)
	verifyca := d.Get("verify_ca").(bool)
	clientCertificatePath := d.Get("client_certificates_folder_path").(string)
	clientCertificateName := d.Get("client_certificate_name").(string)
	clientCertificatePassword := d.Get("client_certificate_password").(string)

	apikey = strings.TrimSpace(apikey)
	url = strings.TrimSpace(url)
	accountname = strings.TrimSpace(accountname)

	var diags diag.Diagnostics

	if apikey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid apiKey",
			Detail:   "Please add a proper Apikey",
		})
		return nil, diags
	}
	if url == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid URL",
			Detail:   "Please add a proper URL",
		})
		return nil, diags
	}

	if accountname == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Account Name",
			Detail:   "Please add a proper Account Name",
		})
		return nil, diags
	}

	apiClient, err := client.NewClient(url, apikey, accountname, verifyca, clientCertificatePath, clientCertificateName, clientCertificatePassword)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return apiClient, diags

}

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

// Read context for getManagedAccount Datasource.
func getManagedAccountReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	apiClient := m.(*client.Client)

	system_name := d.Get("system_name").(string)
	account_name := d.Get("account_name").(string)

	paths := make(map[string]string)
	secret, err := apiClient.ManageAccountFlow(system_name, account_name, paths)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("value", secret)
	d.SetId(hash(secret))

	return diags
}

// Read context for getSecretByPath Datasource.
func getSecretByPathReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	apiClient := m.(*client.Client)

	secretPath := d.Get("path").(string)
	secretTitle := d.Get("title").(string)
	separator := d.Get("separator").(string)

	paths := make(map[string]string)
	secret, err := apiClient.SecretFlow(secretPath, secretTitle, separator, paths)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("value", secret)
	d.SetId(hash(secret))

	return diags

}

// hash function.
func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
