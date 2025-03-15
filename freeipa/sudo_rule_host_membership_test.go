package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPASudoRuleHostMembership_simple(t *testing.T) {
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
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoHostMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"host":  "freeipa_host.host-0.name",
	}
	testSudoHostGrpMembership := map[string]string{
		"index":     "2",
		"name":      "freeipa_sudo_rule.sudorule-1.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-0.name",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "host", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroup", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleHostMembership_simple_CaseInsensitive(t *testing.T) {
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
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoHostMembership := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
		"host":  "freeipa_host.host-0.name",
	}
	testSudoHostGrpMembership := map[string]string{
		"index":     "2",
		"name":      "freeipa_sudo_rule.sudorule-1.name",
		"hostgroup": "freeipa_hostgroup.hostgroup-0.name",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "host", "TestACC-Host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroup", "TestACC-HostGroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleHostMembership_mutiple(t *testing.T) {
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
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"testacc-sudorule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoHostMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"hosts":      "[freeipa_host.host-0.name]",
		"identifier": "\"hostmembers-1\"",
	}
	testSudoHostGrpMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"hostgroups": "[freeipa_hostgroup.hostgroup-0.name]",
		"identifier": "\"hostgroupmembers-2\"",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "hosts.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "hosts.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroups.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "testacc-sudorule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPASudoRuleHostMembership_mutiple_CaseInsensitive(t *testing.T) {
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
	testSudoRule := map[string]string{
		"index":       "1",
		"name":        "\"TestACC-SudoRule\"",
		"description": "\"A sudo rule for acceptance tests\"",
	}
	testSudoHostMembership := map[string]string{
		"index":      "1",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"hosts":      "[freeipa_host.host-0.name]",
		"identifier": "\"hostmembers-1\"",
	}
	testSudoHostGrpMembership := map[string]string{
		"index":      "2",
		"name":       "freeipa_sudo_rule.sudorule-1.name",
		"hostgroups": "[freeipa_hostgroup.hostgroup-0.name]",
		"identifier": "\"hostgroupmembers-2\"",
	}
	testSudoDS := map[string]string{
		"index": "1",
		"name":  "freeipa_sudo_rule.sudorule-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "hosts.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-1", "hosts.0", "TestACC-Host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroups.#", "1"),
					resource.TestCheckResourceAttr("freeipa_sudo_rule_host_membership.sudo-host-membership-2", "hostgroups.0", "TestACC-HostGroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "name", "TestACC-SudoRule"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "description", "A sudo rule for acceptance tests"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_host.0", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_sudo_rule.sudorule-1", "member_hostgroup.0", "testacc-hostgroup"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testMemberHost) + testAccFreeIPAHostGroup_resource(testHostGroup) + testAccFreeIPASudoRule_resource(testSudoRule) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostMembership) + testAccFreeIPASudoRuleHostMembership_resource(testSudoHostGrpMembership) + testAccFreeIPASudoRule_datasource(testSudoDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
