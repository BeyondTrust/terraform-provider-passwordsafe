// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/functional_accounts"
)

var _ datasource.DataSource = &FunctionalAccountDataResource{}

func NewFunctionalAccountDataResource() datasource.DataSource {
	return &FunctionalAccountDataResource{}
}

type FunctionalAccountDataResource struct {
	providerInfo *ProviderData
}

func (d *FunctionalAccountDataResource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_functional_account_datasource"

}

type FunctionalAccountModel struct {
	FunctionalAccountID types.Int32  `tfsdk:"functional_account_id"`
	PlatformID          types.Int32  `tfsdk:"platform_id"`
	DomainName          types.String `tfsdk:"domain_name"`
	AccountName         types.String `tfsdk:"account_name"`
}

type FunctionalDataSourceModel struct {
	Accounts []FunctionalAccountModel `tfsdk:"accounts"`
}

func (d *FunctionalAccountDataResource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Functional Account Datasource, get functional accounts list.",
		Blocks: map[string]schema.Block{
			"accounts": schema.ListNestedBlock{
				Description: "Functional Account Datasource Attibutes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"functional_account_id": schema.Int32Attribute{
							MarkdownDescription: "Functional Account ID",
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
						"domain_name": schema.StringAttribute{
							MarkdownDescription: "Domain Name",
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
					},
				},
			},
		},
	}
}

func (d *FunctionalAccountDataResource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *FunctionalAccountDataResource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FunctionalDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating functional account obj
	functionalAccountObj, _ := functional_accounts.NewFuncionalAccount(*d.providerInfo.authenticationObj, zapLogger)

	// get functional accounts list
	items, err := functionalAccountObj.GetFunctionalAccountsFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting functional accounts list", err.Error())
		return
	}

	var accountsList []FunctionalAccountModel

	for _, acc := range items {
		accountsList = append(accountsList, FunctionalAccountModel{
			FunctionalAccountID: types.Int32Value(int32(acc.FunctionalAccountID)),
			PlatformID:          types.Int32Value(int32(acc.PlatformID)),
			DomainName:          types.StringValue(acc.DomainName),
			AccountName:         types.StringValue(acc.AccountName),
		})
	}

	responseData := FunctionalDataSourceModel{}
	responseData.Accounts = accountsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
