# The import id uses the format: <cmdgroup_name>/<type>/<identifier>
# Use type "msc" for multi-sudocmd membership (non-deprecated).
# Note: slash characters in the cmdgroup name or sudocmd paths must be percent-encoded (%2F).

import {
  to = freeipa_sudo_cmdgroup_membership.terminal_shells
  id = "terminals/msc/terminal_shells"
}

resource "freeipa_sudo_cmdgroup_membership" "terminal_shells" {
  name       = "terminals"
  sudocmds   = ["/bin/bash", "/bin/fish"]
  identifier = "terminal_shells"
}
