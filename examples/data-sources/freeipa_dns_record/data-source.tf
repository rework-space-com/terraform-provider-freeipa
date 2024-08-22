data "freeipa_dns_record" "dns-record-0" {
  record_name = "test-record-A"
  zone_name   = "test.example.lan."
}

data "freeipa_dns_record" "dns-zone-1" {
  record_name = "10"
  zone_name   = "23.168.192.in-addr.arpa."
}
