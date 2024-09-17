resource "freeipa_sudo_cmd" "bash" {
  name        = "/bin/bash"
  description = "The bash shell"
}

resource "freeipa_sudo_cmd" "fish" {
  name        = "/bin/fish"
  description = "The fish shell"
}

resource "freeipa_sudo_cmdgroup" "terminals" {
  name        = "terminals"
  description = "The terminals allowed to be sudoed"
}

resource "freeipa_sudo_cmdgroup_membership" "terminal_bash" {
  name    = freeipa_sudocmdgroup.terminals.id
  sudocmd = freeipa_sudocmd.bash.id
}

resource "freeipa_sudo_cmdgroup_membership" "terminal_fish" {
  name    = freeipa_sudocmdgroup.terminals.id
  sudocmd = freeipa_sudocmd.fish.id
}
