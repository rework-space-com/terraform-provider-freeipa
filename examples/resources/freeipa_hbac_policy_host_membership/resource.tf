resource "freeipa_hbac_policy" "hbac-0" {
  name        = "test-hbac"
  description = "Test HBAC policy"
  enabled     = true
}

resource "freeipa_hbac_policy_host_membership" "hbac-host-1" {
  name = "test-hbac"
  host = "ipaclient1.ipatest.lan"
}

resource "freeipa_hbac_policy_host_membership" "hbac-hosts-1" {
  name       = "test-hbac"
  hosts      = ["ipaclient1.ipatest.lan", "ipaclient2.ipatest.lan"]
  identifier = "hbac-hosts-1"
}

resource "freeipa_hbac_policy_host_membership" "hostgroup-3" {
  name      = "test-hbac"
  hostgroup = "test-hostgroup"
}

resource "freeipa_hbac_policy_host_membership" "hostgroups-3" {
  name       = "test-hbac"
  hostgroups = ["test-hostgroup", "test-hostgroup-2"]
  identifier = "hostgroups-3"
}