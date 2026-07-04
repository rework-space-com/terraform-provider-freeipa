// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASysaccount_full(t *testing.T) {
	testSysaccount := map[string]string{
		"index":       "1",
		"name":        "\"testaccount\"",
		"description": "\"Test system account\"",
		"random":      "true",
	}
	testSysaccountMod := map[string]string{
		"index":       "1",
		"name":        "\"testaccount\"",
		"description": "\"Test system account\"",
		"random":      "true",
		"disabled":    "true",
		"privileged":  "true",
	}
	testSysaccountDS := map[string]string{
		"index": "1",
		"name":  "\"testaccount\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "random", "true"),
					resource.TestCheckNoResourceAttr("freeipa_sysaccount.sysaccount-1", "disabled"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount) + testAccFreeIPASysAccount_datasource(testSysaccountDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "disabled", "false"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccountMod),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "disabled", "true"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "privileged", "true"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccountMod),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccountMod) + testAccFreeIPASysAccount_datasource(testSysaccountDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "disabled", "true"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("freeipa_sysaccount.sysaccount-1", "random", "true"),
					resource.TestCheckNoResourceAttr("freeipa_sysaccount.sysaccount-1", "disabled"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASysAccount_resource(testSysaccount) + testAccFreeIPASysAccount_datasource(testSysaccountDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "name", "testaccount"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "description", "Test system account"),
					resource.TestCheckResourceAttr("data.freeipa_sysaccount.sysaccount-1", "disabled", "false"),
				),
			},
		},
	})
}
