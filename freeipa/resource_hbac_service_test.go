package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPAHBACService(t *testing.T) {
	testHBACService := map[string]string{
		"name":         "/bin/bash",
		"description":  "The bash terminal",
		"description2": "The other bash terminal",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAHBACServiceResource_basic(testHBACService),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service.hbac_svc", "name", testHBACService["name"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceResource_full(testHBACService),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service.hbac_svc", "name", testHBACService["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_service.hbac_svc", "description", testHBACService["description"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceResource_update(testHBACService),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service.hbac_svc", "name", testHBACService["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_service.hbac_svc", "description", testHBACService["description2"]),
				),
			},
		},
	})
}

func testAccFreeIPAHBACServiceResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_service" "hbac_svc" {
		name       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPAHBACServiceResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_service" "hbac_svc" {
		name        = "%s"
		description  = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"])
}

func testAccFreeIPAHBACServiceResource_update(dataset map[string]string) string {
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
	  
	resource "freeipa_hbac_service" "hbac_svc" {
		name        = "%s"
		description  = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description2"])
}
