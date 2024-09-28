package freeipa

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFreeIPADNSGroup_posix(t *testing.T) {
	testGroup := map[string]string{
		"name":        "testgrouppos",
		"description": "User group test",
		"gid_number":  "10001",
		"addattr":     "owner=uid=test",
	}
	testGroup2 := map[string]string{
		"name":        "testgrouppos2",
		"description": "User group test 2",
		"gid_number":  "10002",
		"addattr":     "owner=uid=test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSGroupResource_basic(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_full(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "gid_number", testGroup["gid_number"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "addattr.0", testGroup["addattr"]),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_full(testGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup2["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup2["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "gid_number", testGroup2["gid_number"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "addattr.0", testGroup2["addattr"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSGroupResource_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name       = "%s"
	}
	`, provider_host, provider_user, provider_pass, dataset["name"])
}

func testAccFreeIPADNSGroupResource_full(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name        = "%s"
		description  = "%s"
		gid_number = %s
                addattr = ["%s"]
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["gid_number"], dataset["addattr"])
}

func TestAccFreeIPADNSGroup_noposix(t *testing.T) {
	testGroup := map[string]string{
		"name":        "testgroupnonpos",
		"description": "User group test",
		"nonposix":    "true",
		"addattr":     "owner=uid=test",
	}
	testGroup2 := map[string]string{
		"name":        "testgroupnonpos2",
		"description": "User group test 2",
		"nonposix":    "true",
		"addattr":     "owner=uid=test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSGroupResource_nonposix_basic(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_nonposix_full(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "nonposix", "true"),
					resource.TestCheckResourceAttr("freeipa_group.group", "addattr.0", testGroup["addattr"]),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_nonposix_full(testGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup2["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup2["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "nonposix", "true"),
					resource.TestCheckResourceAttr("freeipa_group.group", "addattr.0", testGroup2["addattr"]),
				),
			},
		},
	})
}

func testAccFreeIPADNSGroupResource_nonposix_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name       = "%s"
		nonposix = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["nonposix"])
}

func testAccFreeIPADNSGroupResource_nonposix_full(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name        = "%s"
		description  = "%s"
		nonposix = %s
		addattr = ["%s"]
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["nonposix"], dataset["addattr"])
}

func TestAccFreeIPADNSGroup_external(t *testing.T) {
	testGroup := map[string]string{
		"name":        "testgroupext",
		"description": "User group test",
		"external":    "true",
	}
	testGroup2 := map[string]string{
		"name":        "testgroupext2",
		"description": "User group test 2",
		"external":    "true",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPADNSGroupResource_external_basic(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_external_full(testGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "external", "true"),
				),
			},
			{
				Config: testAccFreeIPADNSGroupResource_external_full(testGroup2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_group.group", "name", testGroup2["name"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "description", testGroup2["description"]),
					resource.TestCheckResourceAttr("freeipa_group.group", "external", "true"),
				),
			},
		},
	})
}

func testAccFreeIPADNSGroupResource_external_basic(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name       = "%s"
		external = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["external"])
}

func testAccFreeIPADNSGroupResource_external_full(dataset map[string]string) string {
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
	  
	resource "freeipa_group" "group" {
		name        = "%s"
		description  = "%s"
		external = %s
	}
	`, provider_host, provider_user, provider_pass, dataset["name"], dataset["description"], dataset["external"])
}
