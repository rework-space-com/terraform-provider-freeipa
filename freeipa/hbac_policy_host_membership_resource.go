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
var _ resource.Resource = &HbacPolicyHostMembershipResource{}
var _ resource.ResourceWithImportState = &HbacPolicyHostMembershipResource{}

func NewHbacPolicyHostMembershipResource() resource.Resource {
	return &HbacPolicyHostMembershipResource{}
}

// HbacPolicyHostMembershipResource defines the resource implementation.
type HbacPolicyHostMembershipResource struct {
	client *ipa.Client
}

// HbacPolicyHostMembershipResourceModel describes the resource data model.
type HbacPolicyHostMembershipResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Host       types.String `tfsdk:"host"`
	Hosts      types.List   `tfsdk:"hosts"`
	HostGroup  types.String `tfsdk:"hostgroup"`
	HostGroups types.List   `tfsdk:"hostgroups"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *HbacPolicyHostMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hbac_policy_host_membership"
}

func (r *HbacPolicyHostMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hosts"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hostgroup"),
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
	}
}

func (r *HbacPolicyHostMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"host": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Host to add to the HBAC policy",
				DeprecationMessage:  "use hosts instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hosts": schema.ListAttribute{
				MarkdownDescription: "List of hosts to add to the HBAC policy",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"hostgroup": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Hostgroup to add to the HBAC policy",
				DeprecationMessage:  "use hostgroups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostgroups": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups to add to the HBAC policy",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple HBAC policy host membership resources on the same HBAC policy. Manadatory for using hosts/hostgroups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *HbacPolicyHostMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HbacPolicyHostMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HbacPolicyHostMembershipResourceModel
	var id, cmd_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HbacruleAddHostOptionalArgs{}

	args := ipa.HbacruleAddHostArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.Host.IsNull() {
		v := []string{data.Host.ValueString()}
		optArgs.Host = &v
		cmd_id = "h"
	}
	if !data.HostGroup.IsNull() {
		v := []string{data.HostGroup.ValueString()}
		optArgs.Hostgroup = &v
		cmd_id = "hg"
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
				v = append(v, val)
			}
			optArgs.Hostgroup = &v
		}
		cmd_id = "mh"
	}

	_, err := r.client.HbacruleAddHost(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule host membership: %s", err))
		return
	}

	switch cmd_id {
	case "h":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.Host.ValueString())
		data.Id = types.StringValue(id)
	case "hg":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.HostGroup.ValueString())
		data.Id = types.StringValue(id)
	case "mh":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyHostMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HbacPolicyHostMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hbacpolicyid, typeId, policyId, err := parseHBACPolicyHostMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_host_membership: %s", err))
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
	case "h":
		if res.Result.MemberhostHost == nil || !isStringListContainsCaseInsensistive(res.Result.MemberhostHost, &policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy host membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "hg":
		if res.Result.MemberhostHostgroup == nil || !isStringListContainsCaseInsensistive(res.Result.MemberhostHostgroup, &policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy host group membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "mh":
		if !data.Hosts.IsNull() {
			var changedVals []string
			for _, value := range data.Hosts.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy host member failed with error %s", err))
				}
				if res.Result.MemberhostHost != nil && isStringListContainsCaseInsensistive(res.Result.MemberhostHost, &val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy host member %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.Hosts, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if !data.HostGroups.IsNull() {
			var changedVals []string
			for _, value := range data.HostGroups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy member commands failed with error %s", err))
				}
				if res.Result.MemberhostHostgroup != nil && isStringListContainsCaseInsensistive(res.Result.MemberhostHostgroup, &val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.HostGroups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
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

func (r *HbacPolicyHostMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HbacPolicyHostMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.HbacruleAddHostOptionalArgs{}

	memberAddArgs := ipa.HbacruleAddHostArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.HbacruleRemoveHostOptionalArgs{}

	memberDelArgs := ipa.HbacruleRemoveHostArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.Hosts.Equal(state.Hosts) {
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
		var statearr, planarr, addedCmdGrps, deletedCmdGrps []string

		for _, value := range state.HostGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.HostGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedCmdGrps = append(addedCmdGrps, val)
				memberAddOptArgs.Hostgroup = &addedCmdGrps
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedCmdGrps = append(deletedCmdGrps, value)
				memberDelOptArgs.Hostgroup = &deletedCmdGrps
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.HbacruleAddHost(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa hbac policy host membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy host membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy host membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.HbacruleRemoveHost(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa hbac policy host membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy host membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy host membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyHostMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HbacPolicyHostMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	hbacpolicyId, typeId, _, err := parseHBACPolicyHostMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_host_membership: %s", err))
		return
	}

	optArgs := ipa.HbacruleRemoveHostOptionalArgs{}

	args := ipa.HbacruleRemoveHostArgs{
		Cn: hbacpolicyId,
	}

	switch typeId {
	case "h":
		v := []string{data.Host.ValueString()}
		optArgs.Host = &v
	case "hg":
		v := []string{data.HostGroup.ValueString()}
		optArgs.Hostgroup = &v
	case "mh":
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

	_, err = r.client.HbacruleRemoveHost(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa hbac policy host membership: %s", err))
		return
	}
}

func (r *HbacPolicyHostMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseHBACPolicyHostMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine host membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	host := idParts[2]

	return name, _type, host, nil
}
