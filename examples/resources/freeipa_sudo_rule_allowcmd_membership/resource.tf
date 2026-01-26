resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmds" {
  name       = "sudo-rule-editors"
  sudocmds   = ["/bin/bash"]
  identifier = "allowed_bash"
}

resource "freeipa_sudo_rule_allowcmd_membership" "allowed_cmdgrps" {
  name           = "sudo-rule-editors"
  sudocmd_groups = ["allowed-terminals"]
  identifier     = "allowed_terminals"
}

