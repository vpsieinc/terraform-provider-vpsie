package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

var (
	_ resource.Resource              = &reverseDnsResource{}
	_ resource.ResourceWithConfigure = &reverseDnsResource{}
)

type reverseDnsResource struct {
	client DomainAPI
}

type reverseDnsResourceModel struct {
	ID               types.String `tfsdk:"id"`
	VmIdentifier     types.String `tfsdk:"vm_identifier"`
	IP               types.String `tfsdk:"ip"`
	DomainIdentifier types.String `tfsdk:"domain_identifier"`
	HostName         types.String `tfsdk:"hostname"`
}

func NewReverseDnsResource() resource.Resource {
	return &reverseDnsResource{}
}

func (r *reverseDnsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reverse_dns"
}

func (r *reverseDnsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a reverse DNS (PTR) record on the VPSie platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The composite identifier of the reverse DNS record (vm_identifier/ip).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vm_identifier": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The identifier of the virtual machine associated with this reverse DNS record. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The IP address for the reverse DNS record. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"domain_identifier": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The identifier of the domain for this reverse DNS record.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"hostname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The hostname that the IP address resolves to.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *reverseDnsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*govpsie.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configuration Type",
			fmt.Sprintf("Expected *govpsie.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client.Domain
}

func (r *reverseDnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan reverseDnsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reverseReq := &govpsie.ReverseRequest{
		VmIdentifier:     plan.VmIdentifier.ValueString(),
		Ip:               plan.IP.ValueString(),
		DomainIdentifier: plan.DomainIdentifier.ValueString(),
		HostName:         plan.HostName.ValueString(),
	}

	err := r.client.AddReverse(ctx, reverseReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating reverse DNS", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.VmIdentifier.ValueString(), plan.IP.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *reverseDnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state reverseDnsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	records, err := r.client.ListReversePTRRecords(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading reverse DNS records", err.Error())
		return
	}

	found := false
	for _, record := range records {
		if record.Ip == state.IP.ValueString() && record.VmIdentifier == state.VmIdentifier.ValueString() {
			state.HostName = types.StringValue(record.HostName)
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *reverseDnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan reverseDnsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reverseReq := &govpsie.ReverseRequest{
		VmIdentifier:     plan.VmIdentifier.ValueString(),
		Ip:               plan.IP.ValueString(),
		DomainIdentifier: plan.DomainIdentifier.ValueString(),
		HostName:         plan.HostName.ValueString(),
	}

	err := r.client.UpdateReverse(ctx, reverseReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating reverse DNS", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *reverseDnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state reverseDnsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteReverse(ctx, state.IP.ValueString(), state.VmIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting reverse DNS",
			"couldn't delete reverse DNS, unexpected error: "+err.Error(),
		)
		return
	}
}
