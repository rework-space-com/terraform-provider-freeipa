package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPAGroup_posix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testacc-group-1\"",
		"description": "\"Test group 1\"",
		"gid_number":  "10000",
		"addattr":     "[\"owner=uid=test\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}
	testGroup2 := map[string]string{
		"index":       "2",
		"name":        "\"testacc-grouppos-2\"",
		"description": "\"User group test 2\"",
		"gid_number":  "10002",
		"addattr":     "[\"owner=uid=test2\"]",
		"setattr":     "[\"owner=uid=test\"]",
	}
	testGroupDS := map[string]string{
		"index": "1",
		"name":  "freeipa_group.group-2.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_datasource(testGroupDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "description", "User group test 2"),
					resource.TestCheckResourceAttr("data.freeipa_group.group-1", "gid_number", "10002"),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "description", "User group test 2"),
					resource.TestCheckResourceAttr("freeipa_group.group-2", "gid_number", "10002"),
				),
			},
		},
	})
}

func TestAccFreeIPAGroup_noposix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testacc-groupnonpos\"",
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
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testacc-groupnonpos"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "User group test"),
				),
			},
		},
	})
}

func TestAccFreeIPAGroup_external(t *testing.T) {
	testGroup := map[string]string{
		"index":       "1",
		"name":        "\"testacc-groupext\"",
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
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testacc-groupext"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "External user group test"),
				),
			},
		},
	})
}
