resource "freeipa_host" "host-1" {
  name          = "host-1.${freeipa_dns_zone.dns_zone-1.zone_name}"
  ip_address    = "192.168.1.65"
  description   = "FreeIPA client in testzone.ipatest.lan domain"
  operating_system = "Opensuse Leap 15.6"
  mac_addresses = ["00:00:00:AA:AA:AA", "00:00:00:BB:BB:CC"]
  trusted_for_delegation = true
}

resource "freeipa_hostgroup" "hostgroup-1" {
  name        = "hostgroup-1"
  description = "FreeIPA hostgroup"
}