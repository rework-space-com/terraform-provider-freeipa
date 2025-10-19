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
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

func (r *UserResource) CreateActiveUser(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.UserAddOptionalArgs{}
	args := ipa.UserAddArgs{
		Givenname: string(data.FirstName.ValueString()),
		Sn:        string(data.LastName.ValueString()),
	}
	optArgs.All = &all
	optArgs.UID = data.UID.ValueStringPointer()
	if !data.FullName.IsUnknown() {
		optArgs.Cn = data.FullName.ValueStringPointer()
	}
	if !data.DisplayName.IsUnknown() {
		optArgs.Displayname = data.DisplayName.ValueStringPointer()
	}
	if !data.Initials.IsNull() {
		optArgs.Initials = data.Initials.ValueStringPointer()
	}
	if !data.HomeDirectory.IsNull() {
		optArgs.Homedirectory = data.HomeDirectory.ValueStringPointer()
	}
	if !data.RadiusConfig.IsNull() {
		optArgs.Ipatokenradiusconfiglink = data.RadiusConfig.ValueStringPointer()
	}
	if !data.RadiusUser.IsNull() {
		optArgs.Ipatokenradiususername = data.RadiusUser.ValueStringPointer()
	}
	if !data.IdpConfig.IsNull() {
		optArgs.Ipaidpconfiglink = data.IdpConfig.ValueStringPointer()
	}
	if !data.IdpUser.IsNull() {
		optArgs.Ipaidpsub = data.IdpUser.ValueStringPointer()
	}
	if len(data.AuthType.Elements()) > 0 {
		var v []string
		for _, value := range data.AuthType.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipauserauthtype = &v
	}
	if !data.Gecos.IsNull() {
		optArgs.Gecos = data.Gecos.ValueStringPointer()
	}
	if !data.LoginShell.IsNull() {
		optArgs.Loginshell = data.LoginShell.ValueStringPointer()
	}
	if len(data.KrbPrincipalName.Elements()) > 0 {
		var v []string
		for _, value := range data.KrbPrincipalName.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalname = &v
	}
	if !data.UserPassword.IsNull() {
		optArgs.Userpassword = data.UserPassword.ValueStringPointer()
	}
	if len(data.EmailAddress.Elements()) > 0 {
		var v []string
		for _, value := range data.EmailAddress.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mail = &v
	}
	if len(data.TelephoneNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.TelephoneNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Telephonenumber = &v
	}
	if len(data.MobileNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.MobileNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mobile = &v
	}
	if !data.RandomPassword.IsUnknown() {
		optArgs.Random = data.RandomPassword.ValueBoolPointer()
	}
	if !data.UidNumber.IsNull() {
		uid := int(data.UidNumber.ValueInt32())
		optArgs.Uidnumber = &uid
	}
	if !data.GidNumber.IsNull() {
		gid := int(data.GidNumber.ValueInt32())
		optArgs.Gidnumber = &gid
	}
	if !data.StreetAddress.IsNull() {
		optArgs.Street = data.StreetAddress.ValueStringPointer()
	}
	if !data.City.IsNull() {
		optArgs.L = data.City.ValueStringPointer()
	}
	if !data.Province.IsNull() {
		optArgs.St = data.Province.ValueStringPointer()
	}
	if !data.PostalCode.IsNull() {
		optArgs.Postalcode = data.PostalCode.ValueStringPointer()
	}
	if !data.OrganisationUnit.IsNull() {
		optArgs.Ou = data.OrganisationUnit.ValueStringPointer()
	}
	if !data.JobTitle.IsNull() {
		optArgs.Title = data.JobTitle.ValueStringPointer()
	}
	if !data.Manager.IsNull() {
		optArgs.Manager = data.Manager.ValueStringPointer()
	}
	if !data.EmployeeNumber.IsNull() {
		optArgs.Employeenumber = data.EmployeeNumber.ValueStringPointer()
	}
	if !data.EmployeeType.IsNull() {
		optArgs.Employeetype = data.EmployeeType.ValueStringPointer()
	}
	if !data.PreferredLanguage.IsNull() {
		optArgs.Preferredlanguage = data.PreferredLanguage.ValueStringPointer()
	}
	if len(data.SshPublicKeys.Elements()) > 0 {
		var v []string
		for _, value := range data.SshPublicKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if len(data.UserCerts.Elements()) > 0 {
		var v []interface{}
		for _, value := range data.UserCerts.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Usercertificate = &v
	}
	if len(data.CarLicense.Elements()) > 0 {
		var v []string
		for _, value := range data.CarLicense.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Carlicense = &v
	}
	if !data.KrbPrincipalExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPrincipalExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		optArgs.Krbprincipalexpiration = &timestamp
	}
	if !data.KrbPasswordExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPasswordExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_password_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		optArgs.Krbpasswordexpiration = &timestamp
	}
	if len(data.UserClass.Elements()) > 0 {
		var v []string
		for _, value := range data.UserClass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}
	if len(data.AddAttr.Elements()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa group Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if len(data.SetAttr.Elements()) > 0 {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UserAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group: %s", err))
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
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa user %s returns %s", data.UID.String(), res.String()))

	data.State = types.StringValue("active")
	data.Id = types.StringValue(res.Result.UID)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) ReadActiveUser(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, "[DEBUG] User not found")
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading user %s: %s", data.Id.ValueString(), err))
		}
	}

	data.LastName = types.StringValue(res.Result.Sn)
	data.FirstName = types.StringValue(*res.Result.Givenname)
	if res.Result.Cn != nil && !data.FullName.IsNull() {
		data.FullName = types.StringValue(*res.Result.Cn)
	}
	if res.Result.Displayname != nil && !data.DisplayName.IsNull() {
		data.DisplayName = types.StringValue(*res.Result.Displayname)
	}
	if res.Result.Initials != nil && !data.Initials.IsNull() {
		data.Initials = types.StringValue(*res.Result.Initials)
	}
	if res.Result.Homedirectory != nil && !data.HomeDirectory.IsNull() {
		data.HomeDirectory = types.StringValue(*res.Result.Homedirectory)
	}
	if res.Result.Ipauserauthtype != nil && !data.AuthType.IsNull() {
		data.AuthType, _ = types.SetValueFrom(ctx, types.StringType, res.Result.Ipauserauthtype)
	}
	if res.Result.Ipatokenradiusconfiglink != nil && !data.RadiusConfig.IsNull() {
		data.RadiusConfig = types.StringValue(*res.Result.Ipatokenradiusconfiglink)
	}
	if res.Result.Ipatokenradiususername != nil && !data.RadiusUser.IsNull() {
		data.RadiusUser = types.StringValue(*res.Result.Ipatokenradiususername)
	}
	if res.Result.Ipaidpconfiglink != nil && !data.IdpConfig.IsNull() {
		data.IdpConfig = types.StringValue(*res.Result.Ipaidpconfiglink)
	}
	if res.Result.Ipaidpsub != nil && !data.IdpUser.IsNull() {
		data.IdpUser = types.StringValue(*res.Result.Ipaidpsub)
	}
	if res.Result.Gecos != nil && !data.Gecos.IsNull() {
		data.Gecos = types.StringValue(*res.Result.Gecos)
	}
	if res.Result.Loginshell != nil && !data.LoginShell.IsNull() {
		data.LoginShell = types.StringValue(*res.Result.Loginshell)
	}
	if res.Result.Krbprincipalname != nil && !data.KrbPrincipalName.IsNull() {
		data.KrbPrincipalName, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Krbprincipalname)
	}
	if res.Result.Mail != nil && !data.EmailAddress.IsNull() {
		data.EmailAddress, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mail)
	}
	if res.Result.Telephonenumber != nil && !data.TelephoneNumbers.IsNull() {
		data.TelephoneNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Telephonenumber)
	}
	if res.Result.Mobile != nil && !data.MobileNumbers.IsNull() {
		data.MobileNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mobile)
	}
	if res.Result.Random != nil && !data.RandomPassword.IsNull() {
		data.RandomPassword = types.BoolValue(*res.Result.Random)
	}
	if res.Result.Uidnumber != nil && !data.UidNumber.IsNull() {
		data.UidNumber = types.Int32Value(int32(*res.Result.Uidnumber))
	}
	if res.Result.Gidnumber != nil && !data.GidNumber.IsNull() {
		data.GidNumber = types.Int32Value(int32(*res.Result.Gidnumber))
	}
	if res.Result.Street != nil && !data.StreetAddress.IsNull() {
		data.StreetAddress = types.StringValue(*res.Result.Street)
	}
	if res.Result.L != nil && !data.City.IsNull() {
		data.City = types.StringValue(*res.Result.L)
	}
	if res.Result.St != nil && !data.Province.IsNull() {
		data.Province = types.StringValue(*res.Result.St)
	}
	if res.Result.Postalcode != nil && !data.PostalCode.IsNull() {
		data.PostalCode = types.StringValue(*res.Result.Postalcode)
	}
	if res.Result.Ou != nil && !data.OrganisationUnit.IsNull() {
		data.OrganisationUnit = types.StringValue(*res.Result.Ou)
	}
	if res.Result.Title != nil && !data.JobTitle.IsNull() {
		data.JobTitle = types.StringValue(*res.Result.Title)
	}
	if res.Result.Manager != nil && !data.Manager.IsNull() {
		data.Manager = types.StringValue(*res.Result.Manager)
	}
	if res.Result.Employeenumber != nil && !data.EmployeeNumber.IsNull() {
		data.EmployeeNumber = types.StringValue(*res.Result.Employeenumber)
	}
	if res.Result.Employeetype != nil && !data.EmployeeType.IsNull() {
		data.EmployeeType = types.StringValue(*res.Result.Employeetype)
	}
	if res.Result.Preferredlanguage != nil && !data.PreferredLanguage.IsNull() {
		data.PreferredLanguage = types.StringValue(*res.Result.Preferredlanguage)
	}
	if res.Result.Nsaccountlock != nil && !data.AccountDisabled.IsNull() {
		data.AccountDisabled = types.BoolValue(*res.Result.Nsaccountlock)
	}
	if res.Result.Ipasshpubkey != nil && !data.SshPublicKeys.IsNull() {
		data.SshPublicKeys, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ipasshpubkey)
	}
	if res.Result.Usercertificate != nil && !data.UserCerts.IsNull() {
		var resVals []string
		for _, v := range *res.Result.Usercertificate {
			str := v.([]interface{})[0].(map[string]interface{})["__base64__"]
			resVals = append(resVals, str.(string))
		}
		data.UserCerts, _ = types.SetValueFrom(ctx, types.StringType, resVals)
	}
	if res.Result.Carlicense != nil && !data.CarLicense.IsNull() {
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
	if res.Result.Userclass != nil && !data.UserClass.IsNull() {
		data.UserClass, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Userclass)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserResource) UpdateActiveUser(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state, config UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.UserModOptionalArgs{}
	optArgs.All = &all

	if !data.UID.Equal(state.UID) {
		optArgs.UID = data.UID.ValueStringPointer()
	} else {
		optArgs.UID = state.UID.ValueStringPointer()
	}
	if !state.FullName.IsUnknown() && (!data.FullName.Equal(state.FullName) || !data.FullName.Equal(config.FullName)) {
		optArgs.Cn = data.FullName.ValueStringPointer()
	}
	if !data.FirstName.Equal(state.FirstName) {
		optArgs.Givenname = data.FirstName.ValueStringPointer()
	}
	if !data.LastName.Equal(state.LastName) {
		optArgs.Sn = data.LastName.ValueStringPointer()
	}
	if !data.DisplayName.Equal(state.DisplayName) {
		optArgs.Displayname = data.DisplayName.ValueStringPointer()
	}
	if !data.Initials.Equal(state.Initials) {
		optArgs.Initials = data.Initials.ValueStringPointer()
	}
	if !data.HomeDirectory.Equal(state.HomeDirectory) {
		optArgs.Homedirectory = data.HomeDirectory.ValueStringPointer()
	}
	if !data.AuthType.Equal(state.AuthType) {
		var v []string
		for _, value := range data.AuthType.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipauserauthtype = &v
	}
	if !data.RadiusConfig.Equal(state.RadiusConfig) {
		if data.RadiusConfig.IsNull() {
			optArgs.Ipatokenradiusconfiglink = ipa.String("")
		} else {
			optArgs.Ipatokenradiusconfiglink = data.RadiusConfig.ValueStringPointer()
		}
	}
	if !data.RadiusUser.Equal(state.RadiusUser) {
		if data.RadiusUser.IsNull() {
			optArgs.Ipatokenradiususername = ipa.String("")
		} else {
			optArgs.Ipatokenradiususername = data.RadiusUser.ValueStringPointer()
		}
	}
	if !data.IdpConfig.Equal(state.IdpConfig) {
		if data.IdpConfig.IsNull() {
			optArgs.Ipaidpconfiglink = ipa.String("")
		} else {
			optArgs.Ipaidpconfiglink = data.IdpConfig.ValueStringPointer()
		}
	}
	if !data.IdpUser.Equal(state.IdpUser) {
		if data.IdpUser.IsNull() {
			optArgs.Ipaidpsub = ipa.String("")
		} else {
			optArgs.Ipaidpsub = data.IdpUser.ValueStringPointer()
		}
	}
	if !data.Gecos.Equal(state.Gecos) {
		optArgs.Gecos = data.Gecos.ValueStringPointer()
	}
	if !data.LoginShell.Equal(state.LoginShell) {
		optArgs.Loginshell = data.LoginShell.ValueStringPointer()
	}
	if !data.UserPassword.Equal(state.UserPassword) {
		optArgs.Userpassword = data.UserPassword.ValueStringPointer()
	}
	if !data.RandomPassword.Equal(state.RandomPassword) {
		optArgs.Random = data.RandomPassword.ValueBoolPointer()
	}
	if !data.UidNumber.Equal(state.UidNumber) {
		uid := int(data.UidNumber.ValueInt32())
		optArgs.Uidnumber = &uid
	}
	if !data.GidNumber.Equal(state.GidNumber) {
		gid := int(data.GidNumber.ValueInt32())
		optArgs.Gidnumber = &gid
	}
	if !data.StreetAddress.Equal(state.StreetAddress) {
		optArgs.Street = data.StreetAddress.ValueStringPointer()
	}
	if !data.City.Equal(state.City) {
		optArgs.L = data.City.ValueStringPointer()
	}
	if !data.Province.Equal(state.Province) {
		optArgs.St = data.Province.ValueStringPointer()
	}
	if !data.PostalCode.Equal(state.PostalCode) {
		optArgs.Postalcode = data.PostalCode.ValueStringPointer()
	}
	if !data.OrganisationUnit.Equal(state.OrganisationUnit) {
		optArgs.Ou = data.OrganisationUnit.ValueStringPointer()
	}
	if !data.JobTitle.Equal(state.JobTitle) {
		optArgs.Title = data.JobTitle.ValueStringPointer()
	}
	if !data.EmployeeNumber.Equal(state.EmployeeNumber) {
		optArgs.Employeenumber = data.EmployeeNumber.ValueStringPointer()
	}
	if !data.EmployeeType.Equal(state.EmployeeType) {
		optArgs.Employeetype = data.EmployeeType.ValueStringPointer()
	}
	if !data.PreferredLanguage.Equal(state.PreferredLanguage) {
		optArgs.Preferredlanguage = data.PreferredLanguage.ValueStringPointer()
	}
	if !data.TelephoneNumbers.Equal(state.TelephoneNumbers) {
		var v []string
		for _, value := range data.TelephoneNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Telephonenumber = &v
	}
	if !data.MobileNumbers.Equal(state.MobileNumbers) {
		var v []string
		for _, value := range data.MobileNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mobile = &v
	}
	if !data.KrbPrincipalName.Equal(state.KrbPrincipalName) {
		var v []string
		for _, value := range data.KrbPrincipalName.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalname = &v
	}
	if !data.SshPublicKeys.Equal(state.SshPublicKeys) {
		var v []string
		for _, value := range data.SshPublicKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if !data.UserCerts.Equal(state.UserCerts) {
		var v []interface{}
		for _, value := range data.UserCerts.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Usercertificate = &v
	}
	if !data.CarLicense.Equal(state.CarLicense) {
		var v []string
		for _, value := range data.CarLicense.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Carlicense = &v
	}
	if !data.EmailAddress.Equal(state.EmailAddress) {
		var v []string
		for _, value := range data.EmailAddress.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mail = &v
	}
	if !data.KrbPrincipalExpiration.Equal(state.KrbPrincipalExpiration) {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPrincipalExpiration.ValueString())
		if err != nil && !strings.Contains(err.Error(), "EmptyModlist") {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
		}
		optArgs.Krbprincipalexpiration = &timestamp
	}
	if !data.UserClass.Equal(state.UserClass) {
		var v []string
		for _, value := range data.UserClass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}
	if !data.AddAttr.Equal(state.AddAttr) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa group Addattr %s ", data.AddAttr.String()))
		var v []string

		for _, value := range data.AddAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Addattr = &v
	}

	if !data.SetAttr.Equal(state.SetAttr) {
		var v []string
		for _, value := range data.SetAttr.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Setattr = &v
	}

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UserMod(&ipa.UserModArgs{}, &optArgs)
	if err != nil && !strings.Contains(err.Error(), "EmptyModlist") {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	if !data.AccountDisabled.Equal(state.AccountDisabled) {
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
	}
	data.State = types.StringValue("active")
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa user %s returns %s", data.UID.String(), res.String()))

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
