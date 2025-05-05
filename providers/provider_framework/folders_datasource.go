// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &FolderDataSource{}

func NewFolderDataSource() datasource.DataSource {
	return &FolderDataSource{}
}

type FolderDataSource struct {
	providerInfo *ProviderData
}

func (d *FolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder_datasource"

}

type FolderModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ParentId    types.String `tfsdk:"parent_id"`
	UserGroupId types.Int32  `tfsdk:"user_group_id"`
}

type FolderDataSourceModel struct {
	Folders []FolderModel `tfsdk:"folders"`
}

func (d *FolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Folder Datasource, gets folders list.",
		Blocks: map[string]schema.Block{
			"folders": schema.ListNestedBlock{
				Description: "Folder Datasource Attibutes",
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
						"parent_id": schema.StringAttribute{
							MarkdownDescription: "Parent ID (GUID)",
							Required:            true,
						},
						"user_group_id": schema.Int32Attribute{
							MarkdownDescription: "User Group ID",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *FolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *FolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FolderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating secrets obj (contains folder methods)
	secretObj, _ := secrets.NewSecretObj(*d.providerInfo.authenticationObj, zapLogger, maxFileSecretSizeBytes)

	// get folders list
	items, err := secretObj.SecretGetFoldersListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting folders list", err.Error())
		return
	}

	var foldersList []FolderModel

	for _, item := range items {
		foldersList = append(foldersList, FolderModel{
			Id:          types.StringValue(item.Id),
			Name:        types.StringValue(item.Name),
			Description: types.StringValue(item.Description),
		})
	}

	responseData := FolderDataSourceModel{}
	responseData.Folders = foldersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
