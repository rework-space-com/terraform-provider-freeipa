// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
var _ resource.Resource = &HostResource{}
var _ resource.ResourceWithImportState = &HostResource{}

func NewHostResource() resource.Resource {
	return &HostResource{}
}

// HostResource defines the resource implementation.
type HostResource struct {
	client *ipa.Client
}

// HostResourceModel describes the resource data model.
type HostResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	IpAddresses             types.String `tfsdk:"ip_address"`
	Description             types.String `tfsdk:"description"`
	Locality                types.String `tfsdk:"locality"`
	Location                types.String `tfsdk:"location"`
	Platform                types.String `tfsdk:"platform"`
	OperatingSystem         types.String `tfsdk:"operating_system"`
	UserCertificates        types.List   `tfsdk:"user_certificates"`
	MacAddresses            types.List   `tfsdk:"mac_addresses"`
	IpaSshPubKeys           types.List   `tfsdk:"ipasshpubkeys"`
	Userclass               types.List   `tfsdk:"userclass"`
	AssignedIdView          types.String `tfsdk:"assigned_idview"`
	KrbAuthIndicator        types.List   `tfsdk:"krb_auth_indicators"`
	KrbPreAuth              types.Bool   `tfsdk:"krb_preauth"`
	TrustedForDelegation    types.Bool   `tfsdk:"trusted_for_delegation"`
	TrustedToAuthAsDelegate types.Bool   `tfsdk:"trusted_to_auth_as_delegate"`
	Force                   types.Bool   `tfsdk:"force"`
	UserPassword            types.String `tfsdk:"userpassword"`
	RandomPassword          types.Bool   `tfsdk:"random_password"`
	GeneratedPassword       types.String `tfsdk:"generated_password"`
}

func (r *HostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *HostResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *HostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Host name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address of the host",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of this host",
				Optional:            true,
			},
			"locality": schema.StringAttribute{
				MarkdownDescription: "Host locality (e.g. 'Baltimore, MD')",
				Optional:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Host location (e.g. 'Lab 2')",
				Optional:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Host hardware platform (e.g. 'Lenovo T61')",
				Optional:            true,
			},
			"operating_system": schema.StringAttribute{
				MarkdownDescription: "Host operating system and version (e.g. 'Fedora 40')",
				Optional:            true,
			},
			"user_certificates": schema.ListAttribute{
				MarkdownDescription: "Base-64 encoded host certificate",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mac_addresses": schema.ListAttribute{
				MarkdownDescription: "Hardware MAC address(es) on this host",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"ipasshpubkeys": schema.ListAttribute{
				MarkdownDescription: "SSH public keys",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"userclass": schema.ListAttribute{
				MarkdownDescription: "Host category (semantics placed on this attribute are for local interpretation)",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"assigned_idview": schema.StringAttribute{
				MarkdownDescription: "Assigned ID View",
				Optional:            true,
			},
			"krb_auth_indicators": schema.ListAttribute{
				MarkdownDescription: "Defines a whitelist for Authentication Indicators. Use 'otp' to allow OTP-based 2FA authentications. Use 'radius' to allow RADIUS-based 2FA authentications. Other values may be used for custom configurations.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"krb_preauth": schema.BoolAttribute{
				MarkdownDescription: "Pre-authentication is required for the service",
				Optional:            true,
			},
			"trusted_for_delegation": schema.BoolAttribute{
				MarkdownDescription: "Client credentials may be delegated to the service",
				Optional:            true,
			},
			"trusted_to_auth_as_delegate": schema.BoolAttribute{
				MarkdownDescription: "The service is allowed to authenticate on behalf of a client",
				Optional:            true,
			},
			"force": schema.BoolAttribute{
				MarkdownDescription: "Skip host's DNS check (A/AAAA) before adding it",
				Optional:            true,
			},
			"userpassword": schema.StringAttribute{
				MarkdownDescription: "Password used in bulk enrollment",
				Optional:            true,
				Sensitive:           true,
			},
			"random_password": schema.BoolAttribute{
				MarkdownDescription: "Generate a random password to be used in bulk enrollment",
				Optional:            true,
			},
			"generated_password": schema.StringAttribute{
				MarkdownDescription: "Generated random password created at host creation",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *HostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HostAddOptionalArgs{}

	args := ipa.HostAddArgs{
		Fqdn: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		optArgs.Description = data.Description.ValueStringPointer()
	}
	if !data.IpAddresses.IsNull() {
		optArgs.IPAddress = data.IpAddresses.ValueStringPointer()
	}
	if !data.Locality.IsNull() {
		optArgs.L = data.Locality.ValueStringPointer()
	}
	if !data.Location.IsNull() {
		optArgs.Nshostlocation = data.Location.ValueStringPointer()
	}
	if !data.Platform.IsNull() {
		optArgs.Nshardwareplatform = data.Platform.ValueStringPointer()
	}
	if !data.OperatingSystem.IsNull() {
		optArgs.Nsosversion = data.OperatingSystem.ValueStringPointer()
	}
	if !data.UserCertificates.IsNull() {
		var v []interface{}
		for _, value := range data.UserCertificates.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Usercertificate = &v
	}
	if !data.MacAddresses.IsNull() {
		var v []string
		for _, value := range data.MacAddresses.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Macaddress = &v
	}
	if !data.IpaSshPubKeys.IsNull() {
		var v []string
		for _, value := range data.IpaSshPubKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if !data.Userclass.IsNull() {
		var v []string
		for _, value := range data.Userclass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}
	if !data.AssignedIdView.IsNull() {
		optArgs.Ipaassignedidview = data.AssignedIdView.ValueStringPointer()
	}
	if !data.KrbAuthIndicator.IsNull() {
		var v []string
		for _, value := range data.KrbAuthIndicator.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalauthind = &v
	}
	if !data.KrbPreAuth.IsNull() {
		optArgs.Ipakrbrequirespreauth = data.KrbPreAuth.ValueBoolPointer()
	}
	if !data.TrustedForDelegation.IsNull() {
		optArgs.Ipakrbokasdelegate = data.TrustedForDelegation.ValueBoolPointer()
	}
	if !data.TrustedToAuthAsDelegate.IsNull() {
		optArgs.Ipakrboktoauthasdelegate = data.TrustedToAuthAsDelegate.ValueBoolPointer()
	}
	if !data.RandomPassword.IsNull() {
		optArgs.Random = data.RandomPassword.ValueBoolPointer()
	}
	if !data.UserPassword.IsNull() {
		optArgs.Userpassword = data.UserPassword.ValueStringPointer()
	}
	if !data.Force.IsNull() {
		optArgs.Force = data.Force.ValueBoolPointer()
	}

	res, err := r.client.HostAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa host: %s", err))
		return
	}

	if !data.RandomPassword.IsNull() && data.RandomPassword.ValueBool() {
		data.GeneratedPassword = types.StringValue(*res.Result.Randompassword)
	} else {
		data.GeneratedPassword = types.StringValue("")
	}

	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	args := ipa.HostShowArgs{
		Fqdn: data.Name.ValueString(),
	}
	optArgs := ipa.HostShowOptionalArgs{
		All: &all,
	}

	res, err := r.client.HostShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] Host not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa host: %s", err))
			return
		}
	}
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa host %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa host %s", data.Name.ValueString()))
		return
	}

	if res.Result.Description != nil && !data.Description.IsNull() {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.L != nil && !data.Locality.IsNull() {
		data.Locality = types.StringValue(*res.Result.L)
	}
	if res.Result.Nshostlocation != nil && !data.Location.IsNull() {
		data.Location = types.StringValue(*res.Result.Nshostlocation)
	}
	if res.Result.Nshardwareplatform != nil && !data.Platform.IsNull() {
		data.Platform = types.StringValue(*res.Result.Nshardwareplatform)
	}

	if res.Result.Nsosversion != nil && !data.OperatingSystem.IsNull() {
		data.OperatingSystem = types.StringValue(*res.Result.Nsosversion)
	}
	if !data.UserCertificates.IsNull() && res.Result.Usercertificate != nil {
		var changedVals, resVals []string
		for _, v := range *res.Result.Usercertificate {
			resVals = append(resVals, v.(string))
		}
		for _, value := range data.UserCertificates.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(resVals, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.UserCertificates, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if !data.MacAddresses.IsNull() && res.Result.Macaddress != nil {
		var changedVals []string
		for _, value := range data.MacAddresses.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.Macaddress, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.MacAddresses, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if !data.IpaSshPubKeys.IsNull() && res.Result.Ipasshpubkey != nil {
		var changedVals []string

		for _, value := range data.IpaSshPubKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.Ipasshpubkey, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.IpaSshPubKeys, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if !data.Userclass.IsNull() && res.Result.Userclass != nil {
		var changedVals []string
		for _, value := range data.Userclass.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.Userclass, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.Userclass, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ipaassignedidview != nil && !data.AssignedIdView.IsNull() {
		data.AssignedIdView = types.StringValue(*res.Result.Ipaassignedidview)
	}
	if !data.KrbAuthIndicator.IsNull() && res.Result.Krbprincipalauthind != nil {
		var changedVals []string
		for _, value := range data.KrbAuthIndicator.Elements() {
			val, _ := strconv.Unquote(value.String())
			if slices.Contains(*res.Result.Krbprincipalauthind, val) {
				changedVals = append(changedVals, val)
			}
		}
		var diag diag.Diagnostics
		data.KrbAuthIndicator, diag = types.ListValueFrom(ctx, types.StringType, changedVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ipakrbrequirespreauth != nil && !data.KrbPreAuth.IsNull() {
		data.KrbPreAuth = types.BoolValue(*res.Result.Ipakrbrequirespreauth)
	}
	if res.Result.Ipakrbokasdelegate != nil && !data.TrustedForDelegation.IsNull() {
		data.TrustedForDelegation = types.BoolValue(*res.Result.Ipakrbokasdelegate)
	}
	if res.Result.Ipakrboktoauthasdelegate != nil && !data.TrustedToAuthAsDelegate.IsNull() {
		data.TrustedToAuthAsDelegate = types.BoolValue(*res.Result.Ipakrboktoauthasdelegate)
	}
	if res.Result.Randompassword != nil && !data.RandomPassword.IsNull() && data.RandomPassword.ValueBool() {
		data.GeneratedPassword = types.StringValue(*res.Result.Randompassword)
	} else {
		data.GeneratedPassword = types.StringValue("")
	}

	data.Id = data.Name
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa host %s", res.Result.Fqdn))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state HostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.HostModOptionalArgs{}

	args := ipa.HostModArgs{
		Fqdn: data.Name.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa host %s from plan = %v", data.Name.ValueString(), data))
	if !data.Description.Equal(state.Description) {
		if data.Description.ValueStringPointer() != nil {
			optArgs.Description = data.Description.ValueStringPointer()
		} else {
			v := ""
			optArgs.Description = &v
		}
	}
	if !data.Locality.Equal(state.Locality) {
		if data.Locality.ValueStringPointer() != nil {
			optArgs.L = data.Locality.ValueStringPointer()
		} else {
			v := ""
			optArgs.L = &v
		}
	}
	if !data.Location.Equal(state.Location) {
		if data.Location.ValueStringPointer() != nil {
			optArgs.Nshostlocation = data.Location.ValueStringPointer()
		} else {
			v := ""
			optArgs.Nshostlocation = &v
		}
	}
	if !data.Platform.Equal(state.Platform) {
		if data.Platform.ValueStringPointer() != nil {
			optArgs.Nshardwareplatform = data.Platform.ValueStringPointer()
		} else {
			v := ""
			optArgs.Nshardwareplatform = &v
		}
	}
	if !data.OperatingSystem.Equal(state.OperatingSystem) {
		if data.OperatingSystem.ValueStringPointer() != nil {
			optArgs.Nsosversion = data.OperatingSystem.ValueStringPointer()
		} else {
			v := ""
			optArgs.Nsosversion = &v
		}
	}
	if !data.UserCertificates.Equal(state.UserCertificates) {
		var v []interface{}
		for _, value := range data.UserCertificates.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Usercertificate = &v
	}
	if !data.MacAddresses.Equal(state.MacAddresses) {
		var v []string
		for _, value := range data.MacAddresses.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Macaddress = &v
	}
	if !data.IpaSshPubKeys.Equal(state.IpaSshPubKeys) {
		var v []string
		for _, value := range data.IpaSshPubKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if !data.Userclass.Equal(state.Userclass) {
		var v []string
		for _, value := range data.Userclass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}
	if !data.AssignedIdView.Equal(state.AssignedIdView) {
		if data.AssignedIdView.ValueStringPointer() != nil {
			optArgs.Ipaassignedidview = data.AssignedIdView.ValueStringPointer()
		} else {
			v := ""
			optArgs.Ipaassignedidview = &v
		}
	}
	if !data.KrbAuthIndicator.Equal(state.KrbAuthIndicator) {
		var v []string
		for _, value := range data.KrbAuthIndicator.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalauthind = &v
	}
	if !data.KrbPreAuth.Equal(state.KrbPreAuth) {
		if data.KrbPreAuth.ValueBoolPointer() != nil {
			optArgs.Ipakrbrequirespreauth = data.KrbPreAuth.ValueBoolPointer()
		} else {
			v := false
			optArgs.Ipakrbrequirespreauth = &v
		}
	}
	if !data.TrustedForDelegation.Equal(state.TrustedForDelegation) {
		if data.TrustedForDelegation.ValueBoolPointer() != nil {
			optArgs.Ipakrbokasdelegate = data.TrustedForDelegation.ValueBoolPointer()
		} else {
			v := false
			optArgs.Ipakrbokasdelegate = &v
		}
	}
	if !data.TrustedToAuthAsDelegate.Equal(state.TrustedToAuthAsDelegate) {
		if data.TrustedToAuthAsDelegate.ValueBoolPointer() != nil {
			optArgs.Ipakrboktoauthasdelegate = data.TrustedToAuthAsDelegate.ValueBoolPointer()
		} else {
			v := false
			optArgs.Ipakrboktoauthasdelegate = &v
		}
	}
	if !data.RandomPassword.Equal(state.RandomPassword) {
		if data.RandomPassword.ValueBoolPointer() != nil {
			optArgs.Random = data.RandomPassword.ValueBoolPointer()
		} else {
			v := false
			optArgs.Random = &v
		}
	}
	if !data.UserPassword.Equal(state.UserPassword) {
		if data.UserPassword.ValueStringPointer() != nil {
			optArgs.Userpassword = data.UserPassword.ValueStringPointer()
		} else {
			v := ""
			optArgs.Userpassword = &v
		}
	}

	_, err := r.client.HostMod(&args, &optArgs)
	if err != nil && !strings.Contains(err.Error(), "EmptyModlist") {
		resp.Diagnostics.AddWarning("Client Warning", err.Error())
	}

	data.GeneratedPassword = state.GeneratedPassword
	data.Id = data.Name

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa host Id %s", data.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa host Name %s", data.Name.ValueString()))
	args := ipa.HostDelArgs{
		Fqdn: []string{data.Name.ValueString()},
	}
	valTrue := true
	optArgs := ipa.HostDelOptionalArgs{
		Updatedns: &valTrue,
	}
	_, err := r.client.HostDel(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("[DEBUG] Host %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}

}

func (r *HostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
