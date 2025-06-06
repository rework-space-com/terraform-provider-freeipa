# The import id must be exactly the same as the name of the sudo command group.

import {
    to = freeipa_sudo_cmdgroup.testsudocmdgrp
    id = "testsudocmdgrp"
}


resource "freeipa_sudo_cmdgroup" "testsudocmdgrp" {
    name = "testsudocmdgrp"
}