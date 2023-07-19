resource "freeipa_hbac_policy" "hbac-0" {
  name            = "test-hbac"
  description     = "Test HBAC policy"
  enabled         = true
  hostcategory    = "all"
  servicecategory = "all"
}

resource "freeipa_hbac_policy_user_membership" "hbac-group-1" {
  name      = "test-hbac"
  group = "usergroup-1"
}
