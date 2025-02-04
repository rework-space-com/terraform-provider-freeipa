// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &UserGroupDataSource{}

func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

// UserGroupDataSource defines the resource implementation.
type UserGroupDataSource struct {
	client *ipa.Client
}

// UserGroupResourceModel describes the resource data model.
type UserGroupDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	GidNumber                types.Int64  `tfsdk:"gid_number"`
	Ipaexternalmember        types.List   `tfsdk:"member_external"`
	MemberUser               types.List   `tfsdk:"member_user"`
	MemberGroup              types.List   `tfsdk:"member_group"`
	MemberOfGroup            types.List   `tfsdk:"memberof_group"`
	MemberOfSudoRule         types.List   `tfsdk:"memberof_sudorule"`
	MemberOfHBACRule         types.List   `tfsdk:"memberof_hbacrule"`
	MemberIndirectUser       types.List   `tfsdk:"member_indirect_user"`
	MemberIndirectGroup      types.List   `tfsdk:"member_indirect_group"`
	MemberOfIndirectGroup    types.List   `tfsdk:"memberof_indirect_group"`
	MemberOfIndirectSudoRule types.List   `tfsdk:"memberof_indirect_sudorule"`
	MemberOfIndirectHBACRule types.List   `tfsdk:"memberof_indirect_hbacrule"`
}

func (r *UserGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *UserGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User Group data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Group name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Group Description",
				Computed:            true,
			},
			"gid_number": schema.Int64Attribute{
				MarkdownDescription: "GID (use this option to set it manually)",
				Computed:            true,
			},
			"member_external": schema.ListAttribute{
				MarkdownDescription: "List of external users (from trusted domain) that are member of this group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_user": schema.ListAttribute{
				MarkdownDescription: "List of users that are member of this group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_group": schema.ListAttribute{
				MarkdownDescription: "List of groups that are member of this group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_group": schema.ListAttribute{
				MarkdownDescription: "List of groups this group is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this group is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this group is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_indirect_user": schema.ListAttribute{
				MarkdownDescription: "List of users that are is indirectly member of this group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_indirect_group": schema.ListAttribute{
				MarkdownDescription: "List of groups that are is indirectly member of this group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_group": schema.ListAttribute{
				MarkdownDescription: "List of groups this group is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this group is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this group is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *UserGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.GroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.GroupShowArgs{
		Cn: data.Name.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %s", data.Name.ValueString()))
	res, err := r.client.GroupShow(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading group %s: %s", data.Name.ValueString(), err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %v", res))
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa group %s", data.Name.ValueString()))
		return
	}

	data.Name = types.StringValue(res.Result.Cn)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Cn %s", data.Name.ValueString()))
	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Description %s", data.Description.ValueString()))
	}
	if res.Result.Gidnumber != nil {
		data.GidNumber = types.Int64Value(int64(*res.Result.Gidnumber))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group GID %d", data.GidNumber.ValueInt64()))
	}

	if res.Result.Ipaexternalmember != nil {
		var diag diag.Diagnostics
		data.Ipaexternalmember, diag = types.ListValueFrom(ctx, types.StringType, res.Result.Ipaexternalmember)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberUser != nil {
		var diag diag.Diagnostics
		data.MemberUser, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberUser)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberGroup != nil {
		var diag diag.Diagnostics
		data.MemberGroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberGroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofGroup != nil {
		var diag diag.Diagnostics
		data.MemberOfGroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofGroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofSudorule != nil {
		var diag diag.Diagnostics
		data.MemberOfSudoRule, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofSudorule)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofHbacrule != nil {
		var diag diag.Diagnostics
		data.MemberOfHBACRule, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofHbacrule)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberindirectUser != nil {
		var diag diag.Diagnostics
		data.MemberIndirectUser, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberindirectUser)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberindirectGroup != nil {
		var diag diag.Diagnostics
		data.MemberIndirectGroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberindirectGroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofindirectGroup != nil {
		var diag diag.Diagnostics
		data.MemberOfIndirectGroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectGroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofindirectSudorule != nil {
		var diag diag.Diagnostics
		data.MemberOfIndirectSudoRule, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectSudorule)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofindirectHbacrule != nil {
		var diag diag.Diagnostics
		data.MemberOfIndirectHBACRule, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectHbacrule)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	data.Id = types.StringValue(data.Name.ValueString())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
