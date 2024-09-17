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

resource "freeipa_hbac_policy_service_membership" "hbac-svcgrp-1" {
  name         = "test-hbac"
  servicegroup = "Sudo"
}
