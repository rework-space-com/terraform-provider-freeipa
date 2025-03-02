resource "freeipa_hbac_policy" "hbac-0" {
  name            = "test-hbac"
  description     = "Test HBAC policy"
  enabled         = true
  hostcategory    = "all"
  servicecategory = "all"
}

resource "freeipa_hbac_policy_user_membership" "hbac-user-1" {
  name = "test-hbac"
  user = "user-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-users-1" {
  name       = "test-hbac"
  users      = ["user-2", "user-3"]
  identifier = "hbac-users-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-group-1" {
  name  = "test-hbac"
  group = "usergroup-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-groups-1" {
  name       = "test-hbac"
  groups     = ["usergroup-2", "usergroup-3"]
  identifier = "hbac-groups-1"
}
