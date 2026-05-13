package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &coffeeDataSource{}

// coffeeDataSourceModel maps the data source schema data.
type coffeeDataSourceModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Price types.Float64 `tfsdk:"price"`
}

// NewCoffeeDataSource is a helper function to simplify the provider implementation.
func NewCoffeeDataSource() datasource.DataSource {
	return &coffeeDataSource{}
}

// coffeeDataSource is the data source implementation.
type coffeeDataSource struct {
	client string // In a real app, this would be your API client
}

// Metadata returns the data source type name.
func (d *coffeeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffee"
}

// Schema defines the schema for the data source.
func (d *coffeeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the coffee.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Product name of the coffee.",
			},
			"price": schema.Float64Attribute{
				Computed:    true,
				Description: "The price per cup.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data from the API.
func (d *coffeeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data coffeeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Logic to fetch data from your API would go here.
	// For this example, we'll hardcode a "mock" response.
	if data.ID.ValueInt64() == 1 {
		data.Name = types.StringValue("Espresso")
		data.Price = types.Float64Value(2.50)
	} else {
		resp.Diagnostics.AddError("Coffee not found", "Could not find coffee with the provided ID")
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}