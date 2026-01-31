package ip

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

type ipDataSource struct {
	client IPAPI
}

type ipDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	IPs  []ipModel    `tfsdk:"ips"`
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
	CreatedBy     types.String `tfsdk:"created_by"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func NewIPDataSource() datasource.DataSource {
	return &ipDataSource{}
}

func (d *ipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ips"
}

func (d *ipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of IP addresses on the VPSie platform, optionally filtered by type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier for this data source.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by IP type: `all`, `public`, or `private`. Defaults to `all`.",
			},
			"ips": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of IP addresses.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The unique numeric identifier of the IP address.",
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
							MarkdownDescription: "The IP address.",
						},
						"ip_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The IP version (e.g., IPv4, IPv6).",
						},
						"is_primary": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Whether this is the primary IP for the associated server.",
						},
						"hostname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hostname of the server associated with the IP.",
						},
						"box_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the server (box) associated with the IP.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the user who owns the IP.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The category of the IP address.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the IP address (e.g., public, private).",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user who created the IP address.",
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the IP address was last updated.",
						},
					},
				},
			},
		},
	}
}

func (d *ipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client.IP
}

func (d *ipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ipDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipType := "all"
	if !config.Type.IsNull() && !config.Type.IsUnknown() {
		ipType = config.Type.ValueString()
	}

	var ips []govpsie.IP
	var err error

	switch ipType {
	case "public":
		ips, err = d.client.ListPublicIPs(ctx, nil)
	case "private":
		ips, err = d.client.ListPrivateIPs(ctx, nil)
	default:
		ips, err = d.client.ListAllIPs(ctx, nil)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting IPs",
			"Could not get IPs, unexpected error: "+err.Error(),
		)
		return
	}

	var state ipDataSourceModel
	state.Type = config.Type

	for _, ipItem := range ips {
		state.IPs = append(state.IPs, ipModel{
			ID:            types.Int64Value(int64(ipItem.ID)),
			DcName:        types.StringValue(ipItem.DcName),
			DcIdentifier:  types.StringValue(ipItem.DcIdentifier),
			IP:            types.StringValue(ipItem.IP),
			IPVersion:     types.StringValue(ipItem.IPVersion),
			IsPrimary:     types.Int64Value(int64(ipItem.IsPrimary)),
			Hostname:      types.StringValue(ipItem.Hostname),
			BoxIdentifier: types.StringValue(ipItem.BoxIdentifier),
			FullName:      types.StringValue(ipItem.FullName),
			Category:      types.StringValue(ipItem.Category),
			Type:          types.StringValue(ipItem.Type),
			CreatedBy:     types.StringValue(ipItem.CreatedBy),
			UpdatedAt:     types.StringValue(ipItem.UpdatedAt),
		})
	}

	state.ID = types.StringValue("ips")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
