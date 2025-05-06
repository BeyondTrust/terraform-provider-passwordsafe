// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/databases"
)

var _ resource.Resource = &databaseResource{}
var _ resource.ResourceWithImportState = &databaseResource{}

func NewDatabaseResource() resource.Resource {
	return &databaseResource{}
}

type databaseResource struct {
	providerInfo *ProviderData
}

type DatabaseResourceModel struct {
	AssetId           types.String `tfsdk:"asset_id"`
	PlatformID        types.Int32  `tfsdk:"platform_id"`
	InstanceName      types.String `tfsdk:"instance_name"`
	IsDefaultInstance types.Bool   `tfsdk:"is_default_instance"`
	Port              types.Int32  `tfsdk:"port"`
	Version           types.String `tfsdk:"version"`
	Template          types.String `tfsdk:"template"`
	DatabaseID        types.Int32  `tfsdk:"database_id"`
}

func (r *databaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *databaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database Resource, creates database.",

		Attributes: map[string]schema.Attribute{
			"asset_id": schema.StringAttribute{
				MarkdownDescription: "Asset Id",
				Required:            true,
			},
			"database_id": schema.Int32Attribute{
				MarkdownDescription: "Database Id",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"platform_id": schema.Int32Attribute{
				MarkdownDescription: "Platform ID",
				Required:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance Name",
				Required:            true,
			},
			"is_default_instance": schema.BoolAttribute{
				MarkdownDescription: "Indicates if this is the default instance",
				Optional:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: "Port number",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Instance version",
				Optional:            true,
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "Template name",
				Optional:            true,
			},
		},
	}
}

func (r *databaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating database obj
	databaseObj, err := databases.NewDatabaseObj(*r.providerInfo.authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating database object", err.Error())
		return
	}

	databaseDetails := entities.DatabaseDetails{
		PlatformID:        int(data.PlatformID.ValueInt32()),
		InstanceName:      data.InstanceName.ValueString(),
		IsDefaultInstance: data.IsDefaultInstance.ValueBool(),
		Port:              int(data.PlatformID.ValueInt32()),
		Version:           data.Version.ValueString(),
		Template:          data.Template.ValueString(),
	}

	// creating a database.
	createdDataBase, err := databaseObj.CreateDatabaseFlow(data.AssetId.ValueString(), databaseDetails)

	if err != nil {
		resp.Diagnostics.AddError("Error creating database", err.Error())
		return
	}

	data.DatabaseID = types.Int32Value(int32(createdDataBase.DatabaseID))

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// method not implemented
}

func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
