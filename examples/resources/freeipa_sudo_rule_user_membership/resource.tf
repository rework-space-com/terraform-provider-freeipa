resource freeipa_sudo_rule_user_membership "user-0" {
  name = "sudo-rule-test"
  user = "user01"
}

resource freeipa_sudo_rule_user_membership "group-3" {
  name = "sudo-rule-test"
  group = "test-group-0"
}
