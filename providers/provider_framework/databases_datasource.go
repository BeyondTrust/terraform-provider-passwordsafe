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
		Description: "Database Datasource, get databases list.",
		Blocks: map[string]schema.Block{
			"databases": schema.ListNestedBlock{
				Description: "Database Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: d.getDatabaseDataSourceSchemaAttributes(),
				},
			},
		},
	}
}

// getDatabaseDataSourceSchemaAttributes get schema attributes.
func (d *DatabaseDataSource) getDatabaseDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"asset_id":            utils.GetInt32Attribute("Asset ID", false, false, true),
		"database_id":         utils.GetInt32Attribute("Database ID", false, false, true),
		"platform_id":         utils.GetInt32Attribute("Platform ID", false, false, true),
		"instance_name":       utils.GetStringAttribute("Instance Name", false, false, true),
		"is_default_instance": utils.GetBoolAttribute("Is Default Instance", false, false, true),
		"port":                utils.GetInt32Attribute("Port", false, false, true),
		"version":             utils.GetStringAttribute("Version", false, false, true),
		"template":            utils.GetStringAttribute("Template", false, false, true),
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
