package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VBoxDataSourceModel maps the Terraform schema to Go types
// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &vmsDataSource{}

type diskModel struct {
    Path string `tfsdk:"path"`
    Size int64  `tfsdk:"size_bytes"`
    SizeMB int64  `tfsdk:"size_mb"`
    SizeGB int64  `tfsdk:"size_gb"`
    Type string `tfsdk:"type"` // e.g., HardDisk or DVD
}

// In vmModel
type vmModel struct {
    ID     types.String `tfsdk:"id"`
    Name   types.String `tfsdk:"name"`
    Memory types.Int64  `tfsdk:"memory"`
    CPUs   types.Int64  `tfsdk:"cpus"`
    State  types.String `tfsdk:"state"`
    Description types.String `tfsdk:"description"`
    StorageControllers types.List `tfsdk:"storage_controllers"`
    Disks types.List `tfsdk:"disks"`
}

var vmObjectType = types.ObjectType{
    AttrTypes: map[string]attr.Type{
        "id":          types.StringType,
        "name":        types.StringType,
        "memory":      types.Int64Type,
        "cpus":        types.Int64Type,
        "state":       types.StringType,
        "description": types.StringType,
        "storage_controllers": types.ListType{ElemType: types.StringType},
        "disks": types.ListType{ElemType: diskObjectType},
    },
}

var diskObjectType = types.ObjectType{
    AttrTypes: map[string]attr.Type{
        "path":     types.StringType,
        "size_bytes": types.Int64Type,
        "size_mb":    types.Int64Type,
        "size_gb":    types.Int64Type,
        "type":     types.StringType,
    },
}

func getVMSchemaAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        "id":          schema.StringAttribute{Computed: true},
        "name":        schema.StringAttribute{Computed: true},
        "memory":      schema.Int64Attribute{Computed: true},
        "cpus":        schema.Int64Attribute{Computed: true},
        "state":       schema.StringAttribute{Computed: true},
        "description": schema.StringAttribute{Computed: true},
        "storage_controllers": schema.ListAttribute{
            Computed: true,
            ElementType: types.StringType,
            Description: "List of storage controllers attached to the VM.",
        },
        "disks": schema.ListNestedAttribute{
            Computed: true,
            NestedObject: schema.NestedAttributeObject{
                Attributes: map[string]schema.Attribute{
                    "path":       schema.StringAttribute{Computed: true},
                    "size_bytes": schema.Int64Attribute{Computed: true},
                    "size_mb":    schema.Int64Attribute{Computed: true},
                    "size_gb":    schema.Int64Attribute{Computed: true},
                    "type":       schema.StringAttribute{Computed: true},
                },
            },
        },
    }
}
// The main Data Source model
type vmsDataSourceModel struct {
    ID  types.String `tfsdk:"id"`
    VMs types.List   `tfsdk:"vms"` 
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
          "vms": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: getVMSchemaAttributes(),
                },
                Description: "List of VMs with their details.",
            },
        },
    }
}

func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data vmsDataSourceModel

    vms, err := d.client.GetDetailedVMs(ctx)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", err.Error())
        return
    }

    vmsList, diags := types.ListValueFrom(ctx, vmObjectType, vms)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    data.ID = types.StringValue("vbox-vms-list")
    data.VMs = vmsList

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}