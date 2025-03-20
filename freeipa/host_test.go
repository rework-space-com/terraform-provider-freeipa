package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAHost_full(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"testacc.ipatest.lan\"",
	}
	testHost := map[string]string{
		"index":      "0",
		"name":       "\"testacc-host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address": "\"192.168.10.65\"",
	}
	testHostModified := map[string]string{
		"index":                       "0",
		"name":                        "\"testacc-host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address":                  "\"192.168.10.65\"",
		"description":                 "\"FreeIPA client in testacc.ipatest.lan domain\"",
		"locality":                    "\"Some City\"",
		"location":                    "\"lab\"",
		"operating_system":            "\"Fedora 40\"",
		"mac_addresses":               "[\"00:00:00:AA:AA:AA\", \"00:00:00:BB:BB:BB\"]",
		"trusted_for_delegation":      "true",
		"trusted_to_auth_as_delegate": "true",
	}
	testHostModified2 := map[string]string{
		"index":                       "0",
		"name":                        "\"testacc-host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address":                  "\"192.168.10.65\"",
		"description":                 "\"FreeIPA client in testacc.ipatest.lan domain\"",
		"locality":                    "\"Some New City\"",
		"location":                    "\"dc1\"",
		"operating_system":            "\"RHEL 9\"",
		"mac_addresses":               "[\"00:00:00:CC:CC:CC\"]",
		"trusted_for_delegation":      "false",
		"trusted_to_auth_as_delegate": "false",
	}
	testHostDS := map[string]string{
		"index": "0",
		"name":  "freeipa_host.host-0.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHost),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host-0", "name", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "ip_address", "192.168.10.65"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHostModified) + testAccFreeIPAHost_datasource(testHostDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host-0", "name", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "ip_address", "192.168.10.65"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "description", "FreeIPA client in testacc.ipatest.lan domain"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "locality", "Some City"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "location", "lab"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "operating_system", "Fedora 40"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "mac_addresses.#", "2"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "mac_addresses.0", "00:00:00:AA:AA:AA"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "mac_addresses.1", "00:00:00:BB:BB:BB"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "trusted_for_delegation", "true"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "trusted_to_auth_as_delegate", "true"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "description", "FreeIPA client in testacc.ipatest.lan domain"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "locality", "Some City"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "location", "lab"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "operating_system", "Fedora 40"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "mac_addresses.#", "2"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "mac_addresses.0", "00:00:00:AA:AA:AA"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "mac_addresses.1", "00:00:00:BB:BB:BB"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "trusted_for_delegation", "true"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "trusted_to_auth_as_delegate", "true"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHostModified) + testAccFreeIPAHost_datasource(testHostDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHostModified2) + testAccFreeIPAHost_datasource(testHostDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host-0", "name", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "ip_address", "192.168.10.65"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "description", "FreeIPA client in testacc.ipatest.lan domain"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "locality", "Some New City"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "location", "dc1"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "operating_system", "RHEL 9"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "mac_addresses.#", "1"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "mac_addresses.0", "00:00:00:CC:CC:CC"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "trusted_for_delegation", "false"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "trusted_to_auth_as_delegate", "false"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "description", "FreeIPA client in testacc.ipatest.lan domain"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "locality", "Some New City"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "location", "dc1"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "operating_system", "RHEL 9"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "mac_addresses.#", "1"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "mac_addresses.0", "00:00:00:CC:CC:CC"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "trusted_for_delegation", "false"),
					resource.TestCheckResourceAttr("data.freeipa_host.host-0", "trusted_to_auth_as_delegate", "false"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHost) + testAccFreeIPAHost_datasource(testHostDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host-0", "name", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "ip_address", "192.168.10.65"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHost) + testAccFreeIPAHost_datasource(testHostDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAHost_full_CaseInsensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"testacc.ipatest.lan\"",
	}
	testHost := map[string]string{
		"index":      "0",
		"name":       "\"TestACC-Host-1.${freeipa_dns_zone.dns-zone-0.zone_name}\"",
		"ip_address": "\"192.168.10.65\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHost),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host-0", "id", "testacc-host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "name", "TestACC-Host-1.testacc.ipatest.lan"),
					resource.TestCheckResourceAttr("freeipa_host.host-0", "ip_address", "192.168.10.65"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPAHost_resource(testHost),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
