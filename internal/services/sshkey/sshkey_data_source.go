package sshkey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vpsie/govpsie"
)

type sshKeyDataSource struct {
	client *govpsie.Client
}

type sshKeyDataSourceModel struct {
	SshKeys []sshKeysModel `tfsdk:"sshkeys"`
	ID      types.String   `tfsdk:"id"`
}

type sshKeysModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Id         types.Int64  `tfsdk:"id"`
	UserId     types.Int64  `tfsdk:"user_id"`
	Name       types.String `tfsdk:"name"`
	PrivateKey types.String `tfsdk:"private_key"`
	CreatedOn  types.String `tfsdk:"created_on"`
	CreatedBy  types.String `tfsdk:"created_by"`
}

// NewSshKeyDataSource is a helper function to create the data source.
func NewSshKeyDataSource() datasource.DataSource {
	return &sshKeyDataSource{}
}

// Metadata returns the data source type name.
func (s *sshKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sshkeys"
}

// Schema defines the schema for the data source.
func (s *sshKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about all VPSie SSH keys.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source.",
			},
			"sshkeys": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of SSH keys.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the SSH key.",
						},
						"user_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user who owns the SSH key.",
						},
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The numeric ID of the SSH key.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the SSH key.",
						},
						"private_key": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: "The public key content of the SSH key.",
						},
						"created_on": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the SSH key was created.",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user who created the SSH key.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (s *sshKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sshKeyDataSourceModel

	sshKeys, err := s.client.SShKey.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get sshKeys",
			"An unexpected error occurred when getting sshKeys: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Got sshKeys", map[string]interface{}{"sshKeys": sshKeys})

	for _, sshKey := range sshKeys {

		sshKeyState := sshKeysModel{
			Identifier: types.StringValue(sshKey.Identifier),
			CreatedOn:  types.StringValue(sshKey.CreatedOn),
			CreatedBy:  types.StringValue(sshKey.CreatedBy),
			Name:       types.StringValue(sshKey.Name),
			PrivateKey: types.StringValue(sshKey.PrivateKey),
			UserId:     types.Int64Value(sshKey.UserId),
			Id:         types.Int64Value(sshKey.Id),
		}

		state.SshKeys = append(state.SshKeys, sshKeyState)
	}

	state.ID = types.StringValue("ssh_keys")
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *sshKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
