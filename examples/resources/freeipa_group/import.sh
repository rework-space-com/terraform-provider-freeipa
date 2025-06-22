# The import id must be exactly the same as the name of the user group.

# for posix groups, only the name needs to be defined in the resource statement.
import {
  to = freeipa_group.group-posix
  id = "testposix"
}

resource "freeipa_group" "group-posix" {
  name = "testposix"
}

# for external groups, the external attribute must also be defined in the resource statement.
import {
  to = freeipa_group.group-external
  id = "testexternal"
}

resource "freeipa_group" "group-external" {
  name     = "testexternal"
  external = true
}


# for non posix groups, the nonposix attribute must also be defined in the resource statement.
import {
  to = freeipa_group.group-nonposix
  id = "testnonposix"
}

resource "freeipa_group" "group-nonposix" {
  name     = "testnonposix"
  nonposix = true
}

# note that nonposix and external are two mutually exclusive attributes.
# setting the wrong attributes for the group will result in the replacement of the resource (destroy and recreate)
