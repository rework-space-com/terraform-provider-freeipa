---
page_title: "freeipa_host Data Source - freeipa"
description: |-
  FreeIPA Host data source
---

# freeipa_host (Data Source)



## Example Usage

```terraform
data "freeipa_host" "host-0" {
  name = "testhost.example.lan"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Host fully qualified name

	- May contain only letters, numbers, '-'.
	- DNS label may not start or end with '-'

### Optional

- `trusted_for_delegation` (Boolean) Client credentials may be delegated to the service

### Read-Only

- `assigned_idview` (String) Assigned ID View
- `description` (String) A description of this host
- `id` (String) ID of the resource in the terraform state
- `ipasshpubkeys` (List of String) SSH public keys
- `krb_auth_indicators` (List of String) Defines a whitelist for Authentication Indicators. Use 'otp' to allow OTP-based 2FA authentications. Use 'radius' to allow RADIUS-based 2FA authentications. Other values may be used for custom configurations.
- `krb_preauth` (Boolean) Pre-authentication is required for the service
- `locality` (String) Host locality (e.g. 'Baltimore, MD')
- `location` (String) Host location (e.g. 'Lab 2')
- `mac_addresses` (List of String) Hardware MAC address(es) on this host
- `memberof_hbacrule` (List of String) List of HBAC rules this user is member of.
- `memberof_hostgroup` (List of String) List of hostgroups this user is member of.
- `memberof_indirect_hbacrule` (List of String) List of HBAC rules this user is indirectly member of.
- `memberof_indirect_hostgroup` (List of String) List of hostgroups this user is is indirectly member of.
- `memberof_indirect_sudorule` (List of String) List of SUDO rules this user is is indirectly member of.
- `memberof_sudorule` (List of String) List of SUDO rules this user is member of.
- `operating_system` (String) Host operating system and version (e.g. 'Fedora 40')
- `platform` (String) Host hardware platform (e.g. 'Lenovo T61')
- `trusted_to_auth_as_delegate` (Boolean) The service is allowed to authenticate on behalf of a client
- `user_certificates` (List of String) Base-64 encoded host certificate
- `userclass` (List of String) Host category (semantics placed on this attribute are for local interpretation)
