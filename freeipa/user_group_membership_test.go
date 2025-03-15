package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAUserGroupMembership_posix(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group\"",
		"description": "\"User group test\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User\"",
	}
	testMemberGroup := map[string]string{
		"index":       "1",
		"name":        "\"testacc-groupmember\"",
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
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "testacc-group"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "User group test - member of testgroup"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "testacc-groupmember"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "testacc-group"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "user", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "testacc-group"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "group", "testacc-groupmember"),
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

func TestAccFreeIPAUserGroupMembership_simple_posix_CaseInsensitive(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"TestACC-Group\"",
		"description": "\"User group test\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User\"",
	}
	testMemberGroup := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-GroupMember\"",
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
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "TestACC-Group"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "description", "User group test - member of testgroup"),
					resource.TestCheckResourceAttr("freeipa_group.group-1", "name", "TestACC-GroupMember"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "TestACC-Group"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "user", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "TestACC-Group"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "group", "TestACC-GroupMember"),
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

func TestAccFreeIPAUserGroupMembership_multiple_posix(t *testing.T) {
	testGroup1 := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testGroup2 := map[string]string{
		"index":       "1",
		"name":        "\"testacc-group-1\"",
		"description": "\"User group test 1\"",
	}
	testGroup3 := map[string]string{
		"index":       "2",
		"name":        "\"testacc-group-2\"",
		"description": "\"User group test 2\"",
	}
	testMemberUser1 := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testMemberUser2 := map[string]string{
		"index":     "1",
		"login":     "\"testacc-user-1\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User1\"",
	}
	testMembershipGroups1 := map[string]string{
		"index":       "0",
		"name":        "freeipa_group.group-0.name",
		"description": "\"User group test - member of testgroup\"",
		"groups":      "[freeipa_group.group-1.name]",
		"identifier":  "\"groups\"",
	}
	testMembershipGroups2 := map[string]string{
		"index":       "0",
		"name":        "freeipa_group.group-0.name",
		"description": "\"User group test - member of testgroup\"",
		"groups":      "[freeipa_group.group-1.name,freeipa_group.group-2.name]",
		"identifier":  "\"groups\"",
	}
	testMembershipUsers1 := map[string]string{
		"index":      "1",
		"name":       "freeipa_group.group-0.name",
		"users":      "[freeipa_user.user-0.name]",
		"identifier": "\"users\"",
	}
	testMembershipUsers2 := map[string]string{
		"index":      "1",
		"name":       "freeipa_group.group-0.name",
		"users":      "[freeipa_user.user-0.name,freeipa_user.user-1.name]",
		"identifier": "\"users\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_resource(testGroup3) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUser_resource(testMemberUser2) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups1) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-0", "description", "User group test 0"),
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.0", "testacc-group-1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.0", "testacc-user-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_resource(testGroup3) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUser_resource(testMemberUser2) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups1) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_resource(testGroup3) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUser_resource(testMemberUser2) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups2) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-0", "description", "User group test 0"),
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.0", "testacc-group-1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.1", "testacc-group-2"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "testacc-group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.1", "testacc-user-1"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAGroup_resource(testGroup3) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUser_resource(testMemberUser2) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups2) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAUserGroupMembership_multiple_posix_CaseInsensitive(t *testing.T) {
	testGroup1 := map[string]string{
		"index":       "0",
		"name":        "\"TestACC-Group-0\"",
		"description": "\"User group test 0\"",
	}
	testGroup2 := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-Group-1\"",
		"description": "\"User group test 1\"",
	}
	testMemberUser1 := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testMembershipGroups1 := map[string]string{
		"index":       "0",
		"name":        "freeipa_group.group-0.name",
		"description": "\"User group test - member of testgroup\"",
		"groups":      "[freeipa_group.group-1.name]",
		"identifier":  "\"groups\"",
	}
	testMembershipUsers1 := map[string]string{
		"index":      "1",
		"name":       "freeipa_group.group-0.name",
		"users":      "[freeipa_user.user-0.name]",
		"identifier": "\"users\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups1) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group-0", "description", "User group test 0"),
					resource.TestCheckResourceAttr("freeipa_group.group-0", "name", "TestACC-Group-0"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "name", "TestACC-Group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-0", "groups.0", "TestACC-Group-1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "name", "TestACC-Group-0"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.membership-1", "users.0", "TestACC-User-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAGroup_resource(testGroup1) + testAccFreeIPAGroup_resource(testGroup2) + testAccFreeIPAUser_resource(testMemberUser1) + testAccFreeIPAUserGroupMembership_resource(testMembershipGroups1) + testAccFreeIPAUserGroupMembership_resource(testMembershipUsers1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
