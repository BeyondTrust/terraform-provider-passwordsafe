// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ManagedAccountDataSource{}

func NewManagedAccountDataSource() datasource.DataSource {
	return &ManagedAccountDataSource{}
}

type ManagedAccountDataSource struct {
	providerInfo *ProviderData
}

func (d *ManagedAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_account_datasource"

}

type ManagedAccountModel struct {
	PlatformID             types.Int32  `tfsdk:"platform_id"`
	SystemID               types.Int32  `tfsdk:"system_id"`
	SystemName             types.String `tfsdk:"system_name"`
	DomainName             types.String `tfsdk:"domain_name"`
	AccountID              types.Int32  `tfsdk:"account_id"`
	AccountName            types.String `tfsdk:"account_name"`
	InstanceName           types.String `tfsdk:"instance_name"`
	UserPrincipalName      types.String `tfsdk:"user_principal_name"`
	ApplicationID          types.Int32  `tfsdk:"application_id"`
	ApplicationDisplayName types.String `tfsdk:"application_display_name"`
	DefaultReleaseDuration types.Int32  `tfsdk:"default_release_duration"`
	MaximumReleaseDuration types.Int32  `tfsdk:"maximum_release_duration"`
	LastChangeDate         types.String `tfsdk:"last_change_date"`
	NextChangeDate         types.String `tfsdk:"next_change_date"`
	IsChanging             types.Bool   `tfsdk:"is_changing"`
	ChangeState            types.Int32  `tfsdk:"change_state"`
	IsISAAccess            types.Bool   `tfsdk:"is_isa_access"`
	PreferredNodeID        types.String `tfsdk:"preferred_node_id"`
	AccountDescription     types.String `tfsdk:"account_description"`
}

type ManagedAccountDataSourceModel struct {
	ManagedAccounts []ManagedAccountModel `tfsdk:"managed_accounts"`
}

func (d *ManagedAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Managed Account Datasource, get managed accounts list.",
		Blocks: map[string]schema.Block{
			"managed_accounts": schema.ListNestedBlock{
				Description: "Managed Account Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"platform_id": schema.Int32Attribute{
							MarkdownDescription: "Platform ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"system_id": schema.Int32Attribute{
							MarkdownDescription: "System ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"system_name": schema.StringAttribute{
							MarkdownDescription: "System Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							MarkdownDescription: "Domain Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"account_id": schema.Int32Attribute{
							MarkdownDescription: "Account ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"account_name": schema.StringAttribute{
							MarkdownDescription: "Account Name",
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
						"user_principal_name": schema.StringAttribute{
							MarkdownDescription: "User Principal Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"application_id": schema.Int32Attribute{
							MarkdownDescription: "Application ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"application_display_name": schema.StringAttribute{
							MarkdownDescription: "Application Display Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"default_release_duration": schema.Int32Attribute{
							MarkdownDescription: "Default Release Duration",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"maximum_release_duration": schema.Int32Attribute{
							MarkdownDescription: "Maximum Release Duration",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"last_change_date": schema.StringAttribute{
							MarkdownDescription: "Last Change Date (ISO 8601 format)",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"next_change_date": schema.StringAttribute{
							MarkdownDescription: "Next Change Date (ISO 8601 format)",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"is_changing": schema.BoolAttribute{
							MarkdownDescription: "Is Changing",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"change_state": schema.Int32Attribute{
							MarkdownDescription: "Change State",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"is_isa_access": schema.BoolAttribute{
							MarkdownDescription: "ISA Access",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"preferred_node_id": schema.StringAttribute{
							MarkdownDescription: "Preferred Node ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"account_description": schema.StringAttribute{
							MarkdownDescription: "Account Description",
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

func (d *ManagedAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *ManagedAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManagedAccountDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating managed acocunt obj.
	managedAccountObj, _ := managed_accounts.NewManagedAccountObj(*d.providerInfo.authenticationObj, zapLogger)

	// get managed accounts list.
	items, err := managedAccountObj.GetManagedAccountsListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting managed accounts list", err.Error())
		return
	}

	var managedAccountList []ManagedAccountModel

	for _, item := range items {
		managedAccountList = append(managedAccountList, ManagedAccountModel{
			PlatformID:             types.Int32Value(int32(item.PlatformID)),
			SystemID:               types.Int32Value(int32(item.SystemId)),
			SystemName:             types.StringValue(item.SystemName),
			DomainName:             types.StringValue(item.DomainName),
			AccountID:              types.Int32Value(int32(item.AccountId)),
			AccountName:            types.StringValue(item.AccountName),
			InstanceName:           types.StringValue(item.InstanceName),
			UserPrincipalName:      types.StringValue(item.UserPrincipalName),
			ApplicationID:          types.Int32Value(int32(item.ApplicationID)),
			ApplicationDisplayName: types.StringValue(item.ApplicationDisplayName),
			DefaultReleaseDuration: types.Int32Value(int32(item.DefaultReleaseDuration)),
			MaximumReleaseDuration: types.Int32Value(int32(item.MaximumReleaseDuration)),
			LastChangeDate:         types.StringValue(item.LastChangeDate),
			NextChangeDate:         types.StringValue(item.NextChangeDate),
			IsChanging:             types.BoolValue(item.IsChanging),
			ChangeState:            types.Int32Value(int32(item.ChangeState)),
			IsISAAccess:            types.BoolValue(item.IsISAAccess),
			PreferredNodeID:        types.StringValue(item.PreferredNodeID),
			AccountDescription:     types.StringValue(item.AccountDescription),
		})
	}

	responseData := ManagedAccountDataSourceModel{}
	responseData.ManagedAccounts = managedAccountList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
