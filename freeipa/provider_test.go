// This file was originally inspired by the module structure and design patterns
// used in HashiCorp projects, but all code in this file was written from scratch.
//
// Previously licensed under the MPL-2.0.
// This file is now relicensed under the GNU General Public License v3.0 only,
// as permitted by Section 1.10 of the MPL.
//
// Authors:
//	Antoine Gatineau <antoine.gatineau@infra-monkey.com>
//	Roman Butsiy <butsiyroman@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-only

package freeipa

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	ipa "github.com/infra-monkey/go-freeipa/freeipa"
)

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	if v := os.Getenv("FREEIPA_HOST"); v == "" {
		t.Fatal("FREEIPA_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("FREEIPA_USERNAME"); v == "" {
		t.Fatal("FREEIPA_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("FREEIPA_PASSWORD"); v == "" {
		t.Fatal("FREEIPA_PASSWORD must be set for acceptance tests")
	}
}

func testAccFreeIPAClient(t *testing.T) *ipa.Client {
	client, err := (&freeipaProvider{}).NewFreeIPAClient(context.Background(), &freeipaProviderModel{
		Host:               types.StringValue(os.Getenv("FREEIPA_HOST")),
		Username:           types.StringValue(os.Getenv("FREEIPA_USERNAME")),
		Password:           types.StringValue(os.Getenv("FREEIPA_PASSWORD")),
		InsecureSkipVerify: types.BoolValue(true),
	})
	if err != nil {
		t.Fatalf("Error creating FreeIPA client: %s", err)
	}
	return client
}
