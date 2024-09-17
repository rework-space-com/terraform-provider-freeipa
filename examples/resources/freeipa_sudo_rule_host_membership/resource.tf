resource "freeipa_sudo_rule_host_membership" "host-0" {
  name = "sudo-rule-test"
  host = "test.example.test"
}

resource "freeipa_sudo_rule_host_membership" "hostgroup-3" {
  name      = "sudo-rule-test"
  hostgroup = "test-hostgroup"
}
