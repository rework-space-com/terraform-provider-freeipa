// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//   Mixton <maxime.thomas@mtconsulting.tech>
//   Parsa <p.yousefi97@gmail.com>
//   Roman Butsiy <butsiyroman@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}
var _ resource.ResourceWithModifyPlan = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *ipa.Client
}

type UserInterface interface {
	CreateUser(context.Context, resource.CreateRequest, *resource.CreateResponse)
	ReadUser(context.Context, resource.ReadRequest, *resource.ReadResponse)
	UpdateUser(context.Context, resource.UpdateRequest, *resource.UpdateResponse)
}

type ActiveUserResource struct {
	client *ipa.Client
}

type StagedUserResource struct {
	client *ipa.Client
}
type PreservedUserResource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	FirstName              types.String `tfsdk:"first_name"`
	LastName               types.String `tfsdk:"last_name"`
	UID                    types.String `tfsdk:"name"`
	FullName               types.String `tfsdk:"full_name"`
	DisplayName            types.String `tfsdk:"display_name"`
	Initials               types.String `tfsdk:"initials"`
	HomeDirectory          types.String `tfsdk:"home_directory"`
	AuthType               types.Set    `tfsdk:"auth_type"`
	RadiusConfig           types.String `tfsdk:"radius_proxy_config"`
	RadiusUser             types.String `tfsdk:"radius_proxy_username"`
	IdpConfig              types.String `tfsdk:"external_idp_config"`
	IdpUser                types.String `tfsdk:"external_idp_username"`
	Gecos                  types.String `tfsdk:"gecos"`
	LoginShell             types.String `tfsdk:"login_shell"`
	KrbPrincipalName       types.List   `tfsdk:"krb_principal_name"`
	KrbPrincipalExpiration types.String `tfsdk:"krb_principal_expiration"`
	KrbPasswordExpiration  types.String `tfsdk:"krb_password_expiration"`
	UserPassword           types.String `tfsdk:"userpassword"`
	EmailAddress           types.List   `tfsdk:"email_address"`
	TelephoneNumbers       types.List   `tfsdk:"telephone_numbers"`
	MobileNumbers          types.List   `tfsdk:"mobile_numbers"`
	RandomPassword         types.Bool   `tfsdk:"random_password"`
	UidNumber              types.Int32  `tfsdk:"uid_number"`
	GidNumber              types.Int32  `tfsdk:"gid_number"`
	StreetAddress          types.String `tfsdk:"street_address"`
	City                   types.String `tfsdk:"city"`
	Province               types.String `tfsdk:"province"`
	PostalCode             types.String `tfsdk:"postal_code"`
	OrganisationUnit       types.String `tfsdk:"organisation_unit"`
	JobTitle               types.String `tfsdk:"job_title"`
	Manager                types.String `tfsdk:"manager"`
	EmployeeNumber         types.String `tfsdk:"employee_number"`
	EmployeeType           types.String `tfsdk:"employee_type"`
	PreferredLanguage      types.String `tfsdk:"preferred_language"`
	AccountDisabled        types.Bool   `tfsdk:"account_disabled"`
	AccountStaged          types.Bool   `tfsdk:"account_staged"`
	AccountPreserved       types.Bool   `tfsdk:"account_preserved"`
	State                  types.String `tfsdk:"state"`
	SshPublicKeys          types.List   `tfsdk:"ssh_public_key"`
	UserCerts              types.Set    `tfsdk:"user_certificates"`
	CarLicense             types.List   `tfsdk:"car_license"`
	UserClass              types.List   `tfsdk:"userclass"`
	AddAttr                types.List   `tfsdk:"addattr"`
	SetAttr                types.List   `tfsdk:"setattr"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func userSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User resource. The user lifecycle is managed with the attributes:\n\n" +
			"- account_staged\n\n" +
			"- account_preserved\n\n" +
			"(defaults to active)\n\n" +
			"An `active` user can be preserved.\n\n" +
			"A user can be `staged` at the user's creation or from a `preserved`state.\n\n" +
			"A `staged` user can be preserved.\n\n" +
			"A `preserved` or `staged` user can activated.\n\n",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Required:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "UID or Login\n\n	- The name must not exceed 32 characters.\n	- The name must contain only lowercase letters (a-z), digits (0-9), and the characters (. - _).\n	- The name must not start with a special character.\n	- A user and a group cannot have the same name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "Full name",
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Optional:            true,
			},
			"initials": schema.StringAttribute{
				MarkdownDescription: "Initials",
				Optional:            true,
			},
			"home_directory": schema.StringAttribute{
				MarkdownDescription: "Home Directory",
				Optional:            true,
			},
			"auth_type": schema.SetAttribute{
				MarkdownDescription: "User authentication type. Possible values of the elements are (password, radius, otp, pkinit, hardened, idp, passkey)",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("password", "radius", "otp", "pkinit", "hardened", "idp", "passkey")),
				},
			},
			"radius_proxy_config": schema.StringAttribute{
				MarkdownDescription: "RADIUS proxy configuration",
				Optional:            true,
			},
			"radius_proxy_username": schema.StringAttribute{
				MarkdownDescription: "RADIUS proxy username",
				Optional:            true,
			},
			"external_idp_config": schema.StringAttribute{
				MarkdownDescription: "External IdP configuration",
				Optional:            true,
			},
			"external_idp_username": schema.StringAttribute{
				MarkdownDescription: "External IdP user identifier",
				Optional:            true,
			},
			"gecos": schema.StringAttribute{
				MarkdownDescription: "GECOS",
				Optional:            true,
			},
			"login_shell": schema.StringAttribute{
				MarkdownDescription: "Login Shell",
				Optional:            true,
			},
			"krb_principal_name": schema.ListAttribute{
				MarkdownDescription: "Principal alias",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"krb_principal_expiration": schema.StringAttribute{
				MarkdownDescription: "Kerberos principal expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Optional: true,
			},
			"krb_password_expiration": schema.StringAttribute{
				MarkdownDescription: "User password expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Optional: true,
			},
			"userpassword": schema.StringAttribute{
				MarkdownDescription: "Prompt to set the user password",
				Optional:            true,
				Sensitive:           true,
			},
			"email_address": schema.ListAttribute{
				MarkdownDescription: "Email address",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"telephone_numbers": schema.ListAttribute{
				MarkdownDescription: "Telephone Number",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mobile_numbers": schema.ListAttribute{
				MarkdownDescription: "Mobile Number",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"random_password": schema.BoolAttribute{
				MarkdownDescription: "Generate a random user password",
				Optional:            true,
			},
			"uid_number": schema.Int32Attribute{
				MarkdownDescription: "User ID Number (system will assign one if not provided)",
				Optional:            true,
			},
			"gid_number": schema.Int32Attribute{
				MarkdownDescription: "Group ID Number",
				Optional:            true,
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "Street address",
				Optional:            true,
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "City",
				Optional:            true,
			},
			"province": schema.StringAttribute{
				MarkdownDescription: "Province/State/Country",
				Optional:            true,
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "Postal code",
				Optional:            true,
			},
			"organisation_unit": schema.StringAttribute{
				MarkdownDescription: "Org. Unit",
				Optional:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "Job Title",
				Optional:            true,
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "Manager",
				Optional:            true,
			},
			"employee_number": schema.StringAttribute{
				MarkdownDescription: "Employee Number",
				Optional:            true,
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Employee Type",
				Optional:            true,
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "Preferred Language",
				Optional:            true,
			},
			"account_disabled": schema.BoolAttribute{
				MarkdownDescription: "Account disabled.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_staged": schema.BoolAttribute{
				MarkdownDescription: "Account staged.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_preserved": schema.BoolAttribute{
				MarkdownDescription: "Account preserved.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the user, can be `active`, `staged`, or `preserved`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssh_public_key": schema.ListAttribute{
				MarkdownDescription: "List of SSH public keys",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"user_certificates": schema.SetAttribute{
				MarkdownDescription: "List of Base-64 encoded user certificates",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"car_license": schema.ListAttribute{
				MarkdownDescription: "Car Licenses",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"userclass": schema.ListAttribute{
				MarkdownDescription: "User category (semantics placed on this attribute are for local interpretation)",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"addattr": schema.ListAttribute{
				MarkdownDescription: "Add an attribute/value pair. Format is attr=value. The attribute must be part of the LDAP schema.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"setattr": schema.ListAttribute{
				MarkdownDescription: "Set an attribute to a name/value pair. Format is attr=value.",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
		},
		Version: 1,
	}
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userSchema()

}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var accountStage types.String
	req.State.GetAttribute(ctx, path.Root("state"), &accountStage)
	var toPreserved, toStaged types.Bool
	req.Plan.GetAttribute(ctx, path.Root("account_preserved"), &toPreserved)
	req.Plan.GetAttribute(ctx, path.Root("account_staged"), &toStaged)

	if toPreserved.ValueBool() && toStaged.ValueBool() {
		resp.Diagnostics.AddError("User Lifecycle", "account_staged and account_preserved cannot be both true.")
		return
	}
	// on delete
	if req.Plan.Raw.IsNull() {
		return
	}
	// create as preserved
	if req.State.Raw.IsNull() && toPreserved.ValueBool() {
		resp.Diagnostics.AddError("User Lifecycle", "Creating a preserved user is not allowed.")
		return
	}
	if !toStaged.ValueBool() && toPreserved.ValueBool() { // to preserved
		if accountStage.Equal(types.StringValue("staged")) {
			resp.Diagnostics.AddError("User Lifecycle", "Preserving a staged user is not allowed.")
			return
		}
		resp.Plan.SetAttribute(ctx, path.Root("state"), types.StringValue("preserved"))
	} else if toStaged.ValueBool() && !toPreserved.ValueBool() { // to staged
		if accountStage.Equal(types.StringValue("active")) { // active -> staged
			resp.Diagnostics.AddError("User Lifecycle", "Staging an active user is not allowed.")
			return
		}
		resp.Plan.SetAttribute(ctx, path.Root("state"), types.StringValue("staged"))
	} else { // to active
		resp.Plan.SetAttribute(ctx, path.Root("state"), types.StringValue("active"))
	}
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel
	var resource UserInterface
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.State.Equal(types.StringValue("active")) {
		resource = ActiveUserResource{client: r.client}
	} else if data.State.Equal(types.StringValue("staged")) {
		resource = StagedUserResource{client: r.client}
	} else {
		resp.Diagnostics.AddError("User Lifecycle", "User can only be created as active or staged.")
		return
	}
	resource.CreateUser(ctx, req, resp)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel
	var resource UserInterface

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.State.Equal(types.StringValue("active")) {
		resource = ActiveUserResource{client: r.client}
	} else if data.State.Equal(types.StringValue("staged")) {
		resource = StagedUserResource{client: r.client}
	} else if data.State.Equal(types.StringValue("preserved")) {
		resource = PreservedUserResource{client: r.client}
	} else {
		resp.Diagnostics.AddError("User Lifecycle", "User can only be created as active or staged.")
		return
	}

	resource.ReadUser(ctx, req, resp)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state UserResourceModel
	var resource UserInterface

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// update active user
	if state.State.IsNull() || state.State.Equal(types.StringValue("active")) {
		if data.State.Equal(types.StringValue("staged")) {
			resp.Diagnostics.AddError("Client Error", "Staging an active user is not authorized.")
			return
		}
		if data.State.Equal(types.StringValue("preserved")) {
			r.PreserveActiveUser(ctx, req, resp)
		}
	}

	// update staged user
	if state.State.Equal(types.StringValue("staged")) {
		if data.State.Equal(types.StringValue("preserved")) {
			resp.Diagnostics.AddError("Client Error", "Preserving a staged user is not authorized.")
			return
		}
		if data.State.Equal(types.StringValue("active")) {
			r.ActivateStagedUser(ctx, req, resp)
		}
	}

	// update preserved user
	if state.State.Equal(types.StringValue("preserved")) {
		if data.State.Equal(types.StringValue("staged")) {
			r.StagePreservedUser(ctx, req, resp)
			return
		}
		if data.State.Equal(types.StringValue("active")) {
			r.ActivatePreservedUser(ctx, req, resp)
			return
		}
	}

	if data.State.Equal(types.StringValue("active")) {
		resource = ActiveUserResource{client: r.client}
	} else if data.State.Equal(types.StringValue("staged")) {
		resource = StagedUserResource{client: r.client}
	} else if data.State.Equal(types.StringValue("active")) {
		resource = PreservedUserResource{client: r.client}
	} else {
		return
	}
	resource.UpdateUser(ctx, req, resp)

}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// update active user
	if state.State.Equal(types.StringValue("active")) || state.State.Equal(types.StringValue("preserved")) {

		optArgs := ipa.UserDelOptionalArgs{}
		optArgs.UID = &[]string{state.UID.ValueString()}

		_, err := r.client.UserDel(&ipa.UserDelArgs{}, &optArgs)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		return
	} else {
		optArgs := ipa.StageuserDelOptionalArgs{}
		optArgs.UID = &[]string{state.UID.ValueString()}

		_, err := r.client.StageuserDel(&ipa.StageuserDelArgs{}, &optArgs)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	all := true
	var uid, state string
	if strings.Contains(req.ID, ";") {
		idelements := strings.SplitN(req.ID, ";", 2)
		uid = idelements[0]
		state = idelements[1]
	} else {
		uid = req.ID
		state = "active"
	}

	if state == "active" {
		optArgs := ipa.UserShowOptionalArgs{
			All: &all,
		}

		optArgs.UID = &uid

		res, err := r.client.UserShow(&ipa.UserShowArgs{}, &optArgs)
		if err != nil {
			resp.Diagnostics.AddError("Import Error", err.Error())
			return
		}
		if res.Result.UID != uid {
			resp.Diagnostics.AddError("Import Error", "The import ID and the name attribute must be identical")
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("first_name"), res.Result.Givenname)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_name"), res.Result.Sn)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), "active")...)
		return
	}

	if state == "staged" {
		optArgs := ipa.StageuserShowOptionalArgs{
			All: &all,
		}

		optArgs.UID = &uid

		res, err := r.client.StageuserShow(&ipa.StageuserShowArgs{}, &optArgs)
		if err != nil {
			resp.Diagnostics.AddError("Import Error", err.Error())
			return
		}
		if res.Result.UID != uid {
			resp.Diagnostics.AddError("Import Error", "The import ID and the name attribute must be identical")
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("first_name"), res.Result.Givenname)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_name"), res.Result.Sn)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_staged"), true)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), "staged")...)
		return
	}

	if state == "preserved" {
		optArgs := ipa.UserShowOptionalArgs{
			All: &all,
		}

		optArgs.UID = &uid

		res, err := r.client.UserShow(&ipa.UserShowArgs{}, &optArgs)
		if err != nil {
			resp.Diagnostics.AddError("Import Error", err.Error())
			return
		}
		if res.Result.UID != uid {
			resp.Diagnostics.AddError("Import Error", "The import ID and the name attribute must be identical")
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), res.Result.UID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("first_name"), res.Result.Givenname)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_name"), res.Result.Sn)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_preserved"), true)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), "preserved")...)
		return
	}
}

func (r *UserResource) ActivateStagedUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state, config UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	res, err := r.client.StageuserActivate(&ipa.StageuserActivateArgs{}, &ipa.StageuserActivateOptionalArgs{All: &all, UID: data.UID.ValueStringPointer()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !data.AccountDisabled.IsNull() && data.AccountDisabled.Equal(types.BoolValue(true)) {
		_, err := r.client.UserDisable(&ipa.UserDisableArgs{}, &ipa.UserDisableOptionalArgs{UID: data.UID.ValueStringPointer()})
		if err != nil && !strings.Contains(err.Error(), "This entry is already disabled") {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
	} else {
		_, err := r.client.UserEnable(&ipa.UserEnableArgs{}, &ipa.UserEnableOptionalArgs{UID: data.UID.ValueStringPointer()})
		if err != nil && !strings.Contains(err.Error(), "This entry is already enabled") {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
	}

	data.State = types.StringValue("active")
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa user %s returns %s", data.UID.String(), res.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *UserResource) ActivatePreservedUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UserUndel(&ipa.UserUndelArgs{}, &ipa.UserUndelOptionalArgs{UID: data.UID.ValueStringPointer()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if !data.AccountDisabled.IsNull() && data.AccountDisabled.Equal(types.BoolValue(true)) {
		_, err := r.client.UserDisable(&ipa.UserDisableArgs{}, &ipa.UserDisableOptionalArgs{UID: data.UID.ValueStringPointer()})
		if err != nil && !strings.Contains(err.Error(), "This entry is already disabled") {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
	} else {
		_, err := r.client.UserEnable(&ipa.UserEnableArgs{}, &ipa.UserEnableOptionalArgs{UID: data.UID.ValueStringPointer()})
		if err != nil && !strings.Contains(err.Error(), "This entry is already enabled") {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
	}
	data.State = types.StringValue("active")

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *UserResource) StagePreservedUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UserStage(&ipa.UserStageArgs{}, &ipa.UserStageOptionalArgs{UID: &[]string{data.UID.ValueString()}})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	data.State = types.StringValue("staged")

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *UserResource) PreserveActiveUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, config UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	preserve := true
	optArgs := ipa.UserDelOptionalArgs{}
	optArgs.UID = &[]string{data.UID.ValueString()}
	optArgs.Preserve = &preserve

	_, err := r.client.UserDel(&ipa.UserDelArgs{}, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
	}
	data.AccountPreserved = types.BoolValue(true)
	data.AccountDisabled = config.AccountDisabled
	data.State = types.StringValue("preserved")
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
