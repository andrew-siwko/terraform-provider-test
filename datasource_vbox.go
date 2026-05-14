package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VBoxDataSourceModel maps the Terraform schema to Go types
// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &vmsDataSource{}

// vmsDataSourceModel maps the data source schema data.
type vmsDataSourceModel struct {
    ID      types.String `tfsdk:"id"`
    Names   types.List   `tfsdk:"names"` // List of VM names
}

// NewVmsDataSource is a helper function to simplify the provider implementation.
func NewVmsDataSource() datasource.DataSource {
    return &vmsDataSource{}
}

// vmsDataSource is the data source implementation.
type vmsDataSource struct {
    client *VBoxClient // Consuming the shared client defined in provider.go
}

// Metadata returns the data source type name.
// This will result in "mirror_vms" in your HCL.
func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_vms"
}

// Configure adds the provider-level client to the data source.
func (d *vmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*VBoxClient)
    if !ok {
        resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *VBoxClient")
        return
    }

    d.client = client
}

// Schema defines the schema for the data source.
func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:            true,
                Description:         "Identifier for the data source.",
            },
            "names": schema.ListAttribute{
                ElementType: types.StringType,
                Computed:    true,
                Description: "List of all VirtualBox VM names.",
            },
        },
    }
}

// Read refreshes the Terraform state.
func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data vmsDataSourceModel

    // 1. Get the client from the provider
    // 2. Call your GetVMNames()
    names, err := d.client.GetVMNames()
    if err != nil {
        resp.Diagnostics.AddError("Client Error", err.Error())
        return
    }

    // Convert []string to types.List
    namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }    // 3. Map to your Go struct model
    
    data.Names = namesList

    // 4. Set the state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
