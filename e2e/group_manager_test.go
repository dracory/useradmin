package e2e

import (
	"strings"
	"testing"
)

func TestGroupManagerLoadsSeededGroups(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("group-manager"))
	if !strings.Contains(body, "Customers") {
		t.Fatalf("expected seeded 'Customers' in group manager")
	}
}
