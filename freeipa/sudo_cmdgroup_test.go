package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		},
	})
}
