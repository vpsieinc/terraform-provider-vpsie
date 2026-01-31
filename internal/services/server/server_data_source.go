package server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vpsie/govpsie"
)

type serverDataSource struct {
	client *govpsie.Client
}

type serverDataSourceModel struct {
	Servers []serversModel `tfsdk:"servers"`
	ID      types.String   `tfsdk:"id"`
}

type serversModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	Identifier          types.String `tfsdk:"identifier"`
	UserID              types.Int64  `tfsdk:"user_id"`
	BoxSizeID           types.Int64  `tfsdk:"boxsize_id"`
	BoxImageID          types.Int64  `tfsdk:"boximage_id"`
	DataCenterID        types.Int64  `tfsdk:"datacenter_id"`
	NodeID              types.Int64  `tfsdk:"node_id"`
	BoxdIsCountID       types.Int64  `tfsdk:"boxdiscount_id"`
	Hostname            types.String `tfsdk:"hostname"`
	DefaultIP           types.String `tfsdk:"default_ip"`
	DefaultIPv6         types.String `tfsdk:"default_ipv6"`
	PrivateIP           types.String `tfsdk:"private_ip"`
	IsAutoBackup        types.Int64  `tfsdk:"is_autobackup"`
	BoxVirtualization   types.String `tfsdk:"box_virtualization_id"`
	Ram                 types.Int64  `tfsdk:"ram"`
	Cpu                 types.Int64  `tfsdk:"cpu"`
	Ssd                 types.Int64  `tfsdk:"ssd"`
	Traffic             types.Int64  `tfsdk:"traffic"`
	AddedIpAddresses    types.String `tfsdk:"added_ip_addresses"`
	InitialPassword     types.String `tfsdk:"initial_password"`
	Notes               types.String `tfsdk:"notes"`
	CreatedOn           types.String `tfsdk:"created_on"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	DroppedOn           types.String `tfsdk:"dropped_on"`
	IsActive            types.Int64  `tfsdk:"is_active"`
	IsDeleted           types.Int64  `tfsdk:"is_deleted"`
	Power               types.Int64  `tfsdk:"power"`
	ProjectID           types.Int64  `tfsdk:"project_id"`
	IsCustom            types.Int64  `tfsdk:"is_custom"`
	NrAddedIps          types.Int64  `tfsdk:"nr_added_ips"`
	InPcs               types.Int64  `tfsdk:"in_pcs"`
	CustomPrice         types.Int64  `tfsdk:"custom_price"`
	PayableLicense      types.Int64  `tfsdk:"payable_license"`
	LastLicensePay      types.String `tfsdk:"last_license_pay"`
	ScriptID            types.String `tfsdk:"script_id"`
	SshKeyID            types.String `tfsdk:"sshkey_id"`
	IsLocked            types.Int64  `tfsdk:"is_locked"`
	IsWorkWithNew       types.Int64  `tfsdk:"is_work_with_new_version"`
	IsSuspended         types.Int64  `tfsdk:"is_suspended"`
	IsTerminated        types.Int64  `tfsdk:"is_terminated"`
	OldID               types.Int64  `tfsdk:"old_id"`
	CustomIsoID         types.Int64  `tfsdk:"custom_iso_id"`
	IsIsoImageBootAble  types.Int64  `tfsdk:"is_iso_image_bootable"`
	HasSsl              types.Int64  `tfsdk:"has_ssl"`
	LastActionDate      types.String `tfsdk:"last_action_date"`
	IsCreatedFromLegacy types.Int64  `tfsdk:"is_created_from_legacy"`
	IsSmtpAllowed       types.Int64  `tfsdk:"is_smtp_allowed"`
	WeeklyBackup        types.Int64  `tfsdk:"weekly_backup"`
	MonthlyBackup       types.Int64  `tfsdk:"monthly_backup"`
	LibIsoID            types.Int64  `tfsdk:"lib_iso_id"`
	DailySnapshot       types.Int64  `tfsdk:"daily_snapshot"`
	WeeklySnapshot      types.Int64  `tfsdk:"weekly_snapshot"`
	MonthlySnapshot     types.Int64  `tfsdk:"monthly_snap"`
	LastActionInMin     types.Int64  `tfsdk:"last_action_in_min"`
	FirstName           types.String `tfsdk:"firstname"`
	LastName            types.String `tfsdk:"lastname"`
	Username            types.String `tfsdk:"username"`
	State               types.String `tfsdk:"state"`
	IsFipAvailable      types.Int64  `tfsdk:"is_fip_available"`
	IsBucketAvailable   types.Int64  `tfsdk:"is_bucket_available"`
	DcIdentifier        types.String `tfsdk:"dc_identifier"`
	Category            types.String `tfsdk:"category"`
	FullName            types.String `tfsdk:"fullname"`
	VmDescription       types.String `tfsdk:"vm_description"`
	BoxesSuspended      types.Int64  `tfsdk:"boxes_suspended"`
	IsSataAvailable     types.Int64  `tfsdk:"is_sata_available"`
	IsSsdAvailable      types.Int64  `tfsdk:"is_ssd_available"`
	PublicIp            types.String `tfsdk:"public_ip"`
}

// NewServerDataSource is a helper function to create the data source.
func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{}
}

// Metadata returns the data source type name.
func (s *serverDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

// Schema defines the schema for the data source.
func (s *serverDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about all VPSie servers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source.",
			},
			"servers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of servers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The numeric ID of the server.",
						},
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the server.",
						},
						"user_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user who owns the server.",
						},
						"boxsize_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the box size (resource plan) for the server.",
						},
						"boximage_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the box image (OS template) for the server.",
						},
						"datacenter_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the data center where the server is located.",
						},
						"node_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the physical node hosting the server.",
						},
						"boxdiscount_id": schema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The ID of the discount applied to the server.",
						},
						"hostname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname assigned to the server.",
						},
						"default_ip": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default IPv4 address assigned to the server.",
						},
						"default_ipv6": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default IPv6 address assigned to the server.",
						},
						"private_ip": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The private IP address assigned to the server.",
						},
						"is_autobackup": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether automatic backup is enabled for the server.",
						},
						"box_virtualization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The virtualization type identifier for the server.",
						},
						"ram": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The amount of RAM in MB allocated to the server.",
						},
						"cpu": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of CPU cores allocated to the server.",
						},
						"ssd": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The SSD storage size in GB allocated to the server.",
						},
						"traffic": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The traffic bandwidth limit allocated to the server.",
						},
						"added_ip_addresses": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Additional IP addresses added to the server.",
						},
						"initial_password": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: "The initial root password for the server.",
						},
						"notes": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "Optional notes or comments for the server.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the server was created.",
						},
						"last_updated": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the server was last updated.",
						},
						"dropped_on": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The timestamp when the server was dropped or deleted.",
						},
						"is_active": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server is currently active.",
						},
						"is_deleted": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server has been deleted.",
						},
						"power": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The power state of the server (0 = off, 1 = on).",
						},
						"project_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the project to which the server belongs.",
						},
						"is_custom": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server uses a custom configuration.",
						},
						"nr_added_ips": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of additional IP addresses added to the server.",
						},
						"in_pcs": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of processes running on the server.",
						},
						"custom_price": schema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The custom price applied to the server.",
						},
						"payable_license": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The payable license cost for the server.",
						},
						"last_license_pay": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The date of the last license payment.",
						},
						"script_id": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The identifier of a startup script on the server.",
						},
						"sshkey_id": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The identifier of an SSH key on the server.",
						},
						"is_locked": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server is locked from modifications.",
						},
						"is_work_with_new_version": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server is compatible with the new platform version.",
						},
						"is_suspended": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server is currently suspended.",
						},
						"is_terminated": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server has been terminated.",
						},
						"old_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The legacy ID of the server from the previous platform.",
						},
						"custom_iso_id": schema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The ID of a custom ISO image attached to the server.",
						},
						"is_iso_image_bootable": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the attached ISO image is bootable.",
						},
						"has_ssl": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether SSL is enabled for the server.",
						},
						"last_action_date": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The date of the last action performed on the server.",
						},
						"is_created_from_legacy": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether the server was migrated from the legacy platform.",
						},
						"is_smtp_allowed": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether SMTP traffic is allowed on the server.",
						},
						"weekly_backup": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether weekly backups are enabled for the server.",
						},
						"monthly_backup": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether monthly backups are enabled for the server.",
						},
						"lib_iso_id": schema.Int64Attribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The ID of the library ISO image attached to the server.",
						},
						"daily_snapshot": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether daily snapshots are enabled for the server.",
						},
						"weekly_snapshot": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether weekly snapshots are enabled for the server.",
						},
						"monthly_snap": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether monthly snapshots are enabled for the server.",
						},
						"last_action_in_min": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The time in minutes since the last action on the server.",
						},
						"firstname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The first name of the server owner.",
						},
						"lastname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last name of the server owner.",
						},
						"username": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The username of the server owner.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current state of the server.",
						},
						"is_fip_available": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether floating IP is available for the server.",
						},
						"is_bucket_available": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether object storage bucket is available for the server.",
						},
						"dc_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the data center where the server is deployed.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The category of the server.",
						},
						"fullname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the server owner.",
						},
						"vm_description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the virtual machine.",
						},
						"boxes_suspended": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of suspended boxes for the owner.",
						},
						"is_sata_available": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether SATA storage is available for the server.",
						},
						"is_ssd_available": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether SSD storage is available for the server.",
						},
						"public_ip": schema.StringAttribute{
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The public IP address of the server.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverDataSourceModel

	servers, err := s.client.Server.List(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading servers",
			"Could not read servers: "+err.Error(),
		)

		return
	}

	for _, server := range servers {
		serverState := serversModel{
			ID:                  types.Int64Value(server.ID),
			Identifier:          types.StringValue(server.Identifier),
			UserID:              types.Int64Value(server.UserID),
			BoxSizeID:           types.Int64Value(server.BoxSizeID),
			BoxImageID:          types.Int64Value(server.BoxImageID),
			DataCenterID:        types.Int64Value(server.DataCenterID),
			NodeID:              types.Int64Value(server.NodeID),
			BoxdIsCountID:       types.Int64PointerValue(server.BoxdIsCountID),
			Hostname:            types.StringValue(server.Hostname),
			DefaultIP:           types.StringValue(server.DefaultIP),
			DefaultIPv6:         types.StringValue(server.DefaultIPv6),
			PrivateIP:           types.StringValue(server.PrivateIP),
			IsAutoBackup:        types.Int64Value(server.IsAutoBackup),
			BoxVirtualization:   types.StringValue(server.BoxVirtualization),
			Ram:                 types.Int64Value(server.Ram),
			Cpu:                 types.Int64Value(server.Cpu),
			Ssd:                 types.Int64Value(server.Ssd),
			Traffic:             types.Int64Value(server.Traffic),
			AddedIpAddresses:    types.StringPointerValue(server.AddedIpAddresses),
			InitialPassword:     types.StringValue(server.InitialPassword),
			Notes:               types.StringPointerValue(server.Notes),
			CreatedOn:           types.StringValue(server.CreatedOn),
			LastUpdated:         types.StringValue(server.LastUpdated),
			DroppedOn:           types.StringPointerValue(server.DroppedOn),
			IsActive:            types.Int64Value(server.IsActive),
			IsDeleted:           types.Int64Value(server.IsDeleted),
			Power:               types.Int64Value(server.Power),
			ProjectID:           types.Int64Value(server.ProjectID),
			IsCustom:            types.Int64Value(server.IsCustom),
			NrAddedIps:          types.Int64Value(server.NrAddedIps),
			InPcs:               types.Int64Value(server.InPcs),
			CustomPrice:         types.Int64PointerValue(server.CustomPrice),
			PayableLicense:      types.Int64Value(server.PayableLicense),
			LastLicensePay:      types.StringPointerValue(server.LastLicensePay),
			ScriptID:            types.StringPointerValue(server.ScriptID),
			SshKeyID:            types.StringPointerValue(server.SshKeyID),
			IsLocked:            types.Int64Value(server.IsLocked),
			IsWorkWithNew:       types.Int64Value(server.IsWorkWithNew),
			IsSuspended:         types.Int64Value(server.IsSuspended),
			IsTerminated:        types.Int64Value(server.IsTerminated),
			OldID:               types.Int64Value(server.OldID),
			CustomIsoID:         types.Int64PointerValue(server.CustomIsoID),
			IsIsoImageBootAble:  types.Int64Value(server.IsIsoImageBootAble),
			HasSsl:              types.Int64Value(server.HasSsl),
			LastActionDate:      types.StringPointerValue(server.LastActionDate),
			IsCreatedFromLegacy: types.Int64Value(server.IsCreatedFromLegacy),
			IsSmtpAllowed:       types.Int64Value(server.IsSmtpAllowed),
			WeeklyBackup:        types.Int64Value(server.WeeklyBackup),
			MonthlyBackup:       types.Int64Value(server.MonthlyBackup),
			LibIsoID:            types.Int64PointerValue(server.LibIsoID),
			DailySnapshot:       types.Int64Value(server.DailySnapshot),
			WeeklySnapshot:      types.Int64Value(server.WeeklySnapshot),
			MonthlySnapshot:     types.Int64Value(server.MonthlySnapshot),
			LastActionInMin:     types.Int64Value(server.LastActionInMin),
			FirstName:           types.StringValue(server.FirstName),
			LastName:            types.StringValue(server.LastName),
			Username:            types.StringValue(server.Username),
			State:               types.StringValue(server.State),
			IsFipAvailable:      types.Int64Value(server.IsFipAvailable),
			IsBucketAvailable:   types.Int64Value(server.IsBucketAvailable),
			DcIdentifier:        types.StringValue(server.DcIdentifier),
			Category:            types.StringValue(server.Category),
			FullName:            types.StringValue(server.FullName),
			VmDescription:       types.StringValue(server.VmDescription),
			BoxesSuspended:      types.Int64Value(server.BoxesSuspended),
			IsSataAvailable:     types.Int64Value(server.IsSataAvailable),
			IsSsdAvailable:      types.Int64Value(server.IsSsdAvailable),
			PublicIp:            types.StringPointerValue(server.PublicIp),
		}

		state.Servers = append(state.Servers, serverState)
	}

	state.ID = types.StringValue("servers")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *serverDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
