# The import id must be exactly the same as the sudo command.

import {
  to = freeipa_sudo_cmd.testsudocmd
  id = "/bin/bash"
}

resource "freeipa_sudo_cmd" "testsudocmd" {
  name = "/bin/bash"
}