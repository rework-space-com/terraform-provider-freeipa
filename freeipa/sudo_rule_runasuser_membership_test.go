package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoRuleRunAsUserMembership_simple(t *testing.T) {
	testUser := map[string]string{
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
	testSudoRunAsUserMembership := map[string]string{
		"index":     "1",
		"name":      "freeipa_sudo_rule.sudorule-1.name",
		"runasuser": "freeipa_user.user-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasuser", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleRunAsUserMembership_simple_CaseInsensitive(t *testing.T) {
	testUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRunAsUserMembership := map[string]string{
		"index":     "1",
		"name":      "freeipa_sudo_rule.sudorule-1.name",
		"runasuser": "freeipa_user.user-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasuser", "TestACC-User-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleRunAsUserMembership_multiple(t *testing.T) {
	testUser := map[string]string{
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
	testSudoRunAsUserMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"runasusers": "[freeipa_user.user-0.name]",
		"identifier": "\"runasuser0\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasusers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasusers.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleRunAsUserMembership_multiple_CaseInsensitive(t *testing.T) {
	testUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRunAsUserMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"runasusers": "[freeipa_user.user-0.name]",
		"identifier": "\"runasuser0\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasusers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasuser_membership.sudorule-runasuser-membership-1", "runasusers.0", "TestACC-User-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasuser.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsUserMembership_resource(testSudoRunAsUserMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
