package backup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vpsie/govpsie"
)

type backupDataSource struct {
	client BackupAPI
}

type backupDataSourceModel struct {
	Backups []backupsModel `tfsdk:"backups"`
	ID      types.String   `tfsdk:"id"`
}

type backupsModel struct {
	Identifier   types.String `tfsdk:"identifier"`
	CreatedOn    types.String `tfsdk:"created_on"`
	DcIdentifier types.String `tfsdk:"dc_identifier"`
	CreatedBy    types.String `tfsdk:"created_by"`
	HostName     types.String `tfsdk:"hostname"`
	Name         types.String `tfsdk:"name"`
	Note         types.String `tfsdk:"note"`
	BackupKey    types.String `tfsdk:"backup_key"`
	State        types.String `tfsdk:"state"`
	VMIdentifier types.String `tfsdk:"vm_identifier"`
	BoxID        types.Int64  `tfsdk:"box_id"`
	BackupSHA1   types.String `tfsdk:"backupsha1"`
	OSFullName   types.String `tfsdk:"os_full_name"`
	VMCategory   types.String `tfsdk:"vm_category"`
}

// NewBackupDataSource is a helper function to create the data source.
func NewBackupDataSource() datasource.DataSource {
	return &backupDataSource{}
}

// Metadata returns the data source type name.
func (b *backupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups"
}

// Schema defines the schema for the data source.
func (b *backupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of backups on the VPSie platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source.",
			},
			"backups": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of backups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the backup.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the backup was created.",
						},
						"dc_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the data center where the backup is stored.",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user who created the backup.",
						},
						"hostname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname of the server associated with the backup.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the backup.",
						},
						"note": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "An optional note or description for the backup.",
						},
						"backup_key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The key used to identify the backup.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current state of the backup.",
						},
						"vm_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the virtual machine associated with the backup.",
						},
						"box_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The numeric box ID of the backup.",
						},
						"backupsha1": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The SHA1 checksum of the backup.",
						},
						"os_full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the operating system in the backup.",
						},
						"vm_category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The category of the virtual machine associated with the backup.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (b *backupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backupDataSourceModel

	backups, err := b.client.List(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting backups",
			"Could not get backups, unexpected error: "+err.Error(),
		)

		return
	}

	for _, backup := range backups {
		backupState := backupsModel{
			Identifier:   types.StringValue(backup.Identifier),
			CreatedOn:    types.StringValue(backup.CreatedOn),
			DcIdentifier: types.StringValue(backup.DcIdentifier),
			CreatedBy:    types.StringValue(backup.CreatedBy),
			HostName:     types.StringValue(backup.HostName),
			Name:         types.StringValue(backup.Name),
			Note:         types.StringValue(backup.Note),
			BackupKey:    types.StringValue(backup.BackupKey),
			State:        types.StringValue(backup.State),
			VMIdentifier: types.StringValue(backup.VMIdentifier),
			BoxID:        types.Int64Value(int64(backup.BoxID)),
			BackupSHA1:   types.StringValue(backup.BackupSHA1),
			OSFullName:   types.StringValue(backup.OSFullName),
			VMCategory:   types.StringValue(backup.VMCategory),
		}

		state.Backups = append(state.Backups, backupState)
	}

	state.ID = types.StringValue("backups")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
func (b *backupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	b.client = client.Backup
}
