package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// GetCreateManagedSystemCommonAttributes get common attributes to create managed systems by asset, workgroup and database
func GetCreateManagedSystemCommonAttributes() map[string]schema.Attribute {
	commonAttributes := map[string]schema.Attribute{
		"managed_system_id": schema.Int32Attribute{
			MarkdownDescription: "Managed System Id",
			Required:            false,
			Optional:            false,
			Computed:            true,
		},
		"managed_system_name": schema.StringAttribute{
			MarkdownDescription: "Managed System Name",
			Required:            false,
			Optional:            false,
			Computed:            true,
		},
		"contact_email": schema.StringAttribute{
			MarkdownDescription: "Contact Email (max 1000 characters, must be a valid email)",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description (max 255 characters)",
			Optional:            true,
		},
		"timeout": schema.Int32Attribute{
			MarkdownDescription: "Timeout",
			Optional:            true,
		},
		"password_rule_id": schema.Int32Attribute{
			MarkdownDescription: "Password Rule ID",
			Optional:            true,
		},
		"release_duration": schema.Int32Attribute{
			MarkdownDescription: "Release Duration (min: 1, max: 525600)",
			Optional:            true,
			Computed:            true,
			Default:             int32default.StaticInt32(120),
		},
		"max_release_duration": schema.Int32Attribute{
			MarkdownDescription: "Max Release Duration (min: 1, max: 525600)",
			Optional:            true,
			Computed:            true,
			Default:             int32default.StaticInt32(525600),
		},
		"isa_release_duration": schema.Int32Attribute{
			MarkdownDescription: "ISA Release Duration (min: 1, max: 525600)",
			Optional:            true,
			Computed:            true,
			Default:             int32default.StaticInt32(120),
		},
		"auto_management_flag": schema.BoolAttribute{
			MarkdownDescription: "Auto Management Flag",
			Optional:            true,
		},
		"functional_account_id": schema.Int32Attribute{
			MarkdownDescription: "Functional Account ID (required if AutoManagementFlag is true)",
			Optional:            true,
		},
		"check_password_flag": schema.BoolAttribute{
			MarkdownDescription: "Check Password Flag",
			Optional:            true,
		},
		"change_password_after_any_release_flag": schema.BoolAttribute{
			MarkdownDescription: "Change Password After Any Release Flag",
			Optional:            true,
		},
		"reset_password_on_mismatch_flag": schema.BoolAttribute{
			MarkdownDescription: "Reset Password On Mismatch Flag",
			Optional:            true,
		},
		"change_frequency_type": schema.StringAttribute{
			MarkdownDescription: "Change Frequency Type (one of: first, last, xdays)",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("first"),
		},
		"change_frequency_days": schema.Int32Attribute{
			MarkdownDescription: "Change Frequency Days (required if ChangeFrequencyType is xdays)",
			Optional:            true,
		},
		"change_time": schema.StringAttribute{
			MarkdownDescription: "Change Time (format: HH:MM)",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("23:30"),
		},
	}

	return commonAttributes
}

// GetInt32Attribute build and get int32 attribute.
func GetInt32Attribute(description string, required, optional, computed bool) schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: description,
		Required:            required,
		Optional:            optional,
		Computed:            computed,
	}
}

// GetStringAttribute build and get string attribute.
func GetStringAttribute(description string, required, optional, computed bool) schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: description,
		Required:            required,
		Optional:            optional,
		Computed:            computed,
	}
}

// GetBoolAttribute build and get bool attribute.
func GetBoolAttribute(description string, required, optional, computed bool) schema.Attribute {
	return schema.BoolAttribute{
		MarkdownDescription: description,
		Required:            required,
		Optional:            optional,
		Computed:            computed,
	}
}
