// Copyright 2026 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"time"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SecretTotpCodeDataSource{}

func NewSecretTotpCodeDataSource() datasource.DataSource {
	return &SecretTotpCodeDataSource{}
}

type SecretTotpCodeDataSource struct {
	providerInfo *ProviderData
}

type SecretTotpCodeDataSourceModel struct {
	SecretId     types.String `tfsdk:"secret_id"`
	Code         types.String `tfsdk:"code"`
	ValidFromUtc types.String `tfsdk:"valid_from_utc"`
	ValidToUtc   types.String `tfsdk:"valid_to_utc"`
}

func (d *SecretTotpCodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_totp_code"
}

func (d *SecretTotpCodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Secret TOTP Code Datasource, gets the current TOTP code of a credential secret. Requires API version 3.3. The code is stored in the state as a sensitive value and is read again on every plan and refresh.",

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

func (d *SecretTotpCodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c
}

func (d *SecretTotpCodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data SecretTotpCodeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if d.providerInfo.authenticationObj == nil {
		resp.Diagnostics.AddError("Error getting TOTP code", "the provider is not configured")
		return
	}

	secretObj, err := secrets.NewSecretObj(*d.providerInfo.authenticationObj, zapLogger, maxFileSecretSizeBytes, false)

	if err != nil {
		resp.Diagnostics.AddError("Error getting TOTP code", err.Error())
		return
	}

	totpCode, err := secretObj.GetSecretTotpCode(data.SecretId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error getting TOTP code", err.Error())
		return
	}

	data.Code = types.StringValue(totpCode.Code)
	data.ValidFromUtc = types.StringValue(totpCode.ValidFromUtc.UTC().Format(time.RFC3339))
	data.ValidToUtc = types.StringValue(totpCode.ValidToUtc.UTC().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
