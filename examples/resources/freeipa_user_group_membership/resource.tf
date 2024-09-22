resource "freeipa_user_group_membership" "test-0" {
  name = "test-group-2"
  user = "roman"
}

resource "freeipa_user_group_membership" "test-1" {
  name  = "test-group-2"
  group = "test-group"
}

resource "freeipa_user_group_membership" "test-2" {
  name  = "test-group-2"
  external_member = "domain users@adtest.lan"
}
