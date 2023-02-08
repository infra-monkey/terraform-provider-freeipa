---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "freeipa_host_hostgroup_membership Resource - freeipa"
subcategory: ""
description: |-
  
---

# freeipa_host_hostgroup_membership (Resource)

## Example Usage
```hcl
resource freeipa_host_hostgroup_membership "test-0" {
  name = "test-hostgroup-2"
  host = "test.example.test"
}

resource freeipa_host_hostgroup_membership "test-1" {
  name = "test-hostgroup-2"
  hostgroup = "test-hostgroup"
}
```

### Required

- `name` (String) Group name

### Optional

- `host` (String) (Required if `hostgroup` is empty) Host to add
- `hostgroup` (String) (Required if `host` is empty) HostGroup to add

### Read-Only

- `id` (String) The ID of this resource.

