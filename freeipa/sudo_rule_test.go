package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoRule_simple(t *testing.T) {
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRuleModified := map[string]string{
		"index":              "1",
		"name":               "\"testacc-sudorule\"",
		"description":        "\"A new sudo rule for acceptance tests\"",
		"enabled":            "false",
		"usercategory":       "\"all\"",
		"hostcategory":       "\"all\"",
		"commandcategory":    "\"all\"",
		"runasusercategory":  "\"all\"",
		"runasgroupcategory": "\"all\"",
		"order":              "5",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRuleModified),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A new sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "enabled", "false"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "commandcategory", "all"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "runasusercategory", "all"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "runasgroupcategory", "all"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "order", "5"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRuleModified) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A new sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "enabled", "false"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "commandcategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasusercategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasgroupcategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "order", "5"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRuleModified) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRule_simple_CaseInsensitive(t *testing.T) {
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "id", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "id", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
