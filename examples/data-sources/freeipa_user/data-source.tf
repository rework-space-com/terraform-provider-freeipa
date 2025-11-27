# Lookup an active user

data "freeipa_user" "user-0" {
  name = "test-user"
}

# Lookup an staged user

data "freeipa_user" "user-0" {
  name  = "test-user"
  state = "staged"
}
