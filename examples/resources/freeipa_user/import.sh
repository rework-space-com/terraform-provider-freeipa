# The import id of an active must be exactly equal to `username` of the user to import.

# The associated resource in terraform must include the attributes:
# - `name`
# - `first_name`
# - `last_name`

import {
  to = freeipa_user.testuser
  id = "testuser"
}

resource "freeipa_user" "testuser" {
  name       = "testuser"
  first_name = "Test"
  last_name  = "User"
}

# The import id of an staged must be exactly equal to `username;staged` of the user to import.

# The associated resource in terraform must include the attributes:
# - `name`
# - `first_name`
# - `last_name`
# - `account_staged`

import {
  to = freeipa_user.testuser
  id = "testuser;staged"
}

resource "freeipa_user" "testuser" {
  name           = "testuser"
  first_name     = "Test"
  last_name      = "User"
  account_staged = true
}
