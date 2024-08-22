package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPASudoRuleAllowCmdMembership_simple(t *testing.T) {
	testSudoCmd1 := map[string]string{
		"index":       "1",
		"name":        "\"/usr/bin/testacc-bash\"",
		"description": "\"The bash shell\"",
	}
	testSudoCmdGrp := map[string]string{
		"index":       "1",
		"name":        "\"testacc-terminals\"",
		"description": "\"A set of terminals\"",
	}
	testSudoCmdGrpMembership := map[string]string{
		"index":   "1",
		"name":    "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
		"sudocmd": "freeipa_sudo_cmd.sudocmd-1.name",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoAllowCmdMembership := map[string]string{
		"index":   "1",
		"name":    "freeipa_sudo_rule.sudorule-1.name",
		"sudocmd": "freeipa_sudo_cmd.sudocmd-1.name",
	}
	testSudoAllowCmdGrpMembership := map[string]string{
		"index":         "2",
		"name":          "freeipa_sudo_rule.sudorule-1.name",
		"sudocmd_group": "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdMembership) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmd", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-1", "sudocmd", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-2", "sudocmd_group", "testacc-terminals"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdMembership) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmd.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmd.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmdgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmdgroup.0", "testacc-terminals"),
				),
			},
		},
	})
}

func TestAccFreeIPASudoRuleAllowCmdMembership_mutiple(t *testing.T) {
	testSudoCmd1 := map[string]string{
		"index":       "1",
		"name":        "\"/usr/bin/testacc-bash\"",
		"description": "\"The bash shell\"",
	}
	testSudoCmd2 := map[string]string{
		"index":       "2",
		"name":        "\"/usr/bin/testacc-fish\"",
		"description": "\"The fish shell\"",
	}
	testSudoCmd3 := map[string]string{
		"index":       "3",
		"name":        "\"/usr/bin/testacc-zsh\"",
		"description": "\"The zsh shell\"",
	}
	testSudoCmdGrp := map[string]string{
		"index":       "1",
		"name":        "\"testacc-terminals\"",
		"description": "\"A set of terminals\"",
	}
	testSudoCmdGrpMembership := map[string]string{
		"index":       "1",
		"name":        "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
		"sudocmds":    "[freeipa_sudo_cmd.sudocmd-1.name,freeipa_sudo_cmd.sudocmd-2.name,freeipa_sudo_cmd.sudocmd-3.name]",
		"indentifier": "multiplecmds",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoAllowCmdMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"sudocmds":   "[freeipa_sudo_cmd.sudocmd-1.name]",
		"identifier": "\"testacc-allowcmds\"",
	}
	testSudoAllowCmdGrpMembership := map[string]string{
		"index":          "2",
		"name":           "freeipa_sudo_rule.sudorule-1.name",
		"sudocmd_groups": "[freeipa_sudo_cmdgroup.sudocmdgroup-1.name]",
		"identifier":     "\"testacc-allowcmdgroups\"",
	}
	testSudoCmdGrpDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdMembership) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-2", "name", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-2", "description", "The fish shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-3", "name", "/usr/bin/testacc-zsh"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-3", "description", "The zsh shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.#", "3"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.1", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.2", "/usr/bin/testacc-zsh"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-1", "sudocmds.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-1", "sudocmds.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-2", "sudocmd_groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_allowcmd_membership.sudo-allow-membership-2", "sudocmd_groups.0", "testacc-terminals"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdMembership) + testAccFreeIPASudoAllowCmdMembership_resource(testSudoAllowCmdGrpMembership) + testAccFreeIPASudoCmdGroup_datasource(testSudoCmdGrpDS) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.#", "3"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.1", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.2", "/usr/bin/testacc-zsh"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmd.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmd.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmdgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_allow_sudo_cmdgroup.0", "testacc-terminals"),
				),
			},
		},
	})
}
