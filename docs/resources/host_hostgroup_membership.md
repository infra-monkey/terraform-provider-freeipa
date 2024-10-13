---
page_title: "freeipa_host_hostgroup_membership Resource - freeipa"
description: |-
FreeIPA User Group Membership resource
---

# freeipa_host_hostgroup_membership (Resource)



## Example Usage

```terraform
resource "freeipa_host_hostgroup_membership" "test-0" {
  name = "test-hostgroup-2"
  host = "test.example.test"
}

resource "freeipa_host_hostgroup_membership" "test-1" {
  name      = "test-hostgroup-2"
  hostgroup = "test-hostgroup"
}

resource "freeipa_host_hostgroup_membership" "test-2" {
  name       = "test-hostgroup-2"
  hosts      = ["host1","host2"]
  hostgroups = ["test-hostgroup","test-hostgroup2"]
  identifier = "my_unique_identifier"
}
```




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Hostgroup name

### Optional

- `host` (String, Deprecated) **deprecated** Host to add. Will be replaced by hosts.
- `hostgroup` (String, Deprecated) **deprecated** Hostgroup to add. Will be replaced by hostgroups.
- `hostgroups` (List of String) Hostgroups to add as hostgroup members
- `hosts` (List of String) Hosts to add as hostgroup members
- `identifier` (String) Unique identifier to differentiate multiple hostgroup membership resources on the same hostgroup. Manadatory for using hosts/hostgroups configurations.

### Read-Only

- `id` (String) ID of the resource