// Copyright 2025 BeyondTrust. All rights reserved.
// Package provider_framework implements a terraform provider that can talk with Beyondtrust Secret Safe API.
package provider_framework

import (
	"context"
	"terraform-provider-passwordsafe/providers/utils"

	managed_systems "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_systems"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ManagedSystemDataSource{}

func NewManagedSystemDataSource() datasource.DataSource {
	return &ManagedSystemDataSource{}
}

type ManagedSystemDataSource struct {
	providerInfo *ProviderData
}

func (d *ManagedSystemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_system_datasource"

}

type ManagedSystemModel struct {
	WorkgroupID                        types.Int32  `tfsdk:"workgroup_id"`
	HostName                           types.String `tfsdk:"host_name"`
	IPAddress                          types.String `tfsdk:"ip_address"`
	DNSName                            types.String `tfsdk:"dns_name"`
	InstanceName                       types.String `tfsdk:"instance_name"`
	IsDefaultInstance                  types.Bool   `tfsdk:"is_default_instance"`
	Template                           types.String `tfsdk:"template"`
	ForestName                         types.String `tfsdk:"forest_name"`
	UseSSL                             types.Bool   `tfsdk:"use_ssl"`
	ManagedSystemID                    types.Int32  `tfsdk:"managed_system_id"`
	EntityTypeID                       types.Int32  `tfsdk:"entity_type_id"`
	AssetID                            types.Int32  `tfsdk:"asset_id"`
	DatabaseID                         types.Int32  `tfsdk:"database_id"`
	DirectoryID                        types.Int32  `tfsdk:"directory_id"`
	CloudID                            types.Int32  `tfsdk:"cloud_id"`
	SystemName                         types.String `tfsdk:"system_name"`
	Timeout                            types.Int32  `tfsdk:"timeout"`
	PlatformID                         types.Int32  `tfsdk:"platform_id"`
	NetBiosName                        types.String `tfsdk:"net_bios_name"`
	ContactEmail                       types.String `tfsdk:"contact_email"`
	Description                        types.String `tfsdk:"description"`
	Port                               types.Int32  `tfsdk:"port"`
	SshKeyEnforcementMode              types.Int32  `tfsdk:"ssh_key_enforcement_mode"`
	PasswordRuleID                     types.Int32  `tfsdk:"password_rule_id"`
	DSSKeyRuleID                       types.Int32  `tfsdk:"dss_key_rule_id"`
	LoginAccountID                     types.Int32  `tfsdk:"login_account_id"`
	AccountNameFormat                  types.Int32  `tfsdk:"account_name_format"`
	OracleInternetDirectoryID          types.String `tfsdk:"oracle_internet_directory_id"`
	OracleInternetDirectoryServiceName types.String `tfsdk:"oracle_internet_directory_service_name"`
	ReleaseDuration                    types.Int32  `tfsdk:"release_duration"`
	MaxReleaseDuration                 types.Int32  `tfsdk:"max_release_duration"`
	ISAReleaseDuration                 types.Int32  `tfsdk:"isa_release_duration"`
	AutoManagementFlag                 types.Bool   `tfsdk:"auto_management_flag"`
	FunctionalAccountID                types.Int32  `tfsdk:"functional_account_id"`
	ElevationCommand                   types.String `tfsdk:"elevation_command"`
	CheckPasswordFlag                  types.Bool   `tfsdk:"check_password_flag"`
	ChangePasswordAfterAnyReleaseFlag  types.Bool   `tfsdk:"change_password_after_any_release_flag"`
	ResetPasswordOnMismatchFlag        types.Bool   `tfsdk:"reset_password_on_mismatch_flag"`
	ChangeFrequencyType                types.String `tfsdk:"change_frequency_type"`
	ChangeFrequencyDays                types.Int32  `tfsdk:"change_frequency_days"`
	ChangeTime                         types.String `tfsdk:"change_time"`
	RemoteClientType                   types.String `tfsdk:"remote_client_type"`
	ApplicationHostID                  types.Int32  `tfsdk:"application_host_id"`
	IsApplicationHost                  types.Bool   `tfsdk:"is_application_host"`
	AccessURL                          types.String `tfsdk:"access_url"`
}

type ManagedSystemDataSourceModel struct {
	ManagedSystems []ManagedSystemModel `tfsdk:"managed_systems"`
}

func (d *ManagedSystemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Managed System Datasource",
		Blocks: map[string]schema.Block{
			"managed_systems": schema.ListNestedBlock{
				Description: "Managed System Datasource Attibutes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"workgroup_id": schema.Int32Attribute{
							MarkdownDescription: "Workgroup ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"host_name": schema.StringAttribute{
							MarkdownDescription: "Host Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: "IP Address",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"dns_name": schema.StringAttribute{
							MarkdownDescription: "DNS Name",
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
						"is_default_instance": schema.BoolAttribute{
							MarkdownDescription: "Is Default Instance",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"template": schema.StringAttribute{
							MarkdownDescription: "Template",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"forest_name": schema.StringAttribute{
							MarkdownDescription: "Forest Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"use_ssl": schema.BoolAttribute{
							MarkdownDescription: "Use SSL",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"managed_system_id": schema.Int32Attribute{
							MarkdownDescription: "Managed System ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"entity_type_id": schema.Int32Attribute{
							MarkdownDescription: "Entity Type ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"asset_id": schema.Int32Attribute{
							MarkdownDescription: "Asset ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"database_id": schema.Int32Attribute{
							MarkdownDescription: "Database ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"directory_id": schema.Int32Attribute{
							MarkdownDescription: "Directory ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"cloud_id": schema.Int32Attribute{
							MarkdownDescription: "Cloud ID",
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
						"timeout": schema.Int32Attribute{
							MarkdownDescription: "Timeout",
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
						"net_bios_name": schema.StringAttribute{
							MarkdownDescription: "NetBIOS Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"contact_email": schema.StringAttribute{
							MarkdownDescription: "Contact Email",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"port": schema.Int32Attribute{
							MarkdownDescription: "Port",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"ssh_key_enforcement_mode": schema.Int32Attribute{
							MarkdownDescription: "SSH Key Enforcement Mode",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"password_rule_id": schema.Int32Attribute{
							MarkdownDescription: "Password Rule ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"dss_key_rule_id": schema.Int32Attribute{
							MarkdownDescription: "DSS Key Rule ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"login_account_id": schema.Int32Attribute{
							MarkdownDescription: "Login Account ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"account_name_format": schema.Int32Attribute{
							MarkdownDescription: "Account Name Format",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"oracle_internet_directory_id": schema.StringAttribute{
							MarkdownDescription: "Oracle Internet Directory ID (GUID)",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"oracle_internet_directory_service_name": schema.StringAttribute{
							MarkdownDescription: "Oracle Internet Directory Service Name",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"release_duration": schema.Int32Attribute{
							MarkdownDescription: "Release Duration",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"max_release_duration": schema.Int32Attribute{
							MarkdownDescription: "Max Release Duration",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"isa_release_duration": schema.Int32Attribute{
							MarkdownDescription: "ISA Release Duration",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"auto_management_flag": schema.BoolAttribute{
							MarkdownDescription: "Auto Management Flag",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"functional_account_id": schema.Int32Attribute{
							MarkdownDescription: "Functional Account ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"elevation_command": schema.StringAttribute{
							MarkdownDescription: "Elevation Command",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"check_password_flag": schema.BoolAttribute{
							MarkdownDescription: "Check Password Flag",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"change_password_after_any_release_flag": schema.BoolAttribute{
							MarkdownDescription: "Change Password After Any Release Flag",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"reset_password_on_mismatch_flag": schema.BoolAttribute{
							MarkdownDescription: "Reset Password On Mismatch Flag",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"change_frequency_type": schema.StringAttribute{
							MarkdownDescription: "Change Frequency Type",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"change_frequency_days": schema.Int32Attribute{
							MarkdownDescription: "Change Frequency Days",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"change_time": schema.StringAttribute{
							MarkdownDescription: "Change Time",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"remote_client_type": schema.StringAttribute{
							MarkdownDescription: "Remote Client Type",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"application_host_id": schema.Int32Attribute{
							MarkdownDescription: "Application Host ID",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"is_application_host": schema.BoolAttribute{
							MarkdownDescription: "Is Application Host",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"access_url": schema.StringAttribute{
							MarkdownDescription: "Access URL",
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

func (d *ManagedSystemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	d.providerInfo = &c

	if d.providerInfo.userName == "" {
		return
	}

}

func (d *ManagedSystemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManagedSystemDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := utils.Autenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Authentication", err.Error())
		return
	}

	// instantiating managed system obj.
	managedSystemObj, _ := managed_systems.NewManagedSystem(*d.providerInfo.authenticationObj, zapLogger)

	// get managed systems list.
	items, err := managedSystemObj.GetManagedSystemsListFlow()

	if err != nil {
		resp.Diagnostics.AddError("Error getting managed systems list", err.Error())
		return
	}

	var managedAccountList []ManagedSystemModel

	for _, item := range items {
		managedAccountList = append(managedAccountList, ManagedSystemModel{
			ManagedSystemID:                    types.Int32Value(int32(item.ManagedSystemID)),
			EntityTypeID:                       types.Int32Value(int32(item.EntityTypeID)),
			AssetID:                            types.Int32Value(int32(item.AssetID)),
			DatabaseID:                         types.Int32Value(int32(item.DatabaseID)),
			DirectoryID:                        types.Int32Value(int32(item.DirectoryID)),
			CloudID:                            types.Int32Value(int32(item.CloudID)),
			WorkgroupID:                        types.Int32Value(int32(item.WorkgroupID)),
			HostName:                           types.StringValue(item.HostName),
			DNSName:                            types.StringValue(item.DnsName),
			IPAddress:                          types.StringValue(item.IPAddress),
			InstanceName:                       types.StringValue(item.InstanceName),
			IsDefaultInstance:                  types.BoolValue(item.IsDefaultInstance),
			Template:                           types.StringValue(item.Template),
			ForestName:                         types.StringValue(item.ForestName),
			UseSSL:                             types.BoolValue(item.UseSSL),
			OracleInternetDirectoryID:          types.StringValue(item.OracleInternetDirectoryID),
			OracleInternetDirectoryServiceName: types.StringValue(item.OracleInternetDirectoryServiceName),
			SystemName:                         types.StringValue(item.SystemName),
			PlatformID:                         types.Int32Value(int32(item.PlatformID)),
			NetBiosName:                        types.StringValue(item.NetBiosName),
			Port:                               types.Int32Value(int32(item.Port)),
			Timeout:                            types.Int32Value(int32(item.Timeout)),
			Description:                        types.StringValue(item.Description),
			ContactEmail:                       types.StringValue(item.ContactEmail),
			PasswordRuleID:                     types.Int32Value(int32(item.PasswordRuleID)),
			DSSKeyRuleID:                       types.Int32Value(int32(item.DSSKeyRuleID)),
			ReleaseDuration:                    types.Int32Value(int32(item.ReleaseDuration)),
			MaxReleaseDuration:                 types.Int32Value(int32(item.MaxReleaseDuration)),
			ISAReleaseDuration:                 types.Int32Value(int32(item.ISAReleaseDuration)),
			AutoManagementFlag:                 types.BoolValue(item.AutoManagementFlag),
			FunctionalAccountID:                types.Int32Value(int32(item.FunctionalAccountID)),
			LoginAccountID:                     types.Int32Value(int32(item.LoginAccountID)),
			ElevationCommand:                   types.StringValue(item.ElevationCommand),
			SshKeyEnforcementMode:              types.Int32Value(int32(item.SshKeyEnforcementMode)),
			CheckPasswordFlag:                  types.BoolValue(item.CheckPasswordFlag),
			ChangePasswordAfterAnyReleaseFlag:  types.BoolValue(item.ChangePasswordAfterAnyReleaseFlag),
			ResetPasswordOnMismatchFlag:        types.BoolValue(item.ResetPasswordOnMismatchFlag),
			ChangeFrequencyType:                types.StringValue(item.ChangeFrequencyType),
			ChangeFrequencyDays:                types.Int32Value(int32(item.ChangeFrequencyDays)),
			ChangeTime:                         types.StringValue(item.ChangeTime),
			AccountNameFormat:                  types.Int32Value(int32(item.AccountNameFormat)),
			RemoteClientType:                   types.StringValue(item.RemoteClientType),
			ApplicationHostID:                  types.Int32Value(int32(item.ApplicationHostID)),
			IsApplicationHost:                  types.BoolValue(item.IsApplicationHost),
			AccessURL:                          types.StringValue(item.AccessURL),
		})
	}

	responseData := ManagedSystemDataSourceModel{}
	responseData.ManagedSystems = managedAccountList

	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)

}
