// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/assets"
)

var _ datasource.DataSource = &AssetDataSource{}

func NewAssetDataSource() datasource.DataSource {
	return &AssetDataSource{}
}

type AssetDataSource struct {
	providerInfo *ProviderData
}

type AssetDataSourceModel struct {
	Assets    []AssetModel `tfsdk:"assets"`
	Parameter types.String `tfsdk:"parameter"`
}

func (d *AssetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_datasource"

}

type AssetModel struct {
	WorkgroupID     types.Int32  `tfsdk:"workgroup_id"`
	AssetID         types.Int32  `tfsdk:"asset_id"`
	AssetName       types.String `tfsdk:"asset_name"`
	DnsName         types.String `tfsdk:"dns_name"`
	DomainName      types.String `tfsdk:"domain_name"`
	IPAddress       types.String `tfsdk:"ip_address"`
	MacAddress      types.String `tfsdk:"mac_address"`
	AssetType       types.String `tfsdk:"asset_type"`
	OperatingSystem types.String `tfsdk:"operating_system"`
	CreateDate      types.String `tfsdk:"create_date"`
	LastUpdateDate  types.String `tfsdk:"last_update_date"`
	Description     types.String `tfsdk:"description"`
}

func (d *AssetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Asset Datasource, get assets list.",
		Attributes: map[string]schema.Attribute{
			"parameter": schema.StringAttribute{
				MarkdownDescription: "Parameter",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"assets": schema.ListNestedBlock{
				Description: "Asset Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: d.getAssetDataSourceSchemaAttributes(),
				},
			},
		},
	}
}

// getAssetDataSourceSchemaAttributes get schema attributes.
func (d *AssetDataSource) getAssetDataSourceSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"workgroup_id":     utils.GetInt32Attribute("Workgroup ID", false, false, true),
		"asset_id":         utils.GetInt32Attribute("Asset ID", false, false, true),
		"asset_name":       utils.GetStringAttribute("Asset Name", false, false, true),
		"dns_name":         utils.GetStringAttribute("DNS Name", false, false, true),
		"domain_name":      utils.GetStringAttribute("Domain Name", false, false, true),
		"ip_address":       utils.GetStringAttribute("IP Address", false, false, true),
		"mac_address":      utils.GetStringAttribute("MAC Address", false, false, true),
		"asset_type":       utils.GetStringAttribute("Asset Type", false, false, true),
		"operating_system": utils.GetStringAttribute("Operating System", false, false, true),
		"create_date":      utils.GetStringAttribute("Creation Date (ISO 8601 format)", false, false, true),
		"last_update_date": utils.GetStringAttribute("Last Update Date (ISO 8601 format)", false, false, true),
		"description":      utils.GetStringAttribute("Description", false, false, true),
	}
}

func (d *AssetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *AssetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var inputData AssetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &inputData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating asset obj
	asssetObj, _ := assets.NewAssetObj(*d.providerInfo.authenticationObj, zapLogger)

	// get assets list using workgroup id.
	itemsByWorkgroupId, err := asssetObj.GetAssetsListByWorkgroupIdFlow(inputData.Parameter.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error getting assets list by workgroup id", err.Error())
		return
	}

	items := itemsByWorkgroupId

	// if there is not assets by workgroup id so it will search by workgroup name.
	if len(items) == 0 {
		// get assets list using workgroup name.
		itemsByWorkgroupName, err := asssetObj.GetAssetsListByWorkgroupNameFlow(inputData.Parameter.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error getting assets list by workgroup name", err.Error())
			return
		}
		items = itemsByWorkgroupName
	}

	var assetsList []AssetModel

	for _, item := range items {
		assetsList = append(assetsList, AssetModel{
			WorkgroupID:     types.Int32Value(int32(item.WorkgroupID)),
			AssetID:         types.Int32Value(int32(item.AssetID)),
			AssetName:       types.StringValue(item.AssetName),
			AssetType:       types.StringValue(item.AssetType),
			DnsName:         types.StringValue(item.DnsName),
			DomainName:      types.StringValue(item.DomainName),
			IPAddress:       types.StringValue(item.IPAddress),
			OperatingSystem: types.StringValue(item.OperatingSystem),
			CreateDate:      types.StringValue(item.CreateDate),
			LastUpdateDate:  types.StringValue(item.LastUpdateDate),
			Description:     types.StringValue(item.Description),
		})
	}

	responseData := AssetDataSourceModel{}
	responseData.Assets = assetsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
