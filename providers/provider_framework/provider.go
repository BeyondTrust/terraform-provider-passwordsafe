// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/logging"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/utils"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.uber.org/zap"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	maxFileSecretSizeBytes     = 5000000
	clientTimeOutInSeconds     = 45
	separator                  = "/"
	retryMaxElapsedTimeMinutes = 15
)

var signInCount uint64
var mu sync.Mutex
var muOut sync.Mutex

type PasswordSafeProvider struct {
}

// Define the zap configuration
var config = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Encoding:         "console", // You can use "json" for structured logging
	EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	OutputPaths:      []string{"stderr", "providerFramework.log"}, // Logs to both stderr and the file
	ErrorOutputPaths: []string{"stderr"},
}

// Build the logger with the configuration
var logger, _ = config.Build()

// create a zap logger wrapper
var zapLogger = logging.NewZapLogger(logger)

func NewProvider() provider.Provider {
	return &PasswordSafeProvider{}
}

func (p *PasswordSafeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "passwordsafe"
}

type ProviderModel struct {
	APIKey                       types.String `tfsdk:"api_key"`
	ClientId                     types.String `tfsdk:"client_id"`
	ClientSecret                 types.String `tfsdk:"client_secret"`
	Url                          types.String `tfsdk:"url"`
	APIVersion                   types.String `tfsdk:"api_version"`
	APIAccountName               types.String `tfsdk:"api_account_name"`
	VerifyCA                     types.Bool   `tfsdk:"verify_ca"`
	ClientCertificatesFolderPath types.String `tfsdk:"client_certificates_folder_path"`
	ClientCertificateName        types.String `tfsdk:"client_certificate_name"`
	ClientCertificatePassword    types.String `tfsdk:"client_certificate_password"`
}

type ProviderData struct {
	//authenticate authentication.AuthenticationObj
	apiKey                    string
	clientId                  string
	clientSecret              string
	apiVersion                string
	url                       string
	accountname               string
	clientCertificatePath     string
	clientCertificateName     string
	clientCertificatePassword string
	verifyca                  bool
	userName                  string
	authenticationObj         *auth.AuthenticationObj
}

func (p *PasswordSafeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Description: "The API key for making requests to the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "API OAuth Client ID.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "API OAuth Client Secret.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL for the Password Safe instance from which to request a secret.",
			},
			"api_version": schema.StringAttribute{
				Optional:    true,
				Description: "The recommended version is 3.1. If no version is specified, the default API version 3.0 will be used",
			},
			"api_account_name": schema.StringAttribute{
				Required:    true,
				Description: "The user name for the API request to the Password Safe instance. For use when authenticating with an API key.",
			},
			"verify_ca": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether to verify the certificate authority on the Password Safe instance. For use when authenticating to Password Safe.",
			},
			"client_certificates_folder_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to the Client Certificate associated with the Password Safe instance for use when authenticating with an API key using a Client Certificate.",
			},
			"client_certificate_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the Client Certificate for use when authenticating with an API key using a Client Certificate.",
			},
			"client_certificate_password": schema.StringAttribute{
				Optional:    true,
				Description: "The password associated with the Client Certificate. For use when authenticating with an API key using a Client Certificate",
			},
		},
	}
}

// ValidateCredentialsAndConfig make basic validations in credential and config data.
func (p *PasswordSafeProvider) ValidateCredentialsAndConfig(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse, data ProviderData) error {
	var errorSummary string

	if data.apiKey == "" && data.clientId == "" && data.clientSecret == "" {
		errorSummary = "Invalid Authentication method"
		resp.Diagnostics.AddError("Invalid Authentication method", "Please add a valid credential (API Key / Client Credentials)")
		return errors.New(errorSummary)
	}

	if data.url == "" {
		errorSummary = "Invalid URL"
		resp.Diagnostics.AddError(errorSummary, "Please add a proper URL")
		return errors.New(errorSummary)
	}

	if data.apiKey != "" && data.accountname == "" {
		errorSummary = "Invalid Account Name"
		resp.Diagnostics.AddError(errorSummary, "Please add a proper Account Name")
		return errors.New(errorSummary)
	}
	return nil
}

// GetCertificateData decrypt pfx file to get certificate and certificate key data.
func (p *PasswordSafeProvider) GetCertificateData(resp *provider.ConfigureResponse, data ProviderData) (string, string) {
	if data.clientCertificateName != "" {
		certificate, certificateKey, err := utils.GetPFXContent(data.clientCertificatePath, data.clientCertificateName, data.clientCertificatePassword, zapLogger)
		if err != nil {
			resp.Diagnostics.AddError("Error in certificate", err.Error())
		}
		return certificate, certificateKey
	}
	return "", ""
}

func (p *PasswordSafeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data ProviderModel

	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	providerData := ProviderData{
		apiKey:                    strings.TrimSpace(data.APIKey.ValueString()),
		clientId:                  data.ClientId.ValueString(),
		clientSecret:              data.ClientSecret.ValueString(),
		url:                       strings.TrimSpace(data.Url.ValueString()),
		apiVersion:                data.APIVersion.ValueString(),
		accountname:               strings.TrimSpace(data.APIAccountName.ValueString()),
		verifyca:                  data.VerifyCA.ValueBool(),
		clientCertificatePath:     data.ClientCertificatesFolderPath.ValueString(),
		clientCertificateName:     data.ClientCertificateName.ValueString(),
		clientCertificatePassword: data.ClientCertificatePassword.ValueString(),
	}

	// Make basic validations.
	err = p.ValidateCredentialsAndConfig(ctx, req, resp, providerData)
	if err != nil {
		return
	}

	// Get Cerificate and certificate key.
	certificate, certificateKey := p.GetCertificateData(resp, providerData)

	// Create an instance of ValidationParams
	params := utils.ValidationParams{
		ClientID:                   providerData.clientId,
		ClientSecret:               providerData.clientSecret,
		ApiUrl:                     &providerData.url,
		ApiVersion:                 providerData.apiVersion,
		ClientTimeOutInSeconds:     clientTimeOutInSeconds,
		VerifyCa:                   providerData.verifyca,
		Logger:                     zapLogger,
		Certificate:                certificate,
		CertificateKey:             certificateKey,
		RetryMaxElapsedTimeMinutes: &retryMaxElapsedTimeMinutes,
	}

	// validate inputs
	errorsInInputs := utils.ValidateInputs(params)

	if errorsInInputs != nil {
		resp.Diagnostics.AddError("Error in inputs validation", errorsInInputs.Error())
		return
	}

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.InitialInterval = 1 * time.Second
	backoffDefinition.MaxElapsedTime = time.Duration(retryMaxElapsedTimeMinutes) * time.Second
	backoffDefinition.RandomizationFactor = 0.5

	// creating a http client
	httpClientObj, _ := utils.GetHttpClient(clientTimeOutInSeconds, providerData.verifyca, certificate, certificateKey, zapLogger)

	var authenticate *auth.AuthenticationObj

	// authenticate using api key
	if providerData.apiKey != "" {
		authParamsApiKey := &auth.AuthenticationParametersObj{
			HTTPClient:                 *httpClientObj,
			BackoffDefinition:          backoffDefinition,
			EndpointURL:                strings.TrimSpace(data.Url.ValueString()),
			APIVersion:                 data.APIVersion.ValueString(),
			ClientID:                   "",
			ClientSecret:               "",
			ApiKey:                     fmt.Sprintf("%v;runas=%v;", strings.TrimSpace(data.APIKey.ValueString()), strings.TrimSpace(data.APIAccountName.ValueString())),
			Logger:                     zapLogger,
			RetryMaxElapsedTimeSeconds: 30,
		}
		authenticate, err = auth.AuthenticateUsingApiKey(*authParamsApiKey)
		if err != nil {
			resp.Diagnostics.AddError("Error in Provider", err.Error())
		}
	} else {
		// authenticate using client_id and client secret
		authParamsOauth := &auth.AuthenticationParametersObj{
			HTTPClient:                 *httpClientObj,
			BackoffDefinition:          backoffDefinition,
			EndpointURL:                strings.TrimSpace(data.Url.ValueString()),
			APIVersion:                 data.APIVersion.ValueString(),
			ClientID:                   data.ClientId.ValueString(),
			ClientSecret:               data.ClientSecret.ValueString(),
			ApiKey:                     "",
			Logger:                     zapLogger,
			RetryMaxElapsedTimeSeconds: 30,
		}
		authenticate, err = auth.Authenticate(*authParamsOauth)
		if err != nil {
			resp.Diagnostics.AddError("Error in Provider", err.Error())
		}
	}

	// authenticating
	userObject, err := authenticate.GetPasswordSafeAuthentication()
	if err != nil {
		resp.Diagnostics.AddError("Error in Provider", err.Error())
		return
	}

	providerData.userName = userObject.UserName

	// pass authentication obj to ephemeral resources
	providerData.authenticationObj = authenticate

	// pass data to ephemeral resources
	resp.EphemeralResourceData = providerData
	// pass data to ephemeral resources
	resp.ResourceData = providerData
	// pass data to ephemeral resources
	resp.DataSourceData = providerData

}

func (p *PasswordSafeProvider) Functions(_ context.Context) []func() function.Function {
	return nil
}

func (p *PasswordSafeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFunctionalAccountDataResource,
		NewFolderDataSource,
		NewSafesDataSource,
		NewDatabaseDataSource,
		NewWorkgroupDataSource,
		NewPlatformDataSource,
		NewManagedAccountDataSource,
		NewManagedSystemDataSource,
		NewAssetDataSource,
	}
}

func (p *PasswordSafeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkGroupResource,
		NewAssetByWorkgGroypIdResource,
		NewAssetByWorkGroupNameResource,
		NewDatabaseResource,
		NewManagedSytemByAssetResource,
		NewManagedSytemByWorkGroupResource,
		NewManagedSytemByDatabaseResource,
		NewFunctionalAccountResource,
	}
}

func (p *PasswordSafeProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewEphemeralSecret,
		NewEphemeralManagedAccount,
	}
}

var _ provider.Provider = &PasswordSafeProvider{}
var _ provider.ProviderWithFunctions = &PasswordSafeProvider{}
var _ provider.ProviderWithEphemeralResources = &PasswordSafeProvider{}
