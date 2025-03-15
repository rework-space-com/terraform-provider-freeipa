package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoRuleOption_simple(t *testing.T) {
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRuleOption := map[string]string{
		"index":  "1",
		"name":   "freeipa_sudo_rule.sudorule-1.name",
		"option": "\"!authenticate\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_option.sudorule-option-1", "option", "!authenticate"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "option.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "option.0", "!authenticate"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleOption_simple_CaseInsensitive(t *testing.T) {
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRuleOption := map[string]string{
		"index":  "1",
		"name":   "freeipa_sudo_rule.sudorule-1.name",
		"option": "\"!authenticate\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_option.sudorule-option-1", "option", "!authenticate"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "option.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "option.0", "!authenticate"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleOption_resource(testSudoRuleOption) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
