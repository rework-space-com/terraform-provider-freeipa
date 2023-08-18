resource "freeipa_hbac_policy" "hbac-0" {
  name            = "test-hbac"
  description     = "Test HBAC policy"
  enabled         = true
  hostcategory    = "all"
  servicecategory = "all"
}

resource "freeipa_hbac_policy_host_membership" "hbac-host-1" {
  name      = "test-hbac"
  host      = "ipaclient1.ipatest.lan"
}
