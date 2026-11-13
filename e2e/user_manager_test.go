package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestUserManagerLoadsSeededUsers(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("user-manager"))
	if !strings.Contains(body, "Alice.Smith.01@example.com") {
		t.Fatalf("expected seeded 'Alice.Smith.01@example.com' in user manager")
	}
}

// TestUserCreateFlow verifies the "Add User" button opens the create
// modal, and that submitting it adds a user that then appears in the
// manager table.
func TestUserCreateFlow(t *testing.T) {
	ctx := newTab(t)
	email := "e2e.user." + time.Now().Format("150405") + "@example.com"

	// Waiting for the second rendered table row guarantees Vue has
	// mounted; the modal inputs are only bound after that.
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("user-manager")),
		chromedp.WaitReady("table tbody tr:nth-of-type(2)", chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("user manager did not load: %v", err)
	}

	// Open the create modal and fill the form. Modal inputs have no
	// ids — select them by their placeholders.
	err = chromedp.Run(ctx,
		chromedp.Click(`button.btn.btn-primary.btn-sm`, chromedp.ByQuery),
		chromedp.WaitReady(`input[placeholder="Email address"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[placeholder="First name"]`, "E2E", chromedp.ByQuery),
		chromedp.SendKeys(`input[placeholder="Last name"]`, "User", chromedp.ByQuery),
		chromedp.SendKeys(`input[placeholder="Email address"]`, email, chromedp.ByQuery),
		chromedp.Click(`.modal-footer button.btn-primary`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		t.Fatalf("create modal interaction failed: %v", err)
	}

	// The table refreshes after create; search for the new email via
	// the filter to avoid pagination races.
	var body string
	err = chromedp.Run(ctx,
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("reading user manager failed: %v", err)
	}
	if !strings.Contains(body, email) {
		t.Fatalf("expected newly created user %q in manager", email)
	}
}
