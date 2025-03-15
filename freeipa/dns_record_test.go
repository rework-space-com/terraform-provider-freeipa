package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPADNSRecord_A(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"ipa.example.lan\"",
	}
	testRecord := map[string]string{
		"index":     "0",
		"zone_name": "resource.freeipa_dns_zone.dns-zone-0.id",
		"type":      "\"A\"",
		"name":      "\"test-record\"",
		"records":   "[\"192.168.10.10\", \"192.168.10.11\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.dns-record-0", "name", "test-record"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPADNSRecord_A_CaseSensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"ipa.example.lan\"",
	}
	testRecord := map[string]string{
		"index":     "0",
		"zone_name": "resource.freeipa_dns_zone.dns-zone-0.id",
		"type":      "\"A\"",
		"name":      "\"Test-Record\"",
		"records":   "[\"192.168.10.10\", \"192.168.10.11\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.dns-record-0", "name", "Test-Record"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPADNSRecord_CNAME(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"ipa.example.lan\"",
	}
	testRecord := map[string]string{
		"index":     "0",
		"zone_name": "resource.freeipa_dns_zone.dns-zone-0.id",
		"name":      "\"test-cname\"",
		"type":      "\"CNAME\"",
		"records":   "[\"test-record.ipa.example.lan.\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.dns-record-0", "name", "test-cname"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPADNSRecord_CNAME_CaseSensitive(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "\"ipa.example.lan\"",
	}
	testRecord := map[string]string{
		"index":     "0",
		"zone_name": "resource.freeipa_dns_zone.dns-zone-0.id",
		"name":      "\"Test-CNAME\"",
		"type":      "\"CNAME\"",
		"records":   "[\"test-record.ipa.example.lan.\"]",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.dns-record-0", "name", "Test-CNAME"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSRecord_resource(testRecord),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
