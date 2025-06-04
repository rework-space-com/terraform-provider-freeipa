# The import id must be exactly the same as the name of the host group.

import {
    to = freeipa_hostgroup.testhostgroup
    id = "testhostgroup"
}

resource "freeipa_hostgroup" "testhostgroup" {
    name = "testhostgroup"
}