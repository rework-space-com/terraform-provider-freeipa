// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
	"golang.org/x/exp/slices"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HbacPolicyUserMembershipResource{}
var _ resource.ResourceWithImportState = &HbacPolicyUserMembershipResource{}

func NewHbacPolicyUserMembershipResource() resource.Resource {
	return &HbacPolicyUserMembershipResource{}
}

// HbacPolicyUserMembershipResource defines the resource implementation.
type HbacPolicyUserMembershipResource struct {
	client *ipa.Client
}

// HbacPolicyUserMembershipResourceModel describes the resource data model.
type HbacPolicyUserMembershipResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	User       types.String `tfsdk:"user"`
	Users      types.List   `tfsdk:"users"`
	Group      types.String `tfsdk:"group"`
	Groups     types.List   `tfsdk:"groups"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *HbacPolicyUserMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hbac_policy_user_membership"
}

func (r *HbacPolicyUserMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("user"),
			path.MatchRoot("users"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("user"),
			path.MatchRoot("group"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("user"),
			path.MatchRoot("groups"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("group"),
			path.MatchRoot("users"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("group"),
			path.MatchRoot("groups"),
		),
	}
}

func (r *HbacPolicyUserMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA HBAC policy host membership resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "HBAC policy name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User FDQN the policy is applied to",
				DeprecationMessage:  "use users instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"users": schema.ListAttribute{
				MarkdownDescription: "List of user FQDNs to add to the HBAC policy",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User group to add to the HBAC policy",
				DeprecationMessage:  "use hostgroups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "List of user groups to add to the HBAC policy",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple HBAC policy user membership resources on the same HBAC policy. Manadatory for using users/groups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *HbacPolicyUserMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HbacPolicyUserMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HbacPolicyUserMembershipResourceModel
	var id, user_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HbacruleAddUserOptionalArgs{}

	args := ipa.HbacruleAddUserArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.User.IsNull() {
		v := []string{data.User.ValueString()}
		optArgs.User = &v
		user_id = "u"
	}
	if !data.Group.IsNull() {
		v := []string{data.Group.ValueString()}
		optArgs.Group = &v
		user_id = "g"
	}
	if !data.Users.IsNull() || !data.Groups.IsNull() {
		if !data.Users.IsNull() {
			var v []string
			for _, value := range data.Users.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.User = &v
		}
		if !data.Groups.IsNull() {
			var v []string
			for _, value := range data.Groups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Group = &v
		}
		user_id = "mu"
	}

	_, err := r.client.HbacruleAddUser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule user membership: %s", err))
		return
	}

	switch user_id {
	case "u":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.User.ValueString())
		data.Id = types.StringValue(id)
	case "g":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.Group.ValueString())
		data.Id = types.StringValue(id)
	case "mu":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyUserMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HbacPolicyUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hbacpolicyid, typeId, policyId, err := parseHBACPolicyUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_user_membership: %s", err))
		return
	}

	all := true
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.HbacruleShowArgs{
		Cn: hbacpolicyid,
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

	switch typeId {
	case "u":
		if res.Result.MemberuserUser == nil || !isStringListContainsCaseInsensistive(res.Result.MemberuserUser, &policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy user membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "g":
		if res.Result.MemberuserGroup == nil || !isStringListContainsCaseInsensistive(res.Result.MemberuserGroup, &policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy user group membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "mu":
		if !data.Users.IsNull() {
			var changedVals []string
			for _, value := range data.Users.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy user member failed with error %s", err))
				}
				if res.Result.MemberuserUser != nil && isStringListContainsCaseInsensistive(res.Result.MemberuserUser, &val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy user member %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.Users, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if !data.Groups.IsNull() {
			var changedVals []string
			for _, value := range data.Groups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy member commands failed with error %s", err))
				}
				if res.Result.MemberuserGroup != nil && isStringListContainsCaseInsensistive(res.Result.MemberuserGroup, &val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.Groups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HbacPolicyUserMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HbacPolicyUserMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.HbacruleAddUserOptionalArgs{}

	memberAddArgs := ipa.HbacruleAddUserArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.HbacruleRemoveUserOptionalArgs{}

	memberDelArgs := ipa.HbacruleRemoveUserArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.Users.Equal(state.Users) {
		var statearr, planarr, addedUsers, deletedUsers []string

		for _, value := range state.Users.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.Users.Elements() {
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
	if !data.Groups.Equal(state.Groups) {
		var statearr, planarr, addedGroups, deletedGroups []string

		for _, value := range state.Groups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.Groups.Elements() {
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
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.HbacruleAddUser(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa hbac policy user membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy user membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy user membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.HbacruleRemoveUser(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa hbac policy user membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy user membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy user membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyUserMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HbacPolicyUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	hbacpolicyId, typeId, _, err := parseHBACPolicyUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_user_membership: %s", err))
		return
	}

	optArgs := ipa.HbacruleRemoveUserOptionalArgs{}

	args := ipa.HbacruleRemoveUserArgs{
		Cn: hbacpolicyId,
	}

	switch typeId {
	case "u":
		v := []string{data.User.ValueString()}
		optArgs.User = &v
	case "g":
		v := []string{data.Group.ValueString()}
		optArgs.Group = &v
	case "mu":
		if !data.Users.IsNull() {
			var v []string
			for _, value := range data.Users.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.User = &v
		}
		if !data.Groups.IsNull() {
			var v []string
			for _, value := range data.Groups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Group = &v
		}
	}

	_, err = r.client.HbacruleRemoveUser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa hbac policy user membership: %s", err))
		return
	}
}

func (r *HbacPolicyUserMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseHBACPolicyUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine user membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
