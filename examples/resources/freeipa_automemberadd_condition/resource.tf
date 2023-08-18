resource "freeipa_hostgroup" "hostgroup" {
  name        = "my-hostgroup"
  description = "my-hostgroup desc"
}

resource "freeipa_automemberadd" "automember" {
  name = freeipa_hostgroup.hostgroup.name
  type = "hostgroup"
}

resource "freeipa_automemberadd_condition" "automembercondition" {
  name           = freeipa_automemberadd.automember.name
  type           = "hostgroup"
  key            = "fqdn"
  inclusiveregex = ["\\.my\\.first\\.net$", "\\.my\\.second\\.net$"]
}
