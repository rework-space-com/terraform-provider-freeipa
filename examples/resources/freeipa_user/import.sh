# The import id must be exactly equal to `username` of the user to import.

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
