resource freeipa_sudo_rule_denycmd_membership "denied_cmd" {
  name = "sudo-rule-restricted"
  sudocmd = "/usr/bin/systemctl"
}

resource freeipa_sudo_rule_denycmd_membership "denied_cmdgrp" {
  name = "sudo-rule-restricted"
  sudocmd_group = "service-management"
}
