resource "freeipa_sysaccount" "sysaccount-1" {
  name        = "test"
  description = "Test system account"
  random      = true
}

resource "freeipa_sysaccount" "privileged-sysaccount" {
  name        = "test"
  description = "Test system account"
  random      = true
  privileged  = true
}
