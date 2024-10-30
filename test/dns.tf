resource "freeipa_dns_zone" "dns_zone-1" {
  zone_name          = "testzone.ipatest.lan"
  skip_overlap_check = true
  disable_zone       = false
  allow_ptr_sync     = true
  admin_email_address = "admin.testzone.ipatest.lan"
}

resource "freeipa_dns_zone" "dns_zone-2" {
  zone_name          = "10.10.10.0/24"
  is_reverse_zone    = true
  disable_zone       = false
}

resource "freeipa_dns_record" "record-1" {
  zone_name = freeipa_dns_zone.dns_zone-1.id
  name      = "test-record"
  records   = ["10.10.10.10", "10.10.10.11"]
  type      = "A"
}

resource "freeipa_dns_record" "record-2" {
  zone_name = freeipa_dns_zone.dns_zone-1.id
  name      = "test-record2"
  records   = ["test-record"]
  type      = "CNAME"
}