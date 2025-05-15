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
						"platform_id":              utils.GetInt32Attribute("Platform ID", false, false, true),
						"system_id":                utils.GetInt32Attribute("System ID", false, false, true),
						"system_name":              utils.GetStringAttribute("System Name", false, false, true),
						"domain_name":              utils.GetStringAttribute("Domain Name", false, false, true),
						"account_id":               utils.GetInt32Attribute("Account ID", false, false, true),
						"account_name":             utils.GetStringAttribute("Account Name", false, false, true),
						"instance_name":            utils.GetStringAttribute("Instance Name", false, false, true),
						"user_principal_name":      utils.GetStringAttribute("User Principal Name", false, false, true),
						"application_id":           utils.GetInt32Attribute("Application ID", false, false, true),
						"application_display_name": utils.GetStringAttribute("Application Display Name", false, false, true),
						"default_release_duration": utils.GetInt32Attribute("Default Release Duration", false, false, true),
						"maximum_release_duration": utils.GetInt32Attribute("Maximum Release Duration", false, false, true),
						"last_change_date":         utils.GetStringAttribute("Last Change Date (ISO 8601 format)", false, false, true),
						"next_change_date":         utils.GetStringAttribute("Next Change Date (ISO 8601 format)", false, false, true),
						"is_changing":              utils.GetBoolAttribute("Is Changing", false, false, true),
						"change_state":             utils.GetInt32Attribute("Change State", false, false, true),
						"is_isa_access":            utils.GetBoolAttribute("ISA Access", false, false, true),
						"preferred_node_id":        utils.GetStringAttribute("Preferred Node ID", false, false, true),
						"account_description":      utils.GetStringAttribute("Account Description", false, false, true),
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
