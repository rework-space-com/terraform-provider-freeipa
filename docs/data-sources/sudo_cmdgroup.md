---
page_title: "freeipa_sudo_cmdgroup Data Source - freeipa"
description: |-
FreeIPA User sudo command group data source
---

# freeipa_sudo_cmdgroup (Data Source)



## Example Usage

```terraform
data "freeipa_sudo_cmdgroup" "terminals" {
  name = "terminals"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the sudo command group

### Read-Only

- `description` (String) Description of the sudo command group
- `id` (String) ID of the resource in the terraform state
- `member_sudocmd` (List of String) List of sudo commands that are member of the sudo command group