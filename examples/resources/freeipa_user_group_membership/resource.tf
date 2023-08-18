resource freeipa_user_group_membership "test-0" {
  name = "test-group-2"
  user = "roman"
}

resource freeipa_user_group_membership "test-1" {
  name = "test-group-2"
  group = "test-group"
}
