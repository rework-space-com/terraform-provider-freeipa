resource "freeipa_sudo_rule_runasgroup_membership" "group-0" {
  name       = "sudo-rule-test"
  runasgroup = "group01"
}

resource "freeipa_sudo_rule_runasgroup_membership" "groups-0" {
  name        = "sudo-rule-test"
  runasgroups = ["group01", "group02"]
  identifier  = "groups-0"
}