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
var _ datasource.DataSource = &HbacPolicyDataSource{}
var _ datasource.DataSourceWithConfigure = &HbacPolicyDataSource{}

func NewHbacPolicyDataSource() datasource.DataSource {
	return &HbacPolicyDataSource{}
}

// HbacPolicyDataSource defines the resource implementation.
type HbacPolicyDataSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type HbacPolicyDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	UserCategory       types.String `tfsdk:"usercategory"`
	HostCategory       types.String `tfsdk:"hostcategory"`
	ServiceCategory    types.String `tfsdk:"servicecategory"`
	MemberUser         types.List   `tfsdk:"member_user"`
	MemberGroup        types.List   `tfsdk:"member_group"`
	MemberHost         types.List   `tfsdk:"member_host"`
	MemberHostGroup    types.List   `tfsdk:"member_hostgroup"`
	MemberService      types.List   `tfsdk:"member_service"`
	MemberServiceGroup types.List   `tfsdk:"member_servicegroup"`
}

func (r *HbacPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hbac_policy"
}

func (r *HbacPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User hbac policy data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the hbac policy",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the hbac policy",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this hbac policy",
				Computed:            true,
			},
			"usercategory": schema.StringAttribute{
				MarkdownDescription: "User category the hbac policy is applied to (allowed value: all)",
				Computed:            true,
			},
			"hostcategory": schema.StringAttribute{
				MarkdownDescription: "Host category the hbac policy is applied to (allowed value: all)",
				Computed:            true,
			},
			"servicecategory": schema.StringAttribute{
				MarkdownDescription: "Command category the hbac policy is applied to (allowed value: all)",
				Computed:            true,
			},
			"member_user": schema.ListAttribute{
				MarkdownDescription: "List of users member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_group": schema.ListAttribute{
				MarkdownDescription: "List of user groups member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_host": schema.ListAttribute{
				MarkdownDescription: "List of hosts member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of host groups member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_service": schema.ListAttribute{
				MarkdownDescription: "List of services member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_servicegroup": schema.ListAttribute{
				MarkdownDescription: "List of service groups member of this hbac policy.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *HbacPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *HbacPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HbacPolicyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.HbacruleShowArgs{
		Cn: data.Name.ValueString(),
	}

	res, err := r.client.HbacruleShow(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.Ipaenabledflag != nil {
		data.Enabled = types.BoolValue(*res.Result.Ipaenabledflag)
	}
	if res.Result.Usercategory != nil {
		data.UserCategory = types.StringValue(*res.Result.Usercategory)
	}
	if res.Result.Hostcategory != nil {
		data.HostCategory = types.StringValue(*res.Result.Hostcategory)
	}
	if res.Result.Servicecategory != nil {
		data.ServiceCategory = types.StringValue(*res.Result.Servicecategory)
	}
	if res.Result.MemberuserUser != nil {
		data.MemberUser, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberuserUser)
	}
	if res.Result.MemberuserGroup != nil {
		data.MemberGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberuserGroup)
	}
	if res.Result.MemberhostHost != nil {
		data.MemberHost, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberhostHost)
	}
	if res.Result.MemberhostHostgroup != nil {
		data.MemberHostGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberhostHostgroup)
	}
	if res.Result.MemberserviceHbacsvc != nil {
		data.MemberService, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberserviceHbacsvc)
	}
	if res.Result.MemberserviceHbacsvcgroup != nil {
		data.MemberServiceGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberserviceHbacsvcgroup)
	}

	data.Id = data.Name
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy %s", res.Result.Cn))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
