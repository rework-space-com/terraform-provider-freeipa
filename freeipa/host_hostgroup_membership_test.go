package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHostGroupMembership_simple(t *testing.T) {
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
	testMemberHostGroup := map[string]string{
		"index": "1",
		"name":  "\"testacc-groupmember\"",
	}
	testMembershipHost := map[string]string{
		"index": "0",
		"name":  "freeipa_hostgroup.hostgroup-0.name",
		"host":  "freeipa_host.host-0.name",
	}
	testMembershipHostGroup := map[string]string{
		"index":     "1",
		"name":      "freeipa_hostgroup.hostgroup-0.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "host", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroup", "testacc-groupmember"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHostGroupMembership_simple_CaseInsensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"Testacc.ipatest.lan\"",
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
	testMemberHostGroup := map[string]string{
		"index": "1",
		"name":  "\"TestACC-GroupMember\"",
	}
	testMembershipHost := map[string]string{
		"index": "0",
		"name":  "freeipa_hostgroup.hostgroup-0.name",
		"host":  "freeipa_host.host-0.name",
	}
	testMembershipHostGroup := map[string]string{
		"index":     "1",
		"name":      "freeipa_hostgroup.hostgroup-0.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "host", "TestACC-Host-1.Testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroup", "TestACC-GroupMember"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHostGroupMembership_multiple(t *testing.T) {
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
	testMemberHostGroup := map[string]string{
		"index": "1",
		"name":  "\"testacc-groupmember\"",
	}
	testMembershipHost := map[string]string{
		"index":      "0",
		"name":       "freeipa_hostgroup.hostgroup-0.name",
		"hosts":      "[freeipa_host.host-0.name]",
		"identifier": "1",
	}
	testMembershipHostGroup := map[string]string{
		"index":      "1",
		"name":       "freeipa_hostgroup.hostgroup-0.name",
		"hostgroups": "[freeipa_hostgroup.hostgroup-1.name]",
		"identifier": "2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "hosts.#", "1"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "hosts.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroups.0", "testacc-groupmember"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHostGroupMembership_multiple_CaseInsensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"TestACC.ipatest.lan\"",
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
	testMemberHostGroup := map[string]string{
		"index": "1",
		"name":  "\"TestACC-GroupMember\"",
	}
	testMembershipHost := map[string]string{
		"index":      "0",
		"name":       "freeipa_hostgroup.hostgroup-0.name",
		"hosts":      "[freeipa_host.host-0.name]",
		"identifier": "1",
	}
	testMembershipHostGroup := map[string]string{
		"index":      "1",
		"name":       "freeipa_hostgroup.hostgroup-0.name",
		"hostgroups": "[freeipa_hostgroup.hostgroup-1.name]",
		"identifier": "2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "hosts.#", "1"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-0", "hosts.0", "TestACC-Host-1.TestACC.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_host_hostgroup_membership.membership-1", "hostgroups.0", "TestACC-GroupMember"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPAHostGroup_resource(testMemberHostGroup) + testAccFreeIPAHostGroupMembership_resource(testMembershipHost) + testAccFreeIPAHostGroupMembership_resource(testMembershipHostGroup),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
