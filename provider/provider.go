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

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"passwordsafe_secret":          getSecretByPath(),
			"passwordsafe_managed_account": getManagedAccount(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bt_verify_ca": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	apikey := d.Get("api_key").(string)
	url := d.Get("url").(string)
	accountname := d.Get("account_name").(string)
	btVerifyca := d.Get("bt_verify_ca").(bool)

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
	return client.NewClient(url, apikey, accountname, btVerifyca), diags

}

func getSecretByPath() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecret,
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

func getManagedAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceManagedAccount,
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

func dataSourceManagedAccount(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	apiClient := m.(*client.Client)

	system_name := d.Get("system_name").(string)
	account_name := d.Get("account_name").(string)

	secret, err := apiClient.ManageAccountFlow(system_name, account_name)

	if err != nil {
		return diag.FromErr(err)
		return diags
	}

	d.Set("value", secret)
	d.SetId(hash(secret))

	return diags
}

func dataSourceSecret(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	apiClient := m.(*client.Client)

	secretPath := d.Get("path").(string)
	secretTitle := d.Get("title").(string)
	separator := d.Get("separator").(string)

	secret, err := apiClient.SecretFlow(secretPath, secretTitle, separator)

	if err != nil {
		return diag.FromErr(err)
		return diags
	}

	d.Set("value", secret)
	d.SetId(hash(secret))

	return diags

}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
