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
var _ resource.Resource = &SudoCmdGroupMembershipResource{}
var _ resource.ResourceWithImportState = &SudoCmdGroupMembershipResource{}

func NewSudoCmdGroupMembershipResource() resource.Resource {
	return &SudoCmdGroupMembershipResource{}
}

// SudoCmdGroupMembershipResource defines the resource implementation.
type SudoCmdGroupMembershipResource struct {
	client *ipa.Client
}

// SudoCmdGroupMembershipResourceModel describes the resource data model.
type SudoCmdGroupMembershipResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	SudoCmd    types.String `tfsdk:"sudocmd"`
	SudoCmds   types.List   `tfsdk:"sudocmds"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *SudoCmdGroupMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_cmdgroup_membership"
}

func (r *SudoCmdGroupMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("sudocmd"),
			path.MatchRoot("sudocmds"),
		),
	}
}

func (r *SudoCmdGroupMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo command group membership resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the sudo command group",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sudocmd": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Sudo command to add as a member",
				DeprecationMessage:  "use sudocmds instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sudocmds": schema.ListAttribute{
				MarkdownDescription: "List of sudo command to add as a member",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple sudo command group membership resources on the same sudo command group. Manadatory for using sudocmds configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SudoCmdGroupMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoCmdGroupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoCmdGroupMembershipResourceModel
	var id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudocmdgroupAddMemberOptionalArgs{}

	args := ipa.SudocmdgroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}

	if !data.SudoCmd.IsNull() {
		v := []string{data.SudoCmd.ValueString()}
		optArgs.Sudocmd = &v
		id = fmt.Sprintf("%s/sc/%s", encodeSlash(data.Name.ValueString()), data.SudoCmd.ValueString())
	}
	if !data.SudoCmds.IsNull() {
		var v []string
		for _, value := range data.SudoCmds.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Sudocmd = &v
		id = fmt.Sprintf("%s/msc/%s", encodeSlash(data.Name.ValueString()), data.Identifier.ValueString())
	}
	_, err := r.client.SudocmdgroupAddMember(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo command group membership: %s", err))
	}
	data.Id = types.StringValue(id)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoCmdGroupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoCmdGroupMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cmdgrpId, typeId, cmdId, err := parseSudocmdgroupMembershipID(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err))
		return
	}

	all := true
	optArgs := ipa.SudocmdgroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.SudocmdgroupShowArgs{
		Cn: cmdgrpId,
	}

	res, err := r.client.SudocmdgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] Sudo command group not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa Sudo command group: %s", err))
			return
		}
	}

	switch typeId {
	case "sc":
		if res.Result.MemberSudocmd == nil || !slices.Contains(*res.Result.MemberSudocmd, cmdId) {
			tflog.Debug(ctx, "[DEBUG] Sudo command group membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "msc":
		if !data.SudoCmds.IsNull() {
			var changedVals []string
			for _, value := range data.SudoCmds.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group member commands failed with error %s", err))
				}
				if res.Result.MemberSudocmd != nil && slices.Contains(*res.Result.MemberSudocmd, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo command group member commands %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.SudoCmds, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		} else {
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

func (r *SudoCmdGroupMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoCmdGroupMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.SudocmdgroupAddMemberOptionalArgs{}

	memberAddArgs := ipa.SudocmdgroupAddMemberArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.SudocmdgroupRemoveMemberOptionalArgs{}

	memberDelArgs := ipa.SudocmdgroupRemoveMemberArgs{
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
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.SudocmdgroupAddMember(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa sudo command group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo command group membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo command group membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.SudocmdgroupRemoveMember(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa sudo command group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo command group membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo command group membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoCmdGroupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoCmdGroupMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	cmdgrpId, typeId, cmdId, err := parseSudocmdgroupMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudocmdgroup_membership: %s", err))
		return
	}

	optArgs := ipa.SudocmdgroupRemoveMemberOptionalArgs{}

	args := ipa.SudocmdgroupRemoveMemberArgs{
		Cn: cmdgrpId,
	}

	switch typeId {
	case "sc":
		v := []string{cmdId}
		optArgs.Sudocmd = &v
	case "msc":
		if !data.SudoCmds.IsNull() {
			var v []string
			for _, value := range data.SudoCmds.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Sudocmd = &v
		}
	}

	_, err = r.client.SudocmdgroupRemoveMember(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo command group membership: %s", err))
		return
	}
}

func (r *SudoCmdGroupMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseSudocmdgroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine sudo command group membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	sudocmd := idParts[2]

	return name, _type, sudocmd, nil
}
