// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/workgroups"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &WorkGroupResource{}
var _ resource.ResourceWithImportState = &WorkGroupResource{}

func NewWorkGroupResource() resource.Resource {
	return &WorkGroupResource{}
}

type WorkGroupResource struct {
	providerInfo *ProviderData
}

type WorkGroupResorceModel struct {
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Id             types.Int32  `tfsdk:"id"`
}

func (r *WorkGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workgroup"

}

func (r *WorkGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workgroup Resource, creates workgroup. **Note:** Terraform destroy will only remove the workgroup from Terraform state. The actual workgroup must be manually deleted through the BeyondTrust Password Safe web console.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Workgroup Name",
				Required:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				Optional:            true,
				Computed:            false,
			},
			"id": schema.Int32Attribute{
				MarkdownDescription: "Workgroup Id",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
}

func (r *WorkGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

func (r *WorkGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data WorkGroupResorceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Authenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating workgroup obj
	workGroupObj, err := workgroups.NewWorkGroupObj(*r.providerInfo.authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating authentication object", err.Error())
		return
	}

	workGroupDetails := entities.WorkGroupDetails{
		Name:           data.Name.ValueString(),
		OrganizationID: data.OrganizationID.ValueString(),
	}

	// creating a workgroup.
	createdWorkGroup, err := workGroupObj.CreateWorkGroupFlow(workGroupDetails)

	if err != nil {
		resp.Diagnostics.AddError("Error creating workgroup", err.Error())
		return
	}

	if createdWorkGroup.OrganizationID == "" {
		data.OrganizationID = types.StringValue(createdWorkGroup.OrganizationID)
	}

	data.Id = types.Int32Value(int32(createdWorkGroup.ID))

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *WorkGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *WorkGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddAttributeWarning(
		path.Root("id"),
		"Workgroup Deletion Info",
		"Terraform resource deleted. For total deletion of the workgroup, please delete it in the web console as well.",
	)
}

func (r *WorkGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
