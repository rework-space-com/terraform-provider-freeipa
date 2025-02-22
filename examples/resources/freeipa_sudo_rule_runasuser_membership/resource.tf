resource "freeipa_sudo_rule_runasuser_membership" "user-0" {
  name      = "sudo-rule-test"
  runasuser = "user01"
}

resource "freeipa_sudo_rule_runasuser_membership" "users-0" {
  name       = "sudo-rule-test"
  runasusers = ["user01", "user02"]
  identifier = "users-0"
}
