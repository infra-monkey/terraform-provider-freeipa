resource "freeipa_sudo_rule_host_membership" "hosts-0" {
  name       = "sudo-rule-test"
  hosts      = ["test.example.test"]
  identifier = "hosts-0"
}

resource "freeipa_sudo_rule_host_membership" "hostgroups-3" {
  name       = "sudo-rule-test"
  hostgroups = ["test-hostgroup"]
  identifier = "hostgroups-3"
}