package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestRoleViewShowsRole verifies the role view page loads the seeded
// role's details. The page is server-rendered — no Vue mount delay.
func TestRoleViewShowsRole(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("role-view") + "&role_id=" + seedRoleID

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("role view page failed: %v", err)
	}
	if !strings.Contains(body, seedRoleName) {
		t.Fatalf("expected seeded role name %q in role view", seedRoleName)
	}
}
