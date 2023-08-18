resource "freeipa_hostgroup" "hostgroup" {
  name        = "my-hostgroup"
  description = "my-hostgroup desc"
}

resource "freeipa_automemberadd" "automember" {
  name = freeipa_hostgroup.hostgroup.name
  type = "hostgroup"
}
