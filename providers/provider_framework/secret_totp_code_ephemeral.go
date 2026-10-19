// Copyright 2026 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// totpApiVersion is the first Password Safe API version that returns TOTP codes.
const totpApiVersion = "3.3"

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

var _ ephemeral.EphemeralResource = &EphemeralSecretTotpCode{}

// @EphemeralResource(passwordsafe_secret_totp_code_ephemeral, name="Secret TOTP Code")
func NewEphemeralSecretTotpCode() ephemeral.EphemeralResource {
	return &EphemeralSecretTotpCode{}
}

type EphemeralSecretTotpCode struct {
	providerInfo *ProviderData
}

type EphemeralSecretTotpCodeModel struct {
	SecretId     types.String `tfsdk:"secret_id"`
	Code         types.String `tfsdk:"code"`
	ValidFromUtc types.String `tfsdk:"valid_from_utc"`
	ValidToUtc   types.String `tfsdk:"valid_to_utc"`
}

func (e *EphemeralSecretTotpCode) Metadata(ctx context.Context, request ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_secret_totp_code_ephemeral"
}

func (e *EphemeralSecretTotpCode) Schema(ctx context.Context, _ ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{

		MarkdownDescription: "Secret TOTP Code Ephemeral Resource, gets the current TOTP code of a credential secret. Requires API version 3.3.",

		Attributes: map[string]schema.Attribute{
			"secret_id": schema.StringAttribute{
				Description: "ID of a credential secret with TOTP enabled",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex, "must be a valid UUID"),
				},
			},
			"code": schema.StringAttribute{
				Description: "Current TOTP code",
				Computed:    true,
				Sensitive:   true,
			},
			"valid_from_utc": schema.StringAttribute{
				Description: "Start of the code validity window (inclusive), RFC 3339",
				Computed:    true,
			},
			"valid_to_utc": schema.StringAttribute{
				Description: "End of the code validity window (exclusive), RFC 3339",
				Computed:    true,
			},
		},
	}
}

func (e *EphemeralSecretTotpCode) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	e.providerInfo = &c
}

func (e *EphemeralSecretTotpCode) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {

	var data EphemeralSecretTotpCodeModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// The TOTP code endpoint only exists in API v3.3; fail with a clear message instead of a
	// server 404 that reads like "TOTP is not enabled".
	if e.providerInfo.apiVersion != totpApiVersion {
		response.Diagnostics.AddError("Error getting TOTP code", fmt.Sprintf("TOTP codes require api_version %v in the provider configuration, got: %q", totpApiVersion, e.providerInfo.apiVersion))
		return
	}

	if e.providerInfo.authenticationObj == nil {
		response.Diagnostics.AddError("Error getting TOTP code", "the provider is not configured")
		return
	}

	secretObj, err := secrets.NewSecretObj(*e.providerInfo.authenticationObj, zapLogger, maxFileSecretSizeBytes, false)

	if err != nil {
		response.Diagnostics.AddError("Error getting TOTP code", err.Error())
		return
	}

	totpCode, err := secretObj.GetSecretTotpCode(data.SecretId.ValueString())

	if err != nil {
		response.Diagnostics.AddError("Error getting TOTP code", err.Error())
		return
	}

	data.Code = types.StringValue(totpCode.Code)
	data.ValidFromUtc = types.StringValue(totpCode.ValidFromUtc.UTC().Format(time.RFC3339))
	data.ValidToUtc = types.StringValue(totpCode.ValidToUtc.UTC().Format(time.RFC3339))

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)
}
