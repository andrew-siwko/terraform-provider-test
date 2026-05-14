package main

import (
	"context"
	// "fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type VBoxClient struct {
	Endpoint string
	Username string
	Password string
}

type mirrorProvider struct{}

func New() provider.Provider {
	return &mirrorProvider{}
}

func (p *mirrorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mirror"
}

func (p *mirrorProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provschema.Schema{}
}

// Fixed: Added pointer (*) to ConfigureResponse
func (p *mirrorProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *mirrorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeeDataSource,
		NewVmsDataSource,
	}
}

func (p *mirrorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVariableResource,
	}
}
