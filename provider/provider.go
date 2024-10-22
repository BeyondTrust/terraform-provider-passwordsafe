// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	backoff "github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var signInCount uint64
var mu sync.Mutex
var mu_out sync.Mutex

// Define the zap configuration
var config = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Encoding:         "console", // You can use "json" for structured logging
	EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	OutputPaths:      []string{"stderr", "ProviderLogs.log"}, // Logs to both stderr and the file
	ErrorOutputPaths: []string{"stderr"},
}

// Build the logger with the configuration
var logger, _ = config.Build()

// create a zap logger wrapper
var zapLogger = logging.NewZapLogger(logger)

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
				Optional:    true,
				Description: "The api key for making requests to the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API OAuth Client ID.",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: " API OAuth Client Secret.",
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
	client_id := d.Get("client_id").(string)
	client_secret := d.Get("client_secret").(string)
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

	if apikey == "" && client_id == "" && client_secret == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Authentication method",
			Detail:   "Please add a valid credential (API Key / Client Credentials)",
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

	if accountname == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Account Name",
			Detail:   "Please add a proper Account Name",
		})
		return nil, diags
	}

	retryMaxElapsedTimeMinutes := 2

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.InitialInterval = 1 * time.Second
	backoffDefinition.MaxElapsedTime = time.Duration(retryMaxElapsedTimeMinutes) * time.Second
	backoffDefinition.RandomizationFactor = 0.5

	certificate := ""
	certificateKey := ""
	var err error = nil

	if clientCertificateName != "" {
		certificate, certificateKey, err = utils.GetPFXContent(clientCertificatePath, clientCertificateName, clientCertificatePassword, zapLogger)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	httpClientObj, err := utils.GetHttpClient(45, verifyca, certificate, certificateKey, zapLogger)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// If this variable is set, we're using API Key authentication
	// (previous/old authentication method)
	if apikey != "" {
		authenticate, err := auth.AuthenticateUsingApiKey(*httpClientObj, backoffDefinition, d.Get("url").(string), zapLogger, retryMaxElapsedTimeMinutes, fmt.Sprintf("%v;runas=%v;", apikey, accountname))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return authenticate, diags
	}

	authenticate, err := auth.Authenticate(*httpClientObj, backoffDefinition, d.Get("url").(string), client_id, client_secret, zapLogger, retryMaxElapsedTimeMinutes)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return authenticate, diags

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

	authenticationObj := m.(*auth.AuthenticationObj)

	system_name := d.Get("system_name").(string)
	account_name := d.Get("account_name").(string)

	mu.Lock()
	if atomic.LoadUint64(&signInCount) > 0 {
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "Already signed in", atomic.LoadUint64(&signInCount)))
		mu.Unlock()

	} else {
		_, err := authenticationObj.GetPasswordSafeAuthentication()
		if err != nil {
			mu.Unlock()
			zapLogger.Error(err.Error())
			return diag.FromErr(err)
		}
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "signin", atomic.LoadUint64(&signInCount)))
		mu.Unlock()
	}

	manageAccountObj, _ := managed_accounts.NewManagedAccountObj(*authenticationObj, zapLogger)
	gotManagedAccount, err := manageAccountObj.GetSecret(system_name+"/"+account_name, "/")

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("value", gotManagedAccount)
	d.SetId(hash(gotManagedAccount))

	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
	} else {
		authenticationObj.SignOut()
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

	mu.Lock()
	if atomic.LoadUint64(&signInCount) > 0 {
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "Already signed in", atomic.LoadUint64(&signInCount)))
		mu.Unlock()

	} else {
		_, err := authenticationObj.GetPasswordSafeAuthentication()
		if err != nil {
			mu.Unlock()
			zapLogger.Error(err.Error())
			return diag.FromErr(err)
		}
		atomic.AddUint64(&signInCount, 1)
		zapLogger.Debug(fmt.Sprintf("%v %v", "signin", atomic.LoadUint64(&signInCount)))
		mu.Unlock()
	}

	secretObj, _ := secrets.NewSecretObj(*authenticationObj, zapLogger, 5000000)
	secret, err := secretObj.GetSecret(secretPath+separator+secretTitle, separator)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("value", secret)
	d.SetId(hash(secret))

	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		zapLogger.Debug(fmt.Sprintf("%v %v", "Ignore signout", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
	} else {
		authenticationObj.SignOut()
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
