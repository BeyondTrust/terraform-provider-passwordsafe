// Copyright 2023 BeyondTrust. All rights reserved.
// Package Provider implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	backoff "github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var signInCount uint64
var mu sync.Mutex
var muOut sync.Mutex

// Define the zap configuration
var config = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Encoding:         "console", // You can use "json" for structured logging
	EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	OutputPaths:      []string{"stderr", "providerSdkv2.log"}, // Logs to both stderr and the file
	ErrorOutputPaths: []string{"stderr"},
}

// Build the logger with the configuration
var logger, _ = config.Build()

// create a zap logger wrapper
var zapLogger = logging.NewZapLogger(logger)

// Provider Definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"passwordsafe_managed_account":   resourceManagedAccount(),
			"passwordsafe_credential_secret": resourceCredentialSecret(),
			"passwordsafe_text_secret":       resourceTextSecret(),
			"passwordsafe_file_secret":       resourceFileSecret(),
			"passwordsafe_folder":            resourceFolder(),
			"passwordsafe_safe":              resourceSafe(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"passwordsafe_secret":          getSecretByPath(),
			"passwordsafe_managed_account": getManagedAccount(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API key for making requests to the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API OAuth Client ID.",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API OAuth Client Secret.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL for the Password Safe instance from which to request a secret.",
			},
			"api_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The recommended version is 3.1. If no version is specified, the default API version 3.0 will be used",
			},
			"api_account_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name for the API request to the Password Safe instance. For use when authenticating with an API key.",
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
				Description: "The path to the Client Certificate associated with the Password Safe instance for use when authenticating with an API key using a Client Certificate.",
			},
			"client_certificate_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the Client Certificate for use when authenticating with an API key using a Client Certificate.",
			},
			"client_certificate_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The password associated with the Client Certificate. For use when authenticating with an API key using a Client Certificate",
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// ValidateCredentialsAndConfig make basic validations in credential and config data.
func ValidateCredentialsAndConfig(apikey string, clientId string, clientSecret string, url string, accountName string) diag.Diagnostics {
	var diags diag.Diagnostics
	if apikey == "" && clientId == "" && clientSecret == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Authentication method",
			Detail:   "Please add a valid credential (API Key / Client Credentials)",
		})
		return diags
	}

	if url == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid URL",
			Detail:   "Please add a proper URL",
		})
		return diags
	}

	if apikey != "" && accountName == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Account Name",
			Detail:   "Please add a proper Account Name",
		})
		return diags
	}
	return nil
}

// Provider Init Config.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	apikey := d.Get("api_key").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	url := d.Get("url").(string)
	apiVersion := d.Get("api_version").(string)
	accountName := d.Get("api_account_name").(string)
	verifyca := d.Get("verify_ca").(bool)
	clientCertificatePath := d.Get("client_certificates_folder_path").(string)
	clientCertificateName := d.Get("client_certificate_name").(string)
	clientCertificatePassword := d.Get("client_certificate_password").(string)

	apikey = strings.TrimSpace(apikey)
	url = strings.TrimSpace(url)
	accountName = strings.TrimSpace(accountName)

	// Make basic validations.
	diags := ValidateCredentialsAndConfig(apikey, clientId, clientSecret, url, accountName)
	if diags != nil {
		return nil, diags
	}

	retryMaxElapsedTimeMinutes := 2
	clientTimeOutInSeconds := 30

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.InitialInterval = 1 * time.Second
	backoffDefinition.MaxElapsedTime = time.Duration(retryMaxElapsedTimeMinutes) * time.Second
	backoffDefinition.RandomizationFactor = 0.5

	certificate := ""
	certificateKey := ""
	var err error

	if clientCertificateName != "" {
		certificate, certificateKey, err = utils.GetPFXContent(clientCertificatePath, clientCertificateName, clientCertificatePassword, zapLogger)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	// Create an instance of ValidationParams
	params := utils.ValidationParams{
		ClientID:                   clientId,
		ClientSecret:               clientSecret,
		ApiUrl:                     &url,
		ApiVersion:                 apiVersion,
		ClientTimeOutInSeconds:     clientTimeOutInSeconds,
		VerifyCa:                   verifyca,
		Logger:                     zapLogger,
		Certificate:                certificate,
		CertificateKey:             certificateKey,
		RetryMaxElapsedTimeMinutes: &retryMaxElapsedTimeMinutes,
	}

	// validate inputs
	errorsInInputs := utils.ValidateInputs(params)

	if errorsInInputs != nil {
		return nil, diag.FromErr(err)
	}

	httpClientObj, err := utils.GetHttpClient(45, verifyca, certificate, certificateKey, zapLogger)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// If this variable is set, we're using API Key authentication
	// (previous/old authentication method)
	if apikey != "" {
		authParamsApiKey := &auth.AuthenticationParametersObj{
			HTTPClient:                 *httpClientObj,
			BackoffDefinition:          backoffDefinition,
			EndpointURL:                url,
			APIVersion:                 apiVersion,
			ClientID:                   "",
			ClientSecret:               "",
			ApiKey:                     fmt.Sprintf("%v;runas=%v;", apikey, accountName),
			Logger:                     zapLogger,
			RetryMaxElapsedTimeSeconds: retryMaxElapsedTimeMinutes,
		}
		authenticate, err := auth.AuthenticateUsingApiKey(*authParamsApiKey)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return authenticate, diags
	}

	authParamsOauth := &auth.AuthenticationParametersObj{
		HTTPClient:                 *httpClientObj,
		BackoffDefinition:          backoffDefinition,
		EndpointURL:                url,
		APIVersion:                 apiVersion,
		ClientID:                   clientId,
		ClientSecret:               clientSecret,
		ApiKey:                     "",
		Logger:                     zapLogger,
		RetryMaxElapsedTimeSeconds: retryMaxElapsedTimeMinutes,
	}
	authenticate, err := auth.Authenticate(*authParamsOauth)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return authenticate, diags

}
