package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoCmdGrp_simple(t *testing.T) {
	testSudoCmdGrp := map[string]string{
		"index":       "1",
		"name":        "\"testacc-command-group-1\"",
		"description": "\"A set of commands\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "testacc-command-group-1"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of commands"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoCmdGrp_CaseInsensitive(t *testing.T) {
	testSudoCmdGrp := map[string]string{
		"index":       "1",
		"name":        "\"Testacc Command Group 1\"",
		"description": "\"A set of commands\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "name", "Testacc Command Group 1"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmdgroup.sudocmdgroup-1", "description", "A set of commands"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmdGrp_resource(testSudoCmdGrp),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
