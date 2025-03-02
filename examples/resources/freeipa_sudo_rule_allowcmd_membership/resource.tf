resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmd" {
  name    = "sudo-rule-editors"
  sudocmd = "/bin/bash"
}

resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmd" {
  name       = "sudo-rule-editors"
  sudocmds   = ["/bin/bash"]
  identifier = "allowed_bash"
}

resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmdgrp" {
  name          = "sudo-rule-editors"
  sudocmd_group = "allowed-terminals"
}

resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmdgrp" {
  name           = "sudo-rule-editors"
  sudocmd_groups = ["allowed-terminals"]
  identifier     = "allowed_terminals"
}

