package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPASudoCmd_simple(t *testing.T) {
	testSudoCmd := map[string]string{
		"index":       "1",
		"name":        "\"/usr/bin/testacc-bash\"",
		"description": "\"The bash shell\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPASudoCmd_resource(testSudoCmd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "name", "/usr/bin/testacc-bash"),
					resource.TestCheckResourceAttr("freeipa_sudo_cmd.sudocmd-1", "description", "The bash shell"),
				),
			},
		},
	})
}
