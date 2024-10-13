resource "freeipa_sudo_cmd" "sudocmd-1" {
  name          = "/usr/bin/bash"
}
resource "freeipa_sudo_cmd" "sudocmd-2" {
  name          = "/usr/bin/fish"
}
resource "freeipa_sudo_cmd" "sudocmd-3" {
  name          = "/usr/bin/zsh"
}
resource "freeipa_sudo_cmd" "sudocmd-4" {
  name          = "/usr/bin/systemctl"
}

resource "freeipa_sudo_cmdgroup" "services" {
  name        = "services"
  description = "Service management related sudo commands"
}

resource "freeipa_sudo_cmdgroup" "terminals" {
  name        = "terminals"
  description = "Terminal related sudo commands"
}

resource "freeipa_sudo_cmdgroup_membership" "terminals" {
  name    = freeipa_sudo_cmdgroup.terminals.name
  sudocmds= [freeipa_sudo_cmd.sudocmd-1.name,freeipa_sudo_cmd.sudocmd-2.name,freeipa_sudo_cmd.sudocmd-3.name]
  identifier = "terminals"
}

resource "freeipa_sudo_cmdgroup_membership" "services" {
  name    = freeipa_sudo_cmdgroup.services.name
  sudocmds= [freeipa_sudo_cmd.sudocmd-4.name]
  identifier = "services"
}

resource "freeipa_sudo_rule" "sysadmins" {
  name        = "sysadmins"
  description = "Sysadmins have all permissions"
  enabled = true
#   usercategory = "all"
#   hostcategory = "all"
#   commandcategory = "all"
#   runasusercategory = "all"
#   runasgroupcategory ="all"
  # order = 1
}

resource "freeipa_sudo_rule" "limited-admins" {
  name        = "limited-admins"
  description = "Limited admins have some permissions"
  enabled = true
}

resource "freeipa_sudo_rule_denycmd_membership" "sysadmins-cmds" {
  name    = freeipa_sudo_rule.sysadmins.name
  sudocmds = [freeipa_sudo_cmd.sudocmd-4.name]
  sudocmd_groups = [freeipa_sudo_cmdgroup.terminals.name]
  identifier = "sysadmins-cmds"
}
resource "freeipa_sudo_rule_allowcmd_membership" "sysadmins-cmds" {
  name    = freeipa_sudo_rule.sysadmins.name
  sudocmds = [freeipa_sudo_cmd.sudocmd-4.name]
  sudocmd_groups = [freeipa_sudo_cmdgroup.terminals.name]
  identifier = "sysadmins-cmds"
}

resource "freeipa_sudo_rule_host_membership" "hosts-0" {
  name = freeipa_sudo_rule.sysadmins.name
  hosts = [freeipa_host.host-1.name]
  identifier = "hosts-0"
}

resource "freeipa_sudo_rule_host_membership" "hostgroups-3" {
  name      = freeipa_sudo_rule.sysadmins.name
  hostgroups = [freeipa_hostgroup.hostgroup-1.name]
  identifier = "hostgroups-3"
}

resource "freeipa_sudo_rule_runasuser_membership" "users-0" {
  name       = freeipa_sudo_rule.sysadmins.name
  runasusers = [freeipa_user.user-0.name,freeipa_user.user-1.name]
  # runasusers = [freeipa_user.user-0.name]
  identifier = "users-0"
}

resource "freeipa_sudo_rule_runasgroup_membership" "groups-0" {
  name       = freeipa_sudo_rule.sysadmins.name
  runasgroups = [freeipa_group.group-posix.name,freeipa_group.group-0.name]
  # runasgroups = [freeipa_group.group-posix.name]
  identifier = "groups-0"
}

resource "freeipa_sudo_rule_user_membership" "users-0" {
  name      = freeipa_sudo_rule.sysadmins.name
  users = [freeipa_user.user-1.name]
  groups = [freeipa_group.group-posix.name]
  identifier = "users-0"
}