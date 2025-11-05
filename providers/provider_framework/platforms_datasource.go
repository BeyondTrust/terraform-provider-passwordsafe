// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/platforms"
)

var _ datasource.DataSource = &PlatformDataSource{}

func NewPlatformDataSource() datasource.DataSource {
	return &PlatformDataSource{}
}

type PlatformDataSource struct {
	providerInfo *ProviderData
}

func (d *PlatformDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform_datasource"

}

type PlatformModel struct {
	PlatformID              types.Int32  `tfsdk:"platform_id"`
	Name                    types.String `tfsdk:"name"`
	ShortName               types.String `tfsdk:"short_name"`
	PortFlag                types.Bool   `tfsdk:"port_flag"`
	DefaultPort             types.Int32  `tfsdk:"default_port"`
	SupportsElevationFlag   types.Bool   `tfsdk:"supports_elevation_flag"`
	DomainNameFlag          types.Bool   `tfsdk:"domain_name_flag"`
	AutoManagementFlag      types.Bool   `tfsdk:"auto_management_flag"`
	DSSAutoManagementFlag   types.Bool   `tfsdk:"dss_auto_management_flag"`
	ManageableFlag          types.Bool   `tfsdk:"manageable_flag"`
	DSSFlag                 types.Bool   `tfsdk:"dss_flag"`
	LoginAccountFlag        types.Bool   `tfsdk:"login_account_flag"`
	DefaultSessionType      types.String `tfsdk:"default_session_type"`
	ApplicationHostFlag     types.Bool   `tfsdk:"application_host_flag"`
	RequiresApplicationHost types.Bool   `tfsdk:"requires_application_host"`
	RequiresTenantID        types.Bool   `tfsdk:"requires_tenant_id"`
	RequiresObjectID        types.Bool   `tfsdk:"requires_object_id"`
	RequiresSecret          types.Bool   `tfsdk:"requires_secret"`
}

type PlatformDataSourceModel struct {
	Platforms []PlatformModel `tfsdk:"platforms"`
}

func (d *PlatformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Platform Datasource, gets platforms list",
		Blocks: map[string]schema.Block{
			"platforms": schema.ListNestedBlock{
				Description: "Platform Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"platform_id": schema.Int32Attribute{
							MarkdownDescription: "Platform ID",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
						"short_name": schema.StringAttribute{
							MarkdownDescription: "Short Name",
							Required:            true,
						},
						"port_flag": schema.BoolAttribute{
							MarkdownDescription: "Port Flag",
							Required:            true,
						},
						"default_port": schema.Int32Attribute{
							MarkdownDescription: "Default Port (nullable)",
							Optional:            true,
							Computed:            true,
						},
						"supports_elevation_flag": schema.BoolAttribute{
							MarkdownDescription: "Supports Elevation Flag",
							Required:            true,
						},
						"domain_name_flag": schema.BoolAttribute{
							MarkdownDescription: "Domain Name Flag",
							Required:            true,
						},
						"auto_management_flag": schema.BoolAttribute{
							MarkdownDescription: "Auto Management Flag",
							Required:            true,
						},
						"dss_auto_management_flag": schema.BoolAttribute{
							MarkdownDescription: "DSS Auto Management Flag",
							Required:            true,
						},
						"manageable_flag": schema.BoolAttribute{
							MarkdownDescription: "Manageable Flag",
							Required:            true,
						},
						"dss_flag": schema.BoolAttribute{
							MarkdownDescription: "DSS Flag",
							Required:            true,
						},
						"login_account_flag": schema.BoolAttribute{
							MarkdownDescription: "Login Account Flag",
							Required:            true,
						},
						"default_session_type": schema.StringAttribute{
							MarkdownDescription: "Default Session Type (nullable)",
							Optional:            true,
							Computed:            true,
						},
						"application_host_flag": schema.BoolAttribute{
							MarkdownDescription: "Application Host Flag",
							Required:            true,
						},
						"requires_application_host": schema.BoolAttribute{
							MarkdownDescription: "Requires Application Host",
							Required:            true,
						},
						"requires_tenant_id": schema.BoolAttribute{
							MarkdownDescription: "Requires Tenant ID",
							Required:            true,
						},
						"requires_object_id": schema.BoolAttribute{
							MarkdownDescription: "Requires Object ID",
							Required:            true,
						},
						"requires_secret": schema.BoolAttribute{
							MarkdownDescription: "Requires Secret",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *PlatformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *PlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlatformDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Authenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating platform obj
	platformObj, _ := platforms.NewPlatformObj(*d.providerInfo.authenticationObj, zapLogger)

	// get platforms list
	items, err := platformObj.GetPlatformsListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting platforms list", err.Error())
		return
	}

	var platformList []PlatformModel

	for _, item := range items {
		platformList = append(platformList, PlatformModel{
			PlatformID:              types.Int32Value(int32(item.PlatformID)),
			Name:                    types.StringValue(item.Name),
			ShortName:               types.StringValue(item.ShortName),
			PortFlag:                types.BoolValue(item.PortFlag),
			DefaultPort:             types.Int32Value(int32(item.DefaultPort)),
			SupportsElevationFlag:   types.BoolValue(item.SupportsElevationFlag),
			DomainNameFlag:          types.BoolValue(item.DomainNameFlag),
			AutoManagementFlag:      types.BoolValue(item.AutoManagementFlag),
			DSSAutoManagementFlag:   types.BoolValue(item.DSSAutoManagementFlag),
			ManageableFlag:          types.BoolValue(item.ManageableFlag),
			DSSFlag:                 types.BoolValue(item.DSSFlag),
			LoginAccountFlag:        types.BoolValue(item.LoginAccountFlag),
			DefaultSessionType:      types.StringValue(item.DefaultSessionType),
			ApplicationHostFlag:     types.BoolValue(item.ApplicationHostFlag),
			RequiresApplicationHost: types.BoolValue(item.RequiresApplicationHost),
			RequiresTenantID:        types.BoolValue(item.RequiresTenantID),
			RequiresObjectID:        types.BoolValue(item.RequiresObjectID),
			RequiresSecret:          types.BoolValue(item.RequiresSecret),
		})
	}

	responseData := PlatformDataSourceModel{}
	responseData.Platforms = platformList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
