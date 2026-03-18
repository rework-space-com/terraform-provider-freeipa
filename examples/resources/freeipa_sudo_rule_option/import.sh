# The import id uses the format: <sudo_rule_name>/sro/<option_value>
# Note: slash characters in the rule name must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_option.option-0
  id = "sudo-rule-test/sro/!authenticate"
}

resource "freeipa_sudo_rule_option" "option-0" {
  name   = "sudo-rule-test"
  option = "!authenticate"
}
