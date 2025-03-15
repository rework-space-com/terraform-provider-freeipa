package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHbacPolicy_simple(t *testing.T) {
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacPolicyModified := map[string]string{
		"index":           "1",
		"name":            "\"testacc-hbac-policy\"",
		"description":     "\"A new hbac policy for acceptance tests\"",
		"enabled":         "false",
		"usercategory":    "\"all\"",
		"hostcategory":    "\"all\"",
		"servicecategory": "\"all\"",
	}
	testHbacPolicyDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A new hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "enabled", "false"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "servicecategory", "all"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A new hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "enabled", "false"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "servicecategory", "all"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicy_simple_CaseSensitive(t *testing.T) {
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"Testacc HBAC Policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacPolicyModified := map[string]string{
		"index":           "1",
		"name":            "\"Testacc HBAC Policy\"",
		"description":     "\"A new hbac policy for acceptance tests\"",
		"enabled":         "false",
		"usercategory":    "\"all\"",
		"hostcategory":    "\"all\"",
		"servicecategory": "\"all\"",
	}
	testHbacPolicyDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "Testacc HBAC Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "Testacc HBAC Policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "Testacc HBAC Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A new hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "enabled", "false"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "servicecategory", "all"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "Testacc HBAC Policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A new hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "enabled", "false"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "usercategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "hostcategory", "all"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "servicecategory", "all"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicyModified) + testAccFreeIPAHbacPolicy_datasource(testHbacPolicyDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
