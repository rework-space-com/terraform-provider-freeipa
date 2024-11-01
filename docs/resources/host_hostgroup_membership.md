---
page_title: "freeipa_host_hostgroup_membership Resource - freeipa"
description: |-

---

# freeipa_host_hostgroup_membership (Resource)



## Example Usage

```terraform
resource freeipa_host_hostgroup_membership "test-0" {
  name = "test-hostgroup-2"
  host = "test.example.test"
}

resource freeipa_host_hostgroup_membership "test-1" {
  name = "test-hostgroup-2"
  hostgroup = "test-hostgroup"
}
```




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Group name

### Optional

- `host` (String) Host to add
- `hostgroup` (String) HostGroup to add

### Read-Only

- `id` (String) The ID of this resource.