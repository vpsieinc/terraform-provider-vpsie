package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/vpsie/govpsie"
)

type loadbalancerDataSource struct {
	client LoadbalancerAPI
}

type loadbalancerDataSourceModel struct {
	Loadbalancers []loadbalancersModel `tfsdk:"loadbalancers"`
	ID            types.String         `tfsdk:"id"`
}

type loadbalancersModel struct {
	Cpu        types.Int64  `tfsdk:"cpu"`
	Ssd        types.Int64  `tfsdk:"ssd"`
	Ram        types.Int64  `tfsdk:"ram"`
	LBName     types.String `tfsdk:"lb_name"`
	Traffic    types.Int64  `tfsdk:"traffic"`
	BoxsizeID  types.Int64  `tfsdk:"boxsize_id"`
	DefaultIP  types.String `tfsdk:"default_ip"`
	DCName     types.String `tfsdk:"dc_name"`
	Identifier types.String `tfsdk:"identifier"`
	CreatedOn  types.String `tfsdk:"created_on"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	Package    types.String `tfsdk:"package"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UserID     types.Int64  `tfsdk:"user_id"`
}

// NewLoadbalancerDataSource is a helper function to create the data source.
func NewLoadbalancerDataSource() datasource.DataSource {
	return &loadbalancerDataSource{}
}

// Metadata returns the data source type name.
func (l *loadbalancerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_loadbalancers"
}

// Schema defines the schema for the data source.
func (l *loadbalancerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about all VPSie load balancers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source.",
			},
			"loadbalancers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of load balancers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the load balancer.",
						},
						"cpu": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of CPU cores allocated to the load balancer.",
						},
						"ssd": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The SSD storage in GB allocated to the load balancer.",
						},
						"ram": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The RAM in MB allocated to the load balancer.",
						},
						"lb_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the load balancer.",
						},
						"traffic": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The traffic allowance for the load balancer in GB.",
						},
						"boxsize_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The box size ID defining the load balancer resources.",
						},
						"default_ip": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default IP address of the load balancer.",
						},
						"dc_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the data center where the load balancer resides.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the load balancer was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the load balancer was last updated.",
						},
						"package": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The package tier of the load balancer.",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user who created the load balancer.",
						},
						"user_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user who owns the load balancer.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (k *loadbalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadbalancerDataSourceModel

	loadbalancers, err := k.client.ListLBs(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Loadbalancers",
			"Could not get Loadbalancers, unexpected error: "+err.Error(),
		)

		return
	}

	for _, lb := range loadbalancers {
		lbState := loadbalancersModel{
			Identifier: types.StringValue(lb.Identifier),
			UserID:     types.Int64Value(int64(lb.UserID)),
			CreatedOn:  types.StringValue(lb.CreatedOn),
			CreatedBy:  types.StringValue(lb.CreatedBy),
			UpdatedAt:  types.StringValue(lb.UpdatedAt),
			Package:    types.StringValue(lb.Package),
			Cpu:        types.Int64Value(int64(lb.Cpu)),
			Ram:        types.Int64Value(int64(lb.Ram)),
			Ssd:        types.Int64Value(int64(lb.Ssd)),
			LBName:     types.StringValue(lb.LBName),
			Traffic:    types.Int64Value(int64(lb.Traffic)),
			BoxsizeID:  types.Int64Value(int64(lb.BoxsizeID)),
			DefaultIP:  types.StringValue(lb.DefaultIP),
			DCName:     types.StringValue(lb.DCName),
		}

		state.Loadbalancers = append(state.Loadbalancers, lbState)
	}

	state.ID = types.StringValue("loadbalancers")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
func (l *loadbalancerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	l.client = client.LB
}
