// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/workgroups"
)

func NewWorkgroupDataSource() datasource.DataSource {
	return &WorkgroupDataSource{}
}

type WorkgroupDataSource struct {
	providerInfo *ProviderData
}

func (d *WorkgroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workgroup_datasource"

}

type WorkgroupModel struct {
	OrganizationID types.String `tfsdk:"organization_id"`
	ID             types.Int32  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
}

type WorkgroupDataSourceModel struct {
	Workgroups []WorkgroupModel `tfsdk:"workgroups"`
}

func (d *WorkgroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Workgroup Datasource",
		Blocks: map[string]schema.Block{
			"workgroups": schema.ListNestedBlock{
				Description: "Workgroup Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"organization_id": schema.StringAttribute{
							MarkdownDescription: "Organization ID",
							Required:            true,
						},
						"id": schema.Int32Attribute{
							MarkdownDescription: "ID",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *WorkgroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *WorkgroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkgroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating workgroup obj
	workgroupObj, _ := workgroups.NewWorkGroupObj(*d.providerInfo.authenticationObj, zapLogger)

	// get workgroups list
	items, err := workgroupObj.GetWorkgroupListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting workgroups list", err.Error())
		return
	}

	var workgroupList []WorkgroupModel

	for _, item := range items {
		workgroupList = append(workgroupList, WorkgroupModel{
			ID:             types.Int32Value(int32(item.ID)),
			OrganizationID: types.StringValue(item.OrganizationID),
			Name:           types.StringValue(item.Name),
		})
	}

	responseData := WorkgroupDataSourceModel{}
	responseData.Workgroups = workgroupList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
