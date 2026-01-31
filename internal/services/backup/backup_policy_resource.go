package backup

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &backupPolicyResource{}
	_ resource.ResourceWithConfigure   = &backupPolicyResource{}
	_ resource.ResourceWithImportState = &backupPolicyResource{}
)

type backupPolicyResource struct {
	client BackupAPI
}

type backupPolicyResourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
	BackupPlan types.String `tfsdk:"backup_plan"`
	PlanEvery  types.String `tfsdk:"plan_every"`
	Keep       types.String `tfsdk:"keep"`
	CreatedOn  types.String `tfsdk:"created_on"`
	CreatedBy  types.String `tfsdk:"created_by"`
	Disabled   types.Int64  `tfsdk:"disabled"`
	Vms        types.List   `tfsdk:"vms"`
}

func NewBackupPolicyResource() resource.Resource {
	return &backupPolicyResource{}
}

func (b *backupPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_policy"
}

func (b *backupPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a backup policy on the VPSie platform.",
		Attributes: map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the backup policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the backup policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"backup_plan": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of backup plan (e.g., daily, weekly, monthly).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"plan_every": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The frequency interval for the backup plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"keep": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The number of backups to retain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"vms": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "A list of virtual machine identifiers to attach to this backup policy.",
				ElementType:         types.StringType,
			},
			"created_on": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the backup policy was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who created the backup policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Whether the backup policy is disabled (1 for disabled, 0 for enabled).",
			},
		},
	}
}

func (b *backupPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	b.client = client.Backup
}

func (b *backupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan backupPolicyResourceModel
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

	createReq := &govpsie.CreateBackupPolicyReq{
		Name:       plan.Name.ValueString(),
		BackupPlan: plan.BackupPlan.ValueString(),
		PlanEvery:  plan.PlanEvery.ValueString(),
		Keep:       plan.Keep.ValueString(),
		Vms:        vms,
	}

	err := b.client.CreateBackupPolicy(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating backup policy", err.Error())
		return
	}

	policy, err := b.GetPolicyByName(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading backup policy after creation", err.Error())
		return
	}

	plan.Identifier = types.StringValue(policy.Identifier)
	plan.CreatedOn = types.StringValue(policy.CreatedOn)
	plan.CreatedBy = types.StringValue(policy.CreatedBy)
	plan.Disabled = types.Int64Value(int64(policy.Disabled))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (b *backupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backupPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := b.client.GetBackupPolicy(ctx, state.Identifier.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading backup policy",
			"couldn't read backup policy: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(policy.Name)
	state.BackupPlan = types.StringValue(policy.BackupPlan)
	state.PlanEvery = types.StringValue(fmt.Sprintf("%d", policy.PlanEvery))
	state.Keep = types.StringValue(fmt.Sprintf("%d", policy.Keep))
	state.CreatedOn = types.StringValue(policy.CreatedOn)
	state.CreatedBy = types.StringValue(policy.CreatedBy)
	state.Disabled = types.Int64Value(int64(policy.Disabled))

	vmsList, vmsDiags := types.ListValueFrom(ctx, types.StringType, policy.Vms)
	resp.Diagnostics.Append(vmsDiags...)
	if !resp.Diagnostics.HasError() {
		state.Vms = vmsList
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (b *backupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan backupPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state backupPolicyResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle VM attachment changes
	var planVms, stateVms []string
	if !plan.Vms.IsNull() {
		diags = plan.Vms.ElementsAs(ctx, &planVms, false)
		resp.Diagnostics.Append(diags...)
	}
	if !state.Vms.IsNull() {
		diags = state.Vms.ElementsAs(ctx, &stateVms, false)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	stateSet := make(map[string]bool)
	for _, vm := range stateVms {
		stateSet[vm] = true
	}
	planSet := make(map[string]bool)
	for _, vm := range planVms {
		planSet[vm] = true
	}

	// Attach new VMs
	var toAttach []string
	for _, vm := range planVms {
		if !stateSet[vm] {
			toAttach = append(toAttach, vm)
		}
	}
	if len(toAttach) > 0 {
		err := b.client.AttachBackupPolicy(ctx, state.Identifier.ValueString(), toAttach)
		if err != nil {
			resp.Diagnostics.AddError("Error attaching VMs to backup policy", err.Error())
			return
		}
	}

	// Detach removed VMs
	var toDetach []string
	for _, vm := range stateVms {
		if !planSet[vm] {
			toDetach = append(toDetach, vm)
		}
	}
	if len(toDetach) > 0 {
		err := b.client.DetachBackupPolicy(ctx, state.Identifier.ValueString(), toDetach)
		if err != nil {
			resp.Diagnostics.AddError("Error detaching VMs from backup policy", err.Error())
			return
		}
	}

	state.Vms = plan.Vms

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (b *backupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := b.client.DeleteBackupPolicy(ctx, state.Identifier.ValueString(), state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup policy",
			"couldn't delete backup policy, unexpected error: "+err.Error(),
		)
		return
	}
}

func (b *backupPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}

func (b *backupPolicyResource) GetPolicyByName(ctx context.Context, name string) (*govpsie.BackupPolicyListDetail, error) {
	policies, err := b.client.ListBackupPolicies(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if policy.Name == name {
			return &policy, nil
		}
	}

	return nil, fmt.Errorf("backup policy with name %s not found", name)
}
