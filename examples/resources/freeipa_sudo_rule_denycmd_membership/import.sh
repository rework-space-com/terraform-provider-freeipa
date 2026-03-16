# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msrdc" for multi-sudocmd/sudocmd_group membership (non-deprecated).
# Note: slash characters in the rule name or sudocmd paths must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_denycmd_membership.denied_systemctl
  id = "sudo-rule-restricted/msrdc/denied_systemctl"
}

resource "freeipa_sudo_rule_denycmd_membership" "denied_systemctl" {
  name       = "sudo-rule-restricted"
  sudocmds   = ["/usr/bin/systemctl"]
  identifier = "denied_systemctl"
}
