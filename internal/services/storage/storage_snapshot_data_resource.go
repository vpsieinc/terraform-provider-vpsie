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
	_ datasource.DataSource              = &storageSnapshotDataSource{}
	_ datasource.DataSourceWithConfigure = &storageSnapshotDataSource{}
)

// storageSnapshotDataSource is the data source implementation.
type storageSnapshotDataSource struct {
	client *govpsie.Client
}

type storageSnapshotDataSourceModel struct {
	StorageSnapshots []storageSnapshotsModel `tfsdk:"storage_snapshots"`
	ID               types.String            `tfsdk:"id"`
}

type storageSnapshotsModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	StorageID   types.Int64  `tfsdk:"storage_id"`
	Identifier  types.String `tfsdk:"identifier"`
	Size        types.Int64  `tfsdk:"size"`
	CreatedOn   types.String `tfsdk:"created_on"`
	UserID      types.Int64  `tfsdk:"user_id"`
	IsDeleted   types.Int64  `tfsdk:"is_deleted"`
	SnapshotKey types.String `tfsdk:"snapshot_key"`
	StorageName types.String `tfsdk:"storage_name"`
	StorageType types.String `tfsdk:"storage_type"`
	DiskFormat  types.String `tfsdk:"disk_format"`
	BoxID       types.Int64  `tfsdk:"box_id"`
	EntityType  types.String `tfsdk:"entity_type"`
}

// NewstorageSnapshotDataSource is a helper function to create the data source.
func NewStorageSnapshotDataSource() datasource.DataSource {
	return &storageSnapshotDataSource{}
}

// Metadata returns the data source type name.
func (s *storageSnapshotDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_snapshots"
}

// Schema defines the schema for the data source.
func (s *storageSnapshotDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of all storage volume snapshots on the VPSie platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier for this data source.",
			},
			"storage_snapshots": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of storage volume snapshots.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The unique numeric identifier of the storage snapshot.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the storage snapshot.",
						},
						"storage_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The numeric ID of the storage volume this snapshot belongs to.",
						},
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the storage snapshot.",
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the storage snapshot in bytes.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the storage snapshot was created.",
						},
						"user_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user who owns the storage snapshot.",
						},
						"is_deleted": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the storage snapshot has been deleted (1 = deleted, 0 = active).",
						},
						"snapshot_key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The key identifier for the storage snapshot.",
						},
						"storage_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the storage volume this snapshot belongs to.",
						},
						"storage_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the storage volume (e.g., ssd, sata).",
						},
						"disk_format": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The disk format of the storage volume (e.g., EXT4, XFS).",
						},
						"box_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the server (box) the storage volume is attached to.",
						},
						"entity_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The entity type of the storage snapshot.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *storageSnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state storageSnapshotDataSourceModel

	storageSnapshots, err := s.client.Storage.ListSnapshots(ctx, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Storage Snapshots",
			"Could not get storage snapshots for storage: "+err.Error(),
		)

		return
	}

	for _, storageSnap := range storageSnapshots {
		storageSnapState := storageSnapshotsModel{
			ID:          types.Int64Value(int64(storageSnap.ID)),
			Name:        types.StringValue(storageSnap.Name),
			StorageID:   types.Int64Value(int64(storageSnap.StorageID)),
			Identifier:  types.StringValue(storageSnap.Identifier),
			Size:        types.Int64Value(int64(storageSnap.Size)),
			CreatedOn:   types.StringValue(storageSnap.CreatedOn.String()),
			UserID:      types.Int64Value(int64(storageSnap.UserID)),
			IsDeleted:   types.Int64Value(int64(storageSnap.IsDeleted)),
			SnapshotKey: types.StringValue(storageSnap.SnapshotKey),
			StorageName: types.StringValue(storageSnap.StorageName),
			StorageType: types.StringValue(storageSnap.StorageType),
			DiskFormat:  types.StringValue(storageSnap.DiskFormat),
			BoxID:       types.Int64Value(int64(storageSnap.BoxID)),
			EntityType:  types.StringValue(storageSnap.EntityType),
		}

		state.StorageSnapshots = append(state.StorageSnapshots, storageSnapState)

	}

	state.ID = types.StringValue("storage_snapshots")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *storageSnapshotDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	s.client = client
}
