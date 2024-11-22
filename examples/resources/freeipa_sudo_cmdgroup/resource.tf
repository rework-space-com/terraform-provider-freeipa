resource "freeipa_sudo_cmdgroup" "service_management" {
  name        = "service-management"
  description = "Service management related sudo commands"
}
