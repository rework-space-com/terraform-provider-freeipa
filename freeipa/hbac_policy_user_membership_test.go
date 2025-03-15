package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHbacPolicyUserMembership_simple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacUserMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"user":  "freeipa_user.user-0.name",
	}
	testHbacUserGrpMembership := map[string]string{
		"index": "2",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"group": "freeipa_group.group-0.name",
	}
	testHbacDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "user", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "group", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyUserMembership_simple_CaseInsensitive(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"TestACC-Group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-HBAC-Policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacUserMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"user":  "freeipa_user.user-0.name",
	}
	testHbacUserGrpMembership := map[string]string{
		"index": "2",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"group": "freeipa_group.group-0.name",
	}
	testHbacDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "user", "TestACC-User-0"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "group", "TestACC-Group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyUserMembership_mutiple(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"testacc-group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-user-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacUsersMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"users":      "[freeipa_user.user-0.name]",
		"identifier": "\"users-1\"",
	}
	testHbacUserGrpsMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"groups":     "[freeipa_group.group-0.name]",
		"identifier": "\"groups-2\"",
	}
	testHbacDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "users.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "users.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "groups.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyUserMembership_mutiple_CaseInsensitive(t *testing.T) {
	testGroup := map[string]string{
		"index":       "0",
		"name":        "\"TestACC-Group-0\"",
		"description": "\"User group test 0\"",
	}
	testMemberUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User-0\"",
		"firstname": "\"Test\"",
		"lastname":  "\"User0\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-HBAC-Policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacUsersMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"users":      "[freeipa_user.user-0.name]",
		"identifier": "\"users-1\"",
	}
	testHbacUserGrpsMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"groups":     "[freeipa_group.group-0.name]",
		"identifier": "\"groups-2\"",
	}
	testHbacDS := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "users.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-1", "users.0", "TestACC-User-0"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "groups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_user_membership.hbac-user-membership-2", "groups.0", "TestACC-Group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_user.0", "testacc-user-0"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_group.0", "testacc-group-0"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testMemberUser) + testAccFreeIPAGroup_resource(testGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUsersMembership) + testAccFreeIPAHbacPolicyUserMembership_resource(testHbacUserGrpsMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
