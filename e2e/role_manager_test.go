package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestRoleManagerLoadsSeededRoles(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("role-manager"))
	if !strings.Contains(body, "Editor") {
		t.Fatalf("expected seeded 'Editor' in role manager")
	}
}

// TestRoleCreateFlow verifies the create-role modal adds a role that
// then appears in the manager table.
func TestRoleCreateFlow(t *testing.T) {
	ctx := newTab(t)
	name := "E2E Role " + time.Now().Format("150405")

	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("role-manager")),
		chromedp.WaitReady("table tbody tr:nth-of-type(2)", chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("role manager did not load: %v", err)
	}

	err = chromedp.Run(ctx,
		chromedp.Click(`button.btn.btn-primary.btn-sm`, chromedp.ByQuery),
		chromedp.WaitReady(`input[placeholder="Role name"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[placeholder="Role name"]`, name, chromedp.ByQuery),
		chromedp.SendKeys(`input[placeholder="role-handle"]`, "e2e-role-"+time.Now().Format("150405"), chromedp.ByQuery),
		chromedp.Click(`.modal-footer button.btn-primary`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		t.Fatalf("create modal interaction failed: %v", err)
	}

	var body string
	err = chromedp.Run(ctx, chromedp.Text("body", &body))
	if err != nil {
		t.Fatalf("reading role manager failed: %v", err)
	}
	if !strings.Contains(body, name) {
		t.Fatalf("expected newly created role %q in manager", name)
	}
}
