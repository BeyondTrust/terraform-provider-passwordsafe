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
		Description: "Managed System Datasource, get managed systems list.",
		Blocks: map[string]schema.Block{
			"managed_systems": schema.ListNestedBlock{
				Description: "Managed System Datasource Attributes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"workgroup_id":                           utils.GetInt32Attribute("Workgroup ID", false, false, true),
						"host_name":                              utils.GetStringAttribute("Host Name", false, false, true),
						"ip_address":                             utils.GetStringAttribute("IP Address", false, false, true),
						"dns_name":                               utils.GetStringAttribute("DNS Name", false, false, true),
						"instance_name":                          utils.GetStringAttribute("Instance Name", false, false, true),
						"is_default_instance":                    utils.GetBoolAttribute("Is Default Instance", false, false, true),
						"template":                               utils.GetStringAttribute("Template", false, false, true),
						"forest_name":                            utils.GetStringAttribute("Forest Name", false, false, true),
						"use_ssl":                                utils.GetBoolAttribute("Use SSL", false, false, true),
						"managed_system_id":                      utils.GetInt32Attribute("Managed System ID", false, false, true),
						"entity_type_id":                         utils.GetInt32Attribute("Entity Type ID", false, false, true),
						"asset_id":                               utils.GetInt32Attribute("Asset ID", false, false, true),
						"database_id":                            utils.GetInt32Attribute("Database ID", false, false, true),
						"directory_id":                           utils.GetInt32Attribute("Directory ID", false, false, true),
						"cloud_id":                               utils.GetInt32Attribute("Cloud ID", false, false, true),
						"system_name":                            utils.GetStringAttribute("System Name", false, false, true),
						"timeout":                                utils.GetInt32Attribute("Timeout", false, false, true),
						"platform_id":                            utils.GetInt32Attribute("Platform ID", false, false, true),
						"net_bios_name":                          utils.GetStringAttribute("NetBIOS Name", false, false, true),
						"contact_email":                          utils.GetStringAttribute("Contact Email", false, false, true),
						"description":                            utils.GetStringAttribute("Description", false, false, true),
						"port":                                   utils.GetInt32Attribute("Port", false, false, true),
						"ssh_key_enforcement_mode":               utils.GetInt32Attribute("SSH Key Enforcement Mode", false, false, true),
						"password_rule_id":                       utils.GetInt32Attribute("Password Rule ID", false, false, true),
						"dss_key_rule_id":                        utils.GetInt32Attribute("DSS Key Rule ID", false, false, true),
						"login_account_id":                       utils.GetInt32Attribute("Login Account ID", false, false, true),
						"account_name_format":                    utils.GetInt32Attribute("Account Name Format", false, false, true),
						"oracle_internet_directory_id":           utils.GetStringAttribute("Oracle Internet Directory ID (GUID)", false, false, true),
						"oracle_internet_directory_service_name": utils.GetStringAttribute("Oracle Internet Directory Service Name", false, false, true),
						"release_duration":                       utils.GetInt32Attribute("Release Duration", false, false, true),
						"max_release_duration":                   utils.GetInt32Attribute("Max Release Duration", false, false, true),
						"isa_release_duration":                   utils.GetInt32Attribute("ISA Release Duration", false, false, true),
						"auto_management_flag":                   utils.GetBoolAttribute("Auto Management Flag", false, false, true),
						"functional_account_id":                  utils.GetInt32Attribute("Functional Account ID", false, false, true),
						"elevation_command":                      utils.GetStringAttribute("Elevation Command", false, false, true),
						"check_password_flag":                    utils.GetBoolAttribute("Check Password Flag", false, false, true),
						"change_password_after_any_release_flag": utils.GetBoolAttribute("Change Password After Any Release Flag", false, false, true),
						"reset_password_on_mismatch_flag":        utils.GetBoolAttribute("Reset Password On Mismatch Flag", false, false, true),
						"change_frequency_type":                  utils.GetStringAttribute("Change Frequency Type", false, false, true),
						"change_frequency_days":                  utils.GetInt32Attribute("Change Frequency Days", false, false, true),
						"change_time":                            utils.GetStringAttribute("Change Time", false, false, true),
						"remote_client_type":                     utils.GetStringAttribute("Remote Client Type", false, false, true),
						"application_host_id":                    utils.GetInt32Attribute("Application Host ID", false, false, true),
						"is_application_host":                    utils.GetBoolAttribute("Is Application Host", false, false, true),
						"access_url":                             utils.GetStringAttribute("Access URL", false, false, true),
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

	_, err := utils.Authenticate(*d.providerInfo.authenticationObj, &mu, &signInCount, zapLogger)
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
