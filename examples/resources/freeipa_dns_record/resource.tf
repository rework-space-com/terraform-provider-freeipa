resource freeipa_dns_zone "dns_zone-2" {
  zone_name = "test.roman.com.ua."
  skip_overlap_check = true
}

resource freeipa_dns_record "record-8" {
  zone_name = resource.freeipa_dns_zone.dns_zone-2.id
  name = "test-record"
  records = ["192.168.10.10", "192.168.10.11"]
  type = "A"
}

resource freeipa_dns_record "record-7" {
  zone_name = "record.com.ua."
  name = "test-record"
  records = ["2 1 84DE37B22918F76ED66910B47EB440B0A35F4A56"]
  type = "SSHFP"
}
