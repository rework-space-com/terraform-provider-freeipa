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
var _ resource.Resource = &SudoRuleUserMembershipResource{}
var _ resource.ResourceWithImportState = &SudoRuleUserMembershipResource{}

func NewSudoRuleUserMembershipResource() resource.Resource {
	return &SudoRuleUserMembershipResource{}
}

// SudoRuleUserMembershipResource defines the resource implementation.
type SudoRuleUserMembershipResource struct {
	client *ipa.Client
}

// SudoRuleUserMembershipResourceModel describes the resource data model.
type SudoRuleUserMembershipResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	User       types.String `tfsdk:"user"`
	Users      types.List   `tfsdk:"users"`
	Group      types.String `tfsdk:"group"`
	Groups     types.List   `tfsdk:"groups"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *SudoRuleUserMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule_user_membership"
}

func (r *SudoRuleUserMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
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

func (r *SudoRuleUserMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo rule user membership resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Sudo rule name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User to add to the sudo rule",
				DeprecationMessage:  "use users instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"users": schema.ListAttribute{
				MarkdownDescription: "List of users to add to the sudo rule",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User group to add to the sudo rule",
				DeprecationMessage:  "use groups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "List of user groups to add to the sudo rule",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple sudo rule user membership resources on the same sudo rule. Manadatory for using users/groups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SudoRuleUserMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoRuleUserMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoRuleUserMembershipResourceModel
	var id, cmd_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudoruleAddUserOptionalArgs{}

	args := ipa.SudoruleAddUserArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.User.IsNull() {
		v := []string{data.User.ValueString()}
		optArgs.User = &v
		cmd_id = "sru"
	}
	if !data.Group.IsNull() {
		v := []string{data.Group.ValueString()}
		optArgs.Group = &v
		cmd_id = "srug"
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
		cmd_id = "msru"
	}

	_, err := r.client.SudoruleAddUser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule user membership: %s", err))
		return
	}

	switch cmd_id {
	case "sru":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.User.ValueString())
		data.Id = types.StringValue(id)
	case "srug":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.Group.ValueString())
		data.Id = types.StringValue(id)
	case "msru":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleUserMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoRuleUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudorule_user_membership: %s", err))
		return
	}

	all := true
	optArgs := ipa.SudoruleShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudoruleShowArgs{
		Cn: sudoruleId,
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

	switch typeId {
	case "sru":
		if res.Result.MemberuserUser == nil || !slices.Contains(*res.Result.MemberuserUser, cmdId) {
			resp.State.RemoveResource(ctx)
			return
		}
	case "srug":
		if res.Result.MemberuserGroup == nil || !slices.Contains(*res.Result.MemberuserGroup, cmdId) {
			resp.State.RemoveResource(ctx)
			return
		}
	case "msru":
		if res.Result.MemberuserUser == nil && res.Result.MemberuserGroup == nil {
			resp.State.RemoveResource(ctx)
			return
		} else {
			if !data.Users.IsNull() {
				var changedVals []string
				for _, value := range data.Users.Elements() {
					val, err := strconv.Unquote(value.String())
					if err != nil {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo user member failed with error %s", err))
					}
					if res.Result.MemberuserUser != nil && slices.Contains(*res.Result.MemberuserUser, val) {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo user member %s is present in results", val))
						changedVals = append(changedVals, val)
					}
				}
				var diag diag.Diagnostics
				data.Users, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
				if diag.HasError() {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
				}
			}
			if !data.Groups.IsNull() && res.Result.MemberuserGroup == nil {
				var changedVals []string
				for _, value := range data.Groups.Elements() {
					val, err := strconv.Unquote(value.String())
					if err != nil {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo user group member failed with error %s", err))
					}
					if slices.Contains(*res.Result.MemberuserGroup, val) {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo user group member %s is present in results", val))
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
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SudoRuleUserMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoRuleUserMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.SudoruleAddUserOptionalArgs{}

	memberAddArgs := ipa.SudoruleAddUserArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.SudoruleRemoveUserOptionalArgs{}

	memberDelArgs := ipa.SudoruleRemoveUserArgs{
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
		var statearr, planarr, addedGrps, deletedGrps []string

		for _, value := range state.Groups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.Groups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedGrps = append(addedGrps, val)
				memberAddOptArgs.Group = &addedGrps
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedGrps = append(deletedGrps, value)
				memberDelOptArgs.Group = &deletedGrps
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.SudoruleAddUser(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa sudo rule user membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule user membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule user membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.SudoruleRemoveUser(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa sudo rule group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule user membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule user membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleUserMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoRuleUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	sudoRuleId, typeId, _, err := parseSudoRuleUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_user_membership: %s", err))
		return
	}

	optArgs := ipa.SudoruleRemoveUserOptionalArgs{}

	args := ipa.SudoruleRemoveUserArgs{
		Cn: sudoRuleId,
	}

	switch typeId {
	case "sru":
		v := []string{data.User.ValueString()}
		optArgs.User = &v
	case "srug":
		v := []string{data.Group.ValueString()}
		optArgs.Group = &v
	case "msru":
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

	_, err = r.client.SudoruleRemoveUser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo user membership: %s", err))
		return
	}
}

func (r *SudoRuleUserMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseSudoRuleUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine sudo rule user membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
