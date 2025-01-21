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
var _ resource.Resource = &SudoRuleRunAsUserMembershipResource{}
var _ resource.ResourceWithImportState = &SudoRuleRunAsUserMembershipResource{}

func NewSudoRuleRunAsUserMembershipResource() resource.Resource {
	return &SudoRuleRunAsUserMembershipResource{}
}

// SudoRuleRunAsUserMembershipResource defines the resource implementation.
type SudoRuleRunAsUserMembershipResource struct {
	client *ipa.Client
}

// SudoRuleRunAsUserMembershipResourceModel describes the resource data model.
type SudoRuleRunAsUserMembershipResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	RunAsUser  types.String `tfsdk:"runasuser"`
	RunAsUsers types.List   `tfsdk:"runasusers"`
	Identifier types.String `tfsdk:"identifier"`
}

func (r *SudoRuleRunAsUserMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sudo_rule_runasuser_membership"
}

func (r *SudoRuleRunAsUserMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("runasuser"),
			path.MatchRoot("runasusers"),
		),
	}
}

func (r *SudoRuleRunAsUserMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Sudo rule run as user membership resource",

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
			"runasuser": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Run As User to add to the sudo rule. Can be an external user (local user of ipa clients)",
				DeprecationMessage:  "use runasusers instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runasusers": schema.ListAttribute{
				MarkdownDescription: "List of Run As User to add to the sudo rule. Can be an external user (local user of ipa clients)",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple sudo rule runasuser membership resources on the same sudo rule. Manadatory for using runasusers configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SudoRuleRunAsUserMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SudoRuleRunAsUserMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SudoRuleRunAsUserMembershipResourceModel
	var id, usr_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.SudoruleAddRunasuserOptionalArgs{}

	args := ipa.SudoruleAddRunasuserArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.RunAsUser.IsNull() {
		v := []string{data.RunAsUser.ValueString()}
		optArgs.User = &v
		usr_id = "srrau"
	}
	if !data.RunAsUsers.IsNull() {
		var v []string
		for _, value := range data.RunAsUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.User = &v
		usr_id = "msrrau"
	}

	_, err := r.client.SudoruleAddRunasuser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasuser membership: %s", err))
		return
	}

	switch usr_id {
	case "srrau":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), usr_id, data.RunAsUser.ValueString())
		data.Id = types.StringValue(id)
	case "msrrau":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), usr_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleRunAsUserMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SudoRuleRunAsUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sudoruleId, typeId, usrId, err := parseSudoRuleRunAsUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_runasuser_membership: %s", err))
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
	case "srrau":
		if res.Result.IpasudorunasUser == nil || !slices.Contains(*res.Result.IpasudorunasUser, usrId) {
			resp.State.RemoveResource(ctx)
			return
		}
	case "msrrau":
		if !data.RunAsUsers.IsNull() {
			if res.Result.IpasudorunasUser != nil {
				var changedVals []string
				for _, value := range data.RunAsUsers.Elements() {
					val, err := strconv.Unquote(value.String())
					if err != nil {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo runasgroup membership failed with error %s", err))
					}
					if slices.Contains(*res.Result.IpasudorunasUser, val) {
						tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa sudo runasgroup membership %s is present in results", val))
						changedVals = append(changedVals, val)
					}
				}
				var diag diag.Diagnostics
				data.RunAsUsers, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
				if diag.HasError() {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
				}
			} else {
				var diag diag.Diagnostics
				data.RunAsUsers, diag = types.ListValueFrom(ctx, types.StringType, &[]string{})
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

func (r *SudoRuleRunAsUserMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state SudoRuleRunAsUserMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.SudoruleAddRunasuserOptionalArgs{}

	memberAddArgs := ipa.SudoruleAddRunasuserArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.SudoruleRemoveRunasuserOptionalArgs{}

	memberDelArgs := ipa.SudoruleRemoveRunasuserArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.RunAsUsers.Equal(state.RunAsUsers) {
		var statearr, planarr, addedRau, deletedRau []string

		for _, value := range state.RunAsUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.RunAsUsers.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedRau = append(addedRau, val)
				memberAddOptArgs.User = &addedRau
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedRau = append(deletedRau, value)
				memberDelOptArgs.User = &deletedRau
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.SudoruleAddRunasuser(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa sudo rule runasgroup membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasuser membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule runasuser membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.SudoruleRemoveRunasuser(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa sudo command group membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule runasuser membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa sudo rule runasuser membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SudoRuleRunAsUserMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SudoRuleRunAsUserMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	cmdusrId, typeId, _, err := parseSudoRuleRunAsUserMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_sudo_rule_runasuser_membership: %s", err))
		return
	}

	optArgs := ipa.SudoruleRemoveRunasuserOptionalArgs{}

	args := ipa.SudoruleRemoveRunasuserArgs{
		Cn: cmdusrId,
	}

	switch typeId {
	case "srrau":
		v := []string{data.RunAsUser.ValueString()}
		optArgs.User = &v
	case "msrrau":
		if !data.RunAsUsers.IsNull() {
			var v []string
			for _, value := range data.RunAsUsers.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.User = &v
		}
	}

	_, err = r.client.SudoruleRemoveRunasuser(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa sudo runasuser membership: %s", err))
		return
	}
}

func (r *SudoRuleRunAsUserMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseSudoRuleRunAsUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine sudo rule runasuser membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
