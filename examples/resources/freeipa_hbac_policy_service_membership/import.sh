# The import id uses the format: <hbac_policy_name>/<type>/<identifier>
# Use type "ms" for multi-service/servicegroup membership (non-deprecated).

import {
  to = freeipa_hbac_policy_service_membership.hbac-svc-2
  id = "test-hbac/ms/hbac-svc-2"
}

resource "freeipa_hbac_policy_service_membership" "hbac-svc-2" {
  name       = "test-hbac"
  services   = ["sshd"]
  identifier = "hbac-svc-2"
}
