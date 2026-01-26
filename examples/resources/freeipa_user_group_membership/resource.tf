resource "freeipa_user_group_membership" "test-3" {
  name       = "test-group-3"
  users      = ["user1", "user2"]
  groups     = ["group1", "group2"]
  identifier = "my_unique_identifier"
}

resource "freeipa_user_group_membership" "test-3" {
  name             = "test-extgroup-3"
  external_members = ["domain users@adtest.lan"]
  identifier       = "my_unique_identifier"
}