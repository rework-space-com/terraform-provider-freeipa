// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SudoCmdGroupResource{}
var _ resource.ResourceWithImportState = &SudoCmdGroupResource{}

func NewSudoCmdGroupResource() resource.Resource {
	return &SudoCmdGroupResource{}
}

// SudoCmdGroupResource defines the resource implementation.
type SudoCmdGroupResource struct {
	client *ipa.Client
}

// SudoCmdGroupResourceModel describes the resource data model.
type SudoCmdGroupResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *SudoCmdGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_cmdgroup"
}

func (r *SudoCmdGroupResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *SudoCmdGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo command group resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the sudo command group",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Sudo command group description",
				Optional:            true,
			},
		},
	}
}

func (r *SudoCmdGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoCmdGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoCmdGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudocmdgroupAddOptionalArgs{}

	args := ipa.SudocmdgroupAddArgs{
		Cn: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	_, err := r.client.SudocmdgroupAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo command group: %s", err))
	}
	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoCmdGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoCmdGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] Sudo command group not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa sudo command group: %s", err))
			return
		}
	}
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa sudo command group %s", data.Name.ValueString()))
		return
	}

	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
	}
	data.Id = data.Name
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group %s", res.Result.Cn))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SudoCmdGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoCmdGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudocmdgroupModOptionalArgs{}

	args := ipa.SudocmdgroupModArgs{
		Cn: data.Name.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa sudo command group %s from plan = %v", data.Name.ValueString(), data))
	if !data.Description.Equal(state.Description) {
		if data.Description.ValueStringPointer() != nil {
			optArgs.Description = data.Description.ValueStringPointer()
		} else {
			v := ""
			optArgs.Description = &v
		}
	}
	_, err := r.client.SudocmdgroupMod(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating freeipa sudo command group: %s", err))
	}

	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoCmdGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoCmdGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa sudo command group Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa sudo command group Name %s", data.Name.ValueString()))
	args := ipa.SudocmdgroupDelArgs{
		Cn: []string{data.Name.ValueString()},
	}
	optArgs := ipa.SudocmdgroupDelOptionalArgs{}
	_, err := r.client.SudocmdgroupDel(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Sudo command group %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}

}

func (r *SudoCmdGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
