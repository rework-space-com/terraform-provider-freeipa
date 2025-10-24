package freeipa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Other Resource methods are omitted in this example
var _ resource.ResourceWithUpgradeState = &UserResource{}

// UserResourceModelV0 describes the resource data model when upgrading from a Version 0 of the schema
type UserResourceModelV0 struct {
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
	SshPublicKeys          types.List   `tfsdk:"ssh_public_key"`
	UserCerts              types.Set    `tfsdk:"user_certificates"`
	CarLicense             types.List   `tfsdk:"car_license"`
	UserClass              types.List   `tfsdk:"userclass"`
}

func userSchemaV0() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Required: true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"full_name": schema.StringAttribute{
				Optional: true,
			},
			"display_name": schema.StringAttribute{
				Optional: true,
			},
			"initials": schema.StringAttribute{
				Optional: true,
			},
			"home_directory": schema.StringAttribute{
				Optional: true,
			},
			"auth_type": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"radius_proxy_config": schema.StringAttribute{
				Optional: true,
			},
			"radius_proxy_username": schema.StringAttribute{
				Optional: true,
			},
			"external_idp_config": schema.StringAttribute{
				Optional: true,
			},
			"external_idp_username": schema.StringAttribute{
				Optional: true,
			},
			"gecos": schema.StringAttribute{
				Optional: true,
			},
			"login_shell": schema.StringAttribute{
				Optional: true,
			},
			"krb_principal_name": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"krb_principal_expiration": schema.StringAttribute{
				Optional: true,
			},
			"krb_password_expiration": schema.StringAttribute{
				Optional: true,
			},
			"userpassword": schema.StringAttribute{
				Optional: true,
			},
			"email_address": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"telephone_numbers": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"mobile_numbers": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"random_password": schema.BoolAttribute{
				Optional: true,
			},
			"uid_number": schema.Int32Attribute{
				Optional: true,
			},
			"gid_number": schema.Int32Attribute{
				Optional: true,
			},
			"street_address": schema.StringAttribute{
				Optional: true,
			},
			"city": schema.StringAttribute{
				Optional: true,
			},
			"province": schema.StringAttribute{
				Optional: true,
			},
			"postal_code": schema.StringAttribute{
				Optional: true,
			},
			"organisation_unit": schema.StringAttribute{
				Optional: true,
			},
			"job_title": schema.StringAttribute{
				Optional: true,
			},
			"manager": schema.StringAttribute{
				Optional: true,
			},
			"employee_number": schema.StringAttribute{
				Optional: true,
			},
			"employee_type": schema.StringAttribute{
				Optional: true,
			},
			"preferred_language": schema.StringAttribute{
				Optional: true,
			},
			"account_disabled": schema.BoolAttribute{
				Optional: true,
			},
			"ssh_public_key": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_certificates": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"car_license": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"userclass": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *UserResource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := userSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeUserStateV0toV1,
		},
	}
}

func upgradeUserStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {

	var userDataV0 UserResourceModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &userDataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &userDataV0)...)

	if resp.Diagnostics.HasError() {
		return
	}

	upgradedStateData := UserResourceModel{
		Id:                     userDataV0.Id,
		FirstName:              userDataV0.FirstName,
		LastName:               userDataV0.LastName,
		UID:                    userDataV0.UID,
		FullName:               userDataV0.FullName,
		DisplayName:            userDataV0.DisplayName,
		Initials:               userDataV0.Initials,
		HomeDirectory:          userDataV0.HomeDirectory,
		AuthType:               userDataV0.AuthType,
		RadiusConfig:           userDataV0.RadiusConfig,
		RadiusUser:             userDataV0.RadiusUser,
		IdpConfig:              userDataV0.IdpConfig,
		IdpUser:                userDataV0.IdpUser,
		Gecos:                  userDataV0.Gecos,
		LoginShell:             userDataV0.LoginShell,
		KrbPrincipalName:       userDataV0.KrbPrincipalName,
		KrbPrincipalExpiration: userDataV0.KrbPrincipalExpiration,
		KrbPasswordExpiration:  userDataV0.KrbPasswordExpiration,
		UserPassword:           userDataV0.UserPassword,
		EmailAddress:           userDataV0.EmailAddress,
		TelephoneNumbers:       userDataV0.TelephoneNumbers,
		MobileNumbers:          userDataV0.MobileNumbers,
		RandomPassword:         userDataV0.RandomPassword,
		UidNumber:              userDataV0.UidNumber,
		GidNumber:              userDataV0.GidNumber,
		StreetAddress:          userDataV0.StreetAddress,
		City:                   userDataV0.City,
		Province:               userDataV0.Province,
		PostalCode:             userDataV0.PostalCode,
		OrganisationUnit:       userDataV0.OrganisationUnit,
		JobTitle:               userDataV0.JobTitle,
		Manager:                userDataV0.Manager,
		EmployeeNumber:         userDataV0.EmployeeNumber,
		EmployeeType:           userDataV0.EmployeeNumber,
		PreferredLanguage:      userDataV0.PreferredLanguage,
		AccountDisabled:        userDataV0.AccountDisabled,
		SshPublicKeys:          userDataV0.SshPublicKeys,
		UserCerts:              userDataV0.UserCerts,
		CarLicense:             userDataV0.CarLicense,
		UserClass:              userDataV0.UserClass,
	}

	upgradedStateData.State = types.StringValue("active")
	upgradedStateData.AccountPreserved = types.BoolNull()
	upgradedStateData.AccountStaged = types.BoolNull()
	upgradedStateData.SetAttr = types.ListNull(types.StringType)
	upgradedStateData.AddAttr = types.ListNull(types.StringType)

	resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
}
