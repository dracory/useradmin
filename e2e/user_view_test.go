package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestUserViewShowsUser verifies the user view page loads the seeded
// user's details. The page is server-rendered, so waiting for a card
// is enough — no Vue mount delay.
func TestUserViewShowsUser(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("user-view") + "&user_id=" + seedUserID

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("user view page failed: %v", err)
	}
	if !strings.Contains(body, seedUserMail) {
		t.Fatalf("expected seeded user email %q in user view", seedUserMail)
	}
}
