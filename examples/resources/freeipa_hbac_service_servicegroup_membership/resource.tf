resource "freeipa_hbac_service" "hbac_svc_1" {
  name        = "postgresql"
  description = "HBAC Service 1"
}

resource "freeipa_hbac_service" "hbac_svc_2" {
  name        = "mongodb"
  description = "HBAC Service 2"
}

resource "freeipa_hbac_servicegroup" "hbac_svcgroup_1" {
  name        = "database_admins"
  description = "HBAC Service group 1"
}

resource "freeipa_hbac_service_servicegroup_membership" "hbac_svcgroup_member_1" {
  name    = freeipa_hbac_service.hbac_svcgroup_1.name
  service = freeipa_hbac_service.hbac_svc_1.name
}

resource "freeipa_hbac_service_servicegroup_membership" "hbac_svcgroup_member_2" {
  name    = freeipa_hbac_service.hbac_svcgroup_1.id
  service = freeipa_hbac_service.hbac_svc_2.id
}