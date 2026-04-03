package cmd

import "testing"

func TestPullCommandRequiresImage(t *testing.T) {
	if err := pullCommand.Args(pullCommand, []string{}); err == nil {
		t.Fatal("expected missing image to be rejected")
	}
}

func TestXCommandRequiresImage(t *testing.T) {
	if err := xCommand.Args(xCommand, []string{}); err == nil {
		t.Fatal("expected missing image to be rejected")
	}
}
