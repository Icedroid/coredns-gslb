package gslb

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", `gslb`)
	if err := setup(c); err != nil {
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `gslb example.org`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}
}
