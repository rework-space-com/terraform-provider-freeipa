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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &dnsZone{}
var _ resource.ResourceWithImportState = &dnsZone{}

func NewDNSZoneResource() resource.Resource {
	return &dnsZone{}
}

// resourceModel defines the resource implementation.
type dnsZone struct {
	client *ipa.Client
}

// resourceModelModel describes the resource data model.
type dnsZoneModel struct {
	Id                       types.String `tfsdk:"id"`
	ZoneName                 types.String `tfsdk:"zone_name"`
	IsReverseZone            types.Bool   `tfsdk:"is_reverse_zone"`
	DisableZone              types.Bool   `tfsdk:"disable_zone"`
	SkipOverlapCheck         types.Bool   `tfsdk:"skip_overlap_check"`
	AuthoritativeNameserver  types.String `tfsdk:"authoritative_nameserver"`
	SkipNameserverCheck      types.Bool   `tfsdk:"skip_nameserver_check"`
	AdminEmailAddress        types.String `tfsdk:"admin_email_address"`
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
	ComputedZoneName         types.String `tfsdk:"computed_zone_name"`
}

func (r *dnsZone) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *dnsZone) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *dnsZone) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA DNS Zone resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_name": schema.StringAttribute{
				MarkdownDescription: "Zone name (FQDN)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_reverse_zone": schema.BoolAttribute{
				MarkdownDescription: "Allow create the reverse zone",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_zone": schema.BoolAttribute{
				MarkdownDescription: "Allow disabled the zone",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_overlap_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if it will overlap with an existing zone",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"authoritative_nameserver": schema.StringAttribute{
				MarkdownDescription: "Authoritative nameserver domain name",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_nameserver_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if nameserver is not resolvable",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_email_address": schema.StringAttribute{
				MarkdownDescription: "Administrator e-mail address",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soa_refresh": schema.Int64Attribute{
				MarkdownDescription: "SOA record refresh time",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(3600)),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"soa_retry": schema.Int64Attribute{
				MarkdownDescription: "SOA record retry time",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(900)),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"soa_expire": schema.Int64Attribute{
				MarkdownDescription: "SOA record expire time",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(1209600)),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"soa_minimum": schema.Int64Attribute{
				MarkdownDescription: "How long should negative responses be cached",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(3600)),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records at zone apex",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records without explicit TTL definition",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"dynamic_updates": schema.BoolAttribute{
				MarkdownDescription: "Allow dynamic updates",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_update_policy": schema.StringAttribute{
				MarkdownDescription: "BIND update policy",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_query": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to issue queries",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("any;"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_transfer": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to transfer the zone",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("none;"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_forwarders": schema.ListAttribute{
				MarkdownDescription: "Per-zone forwarders. A custom port can be specified for each forwarder using a standard format IP_ADDRESS port PORT",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_ptr_sync": schema.BoolAttribute{
				MarkdownDescription: "Allow synchronization of forward (A, AAAA) and reverse (PTR) records in the zone",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_inline_dnssec_signing": schema.BoolAttribute{
				MarkdownDescription: "Allow inline DNSSEC signing of records in the zone",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"nsec3param_record": schema.StringAttribute{
				MarkdownDescription: "NSEC3PARAM record for zone in format: hash_algorithm flags iterations salt",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"computed_zone_name": schema.StringAttribute{
				MarkdownDescription: "Real zone name compatible with ARPA (ie: `domain.tld.`)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dnsZone) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsZone) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dnsZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.DnszoneAddOptionalArgs{}
	args := ipa.DnszoneAddArgs{}

	if !data.IsReverseZone.IsNull() && data.IsReverseZone.ValueBool() {
		optArgs.NameFromIP = data.ZoneName.ValueStringPointer()
	} else {
		var zone_name interface{} = data.ZoneName.ValueString()
		optArgs.Idnsname = &zone_name
	}
	if !data.SkipOverlapCheck.IsNull() && data.SkipOverlapCheck.ValueBool() {
		optArgs.SkipOverlapCheck = data.SkipOverlapCheck.ValueBoolPointer()
	} else {
		data.SkipOverlapCheck = types.BoolValue(false)
	}
	if !data.SkipNameserverCheck.IsNull() && data.SkipNameserverCheck.ValueBool() {
		optArgs.SkipNameserverCheck = data.SkipNameserverCheck.ValueBoolPointer()
	} else {
		data.SkipNameserverCheck = types.BoolValue(false)
	}
	if !data.AuthoritativeNameserver.IsNull() {
		var auth_nameserver interface{} = data.AuthoritativeNameserver.ValueString()
		optArgs.Idnssoamname = &auth_nameserver
	}
	if !data.AdminEmailAddress.IsNull() {
		var admin_email interface{} = data.AdminEmailAddress.ValueString()
		optArgs.Idnssoarname = &admin_email
	}
	if !data.SoaRefresh.IsNull() {
		soa_refresh := int(data.SoaRefresh.ValueInt64())
		optArgs.Idnssoarefresh = &soa_refresh
	}
	if !data.SoaRetry.IsNull() {
		soa_retry := int(data.SoaRetry.ValueInt64())
		optArgs.Idnssoaretry = &soa_retry
	}
	if !data.SoaExpire.IsNull() {
		soa_expire := int(data.SoaExpire.ValueInt64())
		optArgs.Idnssoaexpire = &soa_expire
	}
	if !data.SoaMinimum.IsNull() {
		soa_min := int(data.SoaMinimum.ValueInt64())
		optArgs.Idnssoaminimum = &soa_min
	}
	if !data.TTL.IsNull() {
		soa_ttl := int(data.TTL.ValueInt64())
		optArgs.Dnsttl = &soa_ttl
	}
	if !data.DefaultTTL.IsNull() {
		soa_default_ttl := int(data.DefaultTTL.ValueInt64())
		optArgs.Dnsdefaultttl = &soa_default_ttl
	}
	if !data.DynamicUpdate.IsNull() {
		optArgs.Idnsallowdynupdate = data.DynamicUpdate.ValueBoolPointer()
	}
	if !data.BindUpdatePolicy.IsNull() {
		optArgs.Idnsupdatepolicy = data.BindUpdatePolicy.ValueStringPointer()
	}
	if !data.AllowQuery.IsNull() {
		optArgs.Idnsallowquery = data.AllowQuery.ValueStringPointer()
	}
	if !data.AllowTransfer.IsNull() {
		optArgs.Idnsallowtransfer = data.AllowTransfer.ValueStringPointer()
	}
	if len(data.ZoneForwarders.Elements()) > 0 {
		var v []string
		for _, value := range data.ZoneForwarders.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Idnsforwarders = &v
	}
	if !data.AllowPtrSync.IsNull() {
		optArgs.Idnsallowsyncptr = data.AllowPtrSync.ValueBoolPointer()
	}
	if !data.AllowInlineDnssecSigning.IsNull() {
		optArgs.Idnssecinlinesigning = data.AllowInlineDnssecSigning.ValueBoolPointer()
	}
	if !data.Nsec3ParamRecord.IsNull() {
		optArgs.Nsec3paramrecord = data.Nsec3ParamRecord.ValueStringPointer()
	}

	res, err := r.client.DnszoneAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa dns zone: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa dns zone %s result: %v", data.ZoneName.ValueString(), res.Result.Idnsname))
	dnsnames := res.Result.Idnsname.([]interface{})
	dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"]
	data.ComputedZoneName = types.StringValue(dnsname.(string))
	data.Id = types.StringValue(dnsname.(string))

	if !data.DisableZone.IsNull() && data.DisableZone.ValueBool() {
		var name interface{} = data.Id.ValueString()
		_, err = r.client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: &name})
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone disable/enable. Something went wrong: %s", err))
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsZone) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dnsZoneModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	var name interface{} = data.Id.ValueString()
	optArgs := ipa.DnszoneShowOptionalArgs{
		All:      &all,
		Rights:   &all,
		Idnsname: &name,
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

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa DNS Zone %s", data.ZoneName.ValueString()))
	if res.Result.Idnszoneactive != nil && !data.DisableZone.IsNull() {
		data.DisableZone = types.BoolValue(!*res.Result.Idnszoneactive)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone disable_zone %s", data.DisableZone.String()))
	}
	if res.Result.Idnssoamname != nil && !data.AuthoritativeNameserver.IsNull() {
		auth_nameservers := (*res.Result.Idnssoamname).([]interface{})
		auth_nameserver := auth_nameservers[0].(map[string]interface{})["__dns_name__"]
		data.AuthoritativeNameserver = types.StringValue(auth_nameserver.(string))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone authoritative_nameserver %s", data.AuthoritativeNameserver.ValueString()))
	}
	if res.Result.Idnssoarname != nil && !data.AdminEmailAddress.IsNull() {
		admin_mails := (*res.Result.Idnssoarname).([]interface{})
		admin_mail := admin_mails[0].(map[string]interface{})["__dns_name__"]
		data.AdminEmailAddress = types.StringValue(admin_mail.(string))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone admin_email_address %s", data.AdminEmailAddress.ValueString()))
	}
	if !data.SoaRefresh.IsNull() {
		data.SoaRefresh = types.Int64Value(int64(res.Result.Idnssoarefresh))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_refresh %d", int(data.SoaRefresh.ValueInt64())))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Remote Read freeipa dns zone soa_refresh %d", res.Result.Idnssoarefresh))
	}
	if !data.SoaRetry.IsNull() {
		data.SoaRetry = types.Int64Value(int64(res.Result.Idnssoaretry))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_retry %d", int(data.SoaRetry.ValueInt64())))
	}
	if !data.SoaExpire.IsNull() {
		data.SoaExpire = types.Int64Value(int64(res.Result.Idnssoaexpire))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_expire %d", int(data.SoaExpire.ValueInt64())))
	}
	if !data.SoaMinimum.IsNull() {
		data.SoaMinimum = types.Int64Value(int64(res.Result.Idnssoaminimum))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone soa_minimum %d", int(data.SoaMinimum.ValueInt64())))
	}
	if res.Result.Dnsttl != nil && !data.TTL.IsNull() {
		data.TTL = types.Int64Value(int64(*res.Result.Dnsttl))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone ttl %d", int(data.TTL.ValueInt64())))
	}
	if res.Result.Dnsdefaultttl != nil && !data.DefaultTTL.IsNull() {
		data.DefaultTTL = types.Int64Value(int64(*res.Result.Dnsdefaultttl))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone default_ttl %d", int(data.DefaultTTL.ValueInt64())))
	}
	if res.Result.Idnsupdatepolicy != nil && !data.BindUpdatePolicy.IsNull() {
		data.BindUpdatePolicy = types.StringValue(*res.Result.Idnsupdatepolicy)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone bind_update_policy %s", data.BindUpdatePolicy.ValueString()))
	}
	if res.Result.Idnsallowdynupdate != nil && !data.DynamicUpdate.IsNull() {
		data.DynamicUpdate = types.BoolValue(*res.Result.Idnsallowdynupdate)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone dynamic_updates %s", data.DynamicUpdate.String()))
	}
	if res.Result.Idnsallowquery != nil && !data.AllowQuery.IsNull() {
		data.AllowQuery = types.StringValue(*res.Result.Idnsallowquery)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_query %s", data.AllowQuery.ValueString()))
	}
	if res.Result.Idnsallowtransfer != nil && !data.AllowTransfer.IsNull() {
		data.AllowTransfer = types.StringValue(*res.Result.Idnsallowtransfer)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_transfer %s", data.AllowTransfer.ValueString()))
	}
	if res.Result.Idnsallowsyncptr != nil && !data.AllowPtrSync.IsNull() {
		data.AllowPtrSync = types.BoolValue(*res.Result.Idnsallowsyncptr)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_ptr_sync %s", data.AllowPtrSync.String()))
	} else {
		data.AllowPtrSync = types.BoolValue(false)
	}
	if res.Result.Idnssecinlinesigning != nil && !data.AllowInlineDnssecSigning.IsNull() {
		data.AllowInlineDnssecSigning = types.BoolValue(*res.Result.Idnssecinlinesigning)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone allow_inline_dnssec_signing %s", data.AllowInlineDnssecSigning.String()))
	}
	if res.Result.Nsec3paramrecord != nil && !data.Nsec3ParamRecord.IsNull() {
		data.Nsec3ParamRecord = types.StringValue(*res.Result.Nsec3paramrecord)
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone nsec3param_record %s", data.Nsec3ParamRecord.ValueString()))
	}

	dnsnames := res.Result.Idnsname.([]interface{})
	dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"]
	data.ComputedZoneName = types.StringValue(dnsname.(string))
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dnsZone) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state dnsZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s: %v", data.ZoneName.ValueString(), data))

	hasChange := false
	var zone_id interface{} = data.Id.ValueString()
	optArgs := ipa.DnszoneModOptionalArgs{
		Idnsname: &zone_id,
	}

	if data.IsReverseZone.IsNull() {
		data.IsReverseZone = state.IsReverseZone
	}
	if !data.AuthoritativeNameserver.Equal(state.AuthoritativeNameserver) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone AuthoritativeNameserver %s: %s", data.ZoneName.ValueString(), data.AuthoritativeNameserver.ValueString()))
		var auth_nameserver interface{} = data.AuthoritativeNameserver.ValueString()
		optArgs.Idnssoamname = &auth_nameserver
		hasChange = true
	}
	if !data.AdminEmailAddress.Equal(state.AdminEmailAddress) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone AdminEmailAddress %s: %s", data.ZoneName.ValueString(), data.AdminEmailAddress.ValueString()))
		var admin_email interface{} = data.AdminEmailAddress.ValueString()
		optArgs.Idnssoarname = &admin_email
		hasChange = true
	}
	if !data.SoaRefresh.Equal(state.SoaRefresh) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.SoaRefresh.ValueInt64())))
		_v := int(data.SoaRefresh.ValueInt64())
		optArgs.Idnssoarefresh = &_v
		hasChange = true
	}
	if !data.SoaRetry.Equal(state.SoaRetry) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.SoaRetry.ValueInt64())))
		_v := int(data.SoaRetry.ValueInt64())
		optArgs.Idnssoaretry = &_v
		hasChange = true
	}
	if !data.SoaExpire.Equal(state.SoaExpire) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.SoaExpire.ValueInt64())))
		_v := int(data.SoaExpire.ValueInt64())
		optArgs.Idnssoaexpire = &_v
		hasChange = true
	}
	if !data.SoaMinimum.Equal(state.SoaMinimum) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.SoaMinimum.ValueInt64())))
		_v := int(data.SoaMinimum.ValueInt64())
		optArgs.Idnssoaminimum = &_v
		hasChange = true
	}
	if !data.TTL.Equal(state.TTL) {
		if !data.TTL.IsNull() {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.TTL.ValueInt64())))
			_v := int(data.TTL.ValueInt64())
			optArgs.Dnsttl = &_v
			hasChange = true
		}
	}
	if !data.DefaultTTL.Equal(state.DefaultTTL) {
		if !data.DefaultTTL.IsNull() {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s change: %d", data.ZoneName.ValueString(), int(data.DefaultTTL.ValueInt64())))
			_v := int(data.DefaultTTL.ValueInt64())
			optArgs.Dnsdefaultttl = &_v
			hasChange = true
		}
	}
	if !data.DynamicUpdate.Equal(state.DynamicUpdate) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s DynamicUpdate has change", data.ZoneName.ValueString()))
		_v := data.DynamicUpdate.ValueBool()
		optArgs.Idnsallowdynupdate = &_v
		hasChange = true
	}
	if !data.AllowPtrSync.Equal(state.AllowPtrSync) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s AllowPtrSync has change", data.ZoneName.ValueString()))
		_v := data.AllowPtrSync.ValueBool()
		optArgs.Idnsallowsyncptr = &_v
		hasChange = true
	}
	if !data.AllowInlineDnssecSigning.Equal(state.AllowInlineDnssecSigning) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s AllowInlineDnssecSigning has change", data.ZoneName.ValueString()))
		_v := data.AllowInlineDnssecSigning.ValueBool()
		optArgs.Idnssecinlinesigning = &_v
		hasChange = true
	}
	if !data.BindUpdatePolicy.Equal(state.BindUpdatePolicy) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s BindUpdatePolicy has change", data.ZoneName.ValueString()))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone BindUpdatePolicy %s: %s", data.ZoneName.ValueString(), data.BindUpdatePolicy.ValueString()))
		optArgs.Idnsupdatepolicy = data.BindUpdatePolicy.ValueStringPointer()
		hasChange = true
	}
	if !data.AllowQuery.Equal(state.AllowQuery) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s AllowQuery has change", data.ZoneName.ValueString()))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone AllowQuery %s: %s", data.ZoneName.ValueString(), data.AllowQuery.ValueString()))
		optArgs.Idnsallowquery = data.AllowQuery.ValueStringPointer()
		hasChange = true
	}
	if !data.AllowTransfer.Equal(state.AllowTransfer) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s AllowTransfer has change", data.ZoneName.ValueString()))
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone AllowTransfer %s: %s", data.ZoneName.ValueString(), data.AllowTransfer.ValueString()))
		optArgs.Idnsallowtransfer = data.AllowTransfer.ValueStringPointer()
		hasChange = true
	}
	if !data.ZoneForwarders.Equal(state.ZoneForwarders) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s ZoneForwarders has change", data.ZoneName.ValueString()))
		var v []string
		for _, value := range data.ZoneForwarders.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Idnsforwarders = &v
		hasChange = true
	}
	if !data.Nsec3ParamRecord.Equal(state.Nsec3ParamRecord) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone Nsec3ParamRecord %s: %s", data.ZoneName.ValueString(), data.Nsec3ParamRecord.ValueString()))
		optArgs.Nsec3paramrecord = data.Nsec3ParamRecord.ValueStringPointer()
		hasChange = true
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s hasChange: %v", data.ZoneName.ValueString(), hasChange))
	var res *ipa.DnszoneModResult
	var err error
	if hasChange {
		res, err = r.client.DnszoneMod(&ipa.DnszoneModArgs{}, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				tflog.Debug(ctx, fmt.Sprintf("EmptyModlist (4202): no modifications to be performed on DNS zone %s", data.ZoneName.ValueString()))
				tflog.Debug(ctx, fmt.Sprintf("EmptyModlist (4202): no modifications to be performed on DNS zone %s", data.ZoneName.ValueString()))
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa dns zone: %s", err))
				return
			}
		}

	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns zone %s plan disabled %s - state disabled %s", data.ZoneName.ValueString(), data.DisableZone.String(), state.DisableZone.String()))
	var name interface{} = data.Id.ValueString()
	if !data.DisableZone.Equal(state.DisableZone) {
		if data.DisableZone.ValueBool() {
			_, err := r.client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: &name})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone disable. Something went wrong: %s", err))
				return
			}
		} else {
			_, err := r.client.DnszoneEnable(&ipa.DnszoneEnableArgs{}, &ipa.DnszoneEnableOptionalArgs{Idnsname: &name})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone enable. Something went wrong: %s", err))
				return
			}
		}
	}

	if res != nil {
		if res.Result.Dnsttl != nil {
			ttl := res.Result.Dnsttl
			data.TTL = types.Int64Value(int64(*ttl))
		} else {
			data.TTL = types.Int64Null()
		}

		if res.Result.Dnsdefaultttl != nil {
			dttl := res.Result.Dnsdefaultttl
			data.DefaultTTL = types.Int64Value(int64(*dttl))
		} else {
			data.DefaultTTL = types.Int64Null()
		}
	} else {
		data.TTL = types.Int64Null()
		data.DefaultTTL = types.Int64Null()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsZone) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dnsZoneModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var id []interface{}
	id = append(id, data.Id.ValueString())
	optArgs := ipa.DnszoneDelOptionalArgs{
		Idnsname: &id,
	}
	_, err := r.client.DnszoneDel(&ipa.DnszoneDelArgs{}, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error delete freeipa dns zone: %s", err))
	}
}

func (r *dnsZone) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
