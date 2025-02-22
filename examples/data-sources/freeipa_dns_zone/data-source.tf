data "freeipa_dns_zone" "dns-zone-0" {
  zone_name = "test.example.lan."
}

data "freeipa_dns_zone" "dns-zone-1" {
  zone_name = "23.168.192.in-addr.arpa."
}
