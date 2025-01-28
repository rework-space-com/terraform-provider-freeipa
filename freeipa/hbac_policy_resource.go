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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HbacPolicyResource{}
var _ resource.ResourceWithImportState = &HbacPolicyResource{}

func NewHbacPolicyResource() resource.Resource {
	return &HbacPolicyResource{}
}

// HbacPolicyResource defines the resource implementation.
type HbacPolicyResource struct {
	client *ipa.Client
}

// HbacPolicyResourceModel describes the resource data model.
type HbacPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	UserCategory    types.String `tfsdk:"usercategory"`
	HostCategory    types.String `tfsdk:"hostcategory"`
	ServiceCategory types.String `tfsdk:"servicecategory"`
}

func (r *HbacPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hbac_policy"
}

func (r *HbacPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *HbacPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA HBAC policy resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the hbac policy",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "HBAC policy description",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this hbac policy",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"usercategory": schema.StringAttribute{
				MarkdownDescription: "User category the hbac policy is applied to (allowed value: all)",
				Optional:            true,
			},
			"hostcategory": schema.StringAttribute{
				MarkdownDescription: "Host category the hbac policy is applied to (allowed value: all)",
				Optional:            true,
			},
			"servicecategory": schema.StringAttribute{
				MarkdownDescription: "Service category the hbac policy is applied to (allowed value: all)",
				Optional:            true,
			},
		},
	}
}

func (r *HbacPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HbacPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HbacPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HbacruleAddOptionalArgs{}

	args := ipa.HbacruleAddArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if !data.Enabled.IsNull() {
		v := data.Enabled.ValueBool()
		optArgs.Ipaenabledflag = &v
	}
	if !data.UserCategory.IsNull() {
		optArgs.Usercategory = data.UserCategory.ValueStringPointer()
	}
	if !data.HostCategory.IsNull() {
		optArgs.Hostcategory = data.HostCategory.ValueStringPointer()
	}
	if !data.ServiceCategory.IsNull() {
		optArgs.Servicecategory = data.ServiceCategory.ValueStringPointer()
	}
	_, err := r.client.HbacruleAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy: %s", err))
		return
	}

	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HbacPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.HbacruleShowArgs{
		Cn: data.Id.ValueString(),
	}

	res, err := r.client.HbacruleShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] Hbac policy not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa hbac policy: %s", err))
			return
		}
	}

	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.Ipaenabledflag != nil && !data.Enabled.IsNull() {
		data.Enabled = types.BoolValue(*res.Result.Ipaenabledflag)
	}
	if res.Result.Usercategory != nil && !data.UserCategory.IsNull() {
		data.UserCategory = types.StringValue(*res.Result.Usercategory)
	}
	if res.Result.Hostcategory != nil && !data.HostCategory.IsNull() {
		data.HostCategory = types.StringValue(*res.Result.Hostcategory)
	}
	if res.Result.Servicecategory != nil && !data.ServiceCategory.IsNull() {
		data.ServiceCategory = types.StringValue(*res.Result.Servicecategory)
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HbacPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HbacPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.HbacruleModArgs{
		Cn: data.Id.ValueString(),
	}
	optArgs := ipa.HbacruleModOptionalArgs{}

	var hasChange = false

	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
		hasChange = true
	}
	if !data.Enabled.Equal(state.Enabled) {
		if !data.Enabled.ValueBool() {
			_, err := r.client.HbacruleDisable(&ipa.HbacruleDisableArgs{Cn: data.Id.ValueString()}, &ipa.HbacruleDisableOptionalArgs{})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling freeipa hbac policy: %s", err))
			}
		} else {
			_, err := r.client.HbacruleEnable(&ipa.HbacruleEnableArgs{Cn: data.Id.ValueString()}, &ipa.HbacruleEnableOptionalArgs{})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error enabling freeipa hbac policy: %s", err))
			}
		}
	}
	if !data.UserCategory.Equal(state.UserCategory) {
		if data.UserCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Usercategory = &v
		} else {
			optArgs.Usercategory = data.UserCategory.ValueStringPointer()
		}
		hasChange = true
	}
	if !data.HostCategory.Equal(state.HostCategory) {
		if data.HostCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Hostcategory = &v
		} else {
			optArgs.Hostcategory = data.HostCategory.ValueStringPointer()
		}
		hasChange = true
	}
	if !data.ServiceCategory.Equal(state.ServiceCategory) {
		if data.ServiceCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Servicecategory = &v
		} else {
			optArgs.Servicecategory = data.ServiceCategory.ValueStringPointer()
		}
		hasChange = true
	}

	if hasChange {
		_, err := r.client.HbacruleMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				resp.Diagnostics.AddError("Client Error", "EmptyModlist (4202): no modifications to be performed")
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa hbac policy: %s", err))
				return
			}
		}
	}

	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HbacPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.HbacruleDelArgs{
		Cn: []string{data.Id.ValueString()},
	}
	_, err := r.client.HbacruleDel(&args, &ipa.HbacruleDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa hbac policy: %s", err))
		return
	}
}

func (r *HbacPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
