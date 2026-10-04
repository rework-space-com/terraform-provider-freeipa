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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

var _ resource.Resource = &dnsGlobalConfig{}
var _ resource.ResourceWithImportState = &dnsGlobalConfig{}

type dnsGlobalConfig struct {
	client *ipa.Client
}

type dnsGlobalConfigModel struct {
	Id              types.String `tfsdk:"id"`
	Forwarders      types.Set    `tfsdk:"forwarders"`
	ForwardPolicy   types.String `tfsdk:"forward_policy"`
	AllowSyncPTR    types.Bool   `tfsdk:"allow_sync_ptr"`
	DNSServers      types.Set    `tfsdk:"dns_servers"`
	DNSSECKeyMaster types.String `tfsdk:"dnssec_key_master"`
	IPADNSVersion   types.Int64  `tfsdk:"ipa_dns_version"`
}

func NewDNSGlobalConfigResource() resource.Resource {
	return &dnsGlobalConfig{}
}

func (r *dnsGlobalConfig) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_global_config"
}

func (r *dnsGlobalConfig) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "FreeIPA global DNS configuration. Singleton resource — only one instance per IPA realm. Create reads existing values then applies modifications; Delete reverts managed attributes to FreeIPA documented defaults (forwarders cleared, forward_policy=first, allow_sync_ptr=false).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource. Always `global` for this singleton.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"forwarders": schema.SetAttribute{
				MarkdownDescription: "Global forwarders. A custom port can be specified for each forwarder using a standard format `IP_ADDRESS port PORT`.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"forward_policy": schema.StringAttribute{
				MarkdownDescription: "Global forwarding policy. One of `only`, `first`, `none`. Set to `none` to disable any configured global forwarders.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("only", "first", "none"),
				},
			},
			"allow_sync_ptr": schema.BoolAttribute{
				MarkdownDescription: "Allow synchronization of forward (A, AAAA) and reverse (PTR) records.",
				Optional:            true,
			},
			"dns_servers": schema.SetAttribute{
				MarkdownDescription: "Set of IPA masters configured as DNS servers. Read-only.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"dnssec_key_master": schema.StringAttribute{
				MarkdownDescription: "IPA server configured as DNSSec key master. Read-only.",
				Computed:            true,
			},
			"ipa_dns_version": schema.Int64Attribute{
				MarkdownDescription: "IPA DNS version. Read-only.",
				Computed:            true,
			},
		},
	}
}

func (r *dnsGlobalConfig) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsGlobalConfig) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dnsGlobalConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "[DEBUG] Create freeipa dns global config")

	optArgs, hasChange := r.buildModOptArgs(ctx, &data, nil)
	if hasChange {
		if _, err := r.client.DnsconfigMod(&ipa.DnsconfigModArgs{}, optArgs); err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				resp.Diagnostics.AddWarning("Client Warning", err.Error())
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error applying freeipa dns global config: %s", err))
				return
			}
		}
	}

	if _, err := r.readGlobalConfig(ctx, &data, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error refreshing freeipa dns global config after create: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsGlobalConfig) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dnsGlobalConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.readGlobalConfig(ctx, &data, &resp.Diagnostics); err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			resp.Diagnostics.AddError("DNS not installed", "dnsconfig_show returned NotFound — IPA DNS subsystem is not installed on this realm")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading freeipa dns global config: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsGlobalConfig) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state dnsGlobalConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "[DEBUG] Update freeipa dns global config")

	optArgs, hasChange := r.buildModOptArgs(ctx, &data, &state)
	if hasChange {
		if _, err := r.client.DnsconfigMod(&ipa.DnsconfigModArgs{}, optArgs); err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				resp.Diagnostics.AddWarning("Client Warning", err.Error())
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating freeipa dns global config: %s", err))
				return
			}
		}
	}

	if _, err := r.readGlobalConfig(ctx, &data, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error refreshing freeipa dns global config after update: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dnsGlobalConfig) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dnsGlobalConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "[DEBUG] Delete freeipa dns global config — reverting managed attrs to defaults")

	defaultPolicy := "first"
	defaultAllowSync := false
	optArgs := ipa.DnsconfigModOptionalArgs{}
	if !data.ForwardPolicy.IsNull() {
		optArgs.Idnsforwardpolicy = &defaultPolicy
	}
	if !data.AllowSyncPTR.IsNull() {
		optArgs.Idnsallowsyncptr = &defaultAllowSync
	}
	if !data.Forwarders.IsNull() {
		empty := []string{}
		optArgs.Idnsforwarders = &empty
	}

	if optArgs.Idnsforwardpolicy == nil && optArgs.Idnsallowsyncptr == nil && optArgs.Idnsforwarders == nil {
		return
	}

	if _, err := r.client.DnsconfigMod(&ipa.DnsconfigModArgs{}, &optArgs); err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddWarning("Client Warning", err.Error())
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reverting freeipa dns global config on delete: %s", err))
	}
}

func (r *dnsGlobalConfig) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "global" {
		resp.Diagnostics.AddError("Import Error", "freeipa_dns_global_config is a singleton — import id must be the literal string \"global\"")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "global")...)
}

func (r *dnsGlobalConfig) buildModOptArgs(ctx context.Context, data *dnsGlobalConfigModel, prev *dnsGlobalConfigModel) (*ipa.DnsconfigModOptionalArgs, bool) {
	optArgs := ipa.DnsconfigModOptionalArgs{}
	hasChange := false

	if prev == nil || !data.Forwarders.Equal(prev.Forwarders) {
		if !data.Forwarders.IsNull() {
			var v []string
			for _, value := range data.Forwarders.Elements() {
				val, _ := strconv.Unquote(value.String())
				v = append(v, val)
			}
			if v == nil {
				v = []string{}
			}
			optArgs.Idnsforwarders = &v
			hasChange = true
		} else if prev != nil {
			empty := []string{}
			optArgs.Idnsforwarders = &empty
			hasChange = true
		}
	}

	if prev == nil || !data.ForwardPolicy.Equal(prev.ForwardPolicy) {
		if !data.ForwardPolicy.IsNull() {
			s := data.ForwardPolicy.ValueString()
			optArgs.Idnsforwardpolicy = &s
			hasChange = true
		}
	}

	if prev == nil || !data.AllowSyncPTR.Equal(prev.AllowSyncPTR) {
		if !data.AllowSyncPTR.IsNull() {
			b := data.AllowSyncPTR.ValueBool()
			optArgs.Idnsallowsyncptr = &b
			hasChange = true
		}
	}

	return &optArgs, hasChange
}

// readGlobalConfig calls dnsconfig_show and copies remote attrs into the model.
// User-managed attrs are only overwritten if they are non-null in the model
// (so Read does not invent values for attrs the user has not declared).
func (r *dnsGlobalConfig) readGlobalConfig(ctx context.Context, data *dnsGlobalConfigModel, diags *diag.Diagnostics) (*ipa.Dnsconfig, error) {
	all := true
	res, err := r.client.DnsconfigShow(&ipa.DnsconfigShowArgs{}, &ipa.DnsconfigShowOptionalArgs{All: &all})
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Read freeipa dns global config: %s", res.String()))

	cfg := &res.Result

	if cfg.DNSServerServer != nil {
		v, d := types.SetValueFrom(ctx, types.StringType, *cfg.DNSServerServer)
		diags.Append(d...)
		data.DNSServers = v
	} else {
		data.DNSServers = types.SetNull(types.StringType)
	}
	if cfg.DnssecKeyMasterServer != nil {
		data.DNSSECKeyMaster = types.StringValue(*cfg.DnssecKeyMasterServer)
	} else {
		data.DNSSECKeyMaster = types.StringNull()
	}
	if cfg.Ipadnsversion != nil {
		data.IPADNSVersion = types.Int64Value(int64(*cfg.Ipadnsversion))
	} else {
		data.IPADNSVersion = types.Int64Null()
	}

	if !data.Forwarders.IsNull() {
		if cfg.Idnsforwarders != nil {
			v, d := types.SetValueFrom(ctx, types.StringType, *cfg.Idnsforwarders)
			diags.Append(d...)
			data.Forwarders = v
		} else {
			data.Forwarders = types.SetValueMust(types.StringType, []attr.Value{})
		}
	}
	if !data.ForwardPolicy.IsNull() {
		if cfg.Idnsforwardpolicy != nil {
			data.ForwardPolicy = types.StringValue(*cfg.Idnsforwardpolicy)
		} else {
			data.ForwardPolicy = types.StringNull()
		}
	}
	if !data.AllowSyncPTR.IsNull() {
		if cfg.Idnsallowsyncptr != nil {
			data.AllowSyncPTR = types.BoolValue(*cfg.Idnsallowsyncptr)
		} else {
			data.AllowSyncPTR = types.BoolNull()
		}
	}

	data.Id = types.StringValue("global")
	return cfg, nil
}
