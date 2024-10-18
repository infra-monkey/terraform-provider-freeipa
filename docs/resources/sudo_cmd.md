---
page_title: "freeipa_sudo_cmd Resource - freeipa"
description: |-
FreeIPA Sudo command resource
---

# freeipa_sudo_cmd (Resource)



## Example Usage

```terraform
resource "freeipa_sudo_cmd" "bash" {
  name        = "/bin/bash"
  description = "The bash terminal"
}
```




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Absolute path of the sudo command

### Optional

- `description` (String) Sudo command description

### Read-Only

- `id` (String) ID of the resource