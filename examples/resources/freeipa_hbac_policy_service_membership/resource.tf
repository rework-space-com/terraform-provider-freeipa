resource "freeipa_hbac_policy" "hbac-0" {
  name            = "test-hbac"
  description     = "Test HBAC policy"
  enabled         = true
  hostcategory    = "all"
  servicecategory = "all"
}

resource "freeipa_hbac_policy_service_membership" "hbac-svc-1" {
  name    = "test-hbac"
  service = "sshd"
}

resource "freeipa_hbac_policy_service_membership" "hbac-svc-2" {
  name       = "test-hbac"
  services   = ["sshd"]
  identifier = "hbac-svc-2"
}

resource "freeipa_hbac_policy_service_membership" "hbac-svcgrp-1" {
  name         = "test-hbac"
  servicegroup = "Sudo"
}

resource "freeipa_hbac_policy_service_membership" "hbac-svcgrp-2" {
  name          = "test-hbac"
  servicegroups = ["Sudo", "ftp"]
  identifier    = "hbac-svcgrp-2"
}
