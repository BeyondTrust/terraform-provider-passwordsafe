// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"maps"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	managed_systems "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_systems"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &managedSystemByDatabaseResource{}
var _ resource.ResourceWithImportState = &managedSystemByDatabaseResource{}

func NewManagedSytemByDatabaseResource() resource.Resource {
	return &managedSystemByDatabaseResource{}
}

type managedSystemByDatabaseResource struct {
	providerInfo *ProviderData
}

type ManagedSystemByDataBaseResourceModel struct {
	DatabaseId                        types.String `tfsdk:"database_id"`
	ManagedSystemID                   types.Int32  `tfsdk:"managed_system_id"`
	ManagedSystemName                 types.String `tfsdk:"managed_system_name"`
	ContactEmail                      types.String `tfsdk:"contact_email"`
	Description                       types.String `tfsdk:"description"`
	Timeout                           types.Int32  `tfsdk:"timeout"`
	PasswordRuleID                    types.Int32  `tfsdk:"password_rule_id"`
	ReleaseDuration                   types.Int32  `tfsdk:"release_duration"`
	MaxReleaseDuration                types.Int32  `tfsdk:"max_release_duration"`
	ISAReleaseDuration                types.Int32  `tfsdk:"isa_release_duration"`
	AutoManagementFlag                types.Bool   `tfsdk:"auto_management_flag"`
	FunctionalAccountID               types.Int32  `tfsdk:"functional_account_id"`
	CheckPasswordFlag                 types.Bool   `tfsdk:"check_password_flag"`
	ChangePasswordAfterAnyReleaseFlag types.Bool   `tfsdk:"change_password_after_any_release_flag"`
	ResetPasswordOnMismatchFlag       types.Bool   `tfsdk:"reset_password_on_mismatch_flag"`
	ChangeFrequencyType               types.String `tfsdk:"change_frequency_type"`
	ChangeFrequencyDays               types.Int32  `tfsdk:"change_frequency_days"`
	ChangeTime                        types.String `tfsdk:"change_time"`
}

func (r *managedSystemByDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_system_by_database"
}

func (r *managedSystemByDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	commonAttributes := utils.GetCreateManagedSystemCommonAttributes()
	databaseAttributes := map[string]schema.Attribute{
		"database_id": schema.StringAttribute{
			MarkdownDescription: "Database Id",
			Required:            true,
		},
	}

	maps.Copy(databaseAttributes, commonAttributes)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Managed System by Database Id Resource, creates managed system by database id.",
		Attributes:          databaseAttributes,
	}
}

func (r *managedSystemByDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

// getManagedSystemObj get managedSystemObj for create manage system by asset, workgroup, database.
func getManagedSystemObj(changeFrequencyType string, changeFrequencyDays int, resp *resource.CreateResponse, authenticationObj authentication.AuthenticationObj) (*managed_systems.ManagedSystemObj, error) {
	err := utils.ValidateChangeFrequencyDays(changeFrequencyType, changeFrequencyDays)

	if err != nil {
		resp.Diagnostics.AddError("Error in inputs", err.Error())
		return nil, err
	}

	_, err = utils.Autenticate(authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return nil, err
	}

	// Instantiating managed system obj
	managedSystemObj, err := managed_systems.NewManagedSystem(authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating managed account object", err.Error())
		return nil, err
	}
	return managedSystemObj, nil
}

func (r *managedSystemByDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ManagedSystemByDataBaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Instantiating managed system obj
	managedSystemObj, err := getManagedSystemObj(data.ChangeFrequencyType.ValueString(), int(data.ChangeFrequencyDays.ValueInt32()), resp, *r.providerInfo.authenticationObj)

	if err != nil {
		return
	}

	// Instantiate object
	databaseDetailsBase := entities.ManagedSystemsByDatabaseIdDetailsBaseConfig{
		ContactEmail:                      data.ContactEmail.ValueString(),
		Description:                       data.Description.ValueString(),
		Timeout:                           int(data.Timeout.ValueInt32()),
		PasswordRuleID:                    int(data.PasswordRuleID.ValueInt32()),
		ReleaseDuration:                   int(data.ReleaseDuration.ValueInt32()),
		MaxReleaseDuration:                int(data.MaxReleaseDuration.ValueInt32()),
		ISAReleaseDuration:                int(data.ISAReleaseDuration.ValueInt32()),
		AutoManagementFlag:                data.AutoManagementFlag.ValueBool(),
		FunctionalAccountID:               int(data.FunctionalAccountID.ValueInt32()),
		CheckPasswordFlag:                 data.CheckPasswordFlag.ValueBool(),
		ChangePasswordAfterAnyReleaseFlag: data.ChangePasswordAfterAnyReleaseFlag.ValueBool(),
		ResetPasswordOnMismatchFlag:       data.ResetPasswordOnMismatchFlag.ValueBool(),
		ChangeFrequencyType:               data.ChangeFrequencyType.ValueString(),
		ChangeFrequencyDays:               int(data.ChangeFrequencyDays.ValueInt32()),
		ChangeTime:                        data.ChangeTime.ValueString(),
	}

	// creating a managed system by database Id.
	createdDataBase, err := managedSystemObj.CreateManagedSystemByDataBaseIdFlow(data.DatabaseId.ValueString(), databaseDetailsBase)

	if err != nil {
		resp.Diagnostics.AddError("Error creating managed system by database Id", err.Error())
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

func (r *managedSystemByDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *managedSystemByDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *managedSystemByDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// method not implemented
}

func (r *managedSystemByDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
