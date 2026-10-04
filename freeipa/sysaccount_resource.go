// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//	Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strconv"
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
var _ resource.Resource = &SysAccountResource{}
var _ resource.ResourceWithModifyPlan = &SysAccountResource{}
var _ resource.ResourceWithImportState = &SysAccountResource{}

func NewSysAccountResource() resource.Resource {
	return &SysAccountResource{}
}

// SysAccountResource defines the resource implementation.
type SysAccountResource struct {
	client *ipa.Client
}

// SysAccountResourceModel describes the resource data model.
type SysAccountResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Password        types.String `tfsdk:"password"`
	RandomPassword  types.String `tfsdk:"random_password"`
	Random          types.Bool   `tfsdk:"random"`
	AccountDisabled types.Bool   `tfsdk:"disabled"`
	Privileged      types.Bool   `tfsdk:"privileged"`
	AddAttr         types.List   `tfsdk:"addattr"`
	SetAttr         types.List   `tfsdk:"setattr"`
}

func (r *SysAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sysaccount"
}

func (r *SysAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA System Account resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "System account name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "System account Description",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The system account's password",
				Optional:            true,
				Sensitive:           true,
			},
			"random_password": schema.StringAttribute{
				MarkdownDescription: "Generated random password",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"random": schema.BoolAttribute{
				MarkdownDescription: "Generate a random user password. Generated random password is returned in `random_password`.",
				Optional:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Disable the system account.",
				Optional:            true,
			},
			"privileged": schema.BoolAttribute{
				MarkdownDescription: "Allow password updates without reset.",
				Optional:            true,
			},
			"addattr": schema.ListAttribute{
				MarkdownDescription: "Add an attribute/value pair. Format is attr=value. The attribute must be part of the LDAP schema.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"setattr": schema.ListAttribute{
				MarkdownDescription: "Set an attribute to a name/value pair. Format is attr=value.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *SysAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r SysAccountResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		return
	}

	var data SysAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if !data.Password.IsNull() && data.Random.ValueBool() {
		resp.Diagnostics.AddError("Atrtibute Validation", "Random password and manual password are in conflict.")
		return
	}

}

func (r *SysAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SysAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.SysaccountAddOptionalArgs{
		All: &all,
	}

	args := ipa.SysaccountAddArgs{
		UID: data.Name.ValueString(),
	}
	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if !data.Password.IsNull() {
		optArgs.Userpassword = data.Password.ValueStringPointer()
	}
	if !data.Random.IsNull() {
		optArgs.Random = data.Random.ValueBoolPointer()
	}
	if !data.Privileged.IsNull() {
		optArgs.Privileged = data.Privileged.ValueBoolPointer()
	}
	if !data.AccountDisabled.IsNull() {
		optArgs.Nsaccountlock = data.AccountDisabled.ValueBoolPointer()
	}
	if len(data.AddAttr.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa system account Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if len(data.SetAttr.Elements()) > 0 {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}

	ret, err := r.client.SysaccountAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa system account: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa system account returned %s ", ret.Result.String()))
	if ret.Result.Randompassword != nil && !data.Random.IsNull() {
		data.RandomPassword = types.StringValue(*ret.Result.Randompassword)
	} else {
		data.RandomPassword = types.StringNull()
	}

	data.Id = types.StringValue(data.Name.ValueString())

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SysAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SysAccountResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
	if res.Result.Random != nil {
		data.Random = types.BoolValue(*res.Result.Random)
	}
	if res.Result.Randompassword != nil {
		data.RandomPassword = types.StringValue(*res.Result.Randompassword)
	}
	if res.Result.Nsaccountlock != nil && !data.AccountDisabled.IsNull() {
		data.AccountDisabled = types.BoolValue(*res.Result.Nsaccountlock)
	}
	data.Id = types.StringValue(data.Name.ValueString())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SysAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SysAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.SysaccountModArgs{
		UID: data.Name.ValueString(),
	}

	all := true
	optArgs := ipa.SysaccountModOptionalArgs{
		All: &all,
	}

	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if !data.Password.Equal(state.Password) {
		optArgs.Userpassword = data.Password.ValueStringPointer()
	}
	if !data.Random.Equal(state.Random) {
		optArgs.Random = data.Random.ValueBoolPointer()
	}

	if !data.AddAttr.Equal(state.AddAttr) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa system account Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if !data.SetAttr.Equal(state.SetAttr) {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}

	_, err := r.client.SysaccountMod(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa system account %s: %s", data.Name.ValueString(), err))
			return
		}
	}

	if !data.Privileged.Equal(state.Privileged) {
		_, err := r.client.SysaccountPolicy(&ipa.SysaccountPolicyArgs{UID: data.Name.ValueString()}, &ipa.SysaccountPolicyOptionalArgs{Privileged: data.Privileged.ValueBoolPointer()})
		if err != nil {
			if !strings.Contains(err.Error(), "EmptyModlist") {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa system account %s: %s", data.Name.ValueString(), err))
				return
			}
		}
	}

	if !data.AccountDisabled.Equal(state.AccountDisabled) {
		if data.AccountDisabled.ValueBool() {
			_, err := r.client.SysaccountDisable(&ipa.SysaccountDisableArgs{UID: data.Name.ValueString()}, &ipa.SysaccountDisableOptionalArgs{})
			if err != nil {
				if !strings.Contains(err.Error(), "EmptyModlist") {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa system account %s: %s", data.Name.ValueString(), err))
					return
				}
			}
		} else {
			_, err := r.client.SysaccountEnable(&ipa.SysaccountEnableArgs{UID: data.Name.ValueString()}, &ipa.SysaccountEnableOptionalArgs{})
			if err != nil {
				if !strings.Contains(err.Error(), "EmptyModlist") {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa system account %s: %s", data.Name.ValueString(), err))
					return
				}
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SysAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SysAccountResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa system account Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa system account Name %s", data.Name.ValueString()))
	args := ipa.SysaccountDelArgs{
		UID: []string{data.Name.ValueString()},
	}
	_, err := r.client.SysaccountDel(&args, &ipa.SysaccountDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] System account %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *SysAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// check if it is a posix group
	optArgs := ipa.SysaccountFindOptionalArgs{
		UID: &req.ID,
	}
	args := ipa.SysaccountFindArgs{}
	res, err := r.client.SysaccountFind(req.ID, &args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Error reading freeipa system account %s", req.ID))
		return
	}

	for _, grp := range res.Result {
		if grp.UID == req.ID {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), grp.UID)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), grp.UID)...)
			return
		}
	}
}
