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
var _ datasource.DataSource = &HostGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &HostGroupDataSource{}

func NewHostGroupDataSource() datasource.DataSource {
	return &HostGroupDataSource{}
}

// HostGroupDataSource defines the resource implementation.
type HostGroupDataSource struct {
	client *ipa.Client
}

// UserGroupResourceModel describes the resource data model.
type HostGroupDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	MemberHost                types.List   `tfsdk:"member_host"`
	MemberHostgroup           types.List   `tfsdk:"member_hostgroup"`
	MemberOfHostgroup         types.List   `tfsdk:"memberof_hostgroup"`
	MemberOfSudoRule          types.List   `tfsdk:"memberof_sudorule"`
	MemberOfHBACRule          types.List   `tfsdk:"memberof_hbacrule"`
	MemberIndirectHost        types.List   `tfsdk:"member_indirect_host"`
	MemberIndirectHostgroup   types.List   `tfsdk:"member_indirect_hostgroup"`
	MemberOfIndirectHostgroup types.List   `tfsdk:"memberof_indirect_hostgroup"`
	MemberOfIndirectSudoRule  types.List   `tfsdk:"memberof_indirect_sudorule"`
	MemberOfIndirectHBACRule  types.List   `tfsdk:"memberof_indirect_hbacrule"`
}

func (r *HostGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hostgroup"
}

func (r *HostGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User Group data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Hostgroup name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Hostgroup Description",
				Computed:            true,
			},
			"member_host": schema.ListAttribute{
				MarkdownDescription: "List of hosts that are member of this hostgroup.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups that are member of this hostgroup.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups this hostgroup is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this hostgroup is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this hostgroup is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_indirect_host": schema.ListAttribute{
				MarkdownDescription: "List of hosts that are is indirectly member of this hostgroup.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_indirect_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups that are is indirectly member of this hostgroup.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups this hostgroup is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this hostgroup is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this hostgroup is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *HostGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *HostGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.HostgroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.HostgroupShowArgs{
		Cn: data.Name.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup %s", data.Name.ValueString()))
	res, err := r.client.HostgroupShow(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa hostgroup %s: %s", data.Name.ValueString(), err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup %v", res))
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa hostgroup %s", data.Name.ValueString()))
		return
	}

	data.Name = types.StringValue(res.Result.Cn)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup Cn %s", data.Name.ValueString()))
	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Description %s", data.Description.ValueString()))
	}
	if res.Result.MemberHost != nil {
		var diag diag.Diagnostics
		data.MemberHost, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberHost)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberHostgroup != nil {
		var diag diag.Diagnostics
		data.MemberHostgroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberHostgroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofHostgroup != nil {
		var diag diag.Diagnostics
		data.MemberOfHostgroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofHostgroup)
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

	if res.Result.MemberindirectHost != nil {
		var diag diag.Diagnostics
		data.MemberIndirectHost, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberindirectHost)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberindirectHostgroup != nil {
		var diag diag.Diagnostics
		data.MemberIndirectHostgroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberindirectHostgroup)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if res.Result.MemberofindirectHostgroup != nil {
		var diag diag.Diagnostics
		data.MemberOfIndirectHostgroup, diag = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectHostgroup)
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
