// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//	Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//	Mixton <maxime.thomas@mtconsulting.tech>
//	Roman Butsiy <butsiyroman@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserGroupResource{}

// var _ resource.ResourceWithImportState = &UserGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

// UserGroupResource defines the resource implementation.
type UserGroupResource struct {
	client *ipa.Client
}

// UserGroupResourceModel describes the resource data model.
type UserGroupResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	GidNumber   types.Int64  `tfsdk:"gid_number"`
	NonPosix    types.Bool   `tfsdk:"nonposix"`
	External    types.Bool   `tfsdk:"external"`
	AddAttr     types.List   `tfsdk:"addattr"`
	SetAttr     types.List   `tfsdk:"setattr"`
}

func (r *UserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *UserGroupResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("gid_number"),
			path.MatchRoot("nonposix"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("gid_number"),
			path.MatchRoot("external"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("external"),
			path.MatchRoot("nonposix"),
		),
	}
}

func (r *UserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User Group resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Group name\n\n	- The name must not exceed 32 characters.\n	- The name must contain only lowercase letters (a-z), digits (0-9), and the characters (. - _).\n	- The name must not start with a special character.\n	- A user and a group cannot have the same name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Group Description",
				Optional:            true,
				Computed:            false,
			},
			"gid_number": schema.Int64Attribute{
				MarkdownDescription: "GID (use this option to set it manually)",
				Optional:            true,
				Computed:            false,
			},
			"nonposix": schema.BoolAttribute{
				MarkdownDescription: "Create as a non-POSIX group",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"external": schema.BoolAttribute{
				MarkdownDescription: "Allow adding external non-IPA members from trusted domains",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"addattr": schema.ListAttribute{
				MarkdownDescription: "Add an attribute/value pair. Format is attr=value. The attribute must be part of the schema.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"setattr": schema.ListAttribute{
				MarkdownDescription: "Set an attribute to a name/value pair. Format is attr=value.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *UserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.GroupAddOptionalArgs{}

	args := ipa.GroupAddArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}

	if !data.GidNumber.IsNull() {
		gid := int(data.GidNumber.ValueInt64())
		optArgs.Gidnumber = &gid
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group returned %d ", gid))
	}

	if !data.NonPosix.IsNull() {
		optArgs.Nonposix = data.NonPosix.ValueBoolPointer()
	}

	if !data.External.IsNull() {
		optArgs.External = data.External.ValueBoolPointer()
	}
	if len(data.AddAttr.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group Addattr %s ", data.AddAttr.String()))
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
	ret, err := r.client.GroupAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group returned %s ", ret.Result.String()))
	if ret.Result.Gidnumber != nil && !data.GidNumber.IsNull() {
		data.GidNumber = types.Int64Value(int64(*ret.Result.Gidnumber))
	}

	tflog.Trace(ctx, "created a user group resource")

	data.Id = types.StringValue(ret.Result.Cn)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserGroupResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.GroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.GroupShowArgs{
		Cn: data.Id.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %s", data.Id.ValueString()))
	res, err := r.client.GroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Group %s not found", data.Id.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Group %s not found: %s", data.Id.ValueString(), err))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %v", res))
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa group %s", data.Name.ValueString()))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Cn %s", data.Name.ValueString()))
	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Description %s", data.Description.ValueString()))
	}
	if res.Result.Gidnumber != nil && !data.GidNumber.IsNull() {
		data.GidNumber = types.Int64Value(int64(*res.Result.Gidnumber))
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group GID %d", data.GidNumber.ValueInt64()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state UserGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.GroupModArgs{
		Cn: data.Name.ValueString(),
	}
	optArgs := ipa.GroupModOptionalArgs{}

	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
	}

	if !data.GidNumber.Equal(state.GidNumber) {
		gid := int(data.GidNumber.ValueInt64())
		optArgs.Gidnumber = &gid
	}

	if !data.AddAttr.Equal(state.AddAttr) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group Addattr %s ", data.AddAttr.String()))
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

	res, err := r.client.GroupMod(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa group %s: %s", res.Result.Cn, err))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa group Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa group Name %s", data.Name.ValueString()))
	args := ipa.GroupDelArgs{
		Cn: []string{data.Name.ValueString()},
	}
	_, err := r.client.GroupDel(&args, &ipa.GroupDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Group %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	valTrue := true
	// check if it is a posix group
	optArgs := ipa.GroupFindOptionalArgs{
		Posix: &valTrue,
		Cn:    &req.ID,
	}
	args := ipa.GroupFindArgs{}
	res, err := r.client.GroupFind(req.ID, &args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Error reading freeipa group %s", req.ID))
		return
	}

	for _, grp := range res.Result {
		if grp.Cn == req.ID {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("nonposix"), false)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("external"), false)...)
			return
		}
	}

	// check if it is a external group
	optArgs = ipa.GroupFindOptionalArgs{
		External: &valTrue,
		Cn:       &req.ID,
	}
	res, err = r.client.GroupFind(req.ID, &args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Error reading freeipa group %s", req.ID))
		return
	}

	for _, grp := range res.Result {
		if grp.Cn == req.ID {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("nonposix"), false)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("external"), true)...)
			return
		}
	}

	// check if it is a non posix group
	optArgs = ipa.GroupFindOptionalArgs{
		Nonposix: &valTrue,
		Cn:       &req.ID,
	}
	res, err = r.client.GroupFind(req.ID, &args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Error reading freeipa group %s", req.ID))
		return
	}

	for _, grp := range res.Result {
		if grp.Cn == req.ID {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), grp.Cn)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("nonposix"), true)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("external"), false)...)
			return
		}
	}

}
