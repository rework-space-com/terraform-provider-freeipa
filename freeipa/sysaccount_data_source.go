// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SysAccountSource{}
var _ datasource.DataSourceWithConfigure = &SysAccountSource{}

func NewSysAccountDataSource() datasource.DataSource {
	return &SysAccountSource{}
}

// SysAccountSource defines the resource implementation.
type SysAccountSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type SysAccountSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	AccountLocked types.Bool   `tfsdk:"disabled"`
	MemberOfRole  types.List   `tfsdk:"memberof_role"`
}

func (r *SysAccountSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sysaccount"
}

func (r *SysAccountSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA System Account data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User ID of the system account",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Is the system account locked",
				Computed:            true,
			},
			"memberof_role": schema.ListAttribute{
				MarkdownDescription: "List of roles this system account is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *SysAccountSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*freeipaProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *freeipaProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p := &freeipaProvider{}
	client, err := p.NewFreeIPAClient(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create FreeIPA API Client",
			"An unexpected error occurred when creating the FreeIPA API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"FreeIPA Client Error: "+err.Error(),
		)
		return
	}

	r.client = client
}

func (r *SysAccountSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SysAccountSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.SysaccountShowOptionalArgs{
		All: &all,
	}

	res, err := r.client.SysaccountShow(&ipa.SysaccountShowArgs{UID: data.Name.ValueString()}, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] System account not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error ", err.Error())
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] System account show returns %v", res))
	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.Nsaccountlock != nil {
		data.AccountLocked = types.BoolValue(*res.Result.Nsaccountlock)
	}
	if res.Result.MemberofRole != nil {
		data.MemberOfRole, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofRole)
	}
	data.Id = types.StringValue(data.Name.ValueString())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
