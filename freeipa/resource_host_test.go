package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSHost(t *testing.T) {
	testHost := map[string]string{
		"name":                        "testhost.testacc.ipatest.lan",
		"description":                 "Host test",
		"ip_address":                  "192.168.1.10",
		"locality":                    "Vienna",
		"location":                    "L3",
		"platform":                    "vSphere 7.0.3",
		"operating_system":            "Debian 11",
		"mac_addresses":               "00:00:00:00:00:00",
		"ipasshpubkeys":               "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDDmMkNHn3R+DzSamQDSW60a0iVlAvzbuC3auu8lNoi3u6lvMemsZqPTuvfY4Xlf7uzm+dya3fTRdPKn8sYgPwQ4saUpCSlegN44PjJMhonR1a7FbpHLWj8CRRfzdUSznQhzFcFff0wMBYAklXlyjvdFM8ahl7zHO08HR6469XOVwO1Tb3OGPrXB2lzStK5PKfk5DO/IKl4vHSKhVNVnsZe52rHiZrxOqdGyCijtvwmW2YfIAGc1k4Seqn/Nn7NxKIFBH3hxaUDqgpZneXzuw9GI/F0M8phnHxXNFVZvIWZVcanEeXtH9Z+vVx1ujNcB2QhiPfLMqkNl9db7uykSGKFM4jD0UjGj5kJ8TOC39Safk7XzpQTnrqvIi158zBHVSgugth+QsE1I9/PL2wlzx1qWV2991JKIOc8m52Iwq02tyO8JaSssFTk9szkLTAHedPnZeBbdnlRYcHqX+NPaUh3hqRTZBIR79Ruk6WAliFkED1L0SgwDfGFlevn1Kde9ok=",
		"userclass":                   "testhost",
		"krb_auth_indicators":         "otp",
		"krb_preauth":                 "false",
		"trusted_for_delegation":      "false",
		"trusted_to_auth_as_delegate": "false",
		"userpassword":                "P@ssword",
		"random_password":             "false",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSHostResource_basic(testHost),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host", "name", testHost["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSHostResource_full(testHost),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_host.host", "name", testHost["name"]),
					resource.TestCheckResourceAttr("freeipa_host.host", "description", testHost["description"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSHostResource_basic(dataset map[string]string) string {
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

	resource "freeipa_dns_zone" "testacc_ipatest_lan" {
		zone_name          = "testacc.ipatest.lan"
	}
	
	  
	resource "freeipa_host" "host" {
		name       = "%s"
		ip_address = "%s"
		depends_on = [
			freeipa_dns_zone.testacc_ipatest_lan
		]
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["ip_address"])
}

func testAccFreeIPADNSHostResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_host" "host" {
		name        = "%s"
		description  = "%s"
		ip_address = "%s"
		locality = "%s"
		location = "%s"
		platform = "%s"
		operating_system = "%s"
		mac_addresses = ["%s"]
		ipasshpubkeys = ["%s"]
		userclass = ["%s"]
		krb_auth_indicators = ["%s"]
		krb_preauth = %s
		trusted_for_delegation = %s
		trusted_to_auth_as_delegate = %s
		userpassword = "%s"
		random_password = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["ip_address"], dataset["locality"],
		dataset["location"], dataset["platform"], dataset["operating_system"], dataset["mac_addresses"], dataset["ipasshpubkeys"], dataset["userclass"],
		dataset["krb_auth_indicators"], dataset["krb_preauth"], dataset["trusted_for_delegation"], dataset["trusted_to_auth_as_delegate"], dataset["userpassword"], dataset["random_password"])
}
