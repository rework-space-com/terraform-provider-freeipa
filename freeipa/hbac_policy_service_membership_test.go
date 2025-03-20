package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHbacPolicyServiceMembership_simple(t *testing.T) {
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacServiceMembership := map[string]string{
		"index":   "1",
		"name":    "freeipa_hbac_policy.hbacpolicy-1.name",
		"service": "\"sshd\"",
	}
	testHbacServiceGroupMembership := map[string]string{
		"index":        "2",
		"name":         "freeipa_hbac_policy.hbacpolicy-1.name",
		"servicegroup": "\"Sudo\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-1", "service", "sshd"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-2", "servicegroup", "Sudo"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_service.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_service.0", "sshd"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_servicegroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_servicegroup.0", "Sudo"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyServiceMembership_mutiple(t *testing.T) {
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacServiceMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"services":   "[\"sshd\"]",
		"identifier": "\"service-1\"",
	}
	testHbacServiceGroupMembership := map[string]string{
		"index":         "2",
		"name":          "freeipa_hbac_policy.hbacpolicy-1.name",
		"servicegroups": "[\"Sudo\"]",
		"identifier":    "\"servicegroup-2\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-1", "services.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-1", "services.0", "sshd"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-2", "servicegroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_service_membership.hbac-service-membership-2", "servicegroups.0", "Sudo"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_service.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_service.0", "sshd"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_servicegroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_servicegroup.0", "Sudo"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceMembership) + testAccFreeIPAHbacPolicyServiceMembership_resource(testHbacServiceGroupMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
