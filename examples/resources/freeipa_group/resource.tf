resource "freeipa_group" "group-posix" {
  name        = "test-group"
  description = "Test group"
  gid_number  = "12345789"
}

resource "freeipa_group" "group-nonposix" {
  name        = "test-group"
  description = "Test group"
  nonposix    = true
}

resource "freeipa_group" "group-external" {
  name        = "test-group"
  description = "Test group"
  external    = true
}