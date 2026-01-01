// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//
// SPDX-License-Identifier: GPL-3.0-only

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
var _ datasource.DataSource = &dnsRecordDataSource{}
var _ datasource.DataSourceWithConfigure = &dnsRecordDataSource{}

func NewDnsRecordDataSource() datasource.DataSource {
	return &dnsRecordDataSource{}
}

// resourceModel defines the resource implementation.
type dnsRecordDataSource struct {
	client *ipa.Client
}

// resourceModelModel describes the resource data model.
type dnsRecordDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	RecordName   types.String `tfsdk:"record_name"`
	ZoneName     types.String `tfsdk:"zone_name"`
	ARecords     types.Set    `tfsdk:"a_records"`
	AAAARecords  types.Set    `tfsdk:"aaaa_records"`
	CnameRecords types.Set    `tfsdk:"cname_records"`
	MxRecords    types.Set    `tfsdk:"mx_records"`
	PtrRecords   types.Set    `tfsdk:"ptr_records"`
	SrvRecords   types.Set    `tfsdk:"srv_records"`
	TxtRecords   types.Set    `tfsdk:"txt_records"`
	SshfpRecords types.Set    `tfsdk:"sshfp_records"`
	NsRecords    types.Set    `tfsdk:"ns_records"`
}

func (r *dnsRecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *dnsRecordDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

func (r *dnsRecordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA DNS Record data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
			},
			"record_name": schema.StringAttribute{
				MarkdownDescription: "Record name",
				Required:            true,
			},
			"zone_name": schema.StringAttribute{
				MarkdownDescription: "Zone name (FQDN)",
				Required:            true,
			},
			"a_records": schema.SetAttribute{
				MarkdownDescription: "List of A records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"aaaa_records": schema.SetAttribute{
				MarkdownDescription: "List of AAAA records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cname_records": schema.SetAttribute{
				MarkdownDescription: "List of CNAME records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"mx_records": schema.SetAttribute{
				MarkdownDescription: "List of MX records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ptr_records": schema.SetAttribute{
				MarkdownDescription: "List of PTR records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"srv_records": schema.SetAttribute{
				MarkdownDescription: "List of SRV records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"txt_records": schema.SetAttribute{
				MarkdownDescription: "List of TXT records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sshfp_records": schema.SetAttribute{
				MarkdownDescription: "List of SSHFP records",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ns_records": schema.SetAttribute{
				MarkdownDescription: "List of NS records",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *dnsRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *dnsRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dnsRecordDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	var zone_name interface{} = data.ZoneName.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	reqArgs := ipa.DnsrecordShowArgs{
		Idnsname: data.RecordName.ValueString(),
	}
	all := true
	optArgs := ipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &zone_name,
		All:             &all,
	}

	res, err := r.client.DnsrecordShow(&reqArgs, &optArgs)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns record %s: %s", data.RecordName.ValueString(), res.String()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone: %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa dns zone %s", data.RecordName.ValueString()))
		return
	}

	if res.Result.Arecord != nil {
		var diag diag.Diagnostics
		data.ARecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Arecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Aaaarecord != nil {
		var diag diag.Diagnostics
		data.AAAARecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Aaaarecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Cnamerecord != nil {
		var diag diag.Diagnostics
		data.CnameRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Cnamerecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Mxrecord != nil {
		var diag diag.Diagnostics
		data.MxRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Mxrecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Ptrrecord != nil {
		var diag diag.Diagnostics
		data.MxRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Ptrrecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Srvrecord != nil {
		var diag diag.Diagnostics
		data.SrvRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Srvrecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Txtrecord != nil {
		var diag diag.Diagnostics
		data.TxtRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Txtrecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Sshfprecord != nil {
		var diag diag.Diagnostics
		data.SshfpRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Sshfprecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}
	if res.Result.Nsrecord != nil {
		var diag diag.Diagnostics
		data.NsRecords, diag = types.SetValueFrom(ctx, types.StringType, res.Result.Nsrecord)
		if diag.HasError() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("diag: %v\n", diag))
		}
	}

	data.Id = types.StringValue(fmt.Sprintf("%s.%s", data.RecordName.ValueString(), data.ZoneName.ValueString()))
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
