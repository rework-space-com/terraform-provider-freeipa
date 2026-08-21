// Authors:
//   Pavel Aksenov <41126916+al-bashkir@users.noreply.github.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

var _ resource.Resource = &dnsForwardZone{}
var _ resource.ResourceWithImportState = &dnsForwardZone{}

type dnsForwardZone struct {
	client *ipa.Client
}

// FreeIPA >= 4.12 returns "managedby": null in dnsforwardzone responses, which
// the go-freeipa library rejects during JSON decoding. Detect that error so we
// can fall back to input-derived values.
func isManagedbyNullErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unexpected value for field Managedby")
}

// canonicalZoneName mirrors FreeIPA's normalization: lowercase + trailing dot.
func canonicalZoneName(name string) string {
	n := strings.ToLower(name)
	if !strings.HasSuffix(n, ".") {
		n += "."
	}
	return n
}

// caseInsensitiveZoneNameModifier suppresses spurious diffs on zone_name when
// the change is only in casing (FreeIPA treats DNS names case-insensitively).
type caseInsensitiveZoneNameModifier struct{}

func (m caseInsensitiveZoneNameModifier) Description(_ context.Context) string {
	return "Suppresses plan diff when zone name differs only in case."
}

func (m caseInsensitiveZoneNameModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m caseInsensitiveZoneNameModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if canonicalZoneName(req.StateValue.ValueString()) == canonicalZoneName(req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

type dnsForwardZoneModel struct {
	Id               types.String `tfsdk:"id"`
	ZoneName         types.String `tfsdk:"zone_name"`
	Forwarders       types.Set    `tfsdk:"forwarders"`
	ForwardPolicy    types.String `tfsdk:"forward_policy"`
	DisableZone      types.Bool   `tfsdk:"disable_zone"`
	SkipOverlapCheck types.Bool   `tfsdk:"skip_overlap_check"`
	ComputedZoneName types.String `tfsdk:"computed_zone_name"`
}

func NewDNSForwardZoneResource() resource.Resource {
	return &dnsForwardZone{}
}

func (r *dnsForwardZone) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_forward_zone"
}

func (r *dnsForwardZone) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*freeipaProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *freeipaProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p := &freeipaProvider{}
	client, err := p.NewFreeIPAClient(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create FreeIPA API Client",
			"An unexpected error occurred when creating the FreeIPA API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"FreeIPA Client Error: "+err.Error(),
		)
		return
	}

	r.client = client
}

func (r *dnsForwardZone) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "FreeIPA DNS forward zone resource. Manages per-zone forwarders, forward policy, and enable/disable state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource (canonical zone name with trailing dot).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_name": schema.StringAttribute{
				MarkdownDescription: "Forward zone name (FQDN).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					caseInsensitiveZoneNameModifier{},
					stringplanmodifier.RequiresReplace(),
				},
			},
			"forwarders": schema.SetAttribute{
				MarkdownDescription: "Per-zone forwarders. A custom port can be specified using a standard format `IP_ADDRESS port PORT`.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"forward_policy": schema.StringAttribute{
				MarkdownDescription: "Per-zone conditional forwarding policy. One of `only`, `first`, `none`. Defaults to `first`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("first"),
				Validators: []validator.String{
					stringvalidator.OneOf("only", "first", "none"),
				},
			},
			"disable_zone": schema.BoolAttribute{
				MarkdownDescription: "Disable the forward zone (inverse of FreeIPA `idnszoneactive`). Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"skip_overlap_check": schema.BoolAttribute{
				MarkdownDescription: "Force creation even if the zone overlaps with an existing zone. Create-time only; ignored on update. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"computed_zone_name": schema.StringAttribute{
				MarkdownDescription: "Canonical zone name returned by FreeIPA (with trailing dot).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dnsForwardZone) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dnsForwardZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa dns forward zone %s", data.ZoneName.ValueString()))

	addOpt := ipa.DnsforwardzoneAddOptionalArgs{}
	var name interface{} = data.ZoneName.ValueString()
	addOpt.Idnsname = &name

	var fwd []string
	for _, value := range data.Forwarders.Elements() {
		val, _ := strconv.Unquote(value.String())
		fwd = append(fwd, val)
	}
	addOpt.Idnsforwarders = &fwd

	if !data.ForwardPolicy.IsNull() && !data.ForwardPolicy.IsUnknown() {
		s := data.ForwardPolicy.ValueString()
		addOpt.Idnsforwardpolicy = &s
	}

	if data.SkipOverlapCheck.ValueBool() {
		b := true
		addOpt.SkipOverlapCheck = &b
	}

	res, err := r.client.DnsforwardzoneAdd(&ipa.DnsforwardzoneAddArgs{}, &addOpt)
	if err != nil {
		if !isManagedbyNullErr(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa dns forward zone: %s", err))
			return
		}
		// FreeIPA >= 4.12 returns managedby:null in Add and Show responses; the zone
		// was created. Try Show first to get FreeIPA's canonical name; if Show also
		// hits the managedby decode error, synthesize the canonical name from input.
		var showName interface{} = data.ZoneName.ValueString()
		showRes, showErr := r.client.DnsforwardzoneShow(
			&ipa.DnsforwardzoneShowArgs{},
			&ipa.DnsforwardzoneShowOptionalArgs{Idnsname: &showName},
		)
		if showErr != nil {
			if !isManagedbyNullErr(showErr) {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error verifying freeipa dns forward zone after create: %s", showErr))
				return
			}
			dnsname := canonicalZoneName(data.ZoneName.ValueString())
			data.Id = types.StringValue(dnsname)
			data.ComputedZoneName = types.StringValue(dnsname)
		} else {
			dnsnames := showRes.Result.Idnsname.([]interface{})
			dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"].(string)
			data.Id = types.StringValue(dnsname)
			data.ComputedZoneName = types.StringValue(dnsname)
		}
	} else {
		dnsnames := res.Result.Idnsname.([]interface{})
		dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"].(string)
		data.Id = types.StringValue(dnsname)
		data.ComputedZoneName = types.StringValue(dnsname)
	}

	if data.DisableZone.ValueBool() {
		var idName interface{} = data.Id.ValueString()
		if _, derr := r.client.DnsforwardzoneDisable(
			&ipa.DnsforwardzoneDisableArgs{},
			&ipa.DnsforwardzoneDisableOptionalArgs{Idnsname: &idName},
		); derr != nil {
			delIDs := []interface{}{data.Id.ValueString()}
			_, _ = r.client.DnsforwardzoneDel(
				&ipa.DnsforwardzoneDelArgs{},
				&ipa.DnsforwardzoneDelOptionalArgs{Idnsname: &delIDs},
			)
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Failed to disable forward zone %s after create; rolled back: %s", data.Id.ValueString(), derr),
			)
			return
		}
	}

	if _, err := r.readZone(ctx, &data, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error refreshing freeipa dns forward zone after create: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsForwardZone) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dnsForwardZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.readZone(ctx, &data, &resp.Diagnostics); err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] dns forward zone %s not found, removing from state", data.ZoneName.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa dns forward zone: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsForwardZone) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state dnsForwardZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa dns forward zone %s", data.ZoneName.ValueString()))

	modOpt := ipa.DnsforwardzoneModOptionalArgs{}
	var name interface{} = data.Id.ValueString()
	modOpt.Idnsname = &name
	hasChange := false

	if !data.Forwarders.Equal(state.Forwarders) {
		var v []string
		for _, value := range data.Forwarders.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		if v == nil {
			v = []string{}
		}
		modOpt.Idnsforwarders = &v
		hasChange = true
	}

	if !data.ForwardPolicy.Equal(state.ForwardPolicy) && !data.ForwardPolicy.IsNull() {
		s := data.ForwardPolicy.ValueString()
		modOpt.Idnsforwardpolicy = &s
		hasChange = true
	}

	if hasChange {
		if _, err := r.client.DnsforwardzoneMod(&ipa.DnsforwardzoneModArgs{}, &modOpt); err != nil {
			switch {
			case strings.Contains(err.Error(), "EmptyModlist"):
				resp.Diagnostics.AddWarning("Client Warning", err.Error())
			case isManagedbyNullErr(err):
				// FreeIPA >= 4.12: Mod succeeded server-side; only the response decode failed.
				resp.Diagnostics.AddWarning("Client Warning", fmt.Sprintf("Mod response decode failed (FreeIPA >= 4.12 managedby:null): %s", err))
			default:
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating freeipa dns forward zone: %s", err))
				return
			}
		}
	}

	if !data.DisableZone.Equal(state.DisableZone) {
		var idName interface{} = data.Id.ValueString()
		if data.DisableZone.ValueBool() {
			if _, err := r.client.DnsforwardzoneDisable(
				&ipa.DnsforwardzoneDisableArgs{},
				&ipa.DnsforwardzoneDisableOptionalArgs{Idnsname: &idName},
			); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling forward zone: %s", err))
				return
			}
		} else {
			if _, err := r.client.DnsforwardzoneEnable(
				&ipa.DnsforwardzoneEnableArgs{},
				&ipa.DnsforwardzoneEnableOptionalArgs{Idnsname: &idName},
			); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error enabling forward zone: %s", err))
				return
			}
		}
	}

	if _, err := r.readZone(ctx, &data, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error refreshing freeipa dns forward zone after update: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsForwardZone) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dnsForwardZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Delete freeipa dns forward zone %s", data.Id.ValueString()))

	ids := []interface{}{data.Id.ValueString()}
	if _, err := r.client.DnsforwardzoneDel(
		&ipa.DnsforwardzoneDelArgs{},
		&ipa.DnsforwardzoneDelOptionalArgs{Idnsname: &ids},
	); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting freeipa dns forward zone: %s", err))
	}
}

func (r *dnsForwardZone) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var name interface{} = req.ID
	res, err := r.client.DnsforwardzoneShow(
		&ipa.DnsforwardzoneShowArgs{},
		&ipa.DnsforwardzoneShowOptionalArgs{Idnsname: &name},
	)
	if err != nil {
		if !isManagedbyNullErr(err) {
			resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Error reading freeipa dns forward zone: %s", err))
			return
		}
		// FreeIPA >= 4.12: synthesize what we can; forward_policy / forwarders /
		// disable_zone stay unknown and the user must declare them in config.
		dnsname := canonicalZoneName(req.ID)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_name"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dnsname)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("computed_zone_name"), dnsname)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skip_overlap_check"), false)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disable_zone"), false)...)
		return
	}
	dnsnames := res.Result.Idnsname.([]interface{})
	dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"].(string)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dnsname)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("computed_zone_name"), dnsname)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skip_overlap_check"), false)...)
	if res.Result.Idnszoneactive != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disable_zone"), !*res.Result.Idnszoneactive)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disable_zone"), false)...)
	}
	if res.Result.Idnsforwardpolicy != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("forward_policy"), *res.Result.Idnsforwardpolicy)...)
	}
	if res.Result.Idnsforwarders != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("forwarders"), *res.Result.Idnsforwarders)...)
	}
}

func (r *dnsForwardZone) readZone(ctx context.Context, data *dnsForwardZoneModel, diags *diag.Diagnostics) (*ipa.Dnsforwardzone, error) {
	var name interface{} = data.Id.ValueString()
	if data.Id.IsNull() || data.Id.ValueString() == "" {
		name = data.ZoneName.ValueString()
	}
	res, err := r.client.DnsforwardzoneShow(
		&ipa.DnsforwardzoneShowArgs{},
		&ipa.DnsforwardzoneShowOptionalArgs{Idnsname: &name},
	)
	if err != nil {
		if isManagedbyNullErr(err) {
			// FreeIPA >= 4.12: cannot decode response. Skip refresh, preserve state.
			diags.AddWarning("Client Warning", fmt.Sprintf("Skipping refresh of dns forward zone %s: %s", data.ZoneName.ValueString(), err))
			return nil, nil
		}
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns forward zone %s: %s", data.ZoneName.ValueString(), res.String()))
	z := &res.Result

	if z.Idnsname != nil {
		dnsnames := z.Idnsname.([]interface{})
		dnsname := dnsnames[0].(map[string]interface{})["__dns_name__"].(string)
		data.ComputedZoneName = types.StringValue(dnsname)
		data.Id = types.StringValue(dnsname)
	}

	if z.Idnsforwarders != nil {
		v, d := types.SetValueFrom(ctx, types.StringType, *z.Idnsforwarders)
		diags.Append(d...)
		data.Forwarders = v
	} else {
		data.Forwarders = types.SetValueMust(types.StringType, []attr.Value{})
	}

	if z.Idnsforwardpolicy != nil {
		data.ForwardPolicy = types.StringValue(*z.Idnsforwardpolicy)
	}

	if z.Idnszoneactive != nil {
		data.DisableZone = types.BoolValue(!*z.Idnszoneactive)
	}

	return z, nil
}
