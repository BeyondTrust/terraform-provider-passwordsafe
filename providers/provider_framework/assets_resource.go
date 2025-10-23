// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/assets"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
)

type assetResource struct {
	providerInfo  *ProviderData
	resourceName  string
	resurceSchema schema.Schema
}

type AssetResorceModel struct {
	AssetID         types.Int32  `tfsdk:"asset_id"`
	IPAddress       types.String `tfsdk:"ip_address"`
	AssetName       types.String `tfsdk:"asset_name"`
	DnsName         types.String `tfsdk:"dns_name"`
	DomainName      types.String `tfsdk:"domain_name"`
	AssetType       types.String `tfsdk:"asset_type"`
	Description     types.String `tfsdk:"description"`
	OperatingSystem types.String `tfsdk:"operating_system"`
}

func (r *assetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.resourceName
}

func (r *assetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.resurceSchema
}

func (r *assetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

func (r *assetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *assetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *assetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssetResorceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Authenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating asset obj
	assetObj, err := assets.NewAssetObj(*r.providerInfo.authenticationObj, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error creating asset object", err.Error())
		return
	}

	// deleting the asset by ID
	err = assetObj.DeleteAssetById(int(data.AssetID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting asset", err.Error())
		return
	}

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}
}

// Shared delete helper function that works with any asset model containing AssetID
func (r *assetResource) deleteAssetByID(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, assetID int32) {
	_, err := utils.Authenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating asset obj
	assetObj, err := assets.NewAssetObj(*r.providerInfo.authenticationObj, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error creating asset object", err.Error())
		return
	}

	// deleting the asset by ID
	err = assetObj.DeleteAssetById(int(assetID))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting asset", err.Error())
		return
	}

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}
}

func (r *assetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// NewAssetByWorkgGroypIdResource

var _ resource.Resource = &assetResourceByWorkGroupId{}
var _ resource.ResourceWithImportState = &assetResourceByWorkGroupId{}

type AssetResorceByWorkGroupIdModel struct {
	AssetResorceModel
	WorkGroupId types.String `tfsdk:"work_group_id"`
}

type assetResourceByWorkGroupId struct {
	assetResource
}

func NewAssetByWorkgGroypIdResource() resource.Resource {
	assetResource := &assetResourceByWorkGroupId{}

	assetResource.resourceName = "_asset_by_workgroup_id"
	assetResource.resurceSchema = schema.Schema{
		MarkdownDescription: "Asset resource, creates asset by workgroup id.",
		Attributes: map[string]schema.Attribute{
			"work_group_id": schema.StringAttribute{
				MarkdownDescription: "Workgroup Id",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP Address",
				Required:            true,
			},
			"asset_id": schema.Int32Attribute{
				MarkdownDescription: "Asset Id",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"asset_name": schema.StringAttribute{
				MarkdownDescription: "Asset Name",
				Optional:            true,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "DNS Name",
				Optional:            true,
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Domain Name",
				Optional:            true,
			},
			"asset_type": schema.StringAttribute{
				MarkdownDescription: "Asset Type",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Optional:            true,
			},
			"operating_system": schema.StringAttribute{
				MarkdownDescription: "Operating System",
				Optional:            true,
			},
		},
	}

	return assetResource
}

func getAssetObj(resp *resource.CreateResponse, authenticationObj authentication.AuthenticationObj, dataInterface interface{}) *assets.AssetObj {

	_, err := utils.Authenticate(authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return nil
	}

	assetGroupObj, err := assets.NewAssetObj(authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating authentication object", err.Error())
		return nil
	}

	return assetGroupObj
}

func (r *assetResourceByWorkGroupId) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data AssetResorceByWorkGroupIdModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	assetGroupObj := getAssetObj(resp, *r.providerInfo.authenticationObj, data)

	assetDetails := entities.AssetDetails{
		IPAddress:       data.IPAddress.ValueString(),
		AssetName:       data.AssetName.ValueString(),
		DnsName:         data.DnsName.ValueString(),
		DomainName:      data.DomainName.ValueString(),
		AssetType:       data.AssetType.ValueString(),
		Description:     data.Description.ValueString(),
		OperatingSystem: data.OperatingSystem.ValueString(),
	}

	createdAsset, err := assetGroupObj.CreateAssetByworkgroupIDFlow(data.WorkGroupId.ValueString(), assetDetails)
	if err != nil {
		resp.Diagnostics.AddError("Error creating asset by WorkGroup Id", err.Error())
		return
	}

	data.AssetID = types.Int32Value(int32(createdAsset.AssetID))

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *assetResourceByWorkGroupId) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssetResorceByWorkGroupIdModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.deleteAssetByID(ctx, req, resp, data.AssetID.ValueInt32())
}

// AssetByWorkGroupNameResource

var _ resource.Resource = &assetResourceByWorkGroupName{}
var _ resource.ResourceWithImportState = &assetResourceByWorkGroupName{}

type AssetResorceByWorkGroupNameModel struct {
	AssetResorceModel
	WorkGroupName types.String `tfsdk:"work_group_name"`
}

type assetResourceByWorkGroupName struct {
	assetResource
}

func NewAssetByWorkGroupNameResource() resource.Resource {
	assetResource := &assetResourceByWorkGroupName{}

	assetResource.resourceName = "_asset_by_workgroup_name"
	assetResource.resurceSchema = schema.Schema{
		MarkdownDescription: "Asset resource, creates asset by workgroup name.",
		Attributes: map[string]schema.Attribute{
			"work_group_name": schema.StringAttribute{
				MarkdownDescription: "Workgroup Name",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP Address",
				Required:            true,
			},
			"asset_id": schema.Int32Attribute{
				MarkdownDescription: "Asset Id",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"asset_name": schema.StringAttribute{
				MarkdownDescription: "Asset Name",
				Optional:            true,
				Computed:            false,
			},
			"dns_name": schema.StringAttribute{
				MarkdownDescription: "Dns Name",
				Optional:            true,
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Domain Name",
				Optional:            true,
			},
			"asset_type": schema.StringAttribute{
				MarkdownDescription: "Asset Type",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Optional:            true,
			},
			"operating_system": schema.StringAttribute{
				MarkdownDescription: "Operating System",
				Optional:            true,
			},
		},
	}

	return assetResource
}

func (r *assetResourceByWorkGroupName) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data AssetResorceByWorkGroupNameModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	assetGroupObj := getAssetObj(resp, *r.providerInfo.authenticationObj, data)

	assetDetails := entities.AssetDetails{
		IPAddress:       data.IPAddress.ValueString(),
		AssetName:       data.AssetName.ValueString(),
		DnsName:         data.DnsName.ValueString(),
		DomainName:      data.DomainName.ValueString(),
		AssetType:       data.AssetType.ValueString(),
		Description:     data.Description.ValueString(),
		OperatingSystem: data.OperatingSystem.ValueString(),
	}

	createdAsset, err := assetGroupObj.CreateAssetByWorkGroupNameFlow(data.WorkGroupName.ValueString(), assetDetails)
	if err != nil {
		resp.Diagnostics.AddError("Error creating asset by WorkGroup Name", err.Error())
		return
	}

	data.AssetID = types.Int32Value(int32(createdAsset.AssetID))

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *assetResourceByWorkGroupName) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssetResorceByWorkGroupNameModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.deleteAssetByID(ctx, req, resp, data.AssetID.ValueInt32())
}
