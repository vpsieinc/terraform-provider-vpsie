package monitoring

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

var (
	_ resource.Resource                = &monitoringRuleResource{}
	_ resource.ResourceWithConfigure   = &monitoringRuleResource{}
	_ resource.ResourceWithImportState = &monitoringRuleResource{}
)

type monitoringRuleResource struct {
	client MonitoringAPI
}

type monitoringRuleResourceModel struct {
	Identifier    types.String `tfsdk:"identifier"`
	RuleName      types.String `tfsdk:"rule_name"`
	MetricType    types.String `tfsdk:"metric_type"`
	Condition     types.String `tfsdk:"condition"`
	ThresholdType types.String `tfsdk:"threshold_type"`
	ThresholdId   types.String `tfsdk:"threshold_id"`
	Threshold     types.String `tfsdk:"threshold"`
	Period        types.String `tfsdk:"period"`
	Frequency     types.String `tfsdk:"frequency"`
	Status        types.String `tfsdk:"status"`
	Email         types.String `tfsdk:"email"`
	ActionKey     types.String `tfsdk:"action_key"`
	ActionName    types.String `tfsdk:"action_name"`
	Vms           types.List   `tfsdk:"vms"`
	CreatedOn     types.String `tfsdk:"created_on"`
}

func NewMonitoringRuleResource() resource.Resource {
	return &monitoringRuleResource{}
}

func (m *monitoringRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_rule"
}

func (m *monitoringRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the monitoring rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the monitoring rule. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"metric_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of metric to monitor (e.g., cpu, memory, disk). Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"condition": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The condition operator for the monitoring rule. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("above", "below"),
				},
			},
			"threshold_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of threshold for the rule (e.g., percentage, absolute). Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"threshold_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The identifier of the threshold configuration. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"threshold": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The threshold value that triggers the alert. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"period": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The evaluation period for the monitoring rule in minutes. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"frequency": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The frequency of checks for the monitoring rule in minutes. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The current status of the monitoring rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address to receive alert notifications. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"action_key": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The key identifying the action to take when the rule triggers. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The display name of the action to take when the rule triggers. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vms": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "The list of VM identifiers to which this monitoring rule applies.",
				ElementType:         types.StringType,
			},
			"created_on": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the monitoring rule was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (m *monitoringRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	m.client = client.Monitoring
}

func (m *monitoringRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan monitoringRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vms []string
	if !plan.Vms.IsNull() {
		diags = plan.Vms.ElementsAs(ctx, &vms, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createReq := &govpsie.CreateMonitoringRuleReq{
		MetricType:    plan.MetricType.ValueString(),
		RuleName:      plan.RuleName.ValueString(),
		Condition:     plan.Condition.ValueString(),
		ThresholdType: plan.ThresholdType.ValueString(),
		ThresholdId:   plan.ThresholdId.ValueString(),
		Period:        plan.Period.ValueString(),
		Frequency:     plan.Frequency.ValueString(),
		Status:        plan.Status.ValueString(),
		Threshold:     plan.Threshold.ValueString(),
		Vms:           vms,
	}
	createReq.Actions.Email = plan.Email.ValueString()
	createReq.Actions.ActionKey = plan.ActionKey.ValueString()
	createReq.Actions.ActionName = plan.ActionName.ValueString()

	err := m.client.CreateRule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating monitoring rule", err.Error())
		return
	}

	rule, err := m.GetRuleByName(ctx, plan.RuleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading monitoring rule after creation", err.Error())
		return
	}

	plan.Identifier = types.StringValue(rule.Identifier)
	plan.CreatedOn = types.StringValue(rule.CreatedOn)
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		plan.Status = types.StringValue(fmt.Sprintf("%d", rule.Status))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (m *monitoringRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitoringRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := m.client.ListMonitoringRule(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading monitoring rules", err.Error())
		return
	}

	found := false
	for _, rule := range rules {
		if rule.Identifier == state.Identifier.ValueString() {
			state.RuleName = types.StringValue(rule.RuleName)
			state.MetricType = types.StringValue(rule.MetricType)
			state.Condition = types.StringValue(rule.Condition)
			state.ThresholdType = types.StringValue(rule.ThresholdType)
			state.Threshold = types.StringValue(fmt.Sprintf("%d", rule.Threshold))
			state.Period = types.StringValue(fmt.Sprintf("%d", rule.Period))
			state.Frequency = types.StringValue(fmt.Sprintf("%d", rule.Frequency))
			state.Status = types.StringValue(fmt.Sprintf("%d", rule.Status))
			state.Email = types.StringValue(rule.Email)
			state.CreatedOn = types.StringValue(rule.CreatedOn)
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

func (m *monitoringRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All Required fields are ForceNew — only status toggle is possible
	var plan monitoringRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state monitoringRuleResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Status.Equal(state.Status) {
		err := m.client.ToggleMonitoringRuleStatus(ctx, plan.Status.ValueString(), state.Identifier.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error toggling monitoring rule status", err.Error())
			return
		}
		state.Status = plan.Status
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (m *monitoringRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}

func (m *monitoringRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitoringRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.DeleteMonitoringRule(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting monitoring rule",
			"couldn't delete monitoring rule, unexpected error: "+err.Error(),
		)
		return
	}
}

func (m *monitoringRuleResource) GetRuleByName(ctx context.Context, name string) (*govpsie.MonitoringRule, error) {
	rules, err := m.client.ListMonitoringRule(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.RuleName == name {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("monitoring rule with name %s not found", name)
}
