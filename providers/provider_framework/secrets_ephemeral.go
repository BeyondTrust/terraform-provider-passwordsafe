package providerv2

import (
	"context"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &EphemeralSecret{}

// @EphemeralResource(passwordsafe_secret_ephemeral, name="Secret Version")
func NewEphemeralSecret() ephemeral.EphemeralResource {
	return &EphemeralSecret{}
}

type EphemeralSecret struct {
	providerInfo *ProviderData
}

type EphemeralSecretModel struct {
	Title     types.String `tfsdk:"title"`
	Path      types.String `tfsdk:"path"`
	Separator types.String `tfsdk:"separator"`
	Value     types.String `tfsdk:"value"`
}

func (e *EphemeralSecret) Metadata(ctx context.Context, _ ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = "passwordsafe_secret_ephemeral"
}

func (e *EphemeralSecret) Schema(ctx context.Context, _ ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{

		MarkdownDescription: "Schema of Secret Retrieval",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Secret path",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1792),
				},
			},
			"title": schema.StringAttribute{
				Description: "Secret title",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"separator": schema.StringAttribute{
				Description: "Separator",
				Optional:    true,
			},
			"value": schema.StringAttribute{
				Description: "Value",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}
func (e *EphemeralSecret) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {

	// getting data from Provider
	c, _ := req.ProviderData.(ProviderData)

	e.providerInfo = &c

	if e.providerInfo.userName == "" {
		return
	}

}

func (e *EphemeralSecret) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {

	var data EphemeralSecretModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// instantiating secret obj
	secretObj, err := secrets.NewSecretObj(*e.providerInfo.authenticationObj, zapLogger, 5000000)

	if err != nil {
		response.Diagnostics.AddError("Error getting secret", err.Error())
		return
	}

	separator := "/"
	if data.Separator.ValueString() != "" {
		separator = data.Separator.ValueString()
	}

	// getting single secret from PS API
	secret, err := secretObj.GetSecret(data.Path.ValueString()+separator+data.Title.ValueString(), separator)

	if err != nil {
		response.Diagnostics.AddError("Error getting secret", err.Error())
		return
	}

	// setting secret to value attribute
	data.Value = types.StringValue(secret)

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)

}
