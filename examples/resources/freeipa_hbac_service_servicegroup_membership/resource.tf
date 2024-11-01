resource "freeipa_hbac_service" "test" {
  name        = "sudo"
  description = "hbacsvc-1"
}

resource "freeipa_hbac_service" "test2" {
  name        = "su -"
  description = "hbacsvc-2"
}

resource "freeipa_hbac_servicegroup" "hbac_svcgroup-1" {
  name        = "admin"
  description = "hbac service group 1"
}


resource "freeipa_hbac_service_servicegroup_membership" "hbac_svcgroup-member-2" {
  name    = freeipa_hbac_service.hbac_svcgroup-1.name
  service = freeipa_hbac_service.test2.name
}

resource "freeipa_hbac_service_servicegroup_membership" "hbac_svcgroup-member-2" {
  name    = freeipa_hbac_service.hbac_svcgroup-1.name
  service = freeipa_hbac_service.test2.name
}