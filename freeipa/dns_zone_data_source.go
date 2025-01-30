// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &dnsZoneDataSource{}
var _ datasource.DataSourceWithConfigure = &dnsZoneDataSource{}

func NewDnsZoneDataSource() datasource.DataSource {
	return &dnsZoneDataSource{}
}

// resourceModel defines the resource implementation.
type dnsZoneDataSource struct {
	client *ipa.Client
}

// resourceModelModel describes the resource data model.
type dnsZoneDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	ZoneName                 types.String `tfsdk:"zone_name"`
	DisableZone              types.Bool   `tfsdk:"disable_zone"`
	SkipOverlapCheck         types.Bool   `tfsdk:"skip_overlap_check"`
	AuthoritativeNameserver  types.String `tfsdk:"authoritative_nameserver"`
	SkipNameserverCheck      types.Bool   `tfsdk:"skip_nameserver_check"`
	AdminEmailAddress        types.String `tfsdk:"admin_email_address"`
	SoaSerialNumber          types.Int64  `tfsdk:"soa_serial_number"`
	SoaRefresh               types.Int64  `tfsdk:"soa_refresh"`
	SoaRetry                 types.Int64  `tfsdk:"soa_retry"`
	SoaExpire                types.Int64  `tfsdk:"soa_expire"`
	SoaMinimum               types.Int64  `tfsdk:"soa_minimum"`
	TTL                      types.Int64  `tfsdk:"ttl"`
	DefaultTTL               types.Int64  `tfsdk:"default_ttl"`
	DynamicUpdate            types.Bool   `tfsdk:"dynamic_updates"`
	BindUpdatePolicy         types.String `tfsdk:"bind_update_policy"`
	AllowQuery               types.String `tfsdk:"allow_query"`
	AllowTransfer            types.String `tfsdk:"allow_transfer"`
	ZoneForwarders           types.List   `tfsdk:"zone_forwarders"`
	AllowPtrSync             types.Bool   `tfsdk:"allow_ptr_sync"`
	AllowInlineDnssecSigning types.Bool   `tfsdk:"allow_inline_dnssec_signing"`
	Nsec3ParamRecord         types.String `tfsdk:"nsec3param_record"`
}

func (r *dnsZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *dnsZoneDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

func (r *dnsZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA DNS Zone resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
			},
			"zone_name": schema.StringAttribute{
				MarkdownDescription: "Zone name (FQDN)",
				Required:            true,
			},
			"disable_zone": schema.BoolAttribute{
				MarkdownDescription: "Allow disabled the zone",
				Computed:            true,
			},
			"skip_overlap_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if it will overlap with an existing zone",
				Computed:            true,
			},
			"authoritative_nameserver": schema.StringAttribute{
				MarkdownDescription: "Authoritative nameserver domain name",
				Computed:            true,
			},
			"skip_nameserver_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if nameserver is not resolvable",
				Computed:            true,
			},
			"admin_email_address": schema.StringAttribute{
				MarkdownDescription: "Administrator e-mail address",
				Computed:            true,
			},
			"soa_serial_number": schema.Int64Attribute{
				MarkdownDescription: "SOA record serial number",
				Computed:            true,
			},
			"soa_refresh": schema.Int64Attribute{
				MarkdownDescription: "SOA record refresh time",
				Computed:            true,
			},
			"soa_retry": schema.Int64Attribute{
				MarkdownDescription: "SOA record retry time",
				Computed:            true,
			},
			"soa_expire": schema.Int64Attribute{
				MarkdownDescription: "SOA record expire time",
				Computed:            true,
			},
			"soa_minimum": schema.Int64Attribute{
				MarkdownDescription: "How long should negative responses be cached",
				Computed:            true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records at zone apex",
				Computed:            true,
			},
			"default_ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records without explicit TTL definition",
				Computed:            true,
			},
			"dynamic_updates": schema.BoolAttribute{
				MarkdownDescription: "Allow dynamic updates",
				Computed:            true,
			},
			"bind_update_policy": schema.StringAttribute{
				MarkdownDescription: "BIND update policy",
				Computed:            true,
			},
			"allow_query": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to issue queries",
				Computed:            true,
			},
			"allow_transfer": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to transfer the zone",
				Computed:            true,
			},
			"zone_forwarders": schema.ListAttribute{
				MarkdownDescription: "Per-zone forwarders. A custom port can be specified for each forwarder using a standard format IP_ADDRESS port PORT",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"allow_ptr_sync": schema.BoolAttribute{
				MarkdownDescription: "Allow synchronization of forward (A, AAAA) and reverse (PTR) records in the zone",
				Computed:            true,
			},
			"allow_inline_dnssec_signing": schema.BoolAttribute{
				MarkdownDescription: "Allow inline DNSSEC signing of records in the zone",
				Computed:            true,
			},
			"nsec3param_record": schema.StringAttribute{
				MarkdownDescription: "NSEC3PARAM record for zone in format: hash_algorithm flags iterations salt",
				Computed:            true,
			},
		},
	}
}

func (r *dnsZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *dnsZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dnsZoneDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	var zone_name interface{} = data.ZoneName.ValueString()

	if resp.Diagnostics.HasError() {
		return
	}

	all := true

	optArgs := ipa.DnszoneShowOptionalArgs{
		All:      &all,
		Rights:   &all,
		Idnsname: &zone_name,
	}

	res, err := r.client.DnszoneShow(&ipa.DnszoneShowArgs{}, &optArgs)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone %s: %s", data.ZoneName.ValueString(), res.String()))
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone %s not found", data.ZoneName.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa DNS zone: %s", err))
			return
		}
	}

	if res != nil {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone: %s", res.Result.String()))
	} else {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa dns zone %s", data.ZoneName.ValueString()))
		return
	}

	dnsnames := res.Result.Idnsname.([]interface{})
	dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"]
	data.ZoneName = types.StringValue(dnsname.(string))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa DNS Zone %s", data.ZoneName.ValueString()))
	if res.Result.Idnszoneactive != nil {
		data.DisableZone = types.BoolValue(!*res.Result.Idnszoneactive)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone disable_zone %s", data.DisableZone.String()))
	}
	if res.Result.Idnssoamname != nil {
		authnames := (*res.Result.Idnssoamname).([]interface{})
		authname := authnames[0].(map[string]interface{})["__dns_name__"]
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone authoritative_nameserver %v", *res.Result.Idnssoamname))
		data.AuthoritativeNameserver = types.StringValue(authname.(string))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone authoritative_nameserver %s", data.AuthoritativeNameserver.ValueString()))
	}
	if res.Result.Idnssoarname != nil {
		adminemails := (*res.Result.Idnssoamname).([]interface{})
		adminemail := adminemails[0].(map[string]interface{})["__dns_name__"]
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone admin_email %v", res.Result.Idnssoarname))
		data.AdminEmailAddress = types.StringValue(adminemail.(string))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone admin_email_address %s", data.AdminEmailAddress.ValueString()))
	}
	if res.Result.Idnssoaserial != nil {
		data.SoaSerialNumber = types.Int64Value(int64(*res.Result.Idnssoaserial))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_serial_number %d", int(data.SoaSerialNumber.ValueInt64())))
	}
	data.SoaRefresh = types.Int64Value(int64(res.Result.Idnssoarefresh))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_refresh %d", int(data.SoaRefresh.ValueInt64())))

	data.SoaRetry = types.Int64Value(int64(res.Result.Idnssoaretry))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_retry %d", int(data.SoaRetry.ValueInt64())))

	data.SoaExpire = types.Int64Value(int64(res.Result.Idnssoaexpire))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_expire %d", int(data.SoaExpire.ValueInt64())))

	data.SoaMinimum = types.Int64Value(int64(res.Result.Idnssoaminimum))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_minimum %d", int(data.SoaMinimum.ValueInt64())))

	if res.Result.Dnsttl != nil {
		data.TTL = types.Int64Value(int64(*res.Result.Dnsttl))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone ttl %d", int(data.TTL.ValueInt64())))
	}
	if res.Result.Dnsdefaultttl != nil {
		data.DefaultTTL = types.Int64Value(int64(*res.Result.Dnsdefaultttl))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone default_ttl %d", int(data.DefaultTTL.ValueInt64())))
	}
	if res.Result.Idnsupdatepolicy != nil {
		data.BindUpdatePolicy = types.StringValue(*res.Result.Idnsupdatepolicy)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone bind_update_policy %s", data.BindUpdatePolicy.ValueString()))
	}
	if res.Result.Idnsallowdynupdate != nil {
		data.DynamicUpdate = types.BoolValue(!*res.Result.Idnsallowdynupdate)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone dynamic_updates %s", data.DynamicUpdate.String()))
	}
	if res.Result.Idnsallowquery != nil {
		data.AllowQuery = types.StringValue(*res.Result.Idnsallowquery)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_query %s", data.AllowQuery.ValueString()))
	}
	if res.Result.Idnsallowtransfer != nil {
		data.AllowTransfer = types.StringValue(*res.Result.Idnsallowtransfer)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_transfer %s", data.AllowTransfer.ValueString()))
	}
	if res.Result.Idnsallowsyncptr != nil {
		data.AllowPtrSync = types.BoolValue(!*res.Result.Idnsallowsyncptr)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_ptr_sync %s", data.AllowPtrSync.String()))
	}
	if res.Result.Idnssecinlinesigning != nil {
		data.AllowInlineDnssecSigning = types.BoolValue(!*res.Result.Idnssecinlinesigning)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_inline_dnssec_signing %s", data.AllowInlineDnssecSigning.String()))
	}
	if res.Result.Nsec3paramrecord != nil {
		data.Nsec3ParamRecord = types.StringValue(*res.Result.Nsec3paramrecord)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone nsec3param_record %s", data.Nsec3ParamRecord.ValueString()))
	}

	data.Id = types.StringValue(data.ZoneName.ValueString())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
