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
var _ resource.Resource = &HostGroupMembership{}
var _ resource.ResourceWithImportState = &HostGroupMembership{}

func NewHostGroupMembershipResource() resource.Resource {
	return &HostGroupMembership{}
}

// resourceModel defines the resource implementation.
type HostGroupMembership struct {
	client *ipa.Client
}

// resourceModelModel describes the resource data model.
type HostGroupMembershipModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Host       types.String `tfsdk:"host"`
	HostGroup  types.String `tfsdk:"hostgroup"`
	Hosts      types.List   `tfsdk:"hosts"`
	HostGroups types.List   `tfsdk:"hostgroups"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *HostGroupMembership) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_hostgroup_membership"
}

func (r *HostGroupMembership) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hostgroup"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hosts"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hostgroups"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("hostgroup"),
			path.MatchRoot("hosts"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("hostgroup"),
			path.MatchRoot("hostgroups"),
		),
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("host"),
			path.MatchRoot("hostgroup"),
			path.MatchRoot("hosts"),
			path.MatchRoot("hostgroups"),
		),
	}
}

func (r *HostGroupMembership) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				MarkdownDescription: "Hostgroup name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Host to add. Will be replaced by hosts.",
				DeprecationMessage:  "use hosts instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostgroup": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Hostgroup to add. Will be replaced by hostgroups.",
				DeprecationMessage:  "use hostgroups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hosts": schema.ListAttribute{
				MarkdownDescription: "Hosts to add as hostgroup members",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"hostgroups": schema.ListAttribute{
				MarkdownDescription: "Hostgroups to add as hostgroup members",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple hostgroup membership resources on the same hostgroup. Manadatory for using hosts/hostgroups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *HostGroupMembership) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostGroupMembership) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostGroupMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create hostgroup membership %s", data.Id.ValueString()))
	optArgs := ipa.HostgroupAddMemberOptionalArgs{}

	args := ipa.HostgroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.Host.IsNull() {
		v := []string{string(data.Host.ValueString())}
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create hostgroup membership host %s", data.Host.ValueString()))
		optArgs.Host = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/h/%s", data.Name.ValueString(), data.Host.ValueString()))

	}
	if !data.HostGroup.IsNull() {
		v := []string{string(data.HostGroup.ValueString())}
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create hostgroup membership hostgroup %s", data.HostGroup.ValueString()))
		optArgs.Hostgroup = &v
		data.Id = types.StringValue(fmt.Sprintf("%s/hg/%s", data.Name.ValueString(), data.HostGroup.ValueString()))
	}
	if !data.Hosts.IsNull() || !data.HostGroups.IsNull() {
		if !data.Hosts.IsNull() {
			var v []string
			for _, value := range data.Hosts.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Host = &v
		}
		if !data.HostGroups.IsNull() {
			var v []string
			for _, value := range data.HostGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				if val == data.Name.ValueString() {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hostgroup membership: %s cannot be membership of itself", data.Name.ValueString()))
					return
				}
				v = append(v, val)
			}
			optArgs.Hostgroup = &v
		}
		data.Id = types.StringValue(fmt.Sprintf("%s/m/%s", data.Name.ValueString(), data.Identifier.ValueString()))
	}

	_v, err := r.client.HostgroupAddMember(&args, &optArgs)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa hostgroup membership: %s", _v.String()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hostgroup membership: %s", err))
		return
	}
	if _v.Completed == 0 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hostgroup membership: %v", _v.Failed))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostGroupMembership) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostGroupMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name, typeId, userId, err := parseHostgroupMembershipID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("State Error", fmt.Sprintf("Unable to parse resource %s: %s", data.Id.ValueString(), err))
	}
	reqArgs := ipa.HostgroupShowArgs{
		Cn: name,
	}
	z := new(bool)
	*z = true
	optArgs := ipa.HostgroupShowOptionalArgs{
		All: z,
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read hostgroup membership %s optArgs %v", data.Id.ValueString(), optArgs))
	res, err := r.client.HostgroupShow(&reqArgs, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound (4001)") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading information on freeipa hostgroup %s: %s", name, err))
		return
	}

	switch typeId {
	case "hg":
		v := []string{userId}
		if res.Result.MemberHostgroup != nil {
			hostgroups := *res.Result.MemberHostgroup
			if slices.Contains(hostgroups, v[0]) {
				data.HostGroup = types.StringValue(v[0])
			} else {
				data.HostGroup = types.StringValue("")
				data.Id = types.StringValue("")
			}
		} else {
			resp.State.RemoveResource(ctx)
			return
		}
	case "h":
		v := []string{userId}
		if res.Result.MemberHost != nil {
			hosts := *res.Result.MemberHost
			if slices.Contains(hosts, v[0]) {
				data.Host = types.StringValue(v[0])
			} else {
				data.Host = types.StringValue("")
				data.Id = types.StringValue("")
			}
		} else {
			resp.State.RemoveResource(ctx)
			return
		}
	case "m":
		if !data.Hosts.IsNull() && res.Result.MemberHost != nil {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hosts %v", *res.Result.MemberHost))
			var changedVals []string
			for _, value := range data.Hosts.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hosts failed with error %s", err))
				}
				if slices.Contains(*res.Result.MemberHost, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hosts %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.Hosts, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if !data.HostGroups.IsNull() && res.Result.MemberHostgroup != nil {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hostgroups %v", *res.Result.MemberHostgroup))
			var changedVals []string
			for _, value := range data.HostGroups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hostgroups failed with error %s", err))
				}
				if slices.Contains(*res.Result.MemberHostgroup, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hostgroup member hostgroups %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.HostGroups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if res.Result.MemberHostgroup == nil && res.Result.MemberHost == nil {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HostGroupMembership) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HostGroupMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.HostgroupAddMemberOptionalArgs{}

	memberAddArgs := ipa.HostgroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.HostgroupRemoveMemberOptionalArgs{}

	memberDelArgs := ipa.HostgroupRemoveMemberArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.Hosts.Equal(state.Hosts) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa hostgroup member hosts %s ", data.Hosts.String()))
		var statearr, planarr, addedHosts, deletedHosts []string

		for _, value := range state.Hosts.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.Hosts.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedHosts = append(addedHosts, val)
				memberAddOptArgs.Host = &addedHosts
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedHosts = append(deletedHosts, value)
				memberDelOptArgs.Host = &deletedHosts
				hasMemberDel = true
			}
		}

	}
	if !data.HostGroups.Equal(state.HostGroups) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa hostgroup member hostgroups %s ", data.HostGroups.String()))
		var statearr, planarr, addedGroups, deletedGroups []string

		for _, value := range state.HostGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.HostGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedGroups = append(addedGroups, val)
				memberAddOptArgs.Hostgroup = &addedGroups
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedGroups = append(deletedGroups, value)
				memberDelOptArgs.Hostgroup = &deletedGroups
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.HostgroupAddMember(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa hostgroup membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hostgroup membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hostgroup membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.HostgroupRemoveMember(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa hostgroup membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hostgroup membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hostgroup membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostGroupMembership) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostGroupMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HostgroupRemoveMemberOptionalArgs{}

	nameId, typeId, userId, err := parseHostgroupMembershipID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_host_hostgroup_membership %s: %s", data.Id.ValueString(), err))

	}

	args := ipa.HostgroupRemoveMemberArgs{
		Cn: nameId,
	}

	switch typeId {
	case "hg":
		v := []string{userId}
		optArgs.Hostgroup = &v
	case "h":
		v := []string{userId}
		optArgs.Host = &v
	case "m":
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa hostgroup member hosts %v ", data.Hosts))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa hostgroup member hostgroups %v ", data.HostGroups))
		if !data.Hosts.IsNull() {
			var v []string
			for _, value := range data.Hosts.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Host = &v
		}
		if !data.HostGroups.IsNull() {
			var v []string
			for _, value := range data.HostGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Hostgroup = &v
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	_, err = r.client.HostgroupRemoveMember(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error remove hostgroup membership %s: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *HostGroupMembership) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseHostgroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
