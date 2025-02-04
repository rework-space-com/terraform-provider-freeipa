// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HostDataSource{}
var _ datasource.DataSourceWithConfigure = &HostDataSource{}

func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

// HostDataSource defines the resource implementation.
type HostDataSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type HostDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Locality                  types.String `tfsdk:"locality"`
	Location                  types.String `tfsdk:"location"`
	Platform                  types.String `tfsdk:"platform"`
	OperatingSystem           types.String `tfsdk:"operating_system"`
	UserCertificates          types.List   `tfsdk:"user_certificates"`
	MacAddresses              types.List   `tfsdk:"mac_addresses"`
	IpaSshPubKeys             types.List   `tfsdk:"ipasshpubkeys"`
	Userclass                 types.List   `tfsdk:"userclass"`
	AssignedIdView            types.String `tfsdk:"assigned_idview"`
	KrbAuthIndicator          types.List   `tfsdk:"krb_auth_indicators"`
	KrbPreAuth                types.Bool   `tfsdk:"krb_preauth"`
	TrustedForDelegation      types.Bool   `tfsdk:"trusted_for_delegation"`
	TrustedToAuthAsDelegate   types.Bool   `tfsdk:"trusted_to_auth_as_delegate"`
	MemberOfHostGroup         types.List   `tfsdk:"memberof_hostgroup"`
	MemberOfSudoRule          types.List   `tfsdk:"memberof_sudorule"`
	MemberOfHBACRule          types.List   `tfsdk:"memberof_hbacrule"`
	MemberOfIndirectHostGroup types.List   `tfsdk:"memberof_indirect_hostgroup"`
	MemberOfIndirectSudoRule  types.List   `tfsdk:"memberof_indirect_sudorule"`
	MemberOfIndirectHBACRule  types.List   `tfsdk:"memberof_indirect_hbacrule"`
}

func (r *HostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *HostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA Host data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Host name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of this host",
				Computed:            true,
			},
			"locality": schema.StringAttribute{
				MarkdownDescription: "Host locality (e.g. 'Baltimore, MD')",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Host location (e.g. 'Lab 2')",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "Host hardware platform (e.g. 'Lenovo T61')",
				Computed:            true,
			},
			"operating_system": schema.StringAttribute{
				MarkdownDescription: "Host operating system and version (e.g. 'Fedora 40')",
				Computed:            true,
			},
			"user_certificates": schema.ListAttribute{
				MarkdownDescription: "Base-64 encoded host certificate",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"mac_addresses": schema.ListAttribute{
				MarkdownDescription: "Hardware MAC address(es) on this host",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ipasshpubkeys": schema.ListAttribute{
				MarkdownDescription: "SSH public keys",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"userclass": schema.ListAttribute{
				MarkdownDescription: "Host category (semantics placed on this attribute are for local interpretation)",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"assigned_idview": schema.StringAttribute{
				MarkdownDescription: "Assigned ID View",
				Computed:            true,
			},
			"krb_auth_indicators": schema.ListAttribute{
				MarkdownDescription: "Defines a whitelist for Authentication Indicators. Use 'otp' to allow OTP-based 2FA authentications. Use 'radius' to allow RADIUS-based 2FA authentications. Other values may be used for custom configurations.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"krb_preauth": schema.BoolAttribute{
				MarkdownDescription: "Pre-authentication is required for the service",
				Computed:            true,
			},
			"trusted_for_delegation": schema.BoolAttribute{
				MarkdownDescription: "Client credentials may be delegated to the service",
				Optional:            true,
			},
			"trusted_to_auth_as_delegate": schema.BoolAttribute{
				MarkdownDescription: "The service is allowed to authenticate on behalf of a client",
				Computed:            true,
			},
			"memberof_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups this user is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this user is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this user is member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_hostgroup": schema.ListAttribute{
				MarkdownDescription: "List of hostgroups this user is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_sudorule": schema.ListAttribute{
				MarkdownDescription: "List of SUDO rules this user is is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"memberof_indirect_hbacrule": schema.ListAttribute{
				MarkdownDescription: "List of HBAC rules this user is indirectly member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *HostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

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
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa host %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa host %s", data.Name.ValueString()))
		return
	}

	if res.Result.Description != nil {
		data.Description = types.StringValue(*res.Result.Description)
	}
	if res.Result.L != nil {
		data.Locality = types.StringValue(*res.Result.L)
	}
	if res.Result.Nshostlocation != nil {
		data.Location = types.StringValue(*res.Result.Nshostlocation)
	}
	if res.Result.Nshardwareplatform != nil {
		data.Platform = types.StringValue(*res.Result.Nshardwareplatform)
	}

	if res.Result.Nsosversion != nil {
		data.OperatingSystem = types.StringValue(*res.Result.Nsosversion)
	}
	if res.Result.Usercertificate != nil {
		var resVals []string
		for _, v := range *res.Result.Usercertificate {
			resVals = append(resVals, v.(string))
		}
		var diag diag.Diagnostics
		data.UserCertificates, diag = types.ListValueFrom(ctx, types.StringType, resVals)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Macaddress != nil {
		var diag diag.Diagnostics
		data.MacAddresses, diag = types.ListValueFrom(ctx, types.StringType, res.Result.Macaddress)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ipasshpubkey != nil {
		var diag diag.Diagnostics
		data.IpaSshPubKeys, diag = types.ListValueFrom(ctx, types.StringType, res.Result.Ipasshpubkey)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Userclass != nil {
		var diag diag.Diagnostics
		data.Userclass, diag = types.ListValueFrom(ctx, types.StringType, res.Result.Userclass)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ipaassignedidview != nil {
		data.AssignedIdView = types.StringValue(*res.Result.Ipaassignedidview)
	}
	if res.Result.Krbprincipalauthind != nil {
		var diag diag.Diagnostics
		data.KrbAuthIndicator, diag = types.ListValueFrom(ctx, types.StringType, res.Result.Krbprincipalauthind)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ipakrbrequirespreauth != nil {
		data.KrbPreAuth = types.BoolValue(*res.Result.Ipakrbrequirespreauth)
	}
	if res.Result.Ipakrbokasdelegate != nil {
		data.TrustedForDelegation = types.BoolValue(*res.Result.Ipakrbokasdelegate)
	}
	if res.Result.Ipakrboktoauthasdelegate != nil {
		data.TrustedToAuthAsDelegate = types.BoolValue(*res.Result.Ipakrboktoauthasdelegate)
	}
	if res.Result.MemberofHostgroup != nil {
		data.MemberOfHostGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofHostgroup)
	}
	if res.Result.MemberofHbacrule != nil {
		data.MemberOfHBACRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofHbacrule)
	}
	if res.Result.MemberofSudorule != nil {
		data.MemberOfSudoRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofSudorule)
	}
	if res.Result.MemberofindirectHostgroup != nil {
		data.MemberOfIndirectHostGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectHostgroup)
	}
	if res.Result.MemberofindirectHbacrule != nil {
		data.MemberOfIndirectHBACRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectHbacrule)
	}
	if res.Result.MemberofindirectSudorule != nil {
		data.MemberOfIndirectSudoRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectSudorule)
	}

	data.Id = data.Name
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa host %s", res.Result.Fqdn))

	data.Id = types.StringValue(data.Name.ValueString())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
