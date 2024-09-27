// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_reverse_zone": schema.BoolAttribute{
				MarkdownDescription: "Allow create the reverse zone",
				Optional:            true,
			},
			"disable_zone": schema.BoolAttribute{
				MarkdownDescription: "Allow disabled the zone",
				Optional:            true,
			},
			"skip_overlap_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if it will overlap with an existing zone",
				Optional:            true,
			},
			"authoritative_nameserver": schema.StringAttribute{
				MarkdownDescription: "Authoritative nameserver domain name",
				Optional:            true,
			},
			"skip_nameserver_check": schema.BoolAttribute{
				MarkdownDescription: "Force DNS zone creation even if nameserver is not resolvable",
				Optional:            true,
			},
			"admin_email_address": schema.StringAttribute{
				MarkdownDescription: "Administrator e-mail address",
				Optional:            true,
			},
			"soa_serial_number": schema.Int64Attribute{
				MarkdownDescription: "SOA record serial number",
				Optional:            true,
			},
			"soa_refresh": schema.Int64Attribute{
				MarkdownDescription: "SOA record refresh time",
				Optional:            true,
			},
			"soa_retry": schema.Int64Attribute{
				MarkdownDescription: "SOA record retry time",
				Optional:            true,
			},
			"soa_expire": schema.Int64Attribute{
				MarkdownDescription: "SOA record expire time",
				Optional:            true,
			},
			"soa_minimum": schema.Int64Attribute{
				MarkdownDescription: "How long should negative responses be cached",
				Optional:            true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records at zone apex",
				Optional:            true,
			},
			"default_ttl": schema.Int64Attribute{
				MarkdownDescription: "Time to live for records without explicit TTL definition",
				Optional:            true,
			},
			"dynamic_updates": schema.BoolAttribute{
				MarkdownDescription: "Allow dynamic updates",
				Optional:            true,
			},
			"bind_update_policy": schema.StringAttribute{
				MarkdownDescription: "BIND update policy",
				Optional:            true,
			},
			"allow_query": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to issue queries",
				Optional:            true,
			},
			"allow_transfer": schema.StringAttribute{
				MarkdownDescription: "Semicolon separated list of IP addresses or networks which are allowed to transfer the zone",
				Optional:            true,
			},
			"zone_forwarders": schema.ListAttribute{
				MarkdownDescription: "Per-zone forwarders. A custom port can be specified for each forwarder using a standard format IP_ADDRESS port PORT",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"allow_ptr_sync": schema.BoolAttribute{
				MarkdownDescription: "Allow synchronization of forward (A, AAAA) and reverse (PTR) records in the zone",
				Optional:            true,
			},
			"allow_inline_dnssec_signing": schema.BoolAttribute{
				MarkdownDescription: "Allow inline DNSSEC signing of records in the zone",
				Optional:            true,
			},
			"nsec3param_record": schema.StringAttribute{
				MarkdownDescription: "NSEC3PARAM record for zone in format: hash_algorithm flags iterations salt",
				Optional:            true,
			},
			"computed_zone_name": schema.StringAttribute{
				MarkdownDescription: "Real zone name compatible with ARPA (ie: `domain.tld.`)",
				Computed:            true,
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

	if !data.SoaSerialNumber.IsUnknown() {
		args.Idnssoaserial = int(data.SoaSerialNumber.ValueInt64())
	}
	if !data.IsReverseZone.IsUnknown() && data.IsReverseZone.ValueBool() {
		if !data.ZoneName.IsUnknown() {
			optArgs.NameFromIP = data.ZoneName.ValueStringPointer()
		}
	} else {
		if !data.ZoneName.IsUnknown() {
			optArgs.Idnsname = data.ZoneName.ValueStringPointer()
		}
	}
	if !data.SkipOverlapCheck.IsUnknown() && data.SkipOverlapCheck.ValueBool() {
		optArgs.SkipOverlapCheck = data.SkipOverlapCheck.ValueBoolPointer()
	}
	if !data.SkipNameserverCheck.IsUnknown() && data.SkipNameserverCheck.ValueBool() {
		optArgs.SkipNameserverCheck = data.SkipNameserverCheck.ValueBoolPointer()
	}
	if !data.AuthoritativeNameserver.IsUnknown() {
		optArgs.Idnssoamname = data.AuthoritativeNameserver.ValueStringPointer()
	}
	if !data.AdminEmailAddress.IsUnknown() {
		optArgs.Idnssoarname = data.AdminEmailAddress.ValueStringPointer()
	}
	if !data.SoaRefresh.IsUnknown() {
		soa_refresh := int(data.SoaRefresh.ValueInt64())
		optArgs.Idnssoarefresh = &soa_refresh
	}
	if !data.SoaRetry.IsUnknown() {
		soa_retry := int(data.SoaRetry.ValueInt64())
		optArgs.Idnssoarefresh = &soa_retry
	}
	if !data.SoaExpire.IsUnknown() {
		soa_expire := int(data.SoaExpire.ValueInt64())
		optArgs.Idnssoaexpire = &soa_expire
	}
	if !data.SoaMinimum.IsUnknown() {
		soa_min := int(data.SoaMinimum.ValueInt64())
		optArgs.Idnssoaminimum = &soa_min
	}
	if !data.TTL.IsUnknown() {
		soa_ttl := int(data.TTL.ValueInt64())
		optArgs.Dnsttl = &soa_ttl
	}
	if !data.DefaultTTL.IsUnknown() {
		soa_default_ttl := int(data.DefaultTTL.ValueInt64())
		optArgs.Dnsdefaultttl = &soa_default_ttl
	}
	if !data.DynamicUpdate.IsUnknown() {
		optArgs.Idnsallowdynupdate = data.DynamicUpdate.ValueBoolPointer()
	}
	if !data.BindUpdatePolicy.IsUnknown() {
		optArgs.Idnsupdatepolicy = data.BindUpdatePolicy.ValueStringPointer()
	}
	if !data.AllowQuery.IsUnknown() {
		optArgs.Idnsallowquery = data.AllowQuery.ValueStringPointer()
	}
	if !data.AllowTransfer.IsUnknown() {
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
	if !data.AllowPtrSync.IsUnknown() {
		optArgs.Idnsallowsyncptr = data.AllowPtrSync.ValueBoolPointer()
	}
	if !data.AllowInlineDnssecSigning.IsUnknown() {
		optArgs.Idnssecinlinesigning = data.AllowInlineDnssecSigning.ValueBoolPointer()
	}
	if !data.Nsec3ParamRecord.IsUnknown() {
		optArgs.Nsec3paramrecord = data.Nsec3ParamRecord.ValueStringPointer()
	}

	res, err := r.client.DnszoneAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa dns zone: %s", err))
		return
	}

	data.ComputedZoneName = types.StringValue(res.Result.Idnsname)
	data.Id = types.StringValue(res.Result.Idnsname)

	if !data.DisableZone.IsUnknown() && data.DisableZone.ValueBool() {
		_, err = r.client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: data.Id.ValueStringPointer()})
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
	optArgs := ipa.DnszoneShowOptionalArgs{
		All:      &all,
		Rights:   &all,
		Idnsname: data.Id.ValueStringPointer(),
	}

	res, err := r.client.DnszoneShow(&ipa.DnszoneShowArgs{}, &optArgs)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns zone %s: %s", data.ZoneName.ValueString(), res.String()))
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone %s not found", data.ZoneName.ValueString()))
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa DNS zone: %s", err))
			return
		}
	}

	if !data.DisableZone.IsNull() {
		data.DisableZone = types.BoolValue(!*res.Result.Idnszoneactive)
	}
	data.ComputedZoneName = types.StringValue(res.Result.Idnsname)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dnsZone) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data dnsZoneModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	optArgs := ipa.DnszoneModOptionalArgs{
		Idnsname: data.ZoneName.ValueStringPointer(),
	}

	if !data.AuthoritativeNameserver.IsUnknown() {
		optArgs.Idnssoamname = data.AuthoritativeNameserver.ValueStringPointer()
	}
	if !data.AdminEmailAddress.IsUnknown() {
		optArgs.Idnssoarname = data.AdminEmailAddress.ValueStringPointer()
	}
	if !data.SoaSerialNumber.IsUnknown() {
		_v := int(data.SoaSerialNumber.ValueInt64())
		optArgs.Idnssoaserial = &_v
	}
	if !data.SoaRefresh.IsUnknown() {
		_v := int(data.SoaRefresh.ValueInt64())
		optArgs.Idnssoarefresh = &_v
	}
	if !data.SoaRetry.IsUnknown() {
		_v := int(data.SoaRetry.ValueInt64())
		optArgs.Idnssoaretry = &_v
	}
	if !data.SoaExpire.IsUnknown() {
		_v := int(data.SoaExpire.ValueInt64())
		optArgs.Idnssoaexpire = &_v
	}
	if !data.SoaMinimum.IsUnknown() {
		_v := int(data.SoaMinimum.ValueInt64())
		optArgs.Idnssoaminimum = &_v
	}
	if !data.TTL.IsUnknown() {
		_v := int(data.TTL.ValueInt64())
		optArgs.Dnsttl = &_v
	}
	if !data.DefaultTTL.IsUnknown() {
		_v := int(data.DefaultTTL.ValueInt64())
		optArgs.Dnsdefaultttl = &_v
	}
	if !data.DynamicUpdate.IsUnknown() {
		optArgs.Idnsallowdynupdate = data.DynamicUpdate.ValueBoolPointer()
	}
	if !data.AllowPtrSync.IsUnknown() {
		optArgs.Idnsallowsyncptr = data.AllowPtrSync.ValueBoolPointer()
	}
	if !data.AllowInlineDnssecSigning.IsUnknown() {
		optArgs.Idnssecinlinesigning = data.AllowInlineDnssecSigning.ValueBoolPointer()
	}
	if !data.BindUpdatePolicy.IsUnknown() {
		optArgs.Idnsupdatepolicy = data.BindUpdatePolicy.ValueStringPointer()
	}
	if !data.AllowQuery.IsUnknown() {
		optArgs.Idnsallowquery = data.AllowQuery.ValueStringPointer()
	}
	if !data.AllowTransfer.IsUnknown() {
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
	if !data.Nsec3ParamRecord.IsUnknown() {
		optArgs.Nsec3paramrecord = data.Nsec3ParamRecord.ValueStringPointer()
	}

	res, err := r.client.DnszoneMod(&ipa.DnszoneModArgs{}, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("EmptyModlist (4202): no modifications to be performed on DNS zone %s", data.ZoneName.ValueString()))
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error update freeipa dns zone: %s", err))
			return
		}
	}

	if !data.DisableZone.IsNull() {
		if data.DisableZone.ValueBool() {
			_, err := r.client.DnszoneDisable(&ipa.DnszoneDisableArgs{}, &ipa.DnszoneDisableOptionalArgs{Idnsname: data.ZoneName.ValueStringPointer()})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone disable/enable. Something went wrong: %s", err))
				return
			}
		} else {
			_, err = r.client.DnszoneEnable(&ipa.DnszoneEnableArgs{}, &ipa.DnszoneEnableOptionalArgs{Idnsname: data.ZoneName.ValueStringPointer()})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("DNS zone disable/enable. Something went wrong: %s", err))
				return
			}
		}
	}
	data.ComputedZoneName = types.StringValue(res.Result.Idnsname)

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
	id := []string{data.Id.ValueString()}
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
