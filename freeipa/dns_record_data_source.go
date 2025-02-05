// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	Id         types.String `tfsdk:"id"`
	RecordName types.String `tfsdk:"record_name"`
	ZoneName   types.String `tfsdk:"zone_name"`
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

	structured := true
	reqArgs := ipa.DnsrecordShowArgs{
		Idnsname: data.RecordName.ValueString(),
	}
	optArgs := ipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &zone_name,
		Structured:      &structured,
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
