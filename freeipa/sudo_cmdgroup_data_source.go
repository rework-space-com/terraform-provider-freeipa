// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SudoCmdGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &SudoCmdGroupDataSource{}

func NewSudoCmdGroupDataSource() datasource.DataSource {
	return &SudoCmdGroupDataSource{}
}

// SudoCmdGroupDataSource defines the resource implementation.
type SudoCmdGroupDataSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type SudoCmdGroupDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	MemberSudocmd types.List   `tfsdk:"member_sudocmd"`
}

func (r *SudoCmdGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_cmdgroup"
}

func (r *SudoCmdGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User sudo command group data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the sudo command group",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the sudo command group",
				Computed:            true,
			},
			"member_sudocmd": schema.ListAttribute{
				MarkdownDescription: "List of sudo commands that are member of the sudo command group",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *SudoCmdGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ipa.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SudoCmdGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SudoCmdGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	args := ipa.SudocmdgroupShowArgs{
		Cn: data.Name.ValueString(),
	}
	optArgs := ipa.SudocmdgroupShowOptionalArgs{
		All: &all,
	}

	res, err := r.client.SudocmdgroupShow(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa sudo command group %s", data.Name.ValueString()))
		return
	}

	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.MemberSudocmd != nil {
		data.MemberSudocmd, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberSudocmd)
	}
	data.Id = types.StringValue(res.Result.Cn)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group %s", res.Result.Cn))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
