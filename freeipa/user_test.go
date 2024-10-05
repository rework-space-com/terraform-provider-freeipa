package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAUser_full(t *testing.T) {
	managerUser := map[string]string{
		"index":     "0",
		"login":     "\"devmanager\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"Manager\"",
	}
	testUser := map[string]string{
		"index":                    "1",
		"login":                    "\"testuser\"",
		"firstname":                "\"Test\"",
		"lastname":                 "\"User\"",
		"account_disabled":         "false",
		"car_license":              "[\"A-111-B\"]",
		"city":                     "\"El Mundo\"",
		"display_name":             "\"Test User\"",
		"email_address":            "[\"testuser@example.com\"]",
		"employee_number":          "\"000001\"",
		"employee_type":            "\"Developer\"",
		"full_name":                "\"Test User\"",
		"gecos":                    "\"Test User\"",
		"gid_number":               "10001",
		"home_directory":           "\"/home/testuser\"",
		"initials":                 "\"TU\"",
		"job_title":                "\"Developer\"",
		"krb_principal_name":       "[\"tuser@IPATEST.LAN\"]",
		"login_shell":              "\"/bin/bash\"",
		"manager":                  "\"devmanager\"",
		"mobile_numbers":           "[\"0123456789\"]",
		"organisation_unit":        "\"Devs\"",
		"postal_code":              "\"12345\"",
		"preferred_language":       "\"English\"",
		"province":                 "\"England\"",
		"random_password":          "false",
		"ssh_public_key":           "[\"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDDmMkNHn3R+DzSamQDSW60a0iVlAvzbuC3auu8lNoi3u6lvMemsZqPTuvfY4Xlf7uzm+dya3fTRdPKn8sYgPwQ4saUpCSlegN44PjJMhonR1a7FbpHLWj8CRRfzdUSznQhzFcFff0wMBYAklXlyjvdFM8ahl7zHO08HR6469XOVwO1Tb3OGPrXB2lzStK5PKfk5DO/IKl4vHSKhVNVnsZe52rHiZrxOqdGyCijtvwmW2YfIAGc1k4Seqn/Nn7NxKIFBH3hxaUDqgpZneXzuw9GI/F0M8phnHxXNFVZvIWZVcanEeXtH9Z+vVx1ujNcB2QhiPfLMqkNl9db7uykSGKFM4jD0UjGj5kJ8TOC39Safk7XzpQTnrqvIi158zBHVSgugth+QsE1I9/PL2wlzx1qWV2991JKIOc8m52Iwq02tyO8JaSssFTk9szkLTAHedPnZeBbdnlRYcHqX+NPaUh3hqRTZBIR79Ruk6WAliFkED1L0SgwDfGFlevn1Kde9ok=\"]",
		"street_address":           "\"1600, Pensylvania av.\"",
		"telephone_numbers":        "[\"1234567890\"]",
		"uid_number":               "10001",
		"userpassword":             "\"P@ssword\"",
		"krb_principal_expiration": "\"2049-12-31T23:59:59Z\"",
		"krb_password_expiration":  "\"2049-12-31T23:59:59Z\"",
		"userclass":                "[\"user-account\"]",
	}
	testUserDS := map[string]string{
		"index": "1",
		"name":  "\"testuser\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "devmanager"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "Manager"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-1", "name", "testuser"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "first_name", "Test"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "last_name", "User"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUser),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUser) + testAccFreeIPAUser_datasource(testUserDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "employee_number", "000001"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "employee_type", "Developer"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "home_directory", "/home/testuser"),
				),
			},
		},
	})
}
