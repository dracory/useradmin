package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestUserUpdatePrefillsForm verifies the user update page loads the
// seeded user into the Vue form (fields populate only after mount +
// AJAX load completes). v-model is compiled away; the email input is
// the only type="email" field in the app.
func TestUserUpdatePrefillsForm(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("user-update") + "&user_id=" + seedUserID

	var email string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Poll(
			`(document.querySelector('#app-user-update input[type="email"]')||{}).value !== ''`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Value(`#app-user-update input[type="email"]`, &email, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("user update page failed: %v", err)
	}
	if email != seedUserMail {
		t.Fatalf("expected prefilled user email %q, got %q", seedUserMail, email)
	}
}
