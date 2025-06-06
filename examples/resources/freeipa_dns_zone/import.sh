# forward zone
# The import id attribute must be the undotted fqdn of the zone to import

import {
    to = freeipa_dns_zone.testzone
    id = "testimport.ipatest.lan"
}

resource "freeipa_dns_zone" "testzone" {
  zone_name          = "testimport.ipatest.lan"
}

# reverse zone
# The import id attribute must be the undotted fqdn of the zone to import

import {
    to = freeipa_dns_zone.reversetestzone
    id = "2.27.172.in-addr.arpa"
}

resource "freeipa_dns_zone" "reversetestzone" {
  zone_name          = "2.27.172.in-addr.arpa"
}

# note that the is_reverse_zone must not be defined, this is only useful for creation