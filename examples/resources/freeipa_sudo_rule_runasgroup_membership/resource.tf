resource "freeipa_sudo_runasrule_user_membership" "user-0" {
  name       = "sudo-rule-test"
  runasgroup = "group01"
}

