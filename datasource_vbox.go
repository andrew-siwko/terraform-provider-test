package main

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &vmsDataSource{}

// vmsDataSourceModel maps the data source schema data.
type vmsDataSourceModel struct {
    Names []types.String `tfsdk:"names"`
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

    // Mocking the data logic. In your RHEL 9 lab, this is where 
    // you'd use d.client to talk to vboxwebsrv.
    vmNames := []string{"rhel9-prod", "ubuntu-test", "ansible-node-01"}
    
    for _, name := range vmNames {
        data.Names = append(data.Names, types.StringValue(name))
    }

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}