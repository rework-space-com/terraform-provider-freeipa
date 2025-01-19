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
var _ resource.Resource = &SudoRuleRunAsGroupMembershipResource{}
var _ resource.ResourceWithImportState = &SudoRuleRunAsGroupMembershipResource{}

func NewSudoRuleRunAsGroupMembershipResource() resource.Resource {
	return &SudoRuleRunAsGroupMembershipResource{}
}

// SudoRuleRunAsGroupMembershipResource defines the resource implementation.
type SudoRuleRunAsGroupMembershipResource struct {
	client *ipa.Client
}

// SudoRuleRunAsGroupMembershipResourceModel describes the resource data model.
type SudoRuleRunAsGroupMembershipResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	RunAsGroup  types.String `tfsdk:"runasgroup"`
	RunAsGroups types.List   `tfsdk:"runasgroups"`
	Identifier  types.String `tfsdk:"identifier"`
}

func (r *SudoRuleRunAsGroupMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule_runasgroup_membership"
}

func (r *SudoRuleRunAsGroupMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("runasgroup"),
			path.MatchRoot("runasgroups"),
		),
	}
}

func (r *SudoRuleRunAsGroupMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo rule run as group membership resource",

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
			"runasgroup": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Run As Group to add to the sudo rule. Can be an external group (local group of ipa clients)",
				DeprecationMessage:  "use runasgroups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runasgroups": schema.ListAttribute{
				MarkdownDescription: "List of Run As Group to add to the sudo rule. Can be an external group (local group of ipa clients)",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple sudo rule runasgroup membership resources on the same sudo rule. Manadatory for using runasgroups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SudoRuleRunAsGroupMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoRuleRunAsGroupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoRuleRunAsGroupMembershipResourceModel
	var id, grp_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudoruleAddRunasgroupOptionalArgs{}

	args := ipa.SudoruleAddRunasgroupArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.RunAsGroup.IsNull() {
		v := []string{data.RunAsGroup.ValueString()}
		optArgs.Group = &v
		grp_id = "srraug"
	}
	if !data.RunAsGroups.IsNull() {
		var v []string
		for _, value := range data.RunAsGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Group = &v
		grp_id = "msrraug"
	}

	_, err := r.client.SudoruleAddRunasgroup(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasgroup membership: %s", err))
		return
	}

	switch grp_id {
	case "srraug":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), grp_id, data.RunAsGroup.ValueString())
		data.Id = types.StringValue(id)
	case "msrraug":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), grp_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleRunAsGroupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoRuleRunAsGroupMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sudoruleId, typeId, grpId, err := parseSudoRuleRunAsGroupMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_runasgroup_membership: %s", err))
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
	case "srraug":
		if res.Result.IpasudorunasgroupGroup == nil || !slices.Contains(*res.Result.IpasudorunasgroupGroup, grpId) {
			tflog.Debug(ctx, "[DEBUG] Sudo rule runasgroup membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "msrraug":
		if !data.RunAsGroups.IsNull() {
			var changedVals []string
			for _, value := range data.RunAsGroups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands failed with error %s", err))
				}
				if res.Result.IpasudorunasgroupGroup != nil && slices.Contains(*res.Result.IpasudorunasgroupGroup, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.RunAsGroups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
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

func (r *SudoRuleRunAsGroupMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoRuleRunAsGroupMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.SudoruleAddRunasgroupOptionalArgs{}

	memberAddArgs := ipa.SudoruleAddRunasgroupArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.SudoruleRemoveRunasgroupOptionalArgs{}

	memberDelArgs := ipa.SudoruleRemoveRunasgroupArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.RunAsGroups.Equal(state.RunAsGroups) {
		var statearr, planarr, addedRag, deletedRag []string

		for _, value := range state.RunAsGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.RunAsGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedRag = append(addedRag, val)
				memberAddOptArgs.Group = &addedRag
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedRag = append(deletedRag, value)
				memberDelOptArgs.Group = &deletedRag
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.SudoruleAddRunasgroup(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa sudo rule runasgroup membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasgroup membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasgroup membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.SudoruleRemoveRunasgroup(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa sudo command group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule runasgroup membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule runasgroup membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleRunAsGroupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoRuleRunAsGroupMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	cmdgrpId, typeId, _, err := parseSudoRuleRunAsGroupMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_runasgroup_membership: %s", err))
		return
	}

	optArgs := ipa.SudoruleRemoveRunasgroupOptionalArgs{}

	args := ipa.SudoruleRemoveRunasgroupArgs{
		Cn: cmdgrpId,
	}

	switch typeId {
	case "srraug":
		v := []string{data.RunAsGroup.ValueString()}
		optArgs.Group = &v
	case "msrraug":
		if !data.RunAsGroups.IsNull() {
			var v []string
			for _, value := range data.RunAsGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Group = &v
		}
	}

	_, err = r.client.SudoruleRemoveRunasgroup(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo runasgroup membership: %s", err))
		return
	}
}

func (r *SudoRuleRunAsGroupMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseSudoRuleRunAsGroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine sudo rule runasgroup membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	group := idParts[2]

	return name, _type, group, nil
}
