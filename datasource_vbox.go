package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vmsDataSource struct {
	client *VBoxClient // Your custom SOAP client wrapper
}

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

func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state struct {
		Names []string `tfsdk:"names"`
	}

	// 1. Connect to vboxwebsrv SOAP endpoint
	// 2. Authenticate and get Session ID
	// 3. Call GetMachines()
	
	// Mocking the data for logic flow:
	state.Names = []string{"rhel9-prod", "ubuntu-test", "ansible-node-01"}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}