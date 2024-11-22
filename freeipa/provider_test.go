// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"os"
	"testing"
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
