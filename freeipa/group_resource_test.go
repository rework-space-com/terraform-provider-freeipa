package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPAGroup_posix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "testgroup1",
		"description": "Test group 1",
		"gid_number":  "10000",
		"addattr":     "owner=uid=test",
		"setattr":     "owner=uid=test",
	}
	testGroup2 := map[string]string{
		"index":       "2",
		"name":        "testgrouppos2",
		"description": "User group test 2",
		"gid_number":  "10002",
		"addattr":     "owner=uid=test2",
		"setattr":     "owner=uid=test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resourcefull(testGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", string(testGroup["description"])),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "gid_number", string(testGroup["gid_number"])),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resourcefull(testGroup) + testAccFreeIPAGroup_resourcefull(testGroup2) + testAccFreeIPAGroup_datasource(testGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "description", string(testGroup["description"])),
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "gid_number", string(testGroup["gid_number"])),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "description", string(testGroup2["description"])),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "gid_number", string(testGroup2["gid_number"])),
				),
			},
		},
	})
}

func TestAccFreeIPAGroup_noposix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "testgroupnonpos",
		"description": "User group test",
		"nonposix":    "true",
		"addattr":     "owner=uid=test",
		"setattr":     "owner=uid=test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroupNonposix_resourcefull(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", testGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", string(testGroup["description"])),
				),
			},
		},
	})
}

func TestAccFreeIPAGroup_external(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "testgroupnonpos",
		"description": "User group test",
		"external":    "true",
		"addattr":     "owner=uid=test",
		"setattr":     "owner=uid=test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroupExternal_resourcefull(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", testGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", string(testGroup["description"])),
				),
			},
		},
	})
}
