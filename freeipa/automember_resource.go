// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AutomemberResource{}
var _ resource.ResourceWithImportState = &AutomemberResource{}

func NewAutomemberResource() resource.Resource {
	return &AutomemberResource{}
}

// AutomemberResource defines the resource implementation.
type AutomemberResource struct {
	client *ipa.Client
}

// AutomemberResourceModel describes the resource data model.
type AutomemberResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	AddAttr     types.List   `tfsdk:"addattr"`
	SetAttr     types.List   `tfsdk:"setattr"`
}

func (r *AutomemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_automemberadd"
}

func (r *AutomemberResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *AutomemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Automember resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Automember rule name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Automember rule description",
				Optional:            true,
				Computed:            false,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Automember rule type",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"addattr": schema.ListAttribute{
				MarkdownDescription: "Add an attribute/value pair. Format is attr=value. The attribute must be part of the schema.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"setattr": schema.ListAttribute{
				MarkdownDescription: "Set an attribute to a name/value pair. Format is attr=value.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *AutomemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AutomemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AutomemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.AutomemberAddOptionalArgs{}

	args := ipa.AutomemberAddArgs{
		Cn:   data.Name.ValueString(),
		Type: data.Type.ValueString(),
	}
	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if len(data.AddAttr.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa automember rule Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if len(data.SetAttr.Elements()) > 0 {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}
	_, err := r.client.AutomemberAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa automember rule: %s", err))
		return
	}

	tflog.Trace(ctx, "created a automember rule resource")

	data.Id = types.StringValue(data.Name.ValueString())

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AutomemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AutomemberResourceModel

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

func (r *AutomemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state AutomemberResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

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
	args := ipa.AutomemberModArgs{
		Cn:   data.Id.ValueString(),
		Type: data.Type.ValueString(),
	}
	optArgs := ipa.AutomemberModOptionalArgs{}

	if !data.Description.Equal(state.Description) {
		optArgs.Description = data.Description.ValueStringPointer()
	}

	if !data.AddAttr.Equal(state.AddAttr) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa automember rule Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if !data.SetAttr.Equal(state.SetAttr) {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}

	res, err := r.client.AutomemberMod(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa automember rule %s: %s", res.Result.Cn, err))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AutomemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AutomemberResourceModel

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
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa automember rule Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa automember rule Name %s", data.Name.ValueString()))
	args := ipa.AutomemberDelArgs{
		Cn:   []string{data.Name.ValueString()},
		Type: data.Type.ValueString(),
	}
	_, err := r.client.AutomemberDel(&args, &ipa.AutomemberDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Automember rule %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *AutomemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
