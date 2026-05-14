package main

import (
	"context"
	// "log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func New() provider.Provider {
	return &hashivarProvider{}
}

type hashivarProvider struct{}

func (p *hashivarProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mirror"
}

func (p *hashivarProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provschema.Schema{}
}

// Fixed: Added pointer (*) to ConfigureResponse
func (p *hashivarProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {}

func (p *hashivarProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewCoffeeDataSource, // This refers to the function in our other file
    }
}
func (p *hashivarProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVariableResource,
	}
}

var _ datasource.DataSource = &vmsDataSource{}

type vmsDataSource struct {
	client *VBoxClient
}

func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

// Configure allows the provider to pass the initialized client to the data source
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

func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state struct {
		Names []string `tfsdk:"names"`
	}

	// This is where you would call d.client.Endpoint and use a SOAP library
	// to fetch the real names from your RHEL 9 VirtualBox instance.
	state.Names = []string{"rhel9-vm-01", "dev-server", "test-node"}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}