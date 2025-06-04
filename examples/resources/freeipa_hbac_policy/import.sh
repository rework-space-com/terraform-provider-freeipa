# The import id must be exactly the same as the name of the HBAC policy.

import {
    to = freeipa_hbac_policy.testhbac
    id = "testhbac"
}

resource "freeipa_hbac_policy" "testhbac" {
    name = "testhbac"
}