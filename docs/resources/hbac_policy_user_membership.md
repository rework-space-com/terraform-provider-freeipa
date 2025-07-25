---
page_title: "freeipa_hbac_policy_user_membership Resource - freeipa"
description: |-
  FreeIPA HBAC policy host membership resource
---

# freeipa_hbac_policy_user_membership (Resource)



## Example Usage

```terraform
resource "freeipa_hbac_policy" "hbac-0" {
  name            = "test-hbac"
  description     = "Test HBAC policy"
  enabled         = true
  hostcategory    = "all"
  servicecategory = "all"
}

resource "freeipa_hbac_policy_user_membership" "hbac-user-1" {
  name = "test-hbac"
  user = "user-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-users-1" {
  name       = "test-hbac"
  users      = ["user-2", "user-3"]
  identifier = "hbac-users-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-group-1" {
  name  = "test-hbac"
  group = "usergroup-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-groups-1" {
  name       = "test-hbac"
  groups     = ["usergroup-2", "usergroup-3"]
  identifier = "hbac-groups-1"
}
```




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) HBAC policy name

### Optional

- `group` (String, Deprecated) **deprecated** User group to add to the HBAC policy
- `groups` (List of String) List of user groups to add to the HBAC policy
- `identifier` (String) Unique identifier to differentiate multiple HBAC policy user membership resources on the same HBAC policy. Manadatory for using users/groups configurations.
- `user` (String, Deprecated) **deprecated** User FDQN the policy is applied to
- `users` (List of String) List of user FQDNs to add to the HBAC policy

### Read-Only

- `id` (String) ID of the resource
