// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
)

var _ datasource.DataSource = &SafesDataSource{}

func NewSafesDataSource() datasource.DataSource {
	return &SafesDataSource{}
}

type SafesDataSource struct {
	providerInfo *ProviderData
}

func (d *SafesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_safe_datasource"

}

type SafeModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type SafesDataSourceModel struct {
	Safes []SafeModel `tfsdk:"safes"`
}

func (d *SafesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Safe Datasource, gets safes list.",
		Blocks: map[string]schema.Block{
			"safes": schema.ListNestedBlock{
				Description: "Safe Datasource, gets safes list.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ID (GUID)",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SafesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *SafesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SafesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating secrets obj (contains safes methods)
	secretObj, _ := secrets.NewSecretObj(*d.providerInfo.authenticationObj, zapLogger, maxFileSecretSizeBytes)

	// get safes list
	items, err := secretObj.SecretGetSafesListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting safes list", err.Error())
		return
	}

	var safesList []SafeModel

	for _, item := range items {
		safesList = append(safesList, SafeModel{
			Id:          types.StringValue(item.Id),
			Name:        types.StringValue(item.Name),
			Description: types.StringValue(item.Description),
		})
	}

	responseData := SafesDataSourceModel{}
	responseData.Safes = safesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
