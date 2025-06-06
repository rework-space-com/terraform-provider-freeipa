// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//   Mixton <maxime.thomas@mtconsulting.tech>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AutomemberConditionResource{}

func NewAutomemberConditionResource() resource.Resource {
	return &AutomemberConditionResource{}
}

// AutomemberConditionResource defines the resource implementation.
type AutomemberConditionResource struct {
	client *ipa.Client
}

// AutomemberConditionResourceModel describes the resource data model.
type AutomemberConditionResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Key            types.String `tfsdk:"key"`
	InclusiveRegex types.List   `tfsdk:"inclusiveregex"`
	ExclusiveRegex types.List   `tfsdk:"exclusiveregex"`
}

func (r *AutomemberConditionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automemberadd_condition"
}

func (r *AutomemberConditionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *AutomemberConditionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Automember conditionresource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Automember rule condition name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Automember rule condition description",
				Optional:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Automember rule condition type",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Automember rule condition key",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"inclusiveregex": schema.ListAttribute{
				MarkdownDescription: "Regex expression for values that should be included.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"exclusiveregex": schema.ListAttribute{
				MarkdownDescription: "Regex expression for values that should be excluded.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AutomemberConditionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AutomemberConditionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AutomemberConditionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.AutomemberAddConditionOptionalArgs{}

	args := ipa.AutomemberAddConditionArgs{
		Cn:   data.Name.ValueString(),
		Key:  data.Key.ValueString(),
		Type: data.Type.ValueString(),
	}
	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if len(data.InclusiveRegex.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa automember rule condition InclusiveRegex %s ", data.InclusiveRegex.String()))
		var v []string

		for _, value := range data.InclusiveRegex.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Automemberinclusiveregex = &v
	}

	if len(data.ExclusiveRegex.Elements()) > 0 {
		var v []string
		for _, value := range data.ExclusiveRegex.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Automemberexclusiveregex = &v
	}
	_, err := r.client.AutomemberAddCondition(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa automember rule condition : %s", err))
		return
	}

	tflog.Trace(ctx, "created a automember rule condition resource")

	data.Id = types.StringValue(data.Name.ValueString())

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AutomemberConditionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AutomemberConditionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.AutomemberShowOptionalArgs{
		All: &all,
	}

	args := ipa.AutomemberShowArgs{
		Cn:   data.Name.ValueString(),
		Type: data.Type.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa automember rule %s", data.Id.ValueString()))
	res, err := r.client.AutomemberShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Automember rule %s not found", data.Id.ValueString()))
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Automember rule %s not found: %s", data.Id.ValueString(), err))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa Automember rule %v", res))
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa automember rule %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa automember rule %s", data.Name.ValueString()))
		return
	}

	data.Name = types.StringValue(res.Result.Cn)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa automember rule Cn %s", data.Name.ValueString()))
	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa automember rule Description %s", data.Description.ValueString()))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AutomemberConditionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state AutomemberConditionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AutomemberConditionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AutomemberConditionResourceModel

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
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa automember rule condition Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa automember rule condition Name %s", data.Name.ValueString()))
	args := ipa.AutomemberRemoveConditionArgs{
		Cn:   data.Name.ValueString(),
		Key:  data.Key.ValueString(),
		Type: data.Type.ValueString(),
	}
	optArgs := ipa.AutomemberRemoveConditionOptionalArgs{}
	if !data.InclusiveRegex.IsNull() {
		var v []string
		for _, value := range data.InclusiveRegex.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Automemberinclusiveregex = &v
	}
	if !data.ExclusiveRegex.IsNull() {
		var v []string
		for _, value := range data.ExclusiveRegex.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Automemberexclusiveregex = &v
	}
	_, err := r.client.AutomemberRemoveCondition(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Automember rule %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}
}
