package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

var (
	_ datasource.DataSource              = &storageDataSource{}
	_ datasource.DataSourceWithConfigure = &storageDataSource{}
)

// storageDataSource is the data source implementation.
type storageDataSource struct {
	client StorageAPI
}

type storageDataSourceModel struct {
	Storages []storagesModel `tfsdk:"storages"`
	ID       types.String    `tfsdk:"id"`
}

type storagesModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	UserID         types.Int64  `tfsdk:"user_id"`
	BoxID          types.Int64  `tfsdk:"box_id"`
	Identifier     types.String `tfsdk:"identifier"`
	UserTemplateID types.Int64  `tfsdk:"user_template_id"`
	StorageType    types.String `tfsdk:"storage_type"`
	DiskFormat     types.String `tfsdk:"disk_format"`
	IsAutomatic    types.Int64  `tfsdk:"is_automatic"`
	Size           types.Int64  `tfsdk:"size"`
	StorageID      types.Int64  `tfsdk:"storage_id"`
	DiskKey        types.String `tfsdk:"disk_key"`
	CreatedOn      types.String `tfsdk:"created_on"`
	VmIdentifier   types.String `tfsdk:"vm_identifier"`
	Hostname       types.String `tfsdk:"hostname"`
	OsIdentifier   types.String `tfsdk:"os_identifier"`
	State          types.String `tfsdk:"state"`
	DcIdentifier   types.String `tfsdk:"dc_identifier"`
	BusDevice      types.String `tfsdk:"bus_device"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
}

// NewStorageDataSource is a helper function to create the data source.
func NewStorageDataSource() datasource.DataSource {
	return &storageDataSource{}
}

// Metadata returns the data source type name.
func (s *storageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storages"
}

// Schema defines the schema for the data source.
func (s *storageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about all VPSie storage volumes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source.",
			},
			"storages": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of storage volumes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The numeric ID of the storage volume.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the storage volume.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description of the storage volume.",
						},
						"user_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user who owns the storage volume.",
						},
						"box_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the server (box) the storage is attached to.",
						},
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the storage volume.",
						},
						"user_template_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user template associated with the storage volume.",
						},
						"storage_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of storage (e.g., ssd, sata).",
						},
						"disk_format": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The disk format of the storage volume.",
						},
						"is_automatic": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the storage volume was created automatically.",
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the storage volume in GB.",
						},
						"storage_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The internal storage ID.",
						},
						"disk_key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The disk key identifier for the storage volume.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the storage volume was created.",
						},
						"vm_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the VM the storage is attached to.",
						},
						"hostname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname of the server the storage is attached to.",
						},
						"os_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS identifier of the server the storage is attached to.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current state of the storage volume.",
						},
						"dc_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the data center where the storage volume resides.",
						},
						"bus_device": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The bus device name of the storage volume.",
						},
						"bus_number": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The bus number of the storage volume.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *storageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state storageDataSourceModel

	storages, err := s.client.ListAll(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Vpsie Storages",
			err.Error(),
		)

		return
	}

	for _, storage := range storages {
		storageState := storagesModel{
			ID:             types.Int64Value(int64(storage.ID)),
			Name:           types.StringValue(storage.Name),
			Description:    types.StringValue(storage.Description),
			UserID:         types.Int64Value(int64(storage.UserID)),
			BoxID:          types.Int64Value(int64(storage.BoxID)),
			Identifier:     types.StringValue(storage.Identifier),
			UserTemplateID: types.Int64Value(int64(storage.UserTemplateID)),
			StorageType:    types.StringValue(storage.StorageType),
			DiskFormat:     types.StringValue(storage.DiskFormat),
			IsAutomatic:    types.Int64Value(int64(storage.IsAutomatic)),
			Size:           types.Int64Value(int64(storage.Size)),
			StorageID:      types.Int64Value(int64(storage.StorageID)),
			DiskKey:        types.StringValue(storage.DiskKey),
			CreatedOn:      types.StringValue(storage.CreatedOn),
			VmIdentifier:   types.StringValue(storage.VmIdentifier),
			Hostname:       types.StringValue(storage.Hostname),
			OsIdentifier:   types.StringValue(storage.OsIdentifier),
			State:          types.StringValue(storage.State),
			DcIdentifier:   types.StringValue(storage.DcIdentifier),
			BusDevice:      types.StringValue(storage.BusDevice),
			BusNumber:      types.Int64Value(int64(storage.BusNumber)),
		}

		state.Storages = append(state.Storages, storageState)
	}

	state.ID = types.StringValue("storages")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *storageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*govpsie.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configuration Type",
			fmt.Sprintf("Expected *govpsie.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	s.client = client.Storage
}
