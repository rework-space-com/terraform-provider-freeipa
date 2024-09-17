package freeipa

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"freeipa": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccFreeIPAProvider() string {
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
	
	`, provider_host, provider_user, provider_pass)
}

func testAccFreeIPAGroup_resourcefull(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_group" "group-%s" {
		name        = "%s"
		description = "%s"
		gid_number  = %s
		addattr = ["%s"]
		setattr = ["%s"]
	}

	`, dataset["index"], dataset["name"], dataset["description"], dataset["gid_number"], dataset["addattr"], dataset["setattr"])
}

func testAccFreeIPAGroupNonposix_resourcefull(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_group" "group-%s" {
		name        = "%s"
		description = "%s"
		nonposix = %s
		addattr = ["%s"]
		setattr = ["%s"]
	}

	`, dataset["index"], dataset["name"], dataset["description"], dataset["nonposix"], dataset["addattr"], dataset["setattr"])
}

func testAccFreeIPAGroupExternal_resourcefull(dataset map[string]string) string {
	return fmt.Sprintf(`
	resource "freeipa_group" "group-%s" {
		name        = "%s"
		description = "%s"
		external = %s
		addattr = ["%s"]
		setattr = ["%s"]
	}

	`, dataset["index"], dataset["name"], dataset["description"], dataset["external"], dataset["addattr"], dataset["setattr"])
}

func testAccFreeIPAGroup_datasource(dataset map[string]string) string {
	return fmt.Sprintf(`
	data "freeipa_group" "group-%s" {
		name        = "%s"
	}

	`, dataset["index"], dataset["name"])
}
