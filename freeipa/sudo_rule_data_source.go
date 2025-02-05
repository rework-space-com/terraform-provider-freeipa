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
var _ datasource.DataSource = &SudoRuleDataSource{}
var _ datasource.DataSourceWithConfigure = &SudoRuleDataSource{}

func NewSudoRuleDataSource() datasource.DataSource {
	return &SudoRuleDataSource{}
}

// SudoRuleDataSource defines the resource implementation.
type SudoRuleDataSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type SudoRuleDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	UserCategory            types.String `tfsdk:"usercategory"`
	HostCategory            types.String `tfsdk:"hostcategory"`
	CommandCategory         types.String `tfsdk:"commandcategory"`
	RunAsUserCategory       types.String `tfsdk:"runasusercategory"`
	RunAsGroupCategory      types.String `tfsdk:"runasgroupcategory"`
	Order                   types.Int32  `tfsdk:"order"`
	Option                  types.List   `tfsdk:"option"`
	MemberUser              types.List   `tfsdk:"member_user"`
	MemberGroup             types.List   `tfsdk:"member_group"`
	MemberHost              types.List   `tfsdk:"member_host"`
	MemberHostGroup         types.List   `tfsdk:"member_hostgroup"`
	MemberAllowSudoCmd      types.List   `tfsdk:"member_allow_sudo_cmd"`
	MemberAllowSudoCmdGroup types.List   `tfsdk:"member_allow_sudo_cmdgroup"`
	MemberDenySudoCmd       types.List   `tfsdk:"member_deny_sudo_cmd"`
	MemberDenySudoCmdGroup  types.List   `tfsdk:"member_deny_sudo_cmdgroup"`
	RunAsUser               types.List   `tfsdk:"runasuser"`
	RunAsGroup              types.List   `tfsdk:"runasgroup"`
}

func (r *SudoRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule"
}

func (r *SudoRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User sudo rule data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the sudo rule",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the sudo rule",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this sudo rule",
				Computed:            true,
			},
			"usercategory": schema.StringAttribute{
				MarkdownDescription: "User category the sudo rule is applied to (allowed value: all)",
				Computed:            true,
			},
			"hostcategory": schema.StringAttribute{
				MarkdownDescription: "Host category the sudo rule is applied to (allowed value: all)",
				Computed:            true,
			},
			"commandcategory": schema.StringAttribute{
				MarkdownDescription: "Command category the sudo rule is applied to (allowed value: all)",
				Computed:            true,
			},
			"runasusercategory": schema.StringAttribute{
				MarkdownDescription: "Run as user category the sudo rule is applied to (allowed value: all)",
				Computed:            true,
			},
			"runasgroupcategory": schema.StringAttribute{
				MarkdownDescription: "Run as group category the sudo rule is applied to (allowed value: all)",
				Computed:            true,
			},
			"order": schema.Int32Attribute{
				MarkdownDescription: "Sudo rule order (must be unique)",
				Computed:            true,
			},
			"option": schema.ListAttribute{
				MarkdownDescription: "List of options defined for this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_user": schema.ListAttribute{
				MarkdownDescription: "List of users member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_group": schema.ListAttribute{
				MarkdownDescription: "List of user groups member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_host": schema.ListAttribute{
				MarkdownDescription: "List of hosts member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of host groups member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_allow_sudo_cmd": schema.ListAttribute{
				MarkdownDescription: "List of allowed sudo commands member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_allow_sudo_cmdgroup": schema.ListAttribute{
				MarkdownDescription: "List of allowed sudo command groups member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_deny_sudo_cmd": schema.ListAttribute{
				MarkdownDescription: "List of denied sudo commands member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"member_deny_sudo_cmdgroup": schema.ListAttribute{
				MarkdownDescription: "List of denied sudo command groups member of this sudo rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"runasuser": schema.ListAttribute{
				MarkdownDescription: "List of users authorised to be run as.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"runasgroup": schema.ListAttribute{
				MarkdownDescription: "List of groups authorised to be run as.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *SudoRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *SudoRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SudoRuleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.SudoruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudoruleShowArgs{
		Cn: data.Name.ValueString(),
	}

	res, err := r.client.SudoruleShow(&args, &optArgs)
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
	if res.Result.Cmdcategory != nil {
		data.CommandCategory = types.StringValue(*res.Result.Cmdcategory)
	}
	if res.Result.Ipasudorunasusercategory != nil {
		data.RunAsUserCategory = types.StringValue(*res.Result.Ipasudorunasusercategory)
	}
	if res.Result.Ipasudorunasgroupcategory != nil {
		data.RunAsGroupCategory = types.StringValue(*res.Result.Ipasudorunasgroupcategory)
	}
	if res.Result.Sudoorder != nil {
		data.Order = types.Int32Value(int32(*res.Result.Sudoorder))
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
	if res.Result.MemberallowcmdSudocmd != nil {
		data.MemberAllowSudoCmd, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberallowcmdSudocmd)
	}
	if res.Result.MemberallowcmdSudocmdgroup != nil {
		data.MemberAllowSudoCmdGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberallowcmdSudocmdgroup)
	}
	if res.Result.MemberdenycmdSudocmd != nil {
		data.MemberDenySudoCmd, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberdenycmdSudocmd)
	}
	if res.Result.MemberdenycmdSudocmdgroup != nil {
		data.MemberDenySudoCmdGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberdenycmdSudocmdgroup)
	}
	if res.Result.IpasudorunasUser != nil {
		data.RunAsUser, _ = types.ListValueFrom(ctx, types.StringType, res.Result.IpasudorunasUser)
	}
	if res.Result.IpasudorunasgroupGroup != nil {
		data.RunAsGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.IpasudorunasgroupGroup)
	}
	if res.Result.Ipasudoopt != nil {
		data.Option, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ipasudoopt)
	}

	data.Id = data.Name
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo rule %s", res.Result.Cn))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
