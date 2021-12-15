package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/iancoleman/strcase"
)

type transformResourceType struct{}

func (t transformResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Transform resource",

		Attributes: map[string]tfsdk.Attribute{
			"source": {
				MarkdownDescription: "Source value for transform",
				Optional:            false,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Type of transform",
				Optional:            true,
				Type:                types.StringType,
			},
			"result": {
				Computed:            true,
				MarkdownDescription: "Transform result",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t transformResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return transformResource{
		provider: provider,
	}, diags
}

type transformResourceData struct {
	SourceString  types.String `tfsdk:"source"`
	TransformType types.String `tfsdk:"type"`
	Result        types.String `tfsdk:"result"`
}

type transformResource struct {
	provider provider
}

func (r transformResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data transformResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// transform, err := d.provider.client.CreateTransform(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create transform, got error: %s", err))
	//     return
	// }

	// For the purposes of this transform code, hardcoding a response value to
	// save into the Terraform state.
	data.TransformType = types.String{Value: "camelCase"}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r transformResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data transformResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	s := ""
	switch data.TransformType.Value {
	case "camel":
		s = strcase.ToLowerCamel(data.SourceString.Value)
		break
	case "snake":
		s = strcase.ToSnake(data.SourceString.Value)
		break
	default:
		s = strcase.ToCamel(data.SourceString.Value)
	}

	data.Result = types.String{Value: s}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// transform, err := d.provider.client.ReadTransform(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read transform, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r transformResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data transformResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// transform, err := d.provider.client.UpdateTransform(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update transform, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r transformResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data transformResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// transform, err := d.provider.client.DeleteTransform(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete transform, got error: %s", err))
	//     return
	// }

	resp.State.RemoveResource(ctx)
}

func (r transformResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
