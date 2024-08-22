resource "freeipa_dns_zone" "dns_zone-2" {
  zone_name          = "test.roman.com.ua"
  skip_overlap_check = true
  disable_zone       = false
}
