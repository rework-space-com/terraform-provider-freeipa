package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSUserGroupMembership(t *testing.T) {
	testDatasetUser := map[string]string{
		"login":     "testuser",
		"firstname": "Test",
		"lastname":  "User",
	}
	testDatasetGroup := map[string]string{
		"name": "testgroup",
	}
	testDatasetGroup2 := map[string]string{
		"name": "testgroup-2",
	}
	testDatasetGroup3 := map[string]string{
		"name": "testgroup-3",
	}

	/*testDatasetExtGroup := map[string]string{
		"name": "testgroup-ext",
	}
	testDatasetExtGroup2 := map[string]string{
		"name": "testgroup2-ext",
	}*/

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSUserGroupMembershipResource_user(testDatasetUser, testDatasetGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership", "name", testDatasetGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership", "user", testDatasetUser["login"]),
				),
			},
			{
				Config: testAccFreeIPADNSUserGroupMembershipResource_group(testDatasetGroup3, testDatasetGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership2", "name", testDatasetGroup3["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership2", "group", testDatasetGroup2["name"]),
				),
			},
			// External users need a valid trust setup to work. Not possible for acceptance test
			/*{
				Config: testAccFreeIPADNSUserGroupMembershipResource_externaluser(testDatasetExtGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership3", "name", testDatasetExtGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership3", "external_member", "user@domain"),
				),
			},
			{
				Config: testAccFreeIPADNSUserGroupMembershipResource_externalgroup(testDatasetExtGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership4", "name", testDatasetExtGroup2["name"]),
					resource.TestCheckResourceAttr("freeipa_user_group_membership.groupmembership4", "external_member", "usergroup@domain"),
				),
			},*/
		},
	})
}

func testAccFreeIPADNSUserGroupMembershipResource_user(dataset_user map[string]string, dataset_group map[string]string) string {
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
	  
	resource "freeipa_user" "user" {
		name       = "%s"
		first_name = "%s"
		last_name  = "%s"
	}

	resource "freeipa_group" "group" {
		name       = "%s"
	}
	resource freeipa_user_group_membership "groupmembership" {
	   name = resource.freeipa_group.group.id
	   user = resource.freeipa_user.user.id
	}
	`, provider_host, provider_user, provider_pass, dataset_user["login"], dataset_user["firstname"], dataset_user["lastname"], dataset_group["name"])
}

func testAccFreeIPADNSUserGroupMembershipResource_group(dataset_group map[string]string, dataset_group2 map[string]string) string {
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
	
	resource "freeipa_group" "group2" {
		name       = "%s"
	}
	resource "freeipa_group" "subgroup" {
		name       = "%s"
	}

	resource freeipa_user_group_membership "groupmembership2" {
	   name = resource.freeipa_group.group2.id
	   group = resource.freeipa_group.subgroup.id
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"], dataset_group2["name"])
}

/*func testAccFreeIPADNSUserGroupMembershipResource_externaluser(dataset_group map[string]string) string {
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

	resource "freeipa_group" "group3" {
		name       = "%s"
		external = true
	}
	resource freeipa_user_group_membership "groupmembership3" {
	   name     = resource.freeipa_group.group3.id
	   external_member = "user@domain"
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"])
}

func testAccFreeIPADNSUserGroupMembershipResource_externalgroup(dataset_group map[string]string) string {
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

	resource "freeipa_group" "group4" {
		name       = "%s"
		external   = true
	}

	resource freeipa_user_group_membership "groupmembership4" {
	   name     = resource.freeipa_group.group4.id
	   external_member = "usergroup@domain"
	}
	}
	`, provider_host, provider_user, provider_pass, dataset_group["name"])
}*/
