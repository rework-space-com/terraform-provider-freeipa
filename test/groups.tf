resource freeipa_group "group-0" {
  name = "test-group2"
  description = "Test group"
  # gid_number = "499600746"
  # nonposix = true
  # external = true
  # addattr = ["owner=foo=bar"]
}

resource "freeipa_group" "group-posix" {
  name        = "test-group-pos"
  description = "Test group posix"
  gid_number  = "12345789"
}

resource "freeipa_group" "group-nonposix" {
  name        = "test-group-nonpos"
  description = "Test group non posix"
  nonposix    = true
}

resource "freeipa_group" "group-external" {
  name        = "test-group-ext"
  description = "Test group external"
  external    = true
}

# resource freeipa_user_group_membership "member-0" {
#   name = freeipa_group.group-posix.name
#   users = [freeipa_user.user-1.name]
# }

resource freeipa_user_group_membership "member-1" {
  name = freeipa_group.group-posix.name
  groups = [freeipa_group.group-nonposix.name,freeipa_group.group-external.name]
  identifier = "posix-members"
}

