package main

import (
	"context"
	// "fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = provschema.Schema{
		Attributes: map[string]provschema.Attribute{
			"vboxwebsrv_endpoint": provschema.StringAttribute{
				Required:    true,
				Description: "The VirtualBox Web Service URL (e.g., http://127.0.0.1:18083).",
			},
			"username": provschema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"password": provschema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *mirrorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Define a local model to match the schema
	var data struct {
		Endpoint types.String `tfsdk:"vboxwebsrv_endpoint"`
		Username types.String `tfsdk:"username"`
		Password types.String `tfsdk:"password"`
	}

	// Read configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize your shared client
	client := &VBoxClient{
		Endpoint: data.Endpoint.ValueString(),
		Username: data.Username.ValueString(),
		Password: data.Password.ValueString(),
	}

	// Pass the client to all Data Sources and Resources
	resp.DataSourceData = client
	resp.ResourceData = client
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
