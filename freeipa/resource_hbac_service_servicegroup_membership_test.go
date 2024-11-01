package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPAHBACServiceGroupMembership(t *testing.T) {
	testHBACServiceGroupMembership := map[string]string{
		"name":     "database_admins",
		"hbacsvc":  "mongodb",
		"hbacsvc2": "postgresql",
	}
	testHBACServiceGroupMembershipWithSlash := map[string]string{
		"name":     "category/database_admins",
		"hbacsvc":  "mongodb",
		"hbacsvc2": "postgresql",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAHBACServiceGroupMembershipResource_basic(testHBACServiceGroupMembership),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member", "name", testHBACServiceGroupMembership["name"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceGroupMembershipResource_full(testHBACServiceGroupMembership),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member", "name", testHBACServiceGroupMembership["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member2", "name", testHBACServiceGroupMembership["name"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceGroupMembershipResource_basic(testHBACServiceGroupMembershipWithSlash),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member", "name", testHBACServiceGroupMembershipWithSlash["name"]),
				),
			},
			{
				Config: testAccFreeIPAHBACServiceGroupMembershipResource_full(testHBACServiceGroupMembershipWithSlash),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member", "name", testHBACServiceGroupMembershipWithSlash["name"]),
					resource.TestCheckResourceAttr("freeipa_hbac_service_servicegroup_membership.hbac_svc_member2", "name", testHBACServiceGroupMembershipWithSlash["name"]),
				),
			},
		},
	})
}

func testAccFreeIPAHBACServiceGroupMembershipResource_basic(dataset map[string]string) string {
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

	resource "freeipa_hbac_service" "hbacsvc" {
		name = "%s"
	}

	resource "freeipa_hbac_servicegroup" "hbacsvcgroup" {
		name = "%s"
	}

	resource "freeipa_hbac_service_servicegroup_membership" "hbac_svc_member" {
		name    = freeipa_hbac_servicegroup.hbacsvcgroup.name
		service = freeipa_hbac_service.hbacsvc.name
	}
	`, provider_host, provider_user, provider_pass, dataset["hbacsvc"], dataset["name"])
}

func testAccFreeIPAHBACServiceGroupMembershipResource_full(dataset map[string]string) string {
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

	resource "freeipa_hbac_service" "hbacsvc" {
		name = "%s"
	}

	resource "freeipa_hbac_service" "hbacsvc2" {
		name = "%s"
	}

	resource "freeipa_hbac_servicegroup" "hbacsvcgroup" {
		name = "%s"
	}

	resource "freeipa_hbac_service_servicegroup_membership" "hbac_svc_member" {
		name    = freeipa_hbac_servicegroup.hbacsvcgroup.name
		service = freeipa_hbac_service.hbacsvc.name
	}

	resource "freeipa_hbac_service_servicegroup_membership" "hbac_svc_member2" {
		name    = freeipa_hbac_servicegroup.hbacsvcgroup.name
		service = freeipa_hbac_service.hbacsvc2.name
	}
	`, provider_host, provider_user, provider_pass, dataset["hbacsvc"], dataset["hbacsvc2"], dataset["name"])
}
