package fip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

type fipDataSource struct {
	client *govpsie.Client
}

type fipDataSourceModel struct {
	ID  types.String `tfsdk:"id"`
	IPs []ipModel    `tfsdk:"ips"`
}

type ipModel struct {
	ID            types.Int64  `tfsdk:"id"`
	DcName        types.String `tfsdk:"dc_name"`
	DcIdentifier  types.String `tfsdk:"dc_identifier"`
	IP            types.String `tfsdk:"ip"`
	IPVersion     types.String `tfsdk:"ip_version"`
	IsPrimary     types.Int64  `tfsdk:"is_primary"`
	Hostname      types.String `tfsdk:"hostname"`
	BoxIdentifier types.String `tfsdk:"box_identifier"`
	FullName      types.String `tfsdk:"full_name"`
	Category      types.String `tfsdk:"category"`
	Type          types.String `tfsdk:"type"`
}

func NewFipDataSource() datasource.DataSource {
	return &fipDataSource{}
}

func (d *fipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_floating_ips"
}

func (d *fipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of all floating IPs on the VPSie platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier for this data source.",
			},
			"ips": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of floating IPs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The unique numeric identifier of the floating IP.",
						},
						"dc_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the data center where the IP is allocated.",
						},
						"dc_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the data center where the IP is allocated.",
						},
						"ip": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The floating IP address.",
						},
						"ip_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The IP version of the floating IP address.",
						},
						"is_primary": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether this is the primary IP for the associated server.",
						},
						"hostname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname of the server associated with the floating IP.",
						},
						"box_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the server (box) associated with the floating IP.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the user who owns the floating IP.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The category of the floating IP.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the floating IP.",
						},
					},
				},
			},
		},
	}
}

func (d *fipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *fipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state fipDataSourceModel

	ips, err := d.client.IP.ListPublicIPs(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Floating IPs",
			"Could not get floating IPs, unexpected error: "+err.Error(),
		)
		return
	}

	for _, ip := range ips {
		state.IPs = append(state.IPs, ipModel{
			ID:            types.Int64Value(int64(ip.ID)),
			DcName:        types.StringValue(ip.DcName),
			DcIdentifier:  types.StringValue(ip.DcIdentifier),
			IP:            types.StringValue(ip.IP),
			IPVersion:     types.StringValue(ip.IPVersion),
			IsPrimary:     types.Int64Value(int64(ip.IsPrimary)),
			Hostname:      types.StringValue(ip.Hostname),
			BoxIdentifier: types.StringValue(ip.BoxIdentifier),
			FullName:      types.StringValue(ip.FullName),
			Category:      types.StringValue(ip.Category),
			Type:          types.StringValue(ip.Type),
		})
	}

	state.ID = types.StringValue("floating_ips")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
