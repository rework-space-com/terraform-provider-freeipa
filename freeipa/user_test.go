// Authors:
//   Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccFreeIPAUser_full(t *testing.T) {
	managerUser := map[string]string{
		"index":     "0",
		"login":     "\"testacc-devmanager\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"Manager\"",
	}
	testUser := map[string]string{
		"index":                    "1",
		"login":                    "\"testacc-user\"",
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
		"manager":                  "\"testacc-devmanager\"",
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
		"userclass":                "[\"user-account\"]",
	}
	testUserModified := map[string]string{
		"index":                    "1",
		"login":                    "\"testacc-user\"",
		"firstname":                "\"Test\"",
		"lastname":                 "\"User\"",
		"account_disabled":         "false",
		"car_license":              "[\"A-222-B\"]",
		"city":                     "\"The World\"",
		"display_name":             "\"Test User Modfied\"",
		"email_address":            "[\"testuser@example.com\",\"testuser2@example.com\"]",
		"employee_number":          "\"000002\"",
		"employee_type":            "\"Tester\"",
		"full_name":                "\"Test User Modfied\"",
		"gecos":                    "\"Test User Modfied\"",
		"gid_number":               "10002",
		"home_directory":           "\"/home/testaccuser\"",
		"initials":                 "\"TUM\"",
		"job_title":                "\"Tester\"",
		"krb_principal_name":       "[\"tuser@IPATEST.LAN\",\"testaccuser@IPATEST.LAN\"]",
		"login_shell":              "\"/bin/fish\"",
		"manager":                  "\"testacc-devmanager\"",
		"mobile_numbers":           "[\"1234567890\"]",
		"organisation_unit":        "\"Testers\"",
		"postal_code":              "\"12340\"",
		"preferred_language":       "\"French\"",
		"province":                 "\"France\"",
		"random_password":          "false",
		"ssh_public_key":           "[]",
		"street_address":           "\"1 Main Street\"",
		"telephone_numbers":        "[\"0123456789\"]",
		"uid_number":               "10001",
		"userpassword":             "\"Password\"",
		"krb_principal_expiration": "\"2050-12-31T23:59:59Z\"",
		"userclass":                "[\"unprivileged-user-account\"]",
	}
	testUserModified2 := map[string]string{
		"index":                    "1",
		"login":                    "\"testacc-user\"",
		"firstname":                "\"Test\"",
		"lastname":                 "\"User\"",
		"account_disabled":         "true",
		"car_license":              "[\"A-222-B\"]",
		"city":                     "\"The World\"",
		"display_name":             "\"Test User Modfied\"",
		"email_address":            "[\"testuser@example.com\",\"testuser2@example.com\"]",
		"employee_number":          "\"000002\"",
		"employee_type":            "\"Tester\"",
		"full_name":                "\"Test User Modfied\"",
		"gecos":                    "\"Test User Modfied\"",
		"gid_number":               "10002",
		"home_directory":           "\"/home/testaccuser\"",
		"initials":                 "\"TUM\"",
		"job_title":                "\"Tester\"",
		"krb_principal_name":       "[\"tuser@IPATEST.LAN\",\"testaccuser@IPATEST.LAN\"]",
		"login_shell":              "\"/bin/fish\"",
		"manager":                  "\"testacc-devmanager\"",
		"mobile_numbers":           "[\"1234567890\"]",
		"organisation_unit":        "\"Testers\"",
		"postal_code":              "\"12340\"",
		"preferred_language":       "\"French\"",
		"province":                 "\"France\"",
		"random_password":          "false",
		"ssh_public_key":           "[]",
		"street_address":           "\"1 Main Street\"",
		"telephone_numbers":        "[\"0123456789\"]",
		"uid_number":               "10001",
		"userpassword":             "\"Password\"",
		"krb_principal_expiration": "\"2050-12-31T23:59:59Z\"",
		"userclass":                "[\"unprivileged-user-account\"]",
	}
	testUserDS := map[string]string{
		"index": "1",
		"name":  "freeipa_user.user-1.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "testacc-devmanager"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "Manager"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-1", "name", "testacc-user"),
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
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUserModified),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-1", "name", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "first_name", "Test"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "last_name", "User"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "account_disabled", "false"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "car_license.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "car_license.0", "A-222-B"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "city", "The World"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "display_name", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.0", "testuser@example.com"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.1", "testuser2@example.com"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "employee_number", "000002"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "employee_type", "Tester"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "full_name", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "gecos", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "gid_number", "10002"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "home_directory", "/home/testaccuser"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "initials", "TUM"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "job_title", "Tester"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.0", "tuser@IPATEST.LAN"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.1", "testaccuser@IPATEST.LAN"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "login_shell", "/bin/fish"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "manager", "testacc-devmanager"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "mobile_numbers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "mobile_numbers.0", "1234567890"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "organisation_unit", "Testers"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "postal_code", "12340"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "preferred_language", "French"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "street_address", "1 Main Street"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "telephone_numbers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "telephone_numbers.0", "0123456789"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "uid_number", "10001"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_expiration", "2050-12-31T23:59:59Z"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "userclass.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "userclass.0", "unprivileged-user-account"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUserModified2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-1", "name", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "first_name", "Test"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "last_name", "User"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "account_disabled", "true"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "car_license.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "car_license.0", "A-222-B"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "city", "The World"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "display_name", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.0", "testuser@example.com"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "email_address.1", "testuser2@example.com"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "employee_number", "000002"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "employee_type", "Tester"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "full_name", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "gecos", "Test User Modfied"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "gid_number", "10002"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "home_directory", "/home/testaccuser"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "initials", "TUM"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "job_title", "Tester"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.#", "2"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.0", "tuser@IPATEST.LAN"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_name.1", "testaccuser@IPATEST.LAN"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "login_shell", "/bin/fish"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "manager", "testacc-devmanager"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "mobile_numbers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "mobile_numbers.0", "1234567890"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "organisation_unit", "Testers"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "postal_code", "12340"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "preferred_language", "French"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "street_address", "1 Main Street"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "telephone_numbers.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "telephone_numbers.0", "0123456789"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "uid_number", "10001"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "krb_principal_expiration", "2050-12-31T23:59:59Z"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "userclass.#", "1"),
					resource.TestCheckResourceAttr("freeipa_user.user-1", "userclass.0", "unprivileged-user-account"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_resource(testUserModified2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAUser_simple_CaseInsensitive(t *testing.T) {
	managerUser := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-DevManager\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"Manager\"",
	}
	testUserDS := map[string]string{
		"index": "1",
		"name":  "freeipa_user.user-0.name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-devmanager"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-DevManager"),
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
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_datasource(testUserDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "name", "TestACC-DevManager"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-1", "last_name", "Manager"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(managerUser) + testAccFreeIPAUser_datasource(testUserDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAUser_staged(t *testing.T) {
	testUser := map[string]string{
		"index":          "0",
		"login":          "\"TestACC-User\"",
		"firstname":      "\"Dev\"",
		"lastname":       "\"User\"",
		"account_staged": "true",
	}
	testUserDS := map[string]string{
		"index":          "0",
		"name":           "freeipa_user.user-0.name",
		"account_staged": "true",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPAUser_datasource(testUserDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "last_name", "User"),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUser) + testAccFreeIPAUser_datasource(testUserDS),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccFreeIPAUser_lifecycle(t *testing.T) {
	testUserStaged := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"User\"",
		"state":     "\"staged\"",
	}
	testUserStagedDS := map[string]string{
		"index": "0",
		"name":  "freeipa_user.user-0.name",
		"state": "\"staged\"",
	}
	testUserActive := map[string]string{
		"index":             "0",
		"login":             "\"TestACC-User\"",
		"state":             "\"active\"",
		"firstname":         "\"Dev\"",
		"lastname":          "\"User\"",
		"organisation_unit": "\"Developers\"",
		"telephone_numbers": "[\"1234567890\"]",
		"mobile_numbers":    "[\"0123456789\"]",
		"car_license":       "[\"A-111-B\"]",
		"street_address":    "\"Mykhaïlo Hrouchevsky Street, 12/2\"",
		"city":              "\"Kyiv\"",
		"province":          "\"Ukraine\"",
		"postal_code":       "\"01008\"",
		"ssh_public_key":    "[\"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDEobivsBY6REElik0hNMLkwcbqIba4c9RWRhD0cy7Kh\"]",
	}
	testUserActiveDS := map[string]string{
		"index": "0",
		"name":  "freeipa_user.user-0.name",
	}
	testUserDisabled := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"User\"",
		"state":     "\"disabled\"",
	}
	testUserPreserved := map[string]string{
		"index":     "0",
		"login":     "\"TestACC-User\"",
		"firstname": "\"Dev\"",
		"lastname":  "\"User\"",
		"state":     "\"preserved\"",
	}
	testUserPreservedDS := map[string]string{
		"index": "0",
		"name":  "freeipa_user.user-0.name",
		"state": "\"preserved\"",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1
			{
				Config:      testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserPreserved),
				ExpectError: regexp.MustCompile("Creating a preserved user is not allowed."),
			},
			// Step 2
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 3
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 4
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged) + testAccFreeIPAUser_datasource(testUserStagedDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 5
			{
				Config:      testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserPreserved),
				ExpectError: regexp.MustCompile("Preserving a staged user is not allowed."),
			},
			// Step 6
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserActive),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 7
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserActive),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 8
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserActive) + testAccFreeIPAUser_datasource(testUserActiveDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "last_name", "User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "organisation_unit", "Developers"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "telephone_numbers.0", "1234567890"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "mobile_numbers.0", "0123456789"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "car_license.0", "A-111-B"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "street_address", "Mykhaïlo Hrouchevsky Street, 12/2"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "city", "Kyiv"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "province", "Ukraine"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "postal_code", "01008"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "ssh_public_key.0", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDEobivsBY6REElik0hNMLkwcbqIba4c9RWRhD0cy7Kh"),
				),
			},
			// Step 9
			{
				Config:      testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged),
				ExpectError: regexp.MustCompile("Staging an active user is not allowed."),
			},
			// Step 10
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 11
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserDisabled),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 12
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserDisabled) + testAccFreeIPAUser_datasource(testUserActiveDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 13
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserPreserved),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 14
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserPreserved),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 15
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserPreserved) + testAccFreeIPAUser_datasource(testUserPreservedDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("data.freeipa_user.user-0", "last_name", "User"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "organisation_unit"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "telephone_numbers.0"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "mobile_numbers.0"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "car_license.0"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "street_address"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "city"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "province"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "postal_code"),
					resource.TestCheckNoResourceAttr("data.freeipa_user.user-0", "ssh_public_key"),
				),
			},
			// Step 16
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 17
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserStaged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 18
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 19
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserDisabled),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 20
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserActive),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_user.user-0", "id", "testacc-user"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "name", "TestACC-User"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "first_name", "Dev"),
					resource.TestCheckResourceAttr("freeipa_user.user-0", "last_name", "User"),
				),
			},
			// Step 21
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPAUser_resource(testUserActive),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
