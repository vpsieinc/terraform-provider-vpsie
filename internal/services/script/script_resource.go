package script

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vpsie/govpsie"
)

var (
	_ resource.Resource                = &scriptResource{}
	_ resource.ResourceWithConfigure   = &scriptResource{}
	_ resource.ResourceWithImportState = &scriptResource{}
)

type scriptResource struct {
	client ScriptAPI
}

type scriptResourceModel struct {
	Identifier    types.String `tfsdk:"identifier"`
	UserID        types.Int64  `tfsdk:"user_id"`
	BoxID         types.Int64  `tfsdk:"box_id"`
	BoxIdentifier types.String `tfsdk:"box_identifier"`
	ScriptName    types.String `tfsdk:"script_name"`
	Script        types.String `tfsdk:"script"`
	CreatedOn     types.String `tfsdk:"created_on"`
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
}

func NewScriptResource() resource.Resource {
	return &scriptResource{}
}

func (s *scriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

func (s *scriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the user who owns the script.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"script_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the script.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"script": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The content of the script to execute.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"created_on": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the script was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"box_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the box associated with the script.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"box_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the box associated with the script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the script.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The resolved name of the script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the script (e.g., bash, cloud-init).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (s *scriptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	s.client = client.Scripts
}

// Create creates the resource and sets the initial Terraform state.
func (s *scriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createScript *govpsie.CreateScriptRequest = &govpsie.CreateScriptRequest{
		Name:          plan.ScriptName.ValueString(),
		ScriptContent: plan.Script.ValueString(),
		ScriptType:    plan.Type.ValueString(),
		Tags:          []string{},
	}
	err := s.client.CreateScript(ctx, createScript)
	if err != nil {
		resp.Diagnostics.AddError("Error creating script", err.Error())
		return
	}

	script, err := s.GetScriptByName(ctx, plan.ScriptName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting script", err.Error())
		return
	}

	// Overwrite items with refreshed state

	plan.Identifier = types.StringValue(script.Identifier)
	plan.CreatedOn = types.StringValue(script.CreatedOn.String())
	plan.BoxID = types.Int64Value(int64(script.BoxID))
	plan.BoxIdentifier = types.StringValue(script.BoxIdentifier)
	plan.ID = types.Int64Value(int64(script.ID))
	plan.Name = types.StringValue(script.Name)
	plan.Type = types.StringValue(script.Type)
	plan.UserID = types.Int64Value(int64(script.UserID))
	plan.Script = types.StringValue(script.Script)
	plan.ScriptName = types.StringValue(script.ScriptName)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *scriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	script, err := s.client.GetScript(ctx, state.Identifier.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading vpsie scripts",
			"couldn't read vpsie scripts identifier "+state.Identifier.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state

	state.Name = types.StringValue(script.Name)
	state.Type = types.StringValue(script.Type)
	state.CreatedOn = types.StringValue(script.CreatedOn.String())
	state.UserID = types.Int64Value(int64(script.UserID))
	state.BoxID = types.Int64Value(int64(script.BoxID))
	state.BoxID = types.Int64Value(int64(script.BoxID))
	state.BoxIdentifier = types.StringValue(script.BoxIdentifier)
	state.BoxIdentifier = types.StringValue(script.BoxIdentifier)
	state.ID = types.Int64Value(int64(script.ID))
	state.Script = types.StringValue(script.Script)
	state.ScriptName = types.StringValue(script.ScriptName)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (s *scriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateScript *govpsie.ScriptUpdateRequest = &govpsie.ScriptUpdateRequest{
		Name:             plan.ScriptName.ValueString(),
		ScriptContent:    plan.Script.ValueString(),
		ScriptType:       plan.Type.ValueString(),
		ScriptIdentifier: plan.Identifier.ValueString(),
	}

	err := s.client.UpdateScript(ctx, updateScript)
	if err != nil {
		resp.Diagnostics.AddError("Error updating script", err.Error())
		return
	}

	script, err := s.client.GetScript(ctx, plan.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vpsie script",
			"Could not read Vpsie script Identifier "+plan.Identifier.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(script.Name)
	plan.Type = types.StringValue(script.Type)
	plan.Script = types.StringValue(script.Script)
	plan.ScriptName = types.StringValue(script.ScriptName)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *scriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.DeleteScript(ctx, state.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting script",
			"couldn't delete script, unexpected error: "+err.Error(),
		)

		return
	}
}

func (s *scriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("identifier"), req, resp)
}

func (s *scriptResource) GetScriptByName(ctx context.Context, scriptName string) (*govpsie.ScriptDetail, error) {
	scripts, err := s.client.GetScripts(ctx)
	if err != nil {
		return nil, err
	}

	for _, script := range scripts {
		if scriptName == strings.Split(script.ScriptName, ".")[0] {
			script, err := s.client.GetScript(ctx, script.Identifier)
			if err != nil {
				return nil, err
			}

			return &script, nil
		}
	}

	return nil, fmt.Errorf("script with name %s not found", scriptName)
}
