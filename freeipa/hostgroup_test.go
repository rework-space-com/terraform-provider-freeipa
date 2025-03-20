package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHostgroup_posix(t *testing.T) {
	testHostgroup := map[string]string{
		"index":       "1",
		"name":        "\"testacc-group-1\"",
		"description": "\"Test hostgroup 1\"",
	}
	testHostgroupModified := map[string]string{
		"index":       "1",
		"name":        "\"testacc-grouppos-1\"",
		"description": "\"Modified description\"",
	}
	testHostgroupDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hostgroup.hostgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "description", "Test hostgroup 1"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroupModified) + testAccFreeIPAHostGroup_datasource(testHostgroupDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "description", "Modified description"),
					resource.TestCheckResourceAttr("data.freeipa_hostgroup.hostgroup-1", "description", "Modified description"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroupModified) + testAccFreeIPAHostGroup_datasource(testHostgroupDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHostgroup_posix_CaseInsensitive(t *testing.T) {
	testHostgroup := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-Group-1\"",
		"description": "\"Test hostgroup 1\"",
	}
	testHostgroupDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hostgroup.hostgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "id", "testacc-group-1"),
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "name", "TestACC-Group-1"),
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "description", "Test hostgroup 1"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroup) + testAccFreeIPAHostGroup_datasource(testHostgroupDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hostgroup.hostgroup-1", "description", "Test hostgroup 1"),
					resource.TestCheckResourceAttr("data.freeipa_hostgroup.hostgroup-1", "name", "TestACC-Group-1"),
					resource.TestCheckResourceAttr("data.freeipa_hostgroup.hostgroup-1", "description", "Test hostgroup 1"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHostGroup_resource(testHostgroup) + testAccFreeIPAHostGroup_datasource(testHostgroupDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
