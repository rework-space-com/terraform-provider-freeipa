package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAGroup_posix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testgroup1\"",
		"description": "\"Test group 1\"",
		"gid_number":  "10000",
		"addattr":     "[\"owner=uid=test\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}
	testGroup2 := map[string]string{
		"index":       "2",
		"name":        "\"testgrouppos2\"",
		"description": "\"User group test 2\"",
		"gid_number":  "10002",
		"addattr":     "[\"owner=uid=test2\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "Test group 1"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "gid_number", "10000"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_datasource(testGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "description", "Test group 1"),
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "gid_number", "10000"),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "description", "User group test 2"),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "gid_number", "10002"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_datasource(testGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAGroup_noposix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testgroupnonpos\"",
		"description": "\"User group test\"",
		"nonposix":    "true",
		"addattr":     "[\"owner=uid=test\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testgroupnonpos"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "User group test"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAGroup_external(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testgroupext\"",
		"description": "\"External user group test\"",
		"external":    "true",
		"addattr":     "[\"owner=uid=test\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testgroupext"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "External user group test"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
