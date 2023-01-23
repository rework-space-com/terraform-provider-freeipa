package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)


func TestAccFreeIPADNSDNSRecord(t *testing.T) {
	var testDnsRecordA map[string]string
	testDnsRecordA = map[string]string{
		"name": "reca",
		"records": "[\"192.168.1.10\"]",
		"type": "A",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordPTR map[string]string
	testDnsRecordPTR = map[string]string{
		"name": "11",
		"records": "[\"default.testacc.ipatest.lan.\"]",
		"type": "PTR",
		"zone_name": "1.168.192.in-addr.arpa.",
	}
	var testDnsRecordAAAA map[string]string
	testDnsRecordAAAA = map[string]string{
		"name": "recaaaa",
		"records": "[\"2001:db8:3333:4444:5555:6666:7777:8888\"]",
		"type": "AAAA",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordCNAME map[string]string
	testDnsRecordCNAME = map[string]string{
		"name": "reccname",
		"records": "[\"reca.testacc.ipatest.lan.\"]",
		"type": "CNAME",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordTXT map[string]string
	testDnsRecordTXT = map[string]string{
		"name": "rectxt",
		"records": "[\"EXAMPLE_TXT_RECORD\"]",
		"type": "TXT",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordSSHFP map[string]string
	testDnsRecordSSHFP = map[string]string{
		"name": "recsshfp",
		"records": "[\"1 1 84DE37B22918F76ED66910B47EB440B0A35F4A56\"]",
		"type": "SSHFP",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordMX map[string]string
	testDnsRecordMX = map[string]string{
		"name": "recmx",
		"records": "[\"0 mail.example.com.\"]",
		"type": "MX",
		"zone_name": "testacc.ipatest.lan.",
	}
	var testDnsRecordSRV map[string]string
	testDnsRecordSRV = map[string]string{
		"name": "recsrv",
		"records": "[\"10 5 443 reccname.testacc.ipatest.lan.\"]",
		"type": "SRV",
		"zone_name": "testacc.ipatest.lan.",
	}


	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_A", "name", testDnsRecordA["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_A", "type", testDnsRecordA["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_A", "zone_name", testDnsRecordA["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordPTR),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_PTR", "name", testDnsRecordPTR["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_PTR", "type", testDnsRecordPTR["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_PTR", "zone_name", testDnsRecordPTR["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordAAAA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_AAAA", "name", testDnsRecordAAAA["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_AAAA", "type", testDnsRecordAAAA["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_AAAA", "zone_name", testDnsRecordAAAA["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordCNAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_CNAME", "name", testDnsRecordCNAME["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_CNAME", "type", testDnsRecordCNAME["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_CNAME", "zone_name", testDnsRecordCNAME["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordTXT),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_TXT", "name", testDnsRecordTXT["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_TXT", "type", testDnsRecordTXT["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_TXT", "zone_name", testDnsRecordTXT["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordSSHFP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SSHFP", "name", testDnsRecordSSHFP["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SSHFP", "type", testDnsRecordSSHFP["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SSHFP", "zone_name", testDnsRecordSSHFP["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordMX),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_MX", "name", testDnsRecordMX["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_MX", "type", testDnsRecordMX["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_MX", "zone_name", testDnsRecordMX["zone_name"]),
				),
			},
			{
				Config: testAccFreeIPADNSDNSRecordResource_basic(testDnsRecordSRV),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SRV", "name", testDnsRecordSRV["name"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SRV", "type", testDnsRecordSRV["type"]),
					resource.TestCheckResourceAttr("freeipa_dns_record.record_SRV", "zone_name", testDnsRecordSRV["zone_name"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSDNSRecordResource_basic(dataset map[string]string) string {
	provider_host := os.Getenv("FREEIPA_HOST")
	provider_user := os.Getenv("FREEIPA_USERNAME")
	provider_pass := os.Getenv("FREEIPA_PASSWORD")
	return fmt.Sprintf(`
	provider "freeipa" {
		host     = "%s"
		username = "%s"
		password = "%s"
		insecure = true
	  }
	  
	resource "freeipa_dns_zone" "zone" {
		zone_name = "testacc.ipatest.lan."
		allow_prt_sync = false
		dynamic_updates = false
		skip_overlap_check = true
	}
	resource "freeipa_dns_zone" "zonereverse" {
		zone_name       = "192.168.1.0/24"
		is_reverse_zone = true
		allow_prt_sync = false
		dynamic_updates = false
		skip_overlap_check = true
	}

	resource "freeipa_dns_record" "record_default" {
		name = "default"
		records = ["192.168.1.11"]
		type = "A"
		zone_name = "testacc.ipatest.lan."
		depends_on = [
			freeipa_dns_zone.zone,
			freeipa_dns_zone.zonereverse
		]
	}

	resource "freeipa_dns_record" "record_%s" {
		name = "%s"
		records = %s
		type = "%s"
		zone_name = "%s"
		depends_on = [
			freeipa_dns_zone.zone,
			freeipa_dns_zone.zonereverse
		]
	}
	`, provider_host, provider_user, provider_pass, dataset["type"], dataset["name"], dataset["records"], dataset["type"], dataset["zone_name"])
}
