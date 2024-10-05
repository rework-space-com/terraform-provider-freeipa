// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &userGroupMembership{}
var _ resource.ResourceWithImportState = &userGroupMembership{}

func NewUserGroupMembershipResource() resource.Resource {
	return &userGroupMembership{}
}

// resourceModel defines the resource implementation.
type userGroupMembership struct {
	client *ipa.Client
}

// resourceModelModel describes the resource data model.
type userGroupMembershipModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	User            types.String `tfsdk:"user"`
	Group           types.String `tfsdk:"group"`
	ExternalMember  types.String `tfsdk:"external_member"`
	Users           types.List   `tfsdk:"users"`
	Groups          types.List   `tfsdk:"groups"`
	ExternalMembers types.List   `tfsdk:"external_members"`
}

func (r *userGroupMembership) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group_membership"
}

func (r *userGroupMembership) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("user"),
			path.MatchRoot("group"),
			path.MatchRoot("external_member"),
			path.MatchRoot("users"),
			path.MatchRoot("groups"),
			path.MatchRoot("external_members"),
		),
	}
}

func (r *userGroupMembership) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User Group Membership resource",

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
			"user": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User to add",
				DeprecationMessage:  "use users instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "**deprecated** User group to add",
				DeprecationMessage:  "use groups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"external_member": schema.StringAttribute{
				MarkdownDescription: "**deprecated** External member to add. name must refer to an external group. (Requires a valid AD Trust configuration).",
				DeprecationMessage:  "use external_members instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"users": schema.ListAttribute{
				MarkdownDescription: "Users to add as group members",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "User groups to add as group members",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"external_members": schema.ListAttribute{
				MarkdownDescription: "External members to add as group members. name must refer to an external group. (Requires a valid AD Trust configuration).",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *userGroupMembership) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userGroupMembership) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userGroupMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create user group membership %s", data.Id.ValueString()))
	optArgs := ipa.GroupAddMemberOptionalArgs{}

	args := ipa.GroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.User.IsNull() {
		v := []string{string(data.User.ValueString())}
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create user group membership user %s", data.User.ValueString()))
		optArgs.User = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/u/%s", data.Name.ValueString(), data.User.ValueString()))

	}
	if !data.Group.IsNull() {
		v := []string{string(data.Group.ValueString())}
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create user group membership group %s", data.Group.ValueString()))
		optArgs.Group = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/g/%s", data.Name.ValueString(), data.Group.ValueString()))
	}
	if !data.ExternalMember.IsNull() {
		v := []string{string(data.ExternalMember.ValueString())}
		optArgs.Ipaexternalmember = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/e/%s", data.Name.ValueString(), data.ExternalMember.ValueString()))
	}
	if len(data.Users.Elements()) > 0 {
		var v []string
		for _, value := range data.Users.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.User = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/mu/%s", data.Name.ValueString(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	}
	if len(data.Groups.Elements()) > 0 {
		var v []string
		for _, value := range data.Groups.Elements() {
			val, _ := strconv.Unquote(value.String())
			if val == data.Name.ValueString() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %s cannot be membership of itself", data.Name.ValueString()))
				return
			}
			v = append(v, val)
		}
		optArgs.Group = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/mg/%s", data.Name.ValueString(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	}
	if len(data.ExternalMembers.Elements()) > 0 {
		var v []string
		for _, value := range data.ExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipaexternalmember = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/me/%s", data.Name.ValueString(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	}

	_v, err := r.client.GroupAddMember(&args, &optArgs)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa user group membership: %s", _v.String()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %s", err))
		return
	}
	if _v.Completed == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group membership: %v", _v.Failed))
		return
	}

	if !data.ExternalMember.IsNull() {
		v := []string{string(data.ExternalMember.ValueString())}
		z := new(bool)
		*z = true
		groupRes, err := r.client.GroupShow(&ipa.GroupShowArgs{Cn: data.Name.ValueString()}, &ipa.GroupShowOptionalArgs{All: z})
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show return is %s", groupRes.Result.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error looking up freeipa user group membership: %s", err))
			return
		}
		if !slices.Contains(*groupRes.Result.Ipaexternalmember, data.ExternalMember.ValueString()) {
			_, err = r.client.GroupRemoveMember(&ipa.GroupRemoveMemberArgs{Cn: data.Name.ValueString()}, &ipa.GroupRemoveMemberOptionalArgs{Ipaexternalmember: &v})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting invalid freeipa user group membership: %s", err))
				return
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("external member is not using the correct format. Use the lowercase upn format (ie: 'domain users@domain.net'): %s", data.ExternalMember.ValueString()))
			return
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show %s is %v", data.Name.ValueString(), groupRes.Result.String()))
		}
	}
	if len(data.ExternalMembers.Elements()) > 0 {
		z := new(bool)
		*z = true
		groupRes, err := r.client.GroupShow(&ipa.GroupShowArgs{Cn: data.Name.ValueString()}, &ipa.GroupShowOptionalArgs{All: z})
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show return is %s", groupRes.Result.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error looking up freeipa user group membership: %s", err))
			return
		}
		for _, value := range data.ExternalMembers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v := []string{val}
			if !slices.Contains(*groupRes.Result.Ipaexternalmember, val) {
				_, err = r.client.GroupRemoveMember(&ipa.GroupRemoveMemberArgs{Cn: data.Name.ValueString()}, &ipa.GroupRemoveMemberOptionalArgs{Ipaexternalmember: &v})
				if err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting invalid freeipa user group membership: %s", err))
					return
				}
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("external member is not using the correct format. Use the lowercase upn format (ie: 'domain users@domain.net'): %s", data.ExternalMember.ValueString()))
				return
			} else {
				tflog.Debug(ctx, fmt.Sprintf("[DEBUG] group show %s is %v", data.Name.ValueString(), groupRes.Result.String()))
			}
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupMembership) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userGroupMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name, typeId, userId, err := parseUserMembershipID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("State Error", fmt.Sprintf("Unable to parse resource %s: %s", data.Id.ValueString(), err))
	}
	reqArgs := ipa.GroupShowArgs{
		Cn: name,
	}
	z := new(bool)
	*z = true
	optArgs := ipa.GroupShowOptionalArgs{
		All: z,
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read user group membership %s optArgs %v", data.Id.ValueString(), optArgs))
	res, err := r.client.GroupShow(&reqArgs, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading information on freeipa user group %s: %s", name, err))
		return
	}

	switch typeId {
	case "g":
		v := []string{userId}
		groups := *res.Result.MemberGroup
		if slices.Contains(groups, v[0]) {
			data.Group = types.StringValue(v[0])
		} else {
			data.Group = types.StringValue("")
			data.Id = types.StringValue("")
		}
	case "u":
		v := []string{userId}
		users := *res.Result.MemberUser
		if slices.Contains(users, v[0]) {
			data.User = types.StringValue(v[0])
		} else {
			data.User = types.StringValue("")
			data.Id = types.StringValue("")
		}
	case "e":
		v := []string{userId}
		extmembers := *res.Result.Ipaexternalmember
		if slices.Contains(extmembers, v[0]) {
			data.ExternalMember = types.StringValue(v[0])
		} else {
			data.ExternalMember = types.StringValue("")
			data.Id = types.StringValue("")
		}
	case "mu":
		if userId == "users" && res.Result.MemberUser != nil {
			data.Users, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberUser)
		}
	case "mg":
		if userId == "groups" && res.Result.MemberGroup != nil {
			data.Groups, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberGroup)
		}
	case "me":
		if userId == "external" && res.Result.Ipaexternalmember != nil {
			data.ExternalMembers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ipaexternalmember)
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userGroupMembership) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userGroupMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupMembership) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userGroupMembershipModel

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
	optArgs := ipa.GroupRemoveMemberOptionalArgs{}

	nameId, typeId, userId, err := parseUserMembershipID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_user_group_membership %s: %s", data.Id.ValueString(), err))

	}

	args := ipa.GroupRemoveMemberArgs{
		Cn: nameId,
	}

	switch typeId {
	case "g":
		v := []string{userId}
		optArgs.Group = &v
	case "u":
		v := []string{userId}
		optArgs.User = &v
	case "e":
		v := []string{userId}
		optArgs.Ipaexternalmember = &v
	case "mu":
		if len(data.Users.Elements()) > 0 {
			var v []string
			for _, value := range data.Users.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.User = &v
		}
	case "mg":
		if len(data.Groups.Elements()) > 0 {
			var v []string
			for _, value := range data.Groups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Group = &v
		}
	case "me":
		if len(data.ExternalMembers.Elements()) > 0 {
			var v []string
			for _, value := range data.ExternalMembers.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Ipaexternalmember = &v
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	_, err = r.client.GroupRemoveMember(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error remove user group membership %s: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *userGroupMembership) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
