package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAUserGroupMembership_posix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testgroup\"",
		"description": "\"User group test\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testuser\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User\"",
	}
	testMemberGroup := map[string]string{
		"index":       "1",
		"name":        "\"testgroupmember\"",
		"description": "\"User group test - member of testgroup\"",
	}
	testMembershipUser := map[string]string{
		"index": "0",
		"name":  "freeipa_group.group-0.name",
		"user":  "freeipa_user.user-0.name",
	}
	testMembershipGroup := map[string]string{
		"index": "1",
		"name":  "freeipa_group.group-0.name",
		"group": "freeipa_group.group-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAGroup_resource(testMemberGroup) + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAUserGroupMembership_resource(testMembershipUser) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-0", "description", "User group test"),
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "testgroup"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "User group test - member of testgroup"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testgroupmember"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "testuser"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "testgroup"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "user", "testuser"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "testgroup"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "group", "testgroupmember"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAGroup_resource(testMemberGroup) + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAUserGroupMembership_resource(testMembershipUser) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
