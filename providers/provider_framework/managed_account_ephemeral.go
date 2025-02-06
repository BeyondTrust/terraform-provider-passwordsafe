package providerv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	managed_accounts "github.com/BeyondTrust/go-client-library-passwordsafe/api/managed_account"
)

// @EphemeralResource(passwordsafe_managed_acccount_ephemeral, name="Secret Version")
func NewEphemeralManagedAccount() ephemeral.EphemeralResource {
	return &EphemeralManagedAccount{}
}

type EphemeralManagedAccount struct {
	providerInfo *ProviderData
}

type EphemeralManagedAccountModel struct {
	SystemName  types.String `tfsdk:"system_name"`
	AccountName types.String `tfsdk:"account_name"`
	Value       types.String `tfsdk:"value"`
}

func (e *EphemeralManagedAccount) Metadata(ctx context.Context, _ ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = "passwordsafe_managed_acccount_ephemeral"
}

func (e *EphemeralManagedAccount) Schema(ctx context.Context, _ ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{

		MarkdownDescription: "Schema of Managed Account Retrieval",

		Attributes: map[string]schema.Attribute{
			"system_name": schema.StringAttribute{
				Description: "System account name",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"account_name": schema.StringAttribute{
				Description: "Managed account name",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 245),
				},
			},
			"value": schema.StringAttribute{
				Description: "Value",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}
func (e *EphemeralManagedAccount) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	e.providerInfo = &c

	if e.providerInfo.userName == "" {
		return
	}

}

func (e *EphemeralManagedAccount) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {

	var data EphemeralManagedAccountModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// instantiating managed account obj
	manageAccountObj, err := managed_accounts.NewManagedAccountObj(*e.providerInfo.authenticationObj, zapLogger)

	if err != nil {
		response.Diagnostics.AddError("Error getting managed account", err.Error())
		return
	}

	// getting single managed account from PS API
	gotManagedAccount, err := manageAccountObj.GetSecret(data.SystemName.ValueString()+"/"+data.AccountName.ValueString(), "/")

	if err != nil {
		response.Diagnostics.AddError("Error getting managed account", err.Error())
		return
	}

	// setting secret to value attribute
	data.Value = types.StringValue(gotManagedAccount)

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)

}
