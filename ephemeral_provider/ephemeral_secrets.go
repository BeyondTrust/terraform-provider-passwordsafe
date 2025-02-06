package ephemeral_provider

import (
	"context"

	"github.com/BeyondTrust/go-client-library-passwordsafe/api/secrets"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ephemeral.EphemeralResource = &EphemeralSecret{}

// @EphemeralResource(passwordsafe_secret_ephemeral_version, name="Secret Version")
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
	response.TypeName = "passwordsafe_secret_ephemeral_version"
}

func (e *EphemeralSecret) Schema(ctx context.Context, _ ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{

		MarkdownDescription: "Schema of Managed Account Retrieval",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Secret path",
				Required:    true,
			},
			"title": schema.StringAttribute{
				Description: "Secret title",
				Required:    true,
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
	tflog.Info(ctx, "Log: ", map[string]interface{}{
		"User Name=>": e.providerInfo.userName,
	})

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
		response.Diagnostics.AddError("Error getting managed account", err.Error())
		return
	}

	// getting single secret from PS API
	secret, err := secretObj.GetSecret(data.Path.ValueString()+data.Separator.ValueString()+data.Title.ValueString(), data.Separator.ValueString())

	if err != nil {
		response.Diagnostics.AddError("Error getting managed account", err.Error())
		return
	}

	// setting secret to value attribute
	data.Value = types.StringValue(secret)

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)

}
