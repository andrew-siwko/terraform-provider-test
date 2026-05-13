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

