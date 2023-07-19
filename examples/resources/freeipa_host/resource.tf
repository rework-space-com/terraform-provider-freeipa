resource freeipa_host "host-1" {
  name = "host-1.example.test"
  ip_address = "192.168.1.65"
  description = "FreeIPA client in example.test domain"
  mac_addresses = ["00:00:00:AA:AA:AA", "00:00:00:BB:BB:BB"]
}
