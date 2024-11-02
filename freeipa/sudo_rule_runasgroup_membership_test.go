package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPASudoRuleRunAsGroupMembership_simple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRunAsGroupMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"runasgroup": "freeipa_group.group-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsGroupMembership_resource(testSudoRunAsGroupMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.sudorule-runasgroup-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.sudorule-runasgroup-membership-1", "runasgroup", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsGroupMembership_resource(testSudoRunAsGroupMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasgroup.0", "testacc-group-0"),
				),
			},
		},
	})
}

func TestAccFreeIPASudoRuleRunAsGroupMembership_multiple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoRunAsGroupMembership := map[string]string{
		"index":       "1",
		"name":        "freeipa_sudo_rule.sudorule-1.name",
		"runasgroups": "[freeipa_group.group-0.name]",
		"identifier":  "\"runasgroup0\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsGroupMembership_resource(testSudoRunAsGroupMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.sudorule-runasgroup-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.sudorule-runasgroup-membership-1", "runasgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_runasgroup_membership.sudorule-runasgroup-membership-1", "runasgroups.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleRunAsGroupMembership_resource(testSudoRunAsGroupMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "runasgroup.0", "testacc-group-0"),
				),
			},
		},
	})
}
