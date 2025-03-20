package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoRuleUserMembership_simple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoUserMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"user":  "freeipa_user.user-0.name",
	}
	testSudoUserGrpMembership := map[string]string{
		"index": "2",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"group": "freeipa_group.group-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "user", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "group", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleUserMembership_simple_CaseInsensitive(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"TestACC-Group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoUserMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"user":  "freeipa_user.user-0.name",
	}
	testSudoUserGrpMembership := map[string]string{
		"index": "2",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"group": "freeipa_group.group-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "id", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "user", "TestACC-User-0"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "group", "TestACC-Group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleUserMembership_mutiple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoUserMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"users":      "[freeipa_user.user-0.name]",
		"identifier": "\"usermembers-1\"",
	}
	testSudoUserGrpMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"groups":     "[freeipa_group.group-0.name]",
		"identifier": "\"usermembers-2\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "users.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-1", "users.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_user_membership.sudo-user-membership-2", "groups.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserMembership) + testAccFreeIPASudoRuleUserMembership_resource(testSudoUserGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
