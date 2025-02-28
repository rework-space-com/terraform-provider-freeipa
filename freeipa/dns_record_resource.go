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
var _ resource.Resource = &DNSRecordResource{}
var _ resource.ResourceWithImportState = &DNSRecordResource{}

func NewDNSRecordResource() resource.Resource {
	return &DNSRecordResource{}
}

// DNSRecordResource defines the resource implementation.
type DNSRecordResource struct {
	client *ipa.Client
}

// DNSRecordResourceModel describes the resource data model.
type DNSRecordResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ZoneName      types.String `tfsdk:"zone_name"`
	Type          types.String `tfsdk:"type"`
	Records       types.List   `tfsdk:"records"`
	TTL           types.Int32  `tfsdk:"ttl"`
	SetIdentifier types.String `tfsdk:"set_identifier"`
}

func (r *DNSRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *DNSRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA DNS Record resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Record name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_name": schema.StringAttribute{
				MarkdownDescription: "Zone name (FQDN)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The record type (A, AAAA, CNAME, MX, PTR, SRV, TXT, SSHP)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"records": schema.ListAttribute{
				MarkdownDescription: "A string list of records",
				Required:            true,
				ElementType:         types.StringType,
			},
			"ttl": schema.Int32Attribute{
				MarkdownDescription: "Time to live",
				Optional:            true,
			},
			"set_identifier": schema.StringAttribute{
				MarkdownDescription: "Unique identifier to differentiate records with routing policies from one another",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *DNSRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DNSRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	var zone_name interface{} = data.ZoneName.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.DnsrecordAddArgs{
		Idnsname: data.Name.ValueString(),
	}

	optArgs := ipa.DnsrecordAddOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := data.Type.ValueString()

	if len(data.Records.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa dns record %s ", data.Name.String()))
		var records []string

		for _, value := range data.Records.Elements() {
			val, _ := strconv.Unquote(value.String())
			records = append(records, val)
		}
		switch _type {
		case "A":
			optArgs.Arecord = &records
		case "AAAA":
			optArgs.Aaaarecord = &records
		case "CNAME":
			optArgs.Cnamerecord = &records
		case "MX":
			optArgs.Mxrecord = &records
		case "NS":
			optArgs.Nsrecord = &records
		case "PTR":
			optArgs.Ptrrecord = &records
		case "SRV":
			optArgs.Srvrecord = &records
		case "TXT":
			optArgs.Txtrecord = &records
		case "SSHFP":
			optArgs.Sshfprecord = &records
		}
	}

	if !data.TTL.IsNull() {
		ttl := int(data.TTL.ValueInt32())
		optArgs.Dnsttl = &ttl
	}

	_, err := r.client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa dns record: %s", err))
		}
	}

	// Generate an ID
	vars := []string{
		data.ZoneName.ValueString(),
		strings.ToLower(data.Name.ValueString()),
		_type,
	}
	if !data.SetIdentifier.IsNull() {
		vars = append(vars, data.SetIdentifier.ValueString())
	}

	data.Id = types.StringValue(strings.Join(vars, "_"))

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DNSRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	var zone_name interface{} = data.ZoneName.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	args := ipa.DnsrecordShowArgs{
		Idnsname: data.Name.ValueString(),
	}

	all := true
	optArgs := ipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &zone_name,
		All:             &all,
	}

	res, err := r.client.DnsrecordShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] DNS record not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa DNS record: %s", err))
			return
		}
	}
	_type := data.Type.ValueString()

	switch _type {
	case "A":
		if res.Result.Arecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Arecord)
		}
	case "AAAA":
		if res.Result.Aaaarecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Aaaarecord)
		}
	case "MX":
		if res.Result.Mxrecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mxrecord)
		}
	case "NS":
		if res.Result.Nsrecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Nsrecord)
		}
	case "PTR":
		if res.Result.Ptrrecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ptrrecord)
		}
	case "SRV":
		if res.Result.Srvrecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Srvrecord)
		}
	case "TXT":
		if res.Result.Txtrecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Txtrecord)
		}
	case "SSHFP":
		if res.Result.Sshfprecord != nil {
			data.Records, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Sshfprecord)
		}
	}

	if res.Result.Dnsttl != nil && !data.TTL.IsNull() {
		data.TTL = types.Int32Value(int32(*res.Result.Dnsttl))
	}

	// Generate an ID
	vars := []string{
		data.ZoneName.ValueString(),
		strings.ToLower(data.Name.ValueString()),
		_type,
	}
	if !data.SetIdentifier.IsNull() {
		vars = append(vars, data.SetIdentifier.ValueString())
	}

	data.Id = types.StringValue(strings.Join(vars, "_"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state DNSRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var zone_name interface{} = data.ZoneName.ValueString()

	args := ipa.DnsrecordModArgs{
		Idnsname: data.Name.ValueString(),
	}

	optArgs := ipa.DnsrecordModOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := data.Type.ValueString()

	if !data.Records.Equal(state.Records) {
		var records []string

		for _, value := range data.Records.Elements() {
			val, _ := strconv.Unquote(value.String())
			records = append(records, val)
		}
		switch _type {
		case "A":
			optArgs.Arecord = &records
		case "AAAA":
			optArgs.Aaaarecord = &records
		case "CNAME":
			optArgs.Cnamerecord = &records
		case "MX":
			optArgs.Mxrecord = &records
		case "NS":
			optArgs.Nsrecord = &records
		case "PTR":
			optArgs.Ptrrecord = &records
		case "SRV":
			optArgs.Srvrecord = &records
		case "TXT":
			optArgs.Txtrecord = &records
		case "SSHFP":
			optArgs.Sshfprecord = &records
		}
	}

	if !data.TTL.Equal(state.TTL) {
		ttl := int(data.TTL.ValueInt32())
		optArgs.Dnsttl = &ttl
	}

	_, err := r.client.DnsrecordMod(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa dns record: %s", err))

		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DNSRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	var zone_name interface{} = data.ZoneName.ValueString()

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
	args := ipa.DnsrecordDelArgs{
		Idnsname: data.Name.ValueString(),
	}

	optArgs := ipa.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := data.Type.ValueString()
	if len(data.Records.Elements()) > 0 {
		var records []string

		for _, value := range data.Records.Elements() {
			val, _ := strconv.Unquote(value.String())
			records = append(records, val)
		}
		switch _type {
		case "A":
			optArgs.Arecord = &records
		case "AAAA":
			optArgs.Aaaarecord = &records
		case "CNAME":
			optArgs.Cnamerecord = &records
		case "MX":
			optArgs.Mxrecord = &records
		case "NS":
			optArgs.Nsrecord = &records
		case "PTR":
			optArgs.Ptrrecord = &records
		case "SRV":
			optArgs.Srvrecord = &records
		case "TXT":
			optArgs.Txtrecord = &records
		case "SSHFP":
			optArgs.Sshfprecord = &records
		}
	}

	_, err := r.client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa dns record: %s", err))
		return
	}
}

func (r *DNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
