package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHbacPolicyHostMembership_simple(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"testacc.ipatest.lan\"",
	}
	testMemberHost := map[string]string{
		"index":      "0",
		"name":       "\"testacc-host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address": "\"192.168.10.65\"",
	}
	testHostGroup := map[string]string{
		"index": "0",
		"name":  "\"testacc-hostgroup\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacHostMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"host":  "freeipa_host.host-0.name",
	}
	testHbacHostGrpMembership := map[string]string{
		"index":     "2",
		"name":      "freeipa_hbac_policy.hbacpolicy-1.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "host", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "hostgroup", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyHostMembership_simple_CaseInsensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"testacc.ipatest.lan\"",
	}
	testMemberHost := map[string]string{
		"index":      "0",
		"name":       "\"TestACC-Host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address": "\"192.168.10.65\"",
	}
	testHostGroup := map[string]string{
		"index": "0",
		"name":  "\"TestACC-HostGroup\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-HBAC-Policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacHostMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_hbac_policy.hbacpolicy-1.name",
		"host":  "freeipa_host.host-0.name",
	}
	testHbacHostGrpMembership := map[string]string{
		"index":     "2",
		"name":      "freeipa_hbac_policy.hbacpolicy-1.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-0.name",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "host", "TestACC-Host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "hostgroup", "TestACC-HostGroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "TestACC-HBAC-Policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHbacPolicyHostMembership_mutiple(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"testacc.ipatest.lan\"",
	}
	testMemberHost := map[string]string{
		"index":      "0",
		"name":       "\"testacc-host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address": "\"192.168.10.65\"",
	}
	testHostGroup := map[string]string{
		"index": "0",
		"name":  "\"testacc-hostgroup\"",
	}
	testHbacPolicy := map[string]string{
		"index":       "1",
		"name":        "\"testacc-hbac-policy\"",
		"description": "\"A hbac policy for acceptance tests\"",
	}
	testHbacHostMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"hosts":      "[freeipa_host.host-0.name]",
		"identifier": "\"hbac-host-1\"",
	}
	testHbacHostGrpMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_hbac_policy.hbacpolicy-1.name",
		"hostgroups": "[freeipa_hostgroup.hostgroup-0.name]",
		"identifier": "\"hbac-hostgrp-2\"",
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
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "hosts.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-1", "hosts.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "hostgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_hbac_policy_host_membership.hbac-host-membership-2", "hostgroups.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "name", "testacc-hbac-policy"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "description", "A hbac policy for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_hbac_policy.hbacpolicy-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHbacPolicy_resource(testHbacPolicy) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostMembership) + testAccFreeIPAHbacPolicyHostMembership_resource(testHbacHostGrpMembership) + testAccFreeIPAHbacPolicy_datasource(testHbacDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
