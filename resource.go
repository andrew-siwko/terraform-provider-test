package main

import (
	"context"
	// "log"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	resp.Schema = resschema.Schema{
		// You can also add a description for the resource itself!
		MarkdownDescription: "Manages a specific configuration variable within the system.",

		Attributes: map[string]resschema.Attribute{
			"id": resschema.StringAttribute{
				MarkdownDescription: "Unique identifier for the variable.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": resschema.StringAttribute{
				MarkdownDescription: "The actual value assigned to the variable. Must be a valid string.",
				Required:            true,
			},
		},
	}
}

func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue("var-123")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
