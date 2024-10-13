data "freeipa_host" "host-1" {
  name = freeipa_host.host-1.name
}

output "host-test" {
  value = data.freeipa_host.host-1
}

data "freeipa_sudo_cmdgroup" "sudocmdgroup" {
  name = freeipa_sudo_cmdgroup.terminals.name
}

output "sudocmdgroup" {
  value = data.freeipa_sudo_cmdgroup.sudocmdgroup
}

data "freeipa_sudo_rule" "sysadmins" {
  name = freeipa_sudo_rule.sysadmins.name
}

output "sudorule" {
  value = data.freeipa_sudo_rule.sysadmins
}