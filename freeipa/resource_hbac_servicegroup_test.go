package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPAHBACServiceGroup(t *testing.T) {
	testHostgroup := map[string]string{
		"name":        "testhbacsvcgroup",
		"description": "hbacscv group test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAHBACServiceGroupResource_basic(testHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_servicegroup.hbac_svcgroup", "name", testHostgroup["name"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceGroupResource_full(testHostgroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_servicegroup.hbac_svcgroup", "name", testHostgroup["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_servicegroup.hbac_svcgroup", "description", testHostgroup["description"]),
				),
			},
		},
	})
}

func testAccFreeIPAHBACServiceGroupResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_servicegroup" "hbac_svcgroup" {
		name       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPAHBACServiceGroupResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_servicegroup" "hbac_svcgroup" {
		name        = "%s"
		description  = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"])
}
