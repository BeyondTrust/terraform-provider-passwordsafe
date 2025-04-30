// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/databases"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

type DatabaseDataSource struct {
	providerInfo *ProviderData
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_datasource"

}

type DatabaseModel struct {
	AssetID           types.Int32  `tfsdk:"asset_id"`
	DatabaseID        types.Int32  `tfsdk:"database_id"`
	PlatformID        types.Int32  `tfsdk:"platform_id"`
	InstanceName      types.String `tfsdk:"instance_name"`
	IsDefaultInstance types.Bool   `tfsdk:"is_default_instance"`
	Port              types.Int32  `tfsdk:"port"`
	Version           types.String `tfsdk:"version"`
	Template          types.String `tfsdk:"template"`
}

type DatabaseDataSourceModel struct {
	Databases []DatabaseModel `tfsdk:"databases"`
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Database Datasource",
		Blocks: map[string]schema.Block{
			"databases": schema.ListNestedBlock{
				Description: "Database Datasource Attibutes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"asset_id": schema.Int32Attribute{
							MarkdownDescription: "Asset ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"database_id": schema.Int32Attribute{
							MarkdownDescription: "Database ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"platform_id": schema.Int32Attribute{
							MarkdownDescription: "Platform ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"instance_name": schema.StringAttribute{
							MarkdownDescription: "Instance Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"is_default_instance": schema.BoolAttribute{
							MarkdownDescription: "Is Default Instance",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"port": schema.Int32Attribute{
							MarkdownDescription: "Port",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "Version",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"template": schema.StringAttribute{
							MarkdownDescription: "Template",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating database obj
	databaseObj, _ := databases.NewDatabaseObj(*d.providerInfo.authenticationObj, zapLogger)

	// get databases list
	items, err := databaseObj.GetDatabasesListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting databases list", err.Error())
		return
	}

	var databasesList []DatabaseModel

	for _, items := range items {
		databasesList = append(databasesList, DatabaseModel{
			AssetID:           types.Int32Value(int32(items.AssetID)),
			DatabaseID:        types.Int32Value(int32(items.DatabaseID)),
			PlatformID:        types.Int32Value(int32(items.PlatformID)),
			InstanceName:      types.StringValue(items.InstanceName),
			IsDefaultInstance: types.BoolValue(items.IsDefaultInstance),
			Port:              types.Int32Value(int32(items.Port)),
			Version:           types.StringValue(items.Version),
			Template:          types.StringValue(items.Template),
		})
	}

	responseData := DatabaseDataSourceModel{}
	responseData.Databases = databasesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
