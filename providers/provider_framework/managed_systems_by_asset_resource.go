// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"fmt"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	managed_systems "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_systems"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &managedSystemResource{}
var _ resource.ResourceWithImportState = &managedSystemResource{}

func NewManagedSytemByAssetResource() resource.Resource {
	return &managedSystemResource{}
}

type managedSystemResource struct {
	providerInfo *ProviderData
}

type ManagedSystemResourceModel struct {
	AssetId                           types.String `tfsdk:"asset_id"`
	ManagedSystemID                   types.Int32  `tfsdk:"managed_system_id"`
	ManagedSystemName                 types.String `tfsdk:"managed_system_name"`
	PlatformID                        types.Int32  `tfsdk:"platform_id"`
	ContactEmail                      types.String `tfsdk:"contact_email"`
	Description                       types.String `tfsdk:"description"`
	Port                              types.Int32  `tfsdk:"port"`
	Timeout                           types.Int32  `tfsdk:"timeout"`
	SshKeyEnforcementMode             types.Int32  `tfsdk:"ssh_key_enforcement_mode"`
	PasswordRuleID                    types.Int32  `tfsdk:"password_rule_id"`
	DSSKeyRuleID                      types.Int32  `tfsdk:"dss_key_rule_id"`
	LoginAccountID                    types.Int32  `tfsdk:"login_account_id"`
	ReleaseDuration                   types.Int32  `tfsdk:"release_duration"`
	MaxReleaseDuration                types.Int32  `tfsdk:"max_release_duration"`
	ISAReleaseDuration                types.Int32  `tfsdk:"isa_release_duration"`
	AutoManagementFlag                types.Bool   `tfsdk:"auto_management_flag"`
	FunctionalAccountID               types.Int32  `tfsdk:"functional_account_id"`
	ElevationCommand                  types.String `tfsdk:"elevation_command"`
	CheckPasswordFlag                 types.Bool   `tfsdk:"check_password_flag"`
	ChangePasswordAfterAnyReleaseFlag types.Bool   `tfsdk:"change_password_after_any_release_flag"`
	ResetPasswordOnMismatchFlag       types.Bool   `tfsdk:"reset_password_on_mismatch_flag"`
	ChangeFrequencyType               types.String `tfsdk:"change_frequency_type"`
	ChangeFrequencyDays               types.Int32  `tfsdk:"change_frequency_days"`
	ChangeTime                        types.String `tfsdk:"change_time"`
	RemoteClientType                  types.String `tfsdk:"remote_client_type"`
	ApplicationHostID                 types.Int32  `tfsdk:"application_host_id"`
	IsApplicationHost                 types.Bool   `tfsdk:"is_application_host"`
}

func (r *managedSystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_system_by_asset"
}

func (r *managedSystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Managed System by Asset Id Resource, creates managed system by asset id.",

		Attributes: map[string]schema.Attribute{
			"asset_id": schema.StringAttribute{
				MarkdownDescription: "Asset Id",
				Required:            true,
			},
			"managed_system_id": schema.Int32Attribute{
				MarkdownDescription: "Managed System Id",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"managed_system_name": schema.StringAttribute{
				MarkdownDescription: "Managed System Name",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"platform_id": schema.Int32Attribute{
				MarkdownDescription: "Platform ID",
				Required:            true,
			},
			"contact_email": schema.StringAttribute{
				MarkdownDescription: "Contact Email (max 1000 characters)",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description (max 255 characters)",
				Optional:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "Port number",
				Optional:            true,
			},

			"timeout": schema.Int32Attribute{
				MarkdownDescription: "Timeout",
				Optional:            true,
			},
			"ssh_key_enforcement_mode": schema.Int32Attribute{
				MarkdownDescription: "SSH Key Enforcement Mode (one of: 0, 1, 2)",
				Optional:            true,
			},
			"password_rule_id": schema.Int32Attribute{
				MarkdownDescription: "Password Rule ID",
				Optional:            true,
			},
			"dss_key_rule_id": schema.Int32Attribute{
				MarkdownDescription: "DSS Key Rule ID",
				Optional:            true,
			},
			"login_account_id": schema.Int32Attribute{
				MarkdownDescription: "Login Account ID",
				Optional:            true,
			},
			"release_duration": schema.Int32Attribute{
				MarkdownDescription: "Release Duration (min: 1, max: 525600)",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(120),
			},
			"max_release_duration": schema.Int32Attribute{
				MarkdownDescription: "Max Release Duration (min: 1, max: 525600)",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(525600),
			},
			"isa_release_duration": schema.Int32Attribute{
				MarkdownDescription: "ISA Release Duration (min: 1, max: 525600)",
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(120),
			},
			"auto_management_flag": schema.BoolAttribute{
				MarkdownDescription: "Auto Management Flag",
				Optional:            true,
			},
			"functional_account_id": schema.Int32Attribute{
				MarkdownDescription: "Functional Account ID (required if AutoManagementFlag is true)",
				Optional:            true,
			},
			"elevation_command": schema.StringAttribute{
				MarkdownDescription: "Elevation Command",
				Optional:            true,
			},
			"check_password_flag": schema.BoolAttribute{
				MarkdownDescription: "Check Password Flag",
				Optional:            true,
			},
			"change_password_after_any_release_flag": schema.BoolAttribute{
				MarkdownDescription: "Change Password After Any Release Flag",
				Optional:            true,
			},
			"reset_password_on_mismatch_flag": schema.BoolAttribute{
				MarkdownDescription: "Reset Password On Mismatch Flag",
				Optional:            true,
			},
			"change_frequency_type": schema.StringAttribute{
				MarkdownDescription: "Change Frequency Type (one of: first, last, xdays)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("first"),
			},
			"change_frequency_days": schema.Int32Attribute{
				MarkdownDescription: "Change Frequency Days (required if ChangeFrequencyType is xdays)",
				Optional:            true,
			},
			"change_time": schema.StringAttribute{
				MarkdownDescription: "Change Time (format: HH:MM)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("23:30"),
			},
			"remote_client_type": schema.StringAttribute{
				MarkdownDescription: "Remote Client Type (one of: None, EPM)",
				Optional:            true,
			},
			"application_host_id": schema.Int32Attribute{
				MarkdownDescription: "Application Host ID",
				Optional:            true,
			},
			"is_application_host": schema.BoolAttribute{
				MarkdownDescription: "Is Application Host",
				Optional:            true,
			},
		},
	}

}

func (r *managedSystemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

func (r *managedSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ManagedSystemResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := utils.ValidateChangeFrequencyDays(data.ChangeFrequencyType.ValueString(), int(data.ChangeFrequencyDays.ValueInt32()))

	if err != nil {
		resp.Diagnostics.AddError("Error in inputs", err.Error())
		return
	}

	_, err = utils.Autenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// Instantiating managed system obj
	managedSystemObj, err := managed_systems.NewManagedSystem(*r.providerInfo.authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating managed account object", err.Error())
		return
	}

	// Instantiate base config object.
	databaseDetailsBase := entities.ManagedSystemsByAssetIdDetailsBaseConfig{
		PlatformID:                        int(data.PlatformID.ValueInt32()),
		ContactEmail:                      data.ContactEmail.ValueString(),
		Description:                       data.Description.ValueString(),
		Port:                              int(data.Port.ValueInt32()),
		Timeout:                           int(data.Timeout.ValueInt32()),
		SshKeyEnforcementMode:             int(data.SshKeyEnforcementMode.ValueInt32()),
		PasswordRuleID:                    int(data.PasswordRuleID.ValueInt32()),
		DSSKeyRuleID:                      int(data.DSSKeyRuleID.ValueInt32()),
		LoginAccountID:                    int(data.LoginAccountID.ValueInt32()),
		ReleaseDuration:                   int(data.ReleaseDuration.ValueInt32()),
		MaxReleaseDuration:                int(data.MaxReleaseDuration.ValueInt32()),
		ISAReleaseDuration:                int(data.ISAReleaseDuration.ValueInt32()),
		AutoManagementFlag:                data.AutoManagementFlag.ValueBool(),
		FunctionalAccountID:               int(data.FunctionalAccountID.ValueInt32()),
		ElevationCommand:                  data.ElevationCommand.ValueString(),
		CheckPasswordFlag:                 data.CheckPasswordFlag.ValueBool(),
		ChangePasswordAfterAnyReleaseFlag: data.ChangePasswordAfterAnyReleaseFlag.ValueBool(),
		ResetPasswordOnMismatchFlag:       data.ResetPasswordOnMismatchFlag.ValueBool(),
		ChangeFrequencyType:               data.ChangeFrequencyType.ValueString(),
		ChangeFrequencyDays:               int(data.ChangeFrequencyDays.ValueInt32()),
		ChangeTime:                        data.ChangeTime.ValueString(),
	}

	// API Version 3.0 input object
	ManagedSystemsByAssetIdDetailsConfig30 := entities.ManagedSystemsByAssetIdDetailsConfig30{
		ManagedSystemsByAssetIdDetailsBaseConfig: databaseDetailsBase,
	}

	// API Version 3.1 input object
	ManagedSystemsByAssetIdDetailsConfig31 := entities.ManagedSystemsByAssetIdDetailsConfig31{
		ManagedSystemsByAssetIdDetailsBaseConfig: databaseDetailsBase,
		RemoteClientType:                         data.RemoteClientType.ValueString(),
	}

	// API Version 3.2 input object
	ManagedSystemsByAssetIdDetailsConfig32 := entities.ManagedSystemsByAssetIdDetailsConfig32{
		ManagedSystemsByAssetIdDetailsBaseConfig: databaseDetailsBase,
		RemoteClientType:                         data.RemoteClientType.ValueString(),
		ApplicationHostID:                        int(data.ApplicationHostID.ValueInt32()),
		IsApplicationHost:                        data.IsApplicationHost.ValueBool(),
	}

	// Configure input object according to API version.
	configMap := map[string]interface{}{
		"3.0": ManagedSystemsByAssetIdDetailsConfig30,
		"3.1": ManagedSystemsByAssetIdDetailsConfig31,
		"3.2": ManagedSystemsByAssetIdDetailsConfig32,
	}

	databaseDetails, exists := configMap[r.providerInfo.apiVersion]

	if !exists {
		resp.Diagnostics.AddError("Invalid API Version", fmt.Sprintf("Unsupported API version: %s", r.providerInfo.apiVersion))
		return
	}

	// creating a managed system.
	createdDataBase, err := managedSystemObj.CreateManagedSystemByAssetIdFlow(data.AssetId.ValueString(), databaseDetails)

	if err != nil {
		resp.Diagnostics.AddError("Error creating managed system by Asset Id", err.Error())
		return
	}

	data.ManagedSystemID = types.Int32Value(int32(createdDataBase.ManagedSystemID))
	data.ManagedSystemName = types.StringValue(createdDataBase.SystemName)

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *managedSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *managedSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *managedSystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// method not implemented
}

func (r *managedSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
