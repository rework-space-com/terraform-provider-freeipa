---
page_title: "freeipa_hbac_policy_user_membership Resource - freeipa"
description: |-

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

resource "freeipa_hbac_policy_user_membership" "hbac-group-1" {
  name      = "test-hbac"
  group = "usergroup-1"
}
```




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) HBAC policy name

### Optional

- `group` (String) Group the policy is applied to
- `user` (String) User FDQN the policy is applied to

### Read-Only

- `id` (String) The ID of this resource.
