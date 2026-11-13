package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestUserImpersonateRedirects verifies the impersonate controller
// calls the configured OnUserImpersonate (sets a session cookie) and
// redirects the browser to the user home URL ("/").
func TestUserImpersonateRedirects(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("user-impersonate") + "&user_id=" + seedUserID

	var path string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second),
		chromedp.Evaluate(`document.location.pathname`, &path),
	)
	if err != nil {
		t.Fatalf("impersonate navigation failed: %v", err)
	}
	if path != "/" {
		t.Fatalf("expected redirect to '/', landed on %q", path)
	}
}
