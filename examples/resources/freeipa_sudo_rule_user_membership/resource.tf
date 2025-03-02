resource "freeipa_sudo_rule_user_membership" "user-0" {
  name = "sudo-rule-test"
  user = "user01"
}

resource "freeipa_sudo_rule_user_membership" "users-1" {
  name       = "sudo-rule-test"
  users      = ["user01"]
  identifier = "users-1"
}

resource "freeipa_sudo_rule_user_membership" "group-3" {
  name  = "sudo-rule-test"
  group = "test-group-0"
}

resource "freeipa_sudo_rule_user_membership" "groups-3" {
  name       = "sudo-rule-test"
  groups     = ["test-group-0"]
  identifier = "groups-3"
}
