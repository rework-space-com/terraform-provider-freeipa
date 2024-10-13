package freeipa

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"freeipa": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccFreeIPAProvider() string {
	provider_host := os.Getenv("FREEIPA_HOST")
	provider_user := os.Getenv("FREEIPA_USERNAME")
	provider_pass := os.Getenv("FREEIPA_PASSWORD")
	return fmt.Sprintf(`
	provider "freeipa" {
		host     = "%s"
		username = "%s"
		password = "%s"
		insecure = true
	  }
	
	`, provider_host, provider_user, provider_pass)
}

func testAccFreeIPAGroup_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_group" "group-%s" {
		name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	if dataset["gid_number"] != "" {
		tf_def += fmt.Sprintf("  gid_number = %s\n", dataset["gid_number"])
	}
	if dataset["external"] != "" {
		tf_def += fmt.Sprintf("  external = %s\n", dataset["external"])
	}
	if dataset["nonposix"] != "" {
		tf_def += fmt.Sprintf("  nonposix = %s\n", dataset["nonposix"])
	}
	if dataset["addattr"] != "" {
		tf_def += fmt.Sprintf("  addattr = %s\n", dataset["addattr"])
	}
	if dataset["setattr"] != "" {
		tf_def += fmt.Sprintf("  setattr = %s\n", dataset["setattr"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPAGroup_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_group" "group-%s" {
		name        = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPAUser_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_user" "user-%s" {
	  name        = %s
	  first_name  = %s
	  last_name   = %s
	`, dataset["index"], dataset["login"], dataset["firstname"], dataset["lastname"])
	if dataset["account_disabled"] != "" {
		tf_def += fmt.Sprintf("  account_disabled = %s\n", dataset["account_disabled"])
	}
	if dataset["car_license"] != "" {
		tf_def += fmt.Sprintf("  car_license = %s\n", dataset["car_license"])
	}
	if dataset["city"] != "" {
		tf_def += fmt.Sprintf("  city = %s\n", dataset["city"])
	}
	if dataset["display_name"] != "" {
		tf_def += fmt.Sprintf("  display_name = %s\n", dataset["display_name"])
	}
	if dataset["email_address"] != "" {
		tf_def += fmt.Sprintf("  email_address = %s\n", dataset["email_address"])
	}
	if dataset["employee_number"] != "" {
		tf_def += fmt.Sprintf("  employee_number = %s\n", dataset["employee_number"])
	}
	if dataset["employee_type"] != "" {
		tf_def += fmt.Sprintf("  employee_type = %s\n", dataset["employee_type"])
	}
	if dataset["full_name"] != "" {
		tf_def += fmt.Sprintf("  full_name = %s\n", dataset["full_name"])
	}
	if dataset["gecos"] != "" {
		tf_def += fmt.Sprintf("  gecos = %s\n", dataset["gecos"])
	}
	if dataset["gid_number"] != "" {
		tf_def += fmt.Sprintf("  gid_number = %s\n", dataset["gid_number"])
	}
	if dataset["home_directory"] != "" {
		tf_def += fmt.Sprintf("  home_directory = %s\n", dataset["home_directory"])
	}
	if dataset["initials"] != "" {
		tf_def += fmt.Sprintf("  initials = %s\n", dataset["initials"])
	}
	if dataset["job_title"] != "" {
		tf_def += fmt.Sprintf("  job_title = %s\n", dataset["job_title"])
	}
	if dataset["krb_principal_name"] != "" {
		tf_def += fmt.Sprintf("  krb_principal_name = %s\n", dataset["krb_principal_name"])
	}
	if dataset["login_shell"] != "" {
		tf_def += fmt.Sprintf("  login_shell = %s\n", dataset["login_shell"])
	}
	if dataset["manager"] != "" {
		tf_def += fmt.Sprintf("  manager = %s\n", dataset["manager"])
	}
	if dataset["mobile_numbers"] != "" {
		tf_def += fmt.Sprintf("  mobile_numbers = %s\n", dataset["mobile_numbers"])
	}
	if dataset["organisation_unit"] != "" {
		tf_def += fmt.Sprintf("  organisation_unit = %s\n", dataset["organisation_unit"])
	}
	if dataset["postal_code"] != "" {
		tf_def += fmt.Sprintf("  postal_code = %s\n", dataset["postal_code"])
	}
	if dataset["preferred_language"] != "" {
		tf_def += fmt.Sprintf("  preferred_language = %s\n", dataset["preferred_language"])
	}
	if dataset["province"] != "" {
		tf_def += fmt.Sprintf("  province = %s\n", dataset["province"])
	}
	if dataset["random_password"] != "" {
		tf_def += fmt.Sprintf("  random_password = %s\n", dataset["random_password"])
	}
	if dataset["ssh_public_key"] != "" {
		tf_def += fmt.Sprintf("  ssh_public_key = %s\n", dataset["ssh_public_key"])
	}
	if dataset["street_address"] != "" {
		tf_def += fmt.Sprintf("  street_address = %s\n", dataset["street_address"])
	}
	if dataset["telephone_numbers"] != "" {
		tf_def += fmt.Sprintf("  telephone_numbers = %s\n", dataset["telephone_numbers"])
	}
	if dataset["uid_number"] != "" {
		tf_def += fmt.Sprintf("  uid_number = %s\n", dataset["uid_number"])
	}
	if dataset["userpassword"] != "" {
		tf_def += fmt.Sprintf("  userpassword = %s\n", dataset["userpassword"])
	}
	if dataset["krb_principal_expiration"] != "" {
		tf_def += fmt.Sprintf("  krb_principal_expiration = %s\n", dataset["krb_principal_expiration"])
	}
	if dataset["krb_password_expiration"] != "" {
		tf_def += fmt.Sprintf("  krb_password_expiration = %s\n", dataset["krb_password_expiration"])
	}
	if dataset["userclass"] != "" {
		tf_def += fmt.Sprintf("  userclass = %s\n", dataset["userclass"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPAUser_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_user" "user-%s" {
		name        = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPAUserGroupMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_user_group_membership" "membership-%s" {
	  name  = %s
	`, dataset["index"], dataset["name"])

	if dataset["user"] != "" {
		tf_def += fmt.Sprintf("  user = %s\n", dataset["user"])
	}
	if dataset["group"] != "" {
		tf_def += fmt.Sprintf("  group = %s\n", dataset["group"])
	}
	if dataset["external_member"] != "" {
		tf_def += fmt.Sprintf("  external_member = %s\n", dataset["external_member"])
	}
	if dataset["users"] != "" {
		tf_def += fmt.Sprintf("  users = %s\n", dataset["users"])
	}
	if dataset["groups"] != "" {
		tf_def += fmt.Sprintf("  groups = %s\n", dataset["groups"])
	}
	if dataset["external_members"] != "" {
		tf_def += fmt.Sprintf("  external_members = %s\n", dataset["external_members"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPADNSZone_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_dns_zone" "dns-zone-%s" {
	  zone_name  = %s
	`, dataset["index"], dataset["zone_name"])

	if dataset["admin_email_address"] != "" {
		tf_def += fmt.Sprintf("  admin_email_address = %s\n", dataset["admin_email_address"])
	}
	if dataset["allow_inline_dnssec_signing"] != "" {
		tf_def += fmt.Sprintf("  allow_inline_dnssec_signing = %s\n", dataset["allow_inline_dnssec_signing"])
	}
	if dataset["allow_ptr_sync"] != "" {
		tf_def += fmt.Sprintf("  allow_prt_sync = %s\n", dataset["allow_prt_sync"])
	}
	if dataset["allow_query"] != "" {
		tf_def += fmt.Sprintf("  allow_query = %s\n", dataset["allow_query"])
	}
	if dataset["allow_transfer"] != "" {
		tf_def += fmt.Sprintf("  allow_transfer = %s\n", dataset["allow_transfer"])
	}
	if dataset["authoritative_nameserver"] != "" {
		tf_def += fmt.Sprintf("  authoritative_nameserver = %s\n", dataset["authoritative_nameserver"])
	}
	if dataset["bind_update_policy"] != "" {
		tf_def += fmt.Sprintf("  bind_update_policy = %s\n", dataset["bind_update_policy"])
	}
	if dataset["default_ttl"] != "" {
		tf_def += fmt.Sprintf("  default_ttl = %s\n", dataset["default_ttl"])
	}
	if dataset["disable_zone"] != "" {
		tf_def += fmt.Sprintf("  disable_zone = %s\n", dataset["disable_zone"])
	}
	if dataset["dynamic_updates"] != "" {
		tf_def += fmt.Sprintf("  dynamic_updates = %s\n", dataset["dynamic_updates"])
	}
	if dataset["is_reverse_zone"] != "" {
		tf_def += fmt.Sprintf("  is_reverse_zone = %s\n", dataset["is_reverse_zone"])
	}
	if dataset["nsec3param_record"] != "" {
		tf_def += fmt.Sprintf("  nsec3param_record = %s\n", dataset["nsec3param_record"])
	}
	if dataset["skip_nameserver_check"] != "" {
		tf_def += fmt.Sprintf("  skip_nameserver_check = %s\n", dataset["skip_nameserver_check"])
	}
	if dataset["skip_overlap_check"] != "" {
		tf_def += fmt.Sprintf("  skip_overlap_check = %s\n", dataset["skip_overlap_check"])
	}
	if dataset["soa_expire"] != "" {
		tf_def += fmt.Sprintf("  soa_expire = %s\n", dataset["soa_expire"])
	}
	if dataset["soa_minimum"] != "" {
		tf_def += fmt.Sprintf("  soa_minimum = %s\n", dataset["soa_minimum"])
	}
	if dataset["soa_refresh"] != "" {
		tf_def += fmt.Sprintf("  soa_refresh = %s\n", dataset["soa_refresh"])
	}
	if dataset["soa_retry"] != "" {
		tf_def += fmt.Sprintf("  soa_retry = %s\n", dataset["soa_retry"])
	}
	if dataset["soa_serial_number"] != "" {
		tf_def += fmt.Sprintf("  soa_serial_number = %s\n", dataset["soa_serial_number"])
	}
	if dataset["ttl"] != "" {
		tf_def += fmt.Sprintf("  ttl = %s\n", dataset["ttl"])
	}
	if dataset["zone_forwarders"] != "" {
		tf_def += fmt.Sprintf("  zone_forwarders = %s\n", dataset["zone_forwarders"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPADNSZone_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_dns_zone" "dns-zone-%s" {
		zone_name       = %s
	}
	`, dataset["index"], dataset["zone_name"])
}

func testAccFreeIPADNSRecord_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_dns_record" "dns-record-%s" {
	  zone_name = %s
	  name      = %s
	  type      = %s
	  records   = %s

	`, dataset["index"], dataset["zone_name"], dataset["name"], dataset["type"], dataset["records"])

	if dataset["ttl"] != "" {
		tf_def += fmt.Sprintf("  ttl = %s\n", dataset["ttl"])
	}
	if dataset["set_identifier"] != "" {
		tf_def += fmt.Sprintf("  set_identifier = %s\n", dataset["set_identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPAHost_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_host" "host-%s" {
	  name        = %s
	  ip_address  = %s
	`, dataset["index"], dataset["name"], dataset["ip_address"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	if dataset["locality"] != "" {
		tf_def += fmt.Sprintf("  locality = %s\n", dataset["locality"])
	}
	if dataset["location"] != "" {
		tf_def += fmt.Sprintf("  location = %s\n", dataset["location"])
	}
	if dataset["platform"] != "" {
		tf_def += fmt.Sprintf("  platform = %s\n", dataset["platform"])
	}
	if dataset["operating_system"] != "" {
		tf_def += fmt.Sprintf("  operating_system = %s\n", dataset["operating_system"])
	}
	if dataset["user_certificates"] != "" {
		tf_def += fmt.Sprintf("  user_certificates = %s\n", dataset["user_certificates"])
	}
	if dataset["mac_addresses"] != "" {
		tf_def += fmt.Sprintf("  mac_addresses = %s\n", dataset["mac_addresses"])
	}
	if dataset["ipasshpubkeys"] != "" {
		tf_def += fmt.Sprintf("  ipasshpubkeys = %s\n", dataset["ipasshpubkeys"])
	}
	if dataset["userclass"] != "" {
		tf_def += fmt.Sprintf("  userclass = %s\n", dataset["userclass"])
	}
	if dataset["assigned_idview"] != "" {
		tf_def += fmt.Sprintf("  assigned_idview = %s\n", dataset["assigned_idview"])
	}
	if dataset["krb_auth_indicators"] != "" {
		tf_def += fmt.Sprintf("  krb_auth_indicators = %s\n", dataset["krb_auth_indicators"])
	}
	if dataset["krb_preauth"] != "" {
		tf_def += fmt.Sprintf("  krb_preauth = %s\n", dataset["krb_preauth"])
	}
	if dataset["trusted_for_delegation"] != "" {
		tf_def += fmt.Sprintf("  trusted_for_delegation = %s\n", dataset["trusted_for_delegation"])
	}
	if dataset["trusted_to_auth_as_delegate"] != "" {
		tf_def += fmt.Sprintf("  trusted_to_auth_as_delegate = %s\n", dataset["trusted_to_auth_as_delegate"])
	}
	if dataset["force"] != "" {
		tf_def += fmt.Sprintf("  force = %s\n", dataset["force"])
	}
	if dataset["userpassword"] != "" {
		tf_def += fmt.Sprintf("  userpassword = %s\n", dataset["userpassword"])
	}
	if dataset["random_password"] != "" {
		tf_def += fmt.Sprintf("  random_password = %s\n", dataset["random_password"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPAHost_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_host" "host-%s" {
		name = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPAHostGroup_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_hostgroup" "hostgroup-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPAHostGroup_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_hostgroup" "hostgroup-%s" {
		name = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPAHostGroupMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_host_hostgroup_membership" "membership-%s" {
	  name  = %s
	`, dataset["index"], dataset["name"])

	if dataset["host"] != "" {
		tf_def += fmt.Sprintf("  host = %s\n", dataset["host"])
	}
	if dataset["hostgroup"] != "" {
		tf_def += fmt.Sprintf("  hostgroup = %s\n", dataset["hostgroup"])
	}
	if dataset["hosts"] != "" {
		tf_def += fmt.Sprintf("  hosts = %s\n", dataset["hosts"])
	}
	if dataset["hostgroups"] != "" {
		tf_def += fmt.Sprintf("  hostgroups = %s\n", dataset["hostgroups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoCmd_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_cmd" "sudocmd-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoCmdGrp_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_cmdgroup" "sudocmdgroup-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoCmdGrpMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_cmdgroup_membership" "sudocmdgroup-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["sudocmd"] != "" {
		tf_def += fmt.Sprintf("  sudocmd = %s\n", dataset["sudocmd"])
	}
	if dataset["sudocmds"] != "" {
		tf_def += fmt.Sprintf("  sudocmds = %s\n", dataset["sudocmds"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoCmdGroup_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_sudo_cmdgroup" "sudocmdgroup-%s" {
		name = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPASudoRule_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule" "sudorule-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["description"] != "" {
		tf_def += fmt.Sprintf("  description = %s\n", dataset["description"])
	}
	if dataset["enabled"] != "" {
		tf_def += fmt.Sprintf("  enabled = %s\n", dataset["enabled"])
	}
	if dataset["usercategory"] != "" {
		tf_def += fmt.Sprintf("  usercategory = %s\n", dataset["usercategory"])
	}
	if dataset["hostcategory"] != "" {
		tf_def += fmt.Sprintf("  hostcategory = %s\n", dataset["hostcategory"])
	}
	if dataset["commandcategory"] != "" {
		tf_def += fmt.Sprintf("  commandcategory = %s\n", dataset["commandcategory"])
	}
	if dataset["runasusercategory"] != "" {
		tf_def += fmt.Sprintf("  runasusercategory = %s\n", dataset["runasusercategory"])
	}
	if dataset["runasgroupcategory"] != "" {
		tf_def += fmt.Sprintf("  runasgroupcategory = %s\n", dataset["runasgroupcategory"])
	}
	if dataset["order"] != "" {
		tf_def += fmt.Sprintf("  order = %s\n", dataset["order"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRuleOption_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_option" "sudorule-option-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["option"] != "" {
		tf_def += fmt.Sprintf("  option = %s\n", dataset["option"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRule_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_sudo_rule" "sudorule-%s" {
		name = %s
	}
	`, dataset["index"], dataset["name"])
}

func testAccFreeIPASudoAllowCmdMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_allowcmd_membership" "sudo-allow-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["sudocmd"] != "" {
		tf_def += fmt.Sprintf("  sudocmd = %s\n", dataset["sudocmd"])
	}
	if dataset["sudocmds"] != "" {
		tf_def += fmt.Sprintf("  sudocmds = %s\n", dataset["sudocmds"])
	}
	if dataset["sudocmd_group"] != "" {
		tf_def += fmt.Sprintf("  sudocmd_group = %s\n", dataset["sudocmd_group"])
	}
	if dataset["sudocmd_groups"] != "" {
		tf_def += fmt.Sprintf("  sudocmd_groups = %s\n", dataset["sudocmd_groups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoDenyCmdMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_denycmd_membership" "sudo-deny-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["sudocmd"] != "" {
		tf_def += fmt.Sprintf("  sudocmd = %s\n", dataset["sudocmd"])
	}
	if dataset["sudocmds"] != "" {
		tf_def += fmt.Sprintf("  sudocmds = %s\n", dataset["sudocmds"])
	}
	if dataset["sudocmd_group"] != "" {
		tf_def += fmt.Sprintf("  sudocmd_group = %s\n", dataset["sudocmd_group"])
	}
	if dataset["sudocmd_groups"] != "" {
		tf_def += fmt.Sprintf("  sudocmd_groups = %s\n", dataset["sudocmd_groups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRuleHostMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_host_membership" "sudo-host-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["host"] != "" {
		tf_def += fmt.Sprintf("  host = %s\n", dataset["host"])
	}
	if dataset["hosts"] != "" {
		tf_def += fmt.Sprintf("  hosts = %s\n", dataset["hosts"])
	}
	if dataset["hostgroup"] != "" {
		tf_def += fmt.Sprintf("  hostgroup = %s\n", dataset["hostgroup"])
	}
	if dataset["hostgroups"] != "" {
		tf_def += fmt.Sprintf("  hostgroups = %s\n", dataset["hostgroups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRuleUserMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_user_membership" "sudo-user-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["user"] != "" {
		tf_def += fmt.Sprintf("  user = %s\n", dataset["user"])
	}
	if dataset["users"] != "" {
		tf_def += fmt.Sprintf("  users = %s\n", dataset["users"])
	}
	if dataset["group"] != "" {
		tf_def += fmt.Sprintf("  group = %s\n", dataset["group"])
	}
	if dataset["groups"] != "" {
		tf_def += fmt.Sprintf("  groups = %s\n", dataset["groups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRuleRunAsGroupMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_runasgroup_membership" "sudorule-runasgroup-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["runasgroup"] != "" {
		tf_def += fmt.Sprintf("  runasgroup = %s\n", dataset["runasgroup"])
	}
	if dataset["runasgroups"] != "" {
		tf_def += fmt.Sprintf("  runasgroups = %s\n", dataset["runasgroups"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}

func testAccFreeIPASudoRuleRunAsUserMembership_resource(dataset map[string]string) string {
	tf_def := fmt.Sprintf(`
	resource "freeipa_sudo_rule_runasuser_membership" "sudorule-runasuser-membership-%s" {
	  name        = %s
	`, dataset["index"], dataset["name"])
	if dataset["runasuser"] != "" {
		tf_def += fmt.Sprintf("  runasuser = %s\n", dataset["runasuser"])
	}
	if dataset["runasusers"] != "" {
		tf_def += fmt.Sprintf("  runasusers = %s\n", dataset["runasusers"])
	}
	if dataset["identifier"] != "" {
		tf_def += fmt.Sprintf("  identifier = %s\n", dataset["identifier"])
	}
	tf_def += "}\n"
	return tf_def
}
