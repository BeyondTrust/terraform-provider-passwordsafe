// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
	"github.com/BeyondTrust/go-client-library-passwordsafe/api/functional_accounts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FunctionalAccountResource{}
var _ resource.ResourceWithImportState = &FunctionalAccountResource{}

func NewFunctionalAccountResource() resource.Resource {
	return &FunctionalAccountResource{}
}

type FunctionalAccountResource struct {
	providerInfo *ProviderData
}

type FunctionalResourceResourceModel struct {
	FunctionalAccountID types.Int32  `tfsdk:"functional_account_id"`
	PlatformID          types.Int32  `tfsdk:"platform_id"`
	DomainName          types.String `tfsdk:"domain_name"`
	AccountName         types.String `tfsdk:"account_name"`
	DisplayName         types.String `tfsdk:"display_name"`
	Password            types.String `tfsdk:"password"`
	PrivateKey          types.String `tfsdk:"private_key"`
	Passphrase          types.String `tfsdk:"passphrase"`
	Description         types.String `tfsdk:"description"`
	ElevationCommand    types.String `tfsdk:"elevation_command"`
	TenantID            types.String `tfsdk:"tenant_id"`
	ObjectID            types.String `tfsdk:"object_id"`
	Secret              types.String `tfsdk:"secret"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	AzureInstance       types.String `tfsdk:"azure_instance"`
}

func (r *FunctionalAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_functional_account"

}

func (r *FunctionalAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Functional Account Resource, creates functional account",

		Attributes: map[string]schema.Attribute{
			"functional_account_id": schema.Int32Attribute{
				MarkdownDescription: "Functional Account ID",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"platform_id": schema.Int32Attribute{
				MarkdownDescription: "Platform ID",
				Required:            true,
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Domain Name",
				Optional:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "Account Name",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display Name",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password",
				Optional:            true,
				Sensitive:           true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Private Key",
				Optional:            true,
				Sensitive:           true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase",
				Optional:            true,
				Sensitive:           true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Optional:            true,
			},
			"elevation_command": schema.StringAttribute{
				MarkdownDescription: "Elevation Command",
				Optional:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "Tenant ID",
				Optional:            true,
			},
			"object_id": schema.StringAttribute{
				MarkdownDescription: "Object ID",
				Optional:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Secret",
				Optional:            true,
				Sensitive:           true,
			},
			"service_account_email": schema.StringAttribute{
				MarkdownDescription: "Service Account Email",
				Optional:            true,
			},
			"azure_instance": schema.StringAttribute{
				MarkdownDescription: "Azure Instance (AzurePublic or AzureUsGovernment)",
				Optional:            true,
			},
		},
	}
}

func (r *FunctionalAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	r.providerInfo = &c

	if r.providerInfo.userName == "" {
		return
	}

}

func (r *FunctionalAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data FunctionalResourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Authenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating functional account obj
	functionalAccountObj, err := functional_accounts.NewFuncionalAccount(*r.providerInfo.authenticationObj, zapLogger)

	if err != nil {
		resp.Diagnostics.AddError("Error creating functional account object", err.Error())
		return
	}

	functionalAccountDetails := entities.FunctionalAccountDetails{
		PlatformID:          int(data.PlatformID.ValueInt32()),
		DomainName:          data.DomainName.ValueString(),
		AccountName:         data.AccountName.ValueString(),
		DisplayName:         data.DisplayName.ValueString(),
		Password:            data.Password.ValueString(),
		PrivateKey:          data.PrivateKey.ValueString(),
		Passphrase:          data.Passphrase.ValueString(),
		Description:         data.Description.ValueString(),
		ElevationCommand:    data.ElevationCommand.ValueString(),
		TenantID:            data.TenantID.ValueString(),
		ObjectID:            data.ObjectID.ValueString(),
		Secret:              data.Secret.ValueString(),
		ServiceAccountEmail: data.ServiceAccountEmail.ValueString(),
		AzureInstance:       data.AzureInstance.ValueString(),
	}

	// creating a functional account.
	createdFunctionalAccount, err := functionalAccountObj.CreateFunctionalAccountFlow(functionalAccountDetails)

	if err != nil {
		resp.Diagnostics.AddError("Error creating functional account object", err.Error())
		return
	}

	data.FunctionalAccountID = types.Int32Value(int32(createdFunctionalAccount.FunctionalAccountID))

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *FunctionalAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// method not implemented
}

func (r *FunctionalAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// method not implemented
}

func (r *FunctionalAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FunctionalResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Authenticate(*r.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating functional account obj
	functionalAccountObj, err := functional_accounts.NewFuncionalAccount(*r.providerInfo.authenticationObj, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error creating functional account object", err.Error())
		return
	}

	// deleting the functional account by ID
	err = functionalAccountObj.DeleteFunctionalAccountById(int(data.FunctionalAccountID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting functional account", err.Error())
		return
	}

	err = utils.SignOut(*r.providerInfo.authenticationObj, &muOut, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error Signing Out", err.Error())
		return
	}
}

func (r *FunctionalAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
