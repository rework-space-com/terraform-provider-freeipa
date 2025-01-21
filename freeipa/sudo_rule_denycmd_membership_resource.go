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
var _ resource.Resource = &SudoRuleDenyCmdMembershipResource{}
var _ resource.ResourceWithImportState = &SudoRuleDenyCmdMembershipResource{}

func NewSudoRuleDenyCmdMembershipResource() resource.Resource {
	return &SudoRuleDenyCmdMembershipResource{}
}

// SudoRuleDenyCmdMembershipResource defines the resource implementation.
type SudoRuleDenyCmdMembershipResource struct {
	client *ipa.Client
}

// SudoRuleDenyCmdMembershipResourceModel describes the resource data model.
type SudoRuleDenyCmdMembershipResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	SudoCmd       types.String `tfsdk:"sudocmd"`
	SudoCmds      types.List   `tfsdk:"sudocmds"`
	SudoCmdGroup  types.String `tfsdk:"sudocmd_group"`
	SudoCmdGroups types.List   `tfsdk:"sudocmd_groups"`
	Identifier    types.String `tfsdk:"identifier"`
}

func (r *SudoRuleDenyCmdMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule_denycmd_membership"
}

func (r *SudoRuleDenyCmdMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd"),
			path.MatchRoot("sudocmds"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd"),
			path.MatchRoot("sudocmd_group"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd"),
			path.MatchRoot("sudocmd_groups"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd_group"),
			path.MatchRoot("sudocmds"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd_group"),
			path.MatchRoot("sudocmd_groups"),
		),
	}
}

func (r *SudoRuleDenyCmdMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo rule deny command membership resource",

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
			"sudocmd": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Sudo command to deny by the sudo rule",
				DeprecationMessage:  "use sudocmds instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sudocmds": schema.ListAttribute{
				MarkdownDescription: "List of Sudo command to deny by the sudo rule",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"sudocmd_group": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Sudo command group to deny by the sudo rule",
				DeprecationMessage:  "use sudocmds instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sudocmd_groups": schema.ListAttribute{
				MarkdownDescription: "List of sudo command group to deny by the sudo rule",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple sudo rule allowed membership resources on the same sudo rule. Manadatory for using sudocmds/sudocmd_groups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SudoRuleDenyCmdMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoRuleDenyCmdMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoRuleDenyCmdMembershipResourceModel
	var id, cmd_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudoruleAddDenyCommandOptionalArgs{}

	args := ipa.SudoruleAddDenyCommandArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.SudoCmd.IsNull() {
		v := []string{data.SudoCmd.ValueString()}
		optArgs.Sudocmd = &v
		cmd_id = "srdc"
	}
	if !data.SudoCmdGroup.IsNull() {
		v := []string{data.SudoCmdGroup.ValueString()}
		optArgs.Sudocmdgroup = &v
		cmd_id = "srdcg"
	}
	if !data.SudoCmds.IsNull() || !data.SudoCmdGroups.IsNull() {
		if !data.SudoCmds.IsNull() {
			var v []string
			for _, value := range data.SudoCmds.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Sudocmd = &v
		}
		if !data.SudoCmdGroups.IsNull() {
			var v []string
			for _, value := range data.SudoCmdGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Sudocmdgroup = &v
		}
		cmd_id = "msrdc"
	}

	_, err := r.client.SudoruleAddDenyCommand(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule denied command membership: %s", err))
		return
	}

	switch cmd_id {
	case "srdc":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.SudoCmd.ValueString())
		data.Id = types.StringValue(id)
	case "srdcg":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.SudoCmdGroup.ValueString())
		data.Id = types.StringValue(id)
	case "msrdc":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), cmd_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleDenyCmdMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoRuleDenyCmdMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sudoruleId, typeId, cmdId, err := parseSudoRuleDenyCommandMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudorule_denycmd_membership: %s", err))
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
	case "srdc":
		if res.Result.MemberdenycmdSudocmd == nil || !slices.Contains(*res.Result.MemberdenycmdSudocmd, cmdId) {
			tflog.Debug(ctx, "[DEBUG] Sudo rule denied command membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "srdcg":
		if res.Result.MemberdenycmdSudocmdgroup == nil || !slices.Contains(*res.Result.MemberdenycmdSudocmdgroup, cmdId) {
			tflog.Debug(ctx, "[DEBUG] Sudo rule denied command membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "msrdc":
		if !data.SudoCmds.IsNull() {
			var changedVals []string
			for _, value := range data.SudoCmds.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands failed with error %s", err))
				}
				if res.Result.MemberdenycmdSudocmd != nil && slices.Contains(*res.Result.MemberdenycmdSudocmd, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.SudoCmds, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if !data.SudoCmdGroups.IsNull() {
			var changedVals []string
			for _, value := range data.SudoCmdGroups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands failed with error %s", err))
				}
				if res.Result.MemberdenycmdSudocmdgroup != nil && slices.Contains(*res.Result.MemberdenycmdSudocmdgroup, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.SudoCmdGroups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
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

func (r *SudoRuleDenyCmdMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoRuleDenyCmdMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.SudoruleAddDenyCommandOptionalArgs{}

	memberAddArgs := ipa.SudoruleAddDenyCommandArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.SudoruleRemoveDenyCommandOptionalArgs{}

	memberDelArgs := ipa.SudoruleRemoveDenyCommandArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.SudoCmds.Equal(state.SudoCmds) {
		var statearr, planarr, addedCmds, deletedCmds []string

		for _, value := range state.SudoCmds.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.SudoCmds.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedCmds = append(addedCmds, val)
				memberAddOptArgs.Sudocmd = &addedCmds
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedCmds = append(deletedCmds, value)
				memberDelOptArgs.Sudocmd = &deletedCmds
				hasMemberDel = true
			}
		}

	}
	if !data.SudoCmdGroups.Equal(state.SudoCmdGroups) {
		var statearr, planarr, addedCmdGrps, deletedCmdGrps []string

		for _, value := range state.SudoCmdGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.SudoCmdGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedCmdGrps = append(addedCmdGrps, val)
				memberAddOptArgs.Sudocmdgroup = &addedCmdGrps
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedCmdGrps = append(deletedCmdGrps, value)
				memberDelOptArgs.Sudocmdgroup = &deletedCmdGrps
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.SudoruleAddDenyCommand(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa sudo rule deny command membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule deny command membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule deny command membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.SudoruleRemoveDenyCommand(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa sudo deny command membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule deny command membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule deny command membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleDenyCmdMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoRuleDenyCmdMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	cmdgrpId, typeId, _, err := parseSudoRuleDenyCommandMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_denycmd_membership: %s", err))
		return
	}

	optArgs := ipa.SudoruleRemoveDenyCommandOptionalArgs{}

	args := ipa.SudoruleRemoveDenyCommandArgs{
		Cn: cmdgrpId,
	}

	switch typeId {
	case "srdc":
		v := []string{data.SudoCmd.ValueString()}
		optArgs.Sudocmd = &v
	case "srdcg":
		v := []string{data.SudoCmdGroup.ValueString()}
		optArgs.Sudocmdgroup = &v
	case "msrdc":
		if !data.SudoCmds.IsNull() {
			var v []string
			for _, value := range data.SudoCmds.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Sudocmd = &v
		}
		if !data.SudoCmdGroups.IsNull() {
			var v []string
			for _, value := range data.SudoCmdGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Sudocmdgroup = &v
		}
	}

	_, err = r.client.SudoruleRemoveDenyCommand(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo command group membership: %s", err))
		return
	}
}

func (r *SudoRuleDenyCmdMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseSudoRuleDenyCommandMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine sudo rule denied command membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	sudocmd := idParts[2]

	return name, _type, sudocmd, nil
}
