package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoCmdGrpMembership_simple(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmd", "/usr/bin/testacc-bash"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoCmdGrpMembership_simple_CaseInsensitive(t *testing.T) {
	testSudoCmd1 := map[string]string{
		"index":       "1",
		"name":        "\"/usr/bin/TestACC-Bash\"",
		"description": "\"The bash shell\"",
	}
	testSudoCmdGrp := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-Terminals\"",
		"description": "\"A set of terminals\"",
	}
	testSudoCmdGrpMembership := map[string]string{
		"index":   "1",
		"name":    "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
		"sudocmd": "freeipa_sudo_cmd.sudocmd-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/TestACC-Bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "TestACC-Terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "name", "TestACC-Terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmd", "/usr/bin/TestACC-Bash"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoCmdGrpMembership_mutiple(t *testing.T) {
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
	testSudoCmdGrpDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
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
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoCmdGroup_datasource(testSudoCmdGrpDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.#", "3"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.0", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.1", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.2", "/usr/bin/testacc-zsh"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoCmdGroup_datasource(testSudoCmdGrpDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoCmdGrpMembership_mutiple_CaseInsensitive(t *testing.T) {
	testSudoCmd1 := map[string]string{
		"index":       "1",
		"name":        "\"/usr/bin/TestACC-Bash\"",
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
		"name":        "\"TestACC-Terminals\"",
		"description": "\"A set of terminals\"",
	}
	testSudoCmdGrpMembership := map[string]string{
		"index":       "1",
		"name":        "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
		"sudocmds":    "[freeipa_sudo_cmd.sudocmd-1.name,freeipa_sudo_cmd.sudocmd-2.name,freeipa_sudo_cmd.sudocmd-3.name]",
		"indentifier": "multiplecmds",
	}
	testSudoCmdGrpDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_cmdgroup.sudocmdgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/TestACC-Bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-2", "name", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-2", "description", "The fish shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-3", "name", "/usr/bin/testacc-zsh"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-3", "description", "The zsh shell"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "TestACC-Terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "name", "TestACC-Terminals"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.#", "3"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.0", "/usr/bin/TestACC-Bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.1", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup_membership.sudocmdgroup-membership-1", "sudocmds.2", "/usr/bin/testacc-zsh"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoCmdGroup_datasource(testSudoCmdGrpDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "id", "testacc-terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "TestACC-Terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of terminals"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.#", "3"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.0", "/usr/bin/TestACC-Bash"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.1", "/usr/bin/testacc-fish"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_cmdgroup.sudocmdgroup-1", "member_sudocmd.2", "/usr/bin/testacc-zsh"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd1) + testAccFreeIPASudoCmd_resource(testSudoCmd2) + testAccFreeIPASudoCmd_resource(testSudoCmd3) + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp) + testAccFreeIPASudoCmdGrpMembership_resource(testSudoCmdGrpMembership) + testAccFreeIPASudoCmdGroup_datasource(testSudoCmdGrpDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
