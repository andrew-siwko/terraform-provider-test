package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func main() {
	opts := provider.ServeOpts{
		Address: "registry.terraform.io/andrew/hashivar",
	}

	err := provider.Serve(context.Background(), New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func New() provider.Provider {
	return &hashivarProvider{}
}

type hashivarProvider struct{}

func (p *hashivarProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hashivar"
}

func (p *hashivarProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *hashivarProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ provider.ConfigureResponse) {}

func (p *hashivarProvider) DataSources(_ context.Context) []func() datasource.DataSource { return nil }

func (p *hashivarProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVariableResource,
	}
}

// --- Resource Logic ---

type variableResource struct{}

func NewVariableResource() resource.Resource { return &variableResource{} }

type VariableResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

func (r *variableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *variableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() { return }

	data.ID = types.StringValue("var-123") // Hardcoded ID for this simple example
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() { return }
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() { return }
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Delete(_ context.Context, _ resource.DeleteRequest, _ resource.DeleteResponse) {}