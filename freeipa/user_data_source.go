// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}
var _ datasource.DataSourceWithConfigure = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the resource implementation.
type UserDataSource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type UserDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	FirstName                types.String `tfsdk:"first_name"`
	LastName                 types.String `tfsdk:"last_name"`
	UID                      types.String `tfsdk:"name"`
	FullName                 types.String `tfsdk:"full_name"`
	DisplayName              types.String `tfsdk:"display_name"`
	Initials                 types.String `tfsdk:"initials"`
	HomeDirectory            types.String `tfsdk:"home_directory"`
	Gecos                    types.String `tfsdk:"gecos"`
	LoginShell               types.String `tfsdk:"login_shell"`
	KrbPrincipalName         types.List   `tfsdk:"krb_principal_name"`
	KrbPrincipalExpiration   types.String `tfsdk:"krb_principal_expiration"`
	KrbPasswordExpiration    types.String `tfsdk:"krb_password_expiration"`
	EmailAddress             types.List   `tfsdk:"email_address"`
	TelephoneNumbers         types.List   `tfsdk:"telephone_numbers"`
	MobileNumbers            types.List   `tfsdk:"mobile_numbers"`
	RandomPassword           types.Bool   `tfsdk:"random_password"`
	UidNumber                types.Int32  `tfsdk:"uid_number"`
	GidNumber                types.Int32  `tfsdk:"gid_number"`
	StreetAddress            types.String `tfsdk:"street_address"`
	City                     types.String `tfsdk:"city"`
	Province                 types.String `tfsdk:"province"`
	PostalCode               types.String `tfsdk:"postal_code"`
	OrganisationUnit         types.String `tfsdk:"organisation_unit"`
	JobTitle                 types.String `tfsdk:"job_title"`
	Manager                  types.String `tfsdk:"manager"`
	EmployeeNumber           types.String `tfsdk:"employee_number"`
	EmployeeType             types.String `tfsdk:"employee_type"`
	PreferredLanguage        types.String `tfsdk:"preferred_language"`
	AccountDisabled          types.Bool   `tfsdk:"account_disabled"`
	SshPublicKeys            types.List   `tfsdk:"ssh_public_key"`
	CarLicense               types.List   `tfsdk:"car_license"`
	UserClass                types.List   `tfsdk:"userclass"`
	UserPassword             types.String `tfsdk:"user_password"`
	MemberOfGroup            types.List   `tfsdk:"memberof_group"`
	MemberOfSudoRule         types.List   `tfsdk:"memberof_sudorule"`
	MemberOfHBACRule         types.List   `tfsdk:"memberof_hbacrule"`
	MemberOfIndirectGroup    types.List   `tfsdk:"memberof_indirect_group"`
	MemberOfIndirectSudoRule types.List   `tfsdk:"memberof_indirect_sudorule"`
	MemberOfIndirectHBACRule types.List   `tfsdk:"memberof_indirect_hbacrule"`
}

func (r *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource in the terraform state",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "UID or login",
				Required:            true,
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "Full name",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Computed:            true,
			},
			"initials": schema.StringAttribute{
				MarkdownDescription: "Initials",
				Computed:            true,
			},
			"home_directory": schema.StringAttribute{
				MarkdownDescription: "Home Directory",
				Computed:            true,
			},
			"gecos": schema.StringAttribute{
				MarkdownDescription: "GECOS",
				Computed:            true,
			},
			"login_shell": schema.StringAttribute{
				MarkdownDescription: "Login Shell",
				Computed:            true,
			},
			"krb_principal_name": schema.ListAttribute{
				MarkdownDescription: "Principal alias",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"krb_principal_expiration": schema.StringAttribute{
				MarkdownDescription: "Kerberos principal expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Computed: true,
			},
			"krb_password_expiration": schema.StringAttribute{
				MarkdownDescription: "User password expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Computed: true,
			},
			"email_address": schema.ListAttribute{
				MarkdownDescription: "Email address",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"telephone_numbers": schema.ListAttribute{
				MarkdownDescription: "Telephone Number",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"mobile_numbers": schema.ListAttribute{
				MarkdownDescription: "Mobile Number",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"random_password": schema.BoolAttribute{
				MarkdownDescription: "Generate a random user password",
				Computed:            true,
			},
			"uid_number": schema.Int32Attribute{
				MarkdownDescription: "User ID Number (system will assign one if not provided)",
				Computed:            true,
			},
			"gid_number": schema.Int32Attribute{
				MarkdownDescription: "Group ID Number",
				Computed:            true,
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "Street address",
				Computed:            true,
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "City",
				Computed:            true,
			},
			"province": schema.StringAttribute{
				MarkdownDescription: "Province/State/Country",
				Computed:            true,
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "Postal code",
				Computed:            true,
			},
			"organisation_unit": schema.StringAttribute{
				MarkdownDescription: "Org. Unit",
				Computed:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "Job Title",
				Computed:            true,
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "Manager",
				Computed:            true,
			},
			"employee_number": schema.StringAttribute{
				MarkdownDescription: "Employee Number",
				Computed:            true,
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Employee Type",
				Computed:            true,
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "Preferred Language",
				Computed:            true,
			},
			"account_disabled": schema.BoolAttribute{
				MarkdownDescription: "Account disabled",
				Computed:            true,
			},
			"ssh_public_key": schema.ListAttribute{
				MarkdownDescription: "List of SSH public keys",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"car_license": schema.ListAttribute{
				MarkdownDescription: "Car Licenses",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"userclass": schema.ListAttribute{
				MarkdownDescription: "User category (semantics placed on this attribute are for local interpretation)",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"user_password": schema.StringAttribute{
				MarkdownDescription: "User password",
				Computed:            true,
				Sensitive:           true,
			},
			"memberof_group": schema.ListAttribute{
				MarkdownDescription: "List of groups this user is member of.",
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
			"memberof_indirect_group": schema.ListAttribute{
				MarkdownDescription: "List of groups this user is is indirectly member of.",
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

func (r *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.UserShowOptionalArgs{
		All: &all,
	}

	optArgs.UID = data.UID.ValueStringPointer()

	res, err := r.client.UserShow(&ipa.UserShowArgs{}, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading user %s: %s", data.UID.ValueString(), err))
		return
	}

	if res.Result.Cn != nil {
		data.FullName = types.StringValue(*res.Result.Cn)
	}
	if res.Result.Displayname != nil {
		data.DisplayName = types.StringValue(*res.Result.Displayname)
	}
	if res.Result.Initials != nil {
		data.Initials = types.StringValue(*res.Result.Initials)
	}
	if res.Result.Homedirectory != nil {
		data.HomeDirectory = types.StringValue(*res.Result.Homedirectory)
	}
	if res.Result.Gecos != nil {
		data.Gecos = types.StringValue(*res.Result.Gecos)
	}
	if res.Result.Initials != nil {
		data.LoginShell = types.StringValue(*res.Result.Loginshell)
	}
	if res.Result.Krbprincipalname != nil {
		data.KrbPrincipalName, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Krbprincipalname)
	}
	if res.Result.Mail != nil {
		data.EmailAddress, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mail)
	}
	if res.Result.Telephonenumber != nil {
		data.TelephoneNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Telephonenumber)
	}
	if res.Result.Mobile != nil {
		data.MobileNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mobile)
	}
	if res.Result.Random != nil {
		data.RandomPassword = types.BoolValue(*res.Result.Random)
	}
	if res.Result.Uidnumber != nil {
		data.UidNumber = types.Int32Value(int32(*res.Result.Uidnumber))
	}
	if res.Result.Gidnumber != nil {
		data.GidNumber = types.Int32Value(int32(*res.Result.Gidnumber))
	}
	if res.Result.Street != nil {
		data.StreetAddress = types.StringValue(*res.Result.Street)
	}
	if res.Result.L != nil {
		data.City = types.StringValue(*res.Result.L)
	}
	if res.Result.St != nil {
		data.Province = types.StringValue(*res.Result.St)
	}
	if res.Result.Postalcode != nil {
		data.PostalCode = types.StringValue(*res.Result.Postalcode)
	}
	if res.Result.Ou != nil {
		data.OrganisationUnit = types.StringValue(*res.Result.Ou)
	}
	if res.Result.Title != nil {
		data.JobTitle = types.StringValue(*res.Result.Title)
	}
	if res.Result.Manager != nil {
		data.Manager = types.StringValue(*res.Result.Manager)
	}
	if res.Result.Employeenumber != nil {
		data.EmployeeNumber = types.StringValue(*res.Result.Employeenumber)
	}
	if res.Result.Employeetype != nil {
		data.EmployeeType = types.StringValue(*res.Result.Employeetype)
	}
	if res.Result.Preferredlanguage != nil {
		data.PreferredLanguage = types.StringValue(*res.Result.Preferredlanguage)
	}
	if res.Result.Nsaccountlock != nil {
		data.AccountDisabled = types.BoolValue(*res.Result.Nsaccountlock)
	}
	if res.Result.Ipasshpubkey != nil && !data.SshPublicKeys.IsNull() {
		data.SshPublicKeys, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ipasshpubkey)
	}
	if res.Result.Carlicense != nil {
		data.CarLicense, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Carlicense)
	}
	if res.Result.Krbprincipalexpiration != nil && !data.KrbPrincipalExpiration.IsNull() {
		timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", res.Result.Krbprincipalexpiration.String())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		data.KrbPrincipalExpiration = types.StringValue(timestamp.Format(time.RFC3339))
	}
	if res.Result.Krbpasswordexpiration != nil && !data.KrbPasswordExpiration.IsNull() {
		timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", res.Result.Krbpasswordexpiration.String())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		data.KrbPasswordExpiration = types.StringValue(timestamp.Format(time.RFC3339))
	}
	if res.Result.Userpassword != nil {
		data.UserPassword = types.StringValue(*res.Result.Userpassword)
	}
	if res.Result.Userclass != nil {
		data.UserClass, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Userclass)
	}
	if res.Result.MemberofGroup != nil {
		data.MemberOfGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofGroup)
	}
	if res.Result.MemberofHbacrule != nil {
		data.MemberOfHBACRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofHbacrule)
	}
	if res.Result.MemberofSudorule != nil {
		data.MemberOfSudoRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofSudorule)
	}
	if res.Result.MemberofindirectGroup != nil {
		data.MemberOfIndirectGroup, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectGroup)
	}
	if res.Result.MemberofindirectHbacrule != nil {
		data.MemberOfIndirectHBACRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectHbacrule)
	}
	if res.Result.MemberofindirectSudorule != nil {
		data.MemberOfIndirectSudoRule, _ = types.ListValueFrom(ctx, types.StringType, res.Result.MemberofindirectSudorule)
	}

	data.Id = types.StringValue(data.UID.ValueString())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
