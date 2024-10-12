// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

// UserGroupResource defines the resource implementation.
type UserGroupResource struct {
	client *ipa.Client
}

// UserGroupResourceModel describes the resource data model.
type UserGroupResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	GidNumber             types.Int64  `tfsdk:"gid_number"`
	NonPosix              types.Bool   `tfsdk:"nonposix"`
	External              types.Bool   `tfsdk:"external"`
	MemberUsers           types.List   `tfsdk:"member_users"`
	MemberGroups          types.List   `tfsdk:"member_groups"`
	MemberExternalMembers types.List   `tfsdk:"member_external_members"`
	AddAttr               types.List   `tfsdk:"addattr"`
	SetAttr               types.List   `tfsdk:"setattr"`
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
				MarkdownDescription: "Group name",
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
				Computed:            false,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"external": schema.BoolAttribute{
				MarkdownDescription: "Allow adding external non-IPA members from trusted domains",
				Optional:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"member_users": schema.ListAttribute{
				MarkdownDescription: "Users to add as group members",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"member_groups": schema.ListAttribute{
				MarkdownDescription: "User groups to add as group members",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"member_external_members": schema.ListAttribute{
				MarkdownDescription: "External members to add as group members. name must refer to an external group. (Requires a valid AD Trust configuration).",
				Optional:            true,
				ElementType:         types.StringType,
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

	memberOptArgs := ipa.GroupAddMemberOptionalArgs{}

	memberArgs := ipa.GroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}

	hasMembership := false
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
	// member users, group and external members are added when all inputs are processed.
	// This allows to make one single api call for all members.
	if len(data.MemberUsers.Elements()) > 0 {
		var v []string
		for _, value := range data.MemberUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		memberOptArgs.User = &v
		hasMembership = true
	}
	if len(data.MemberGroups.Elements()) > 0 {
		var v []string
		for _, value := range data.MemberGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			if val == data.Name.ValueString() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %s cannot be membership of itself", data.Name.ValueString()))
				return
			}
			v = append(v, val)
		}
		memberOptArgs.Group = &v
		hasMembership = true
	}
	if len(data.MemberExternalMembers.Elements()) > 0 {
		var v []string
		for _, value := range data.MemberExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		memberOptArgs.Ipaexternalmember = &v
		hasMembership = true
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

	// Add membership if any is detected in the resource configuration.
	if hasMembership {
		_v, err := r.client.GroupAddMember(&memberArgs, &memberOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa user group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %v", _v.Failed))
			return
		}
		// If the external members from configuration are not reported in the list of external members from groupshow
		// Then FreeIPA transformed the name of the external member during creation.
		// If group add doesn't fail, remove the external member missing from the list and throw an error.
		if len(data.MemberExternalMembers.Elements()) > 0 {
			z := new(bool)
			*z = true
			groupRes, err := r.client.GroupShow(&ipa.GroupShowArgs{Cn: data.Name.ValueString()}, &ipa.GroupShowOptionalArgs{All: z})
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show return is %s", groupRes.Result.String()))
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error looking up freeipa user group membership: %s", err))
			}
			for _, value := range data.MemberExternalMembers.Elements() {
				val, _ := strconv.Unquote(value.String())
				v := []string{val}
				if !slices.Contains(*groupRes.Result.Ipaexternalmember, val) {
					_, err = r.client.GroupRemoveMember(&ipa.GroupRemoveMemberArgs{Cn: data.Name.ValueString()}, &ipa.GroupRemoveMemberOptionalArgs{Ipaexternalmember: &v})
					if err != nil {
						resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting invalid freeipa user group membership: %s", err))
					}
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("external member is not using the correct format. Use the lowercase upn format (ie: 'domain users@domain.net'): %s", val))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show %s is %v", data.Name.ValueString(), groupRes.Result.String()))
				}
			}
		}

	}

	data.Id = types.StringValue(data.Name.ValueString())

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

	data.Name = types.StringValue(res.Result.Cn)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Cn %s", data.Name.ValueString()))
	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group Description %s", data.Description.ValueString()))
	}
	if res.Result.Gidnumber != nil && !data.GidNumber.IsNull() {
		data.GidNumber = types.Int64Value(int64(*res.Result.Gidnumber))
	}
	// This read function will only keep members that are defined in the current state AND in the real member list of the group
	// This avoids conflicts with membership that would have been set outside of this resource.
	if !data.MemberExternalMembers.IsNull() && res.Result.Ipaexternalmember != nil {
		var changedVals []string
		for _, value := range data.MemberExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.Ipaexternalmember, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.MemberExternalMembers, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if !data.MemberUsers.IsNull() && res.Result.MemberUser != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group member users %v", *res.Result.MemberUser))
		var changedVals []string
		for _, value := range data.MemberUsers.Elements() {
			val, err := strconv.Unquote(value.String())
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group member users failed with error %s", err))
			}
			if slices.Contains(*res.Result.MemberUser, val) {
				tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa group member users %s is present in results", val))
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.MemberUsers, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	if !data.MemberGroups.IsNull() && res.Result.MemberGroup != nil {
		var changedVals []string
		for _, value := range data.MemberGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.MemberGroup, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.MemberGroups, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }
	args := ipa.GroupModArgs{
		Cn: data.Name.ValueString(),
	}
	optArgs := ipa.GroupModOptionalArgs{}

	memberAddOptArgs := ipa.GroupAddMemberOptionalArgs{}

	memberAddArgs := ipa.GroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.GroupRemoveMemberOptionalArgs{}

	memberDelArgs := ipa.GroupRemoveMemberArgs{
		Cn: data.Name.ValueString(),
	}
	hasGroupMod := false
	hasMemberAdd := false
	hasMemberDel := false
	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
		hasGroupMod = true
	}

	if !data.GidNumber.Equal(state.GidNumber) {
		gid := int(data.GidNumber.ValueInt64())
		optArgs.Gidnumber = &gid
		hasGroupMod = true
	}
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.MemberUsers.Equal(state.MemberUsers) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group member users %s ", data.MemberUsers.String()))
		var statearr, planarr, addedUsers, deletedUsers []string

		for _, value := range state.MemberUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.MemberUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedUsers = append(addedUsers, val)
				memberAddOptArgs.User = &addedUsers
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedUsers = append(deletedUsers, value)
				memberDelOptArgs.User = &deletedUsers
				hasMemberDel = true
			}
		}

	}
	if !data.MemberGroups.Equal(state.MemberGroups) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group member users %s ", data.MemberGroups.String()))
		var statearr, planarr, addedGroups, deletedGroups []string

		for _, value := range state.MemberGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.MemberGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedGroups = append(addedGroups, val)
				memberAddOptArgs.Group = &addedGroups
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedGroups = append(deletedGroups, value)
				memberDelOptArgs.Group = &deletedGroups
				hasMemberDel = true
			}
		}

	}
	if !data.MemberExternalMembers.Equal(state.MemberExternalMembers) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group member users %s ", data.MemberExternalMembers.String()))
		var statearr, planarr, addedExt, deletedExt []string

		for _, value := range state.MemberExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.MemberExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedExt = append(addedExt, val)
				memberAddOptArgs.Ipaexternalmember = &addedExt
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedExt = append(deletedExt, value)
				memberDelOptArgs.Ipaexternalmember = &deletedExt
				hasMemberAdd = true
			}
		}

	}
	if !data.AddAttr.Equal(state.AddAttr) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
		hasGroupMod = true
	}

	if !data.SetAttr.Equal(state.SetAttr) {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
		hasGroupMod = true
	}

	if hasGroupMod {
		res, err := r.client.GroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				tflog.Debug(ctx, "[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa group %s: %s", res.Result.Cn, err))
				return
			}
		}
	}

	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.GroupAddMember(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa user group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %v", _v.Failed))
			return
		}
		if memberAddOptArgs.Ipaexternalmember != nil {
			z := new(bool)
			*z = true
			groupRes, err := r.client.GroupShow(&ipa.GroupShowArgs{Cn: data.Name.ValueString()}, &ipa.GroupShowOptionalArgs{All: z})
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show return is %s", groupRes.Result.String()))
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error looking up freeipa user group membership: %s", err))
				return
			}
			for _, value := range *memberAddOptArgs.Ipaexternalmember {
				v := []string{value}
				if !slices.Contains(*groupRes.Result.Ipaexternalmember, value) {
					_, err = r.client.GroupRemoveMember(&ipa.GroupRemoveMemberArgs{Cn: data.Name.ValueString()}, &ipa.GroupRemoveMemberOptionalArgs{Ipaexternalmember: &v})
					if err != nil {
						resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting invalid freeipa user group membership: %s", err))
					}
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("external member is not using the correct format. Use the lowercase upn format (ie: 'domain users@domain.net'): %s", value))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show %s is %v", data.Name.ValueString(), groupRes.Result.String()))
				}
			}
		}
	}
	if hasMemberDel {
		_v, err := r.client.GroupRemoveMember(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa user group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa user group membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa user group membership: %v", _v.Failed))
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
