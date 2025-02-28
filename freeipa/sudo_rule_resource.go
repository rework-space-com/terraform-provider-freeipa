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
var _ resource.Resource = &SudoRuleResource{}
var _ resource.ResourceWithImportState = &SudoRuleResource{}

func NewSudoRuleResource() resource.Resource {
	return &SudoRuleResource{}
}

// SudoRuleResource defines the resource implementation.
type SudoRuleResource struct {
	client *ipa.Client
}

// SudoRuleResourceModel describes the resource data model.
type SudoRuleResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	UserCategory       types.String `tfsdk:"usercategory"`
	HostCategory       types.String `tfsdk:"hostcategory"`
	CommandCategory    types.String `tfsdk:"commandcategory"`
	RunAsUserCategory  types.String `tfsdk:"runasusercategory"`
	RunAsGroupCategory types.String `tfsdk:"runasgroupcategory"`
	Order              types.Int32  `tfsdk:"order"`
}

func (r *SudoRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule"
}

func (r *SudoRuleResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *SudoRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo rule resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the sudo rule",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Sudo rule description",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this sudo rule",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"usercategory": schema.StringAttribute{
				MarkdownDescription: "User category the sudo rule is applied to (allowed value: all)",
				Optional:            true,
			},
			"hostcategory": schema.StringAttribute{
				MarkdownDescription: "Host category the sudo rule is applied to (allowed value: all)",
				Optional:            true,
			},
			"commandcategory": schema.StringAttribute{
				MarkdownDescription: "Command category the sudo rule is applied to (allowed value: all)",
				Optional:            true,
			},
			"runasusercategory": schema.StringAttribute{
				MarkdownDescription: "Run as user category the sudo rule is applied to (allowed value: all)",
				Optional:            true,
			},
			"runasgroupcategory": schema.StringAttribute{
				MarkdownDescription: "Run as group category the sudo rule is applied to (allowed value: all)",
				Optional:            true,
			},
			"order": schema.Int32Attribute{
				MarkdownDescription: "Sudo rule order (must be unique)",
				Optional:            true,
			},
		},
	}
}

func (r *SudoRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudoruleAddOptionalArgs{}

	args := ipa.SudoruleAddArgs{
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
	if !data.RunAsUserCategory.IsNull() {
		optArgs.Ipasudorunasusercategory = data.RunAsUserCategory.ValueStringPointer()
	}
	if !data.CommandCategory.IsNull() {
		optArgs.Cmdcategory = data.CommandCategory.ValueStringPointer()
	}
	if !data.RunAsGroupCategory.IsNull() {
		optArgs.Ipasudorunasgroupcategory = data.RunAsGroupCategory.ValueStringPointer()
	}
	if !data.Order.IsNull() {
		v := int(data.Order.ValueInt32())
		optArgs.Sudoorder = &v
	}
	_, err := r.client.SudoruleAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule: %s", err))
		return
	}

	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.SudoruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudoruleShowArgs{
		Cn: data.Id.ValueString(),
	}

	res, err := r.client.SudoruleShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] Sudo rule not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa sudo rule: %s", err))
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
	if res.Result.Cmdcategory != nil && !data.CommandCategory.IsNull() {
		data.CommandCategory = types.StringValue(*res.Result.Cmdcategory)
	}
	if res.Result.Ipasudorunasusercategory != nil && !data.RunAsUserCategory.IsNull() {
		data.RunAsUserCategory = types.StringValue(*res.Result.Ipasudorunasusercategory)
	}
	if res.Result.Ipasudorunasgroupcategory != nil && !data.RunAsUserCategory.IsNull() {
		data.RunAsUserCategory = types.StringValue(*res.Result.Ipasudorunasgroupcategory)
	}
	if res.Result.Sudoorder != nil && !data.Order.IsNull() {
		data.Order = types.Int32Value(int32(*res.Result.Sudoorder))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SudoRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.SudoruleModArgs{
		Cn: data.Id.ValueString(),
	}
	optArgs := ipa.SudoruleModOptionalArgs{}

	var hasChange = false

	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
		hasChange = true
	}
	if !data.Enabled.Equal(state.Enabled) {
		if !data.Enabled.ValueBool() {
			_, err := r.client.SudoruleDisable(&ipa.SudoruleDisableArgs{Cn: data.Id.ValueString()}, &ipa.SudoruleDisableOptionalArgs{})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling freeipa sudo rule: %s", err))
			}
		} else {
			_, err := r.client.SudoruleEnable(&ipa.SudoruleEnableArgs{Cn: data.Id.ValueString()}, &ipa.SudoruleEnableOptionalArgs{})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error enabling freeipa sudo rule: %s", err))
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
	if !data.RunAsUserCategory.Equal(state.RunAsUserCategory) {
		if data.RunAsUserCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Ipasudorunasusercategory = &v
		} else {
			optArgs.Ipasudorunasusercategory = data.RunAsUserCategory.ValueStringPointer()
		}
		hasChange = true
	}
	if !data.CommandCategory.Equal(state.CommandCategory) {
		if data.CommandCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Cmdcategory = &v
		} else {
			optArgs.Cmdcategory = data.CommandCategory.ValueStringPointer()
		}
		hasChange = true
	}
	if !data.RunAsGroupCategory.Equal(state.RunAsGroupCategory) {
		if data.RunAsGroupCategory.ValueStringPointer() == nil {
			v := ""
			optArgs.Ipasudorunasgroupcategory = &v
		} else {
			optArgs.Ipasudorunasgroupcategory = data.RunAsGroupCategory.ValueStringPointer()
		}
		hasChange = true
	}
	// TODO update rule when order is removed
	if !data.Order.Equal(state.Order) {
		if data.Order.ValueInt32Pointer() != nil {
			v := int(data.Order.ValueInt32())
			optArgs.Sudoorder = &v
		} else {
			optArgs.Sudoorder = nil
		}
		hasChange = true
	}

	if hasChange {
		_, err := r.client.SudoruleMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				resp.Diagnostics.AddWarning("Client Warning", err.Error())
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa sudo rule: %s", err))
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

func (r *SudoRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.SudoruleDelArgs{
		Cn: []string{data.Id.ValueString()},
	}
	_, err := r.client.SudoruleDel(&args, &ipa.SudoruleDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo rule: %s", err))
		return
	}
}

func (r *SudoRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
