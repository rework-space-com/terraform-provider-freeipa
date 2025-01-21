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
var _ resource.Resource = &HbacPolicyServiceMembershipResource{}
var _ resource.ResourceWithImportState = &HbacPolicyServiceMembershipResource{}

func NewHbacPolicyServiceMembershipResource() resource.Resource {
	return &HbacPolicyServiceMembershipResource{}
}

// HbacPolicyServiceMembershipResource defines the resource implementation.
type HbacPolicyServiceMembershipResource struct {
	client *ipa.Client
}

// HbacPolicyServiceMembershipResourceModel describes the resource data model.
type HbacPolicyServiceMembershipResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Service       types.String `tfsdk:"service"`
	Services      types.List   `tfsdk:"services"`
	ServiceGroup  types.String `tfsdk:"servicegroup"`
	ServiceGroups types.List   `tfsdk:"servicegroups"`
	Identifier    types.String `tfsdk:"identifier"`
}

func (r *HbacPolicyServiceMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hbac_policy_service_membership"
}

func (r *HbacPolicyServiceMembershipResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("service"),
			path.MatchRoot("services"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("service"),
			path.MatchRoot("servicegroup"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("service"),
			path.MatchRoot("servicegroups"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("servicegroup"),
			path.MatchRoot("services"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("servicegroup"),
			path.MatchRoot("servicegroups"),
		),
	}
}

func (r *HbacPolicyServiceMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA HBAC policy service membership resource",

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
			"service": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Service name the policy is applied t",
				DeprecationMessage:  "use services instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"services": schema.ListAttribute{
				MarkdownDescription: "List of service name the policy is applied t",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"servicegroup": schema.StringAttribute{
				MarkdownDescription: "**deprecated** Service group name the policy is applied to",
				DeprecationMessage:  "use servicegroups instead",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"servicegroups": schema.ListAttribute{
				MarkdownDescription: "List of service group name the policy is applied to",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate multiple HBAC policy service membership resources on the same HBAC policy. Manadatory for using services/servicegroups configurations.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *HbacPolicyServiceMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HbacPolicyServiceMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HbacPolicyServiceMembershipResourceModel
	var id, user_id string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HbacruleAddServiceOptionalArgs{}

	args := ipa.HbacruleAddServiceArgs{
		Cn: data.Name.ValueString(),
	}
	if !data.Service.IsNull() {
		v := []string{data.Service.ValueString()}
		optArgs.Hbacsvc = &v
		user_id = "s"
	}
	if !data.ServiceGroup.IsNull() {
		v := []string{data.ServiceGroup.ValueString()}
		optArgs.Hbacsvcgroup = &v
		user_id = "sg"
	}
	if !data.Services.IsNull() || !data.ServiceGroups.IsNull() {
		if !data.Services.IsNull() {
			var v []string
			for _, value := range data.Services.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Hbacsvc = &v
		}
		if !data.ServiceGroups.IsNull() {
			var v []string
			for _, value := range data.ServiceGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Hbacsvcgroup = &v
		}
		user_id = "ms"
	}

	_, err := r.client.HbacruleAddService(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa sudo rule service membership: %s", err))
		return
	}

	switch user_id {
	case "s":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.Service.ValueString())
		data.Id = types.StringValue(id)
	case "sg":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.ServiceGroup.ValueString())
		data.Id = types.StringValue(id)
	case "ms":
		id = fmt.Sprintf("%s/%s/%s", encodeSlash(data.Name.ValueString()), user_id, data.Identifier.ValueString())
		data.Id = types.StringValue(id)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyServiceMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HbacPolicyServiceMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hbacpolicyid, typeId, policyId, err := parseHBACPolicyServiceMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_service_membership: %s", err))
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
	case "s":
		if res.Result.MemberserviceHbacsvc == nil || !slices.Contains(*res.Result.MemberserviceHbacsvc, policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy service membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "sg":
		if res.Result.MemberserviceHbacsvcgroup == nil || !slices.Contains(*res.Result.MemberserviceHbacsvcgroup, policyId) {
			tflog.Debug(ctx, "[DEBUG] HBAC policy service group membership does not exist")
			resp.State.RemoveResource(ctx)
			return
		}
	case "ms":
		if !data.Services.IsNull() {
			var changedVals []string
			for _, value := range data.Services.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy service member failed with error %s", err))
				}
				if res.Result.MemberserviceHbacsvc != nil && slices.Contains(*res.Result.MemberserviceHbacsvc, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy service member %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.Services, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
			if diag.HasError() {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
			}
		}
		if !data.ServiceGroups.IsNull() {
			var changedVals []string
			for _, value := range data.ServiceGroups.Elements() {
				val, err := strconv.Unquote(value.String())
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy service member failed with error %s", err))
				}
				if res.Result.MemberserviceHbacsvcgroup != nil && slices.Contains(*res.Result.MemberserviceHbacsvcgroup, val) {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa hbac policy service member %s is present in results", val))
					changedVals = append(changedVals, val)
				}
			}
			var diag diag.Diagnostics
			data.ServiceGroups, diag = types.ListValueFrom(ctx, types.StringType, &changedVals)
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

func (r *HbacPolicyServiceMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HbacPolicyServiceMembershipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	memberAddOptArgs := ipa.HbacruleAddServiceOptionalArgs{}

	memberAddArgs := ipa.HbacruleAddServiceArgs{
		Cn: data.Name.ValueString(),
	}

	memberDelOptArgs := ipa.HbacruleRemoveServiceOptionalArgs{}

	memberDelArgs := ipa.HbacruleRemoveServiceArgs{
		Cn: data.Name.ValueString(),
	}
	hasMemberAdd := false
	hasMemberDel := false
	// Memberships can be added or removed, comparing the current state and the plan allows us to define 2 lists of members to add or remove.
	if !data.Services.Equal(state.Services) {
		var statearr, planarr, addedSvc, deletedSvc []string

		for _, value := range state.Services.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.Services.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedSvc = append(addedSvc, val)
				memberAddOptArgs.Hbacsvc = &addedSvc
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedSvc = append(deletedSvc, value)
				memberDelOptArgs.Hbacsvc = &deletedSvc
				hasMemberDel = true
			}
		}

	}
	if !data.ServiceGroups.Equal(state.ServiceGroups) {
		var statearr, planarr, addedGroups, deletedGroups []string

		for _, value := range state.ServiceGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			statearr = append(statearr, val)
		}
		for _, value := range data.ServiceGroups.Elements() {
			val, _ := strconv.Unquote(value.String())
			planarr = append(planarr, val)
			if !slices.Contains(statearr, val) {
				addedGroups = append(addedGroups, val)
				memberAddOptArgs.Hbacsvcgroup = &addedGroups
				hasMemberAdd = true
			}
		}
		for _, value := range statearr {
			if !slices.Contains(planarr, value) {
				deletedGroups = append(deletedGroups, value)
				memberDelOptArgs.Hbacsvcgroup = &deletedGroups
				hasMemberDel = true
			}
		}

	}
	// The api provides a add and a remove function for membership. Therefore we need to call the right one when appropriate.
	if hasMemberAdd {
		_v, err := r.client.HbacruleAddService(&memberAddArgs, &memberAddOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error creating freeipa hbac policy service membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy service membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa hbac policy service membership: %v", _v.Failed))
			return
		}
	}
	if hasMemberDel {
		_v, err := r.client.HbacruleRemoveService(&memberDelArgs, &memberDelOptArgs)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Error removing freeipa hbac policy service membership: %s", _v.String()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy service membership: %s", err))
			return
		}
		if _v.Completed == 0 {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error removing freeipa hbac policy service membership: %v", _v.Failed))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HbacPolicyServiceMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HbacPolicyServiceMembershipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	hbacpolicyId, typeId, _, err := parseHBACPolicyServiceMembershipID(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error parsing ID of freeipa_hbac_policy_user_membership: %s", err))
		return
	}

	optArgs := ipa.HbacruleRemoveServiceOptionalArgs{}

	args := ipa.HbacruleRemoveServiceArgs{
		Cn: hbacpolicyId,
	}

	switch typeId {
	case "s":
		v := []string{data.Service.ValueString()}
		optArgs.Hbacsvc = &v
	case "sg":
		v := []string{data.ServiceGroup.ValueString()}
		optArgs.Hbacsvcgroup = &v
	case "ms":
		if !data.Services.IsNull() {
			var v []string
			for _, value := range data.Services.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Hbacsvc = &v
		}
		if !data.ServiceGroups.IsNull() {
			var v []string
			for _, value := range data.ServiceGroups.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			optArgs.Hbacsvcgroup = &v
		}
	}

	_, err = r.client.HbacruleRemoveService(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa hbac policy servie membership: %s", err))
		return
	}
}

func (r *HbacPolicyServiceMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func parseHBACPolicyServiceMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine service membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	svc := idParts[2]

	return name, _type, svc, nil
}
