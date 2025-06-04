import {
    to = freeipa_sudo_rule.testsudorule
    id = "testsudorule"
}

resource "freeipa_sudo_rule" "testsudorule" {
    name = "testsudorule"
}