package main

import "testing"

func TestDummy(t *testing.T) {
	// A simple dummy test to ensure the CI test runner executes correctly
	if 1 != 1 {
		t.Errorf("Math is broken")
	}
}
