resource "freeipa_host_hostgroup_membership" "test-0" {
  name = "test-hostgroup-2"
  host = "test.example.test"
}

resource "freeipa_host_hostgroup_membership" "test-1" {
  name      = "test-hostgroup-2"
  hostgroup = "test-hostgroup"
}
