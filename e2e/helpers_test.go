package e2e

import (
	"context"
	"testing"

	"github.com/chromedp/chromedp"
)

// adminURL builds a controller URL, e.g. /admin/users?controller=user-manager.
func adminURL(controller string) string {
	return baseURL + adminPath + "?controller=" + controller
}

// waitTable navigates to a manager page and waits until Vue has loaded
// data and rendered the results table. Returns the page's text content.
//
// Waits for the SECOND table row: the raw v-for template row exists in
// the DOM before Vue mounts, so `table tr` alone can resolve while the
// table is still a template. A second row only exists after v-for has
// expanded with real data.
func waitTable(t *testing.T, ctx context.Context, url string) string {
	t.Helper()
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("table tbody tr:nth-of-type(2)", chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("page %s failed: %v", url, err)
	}
	return body
}
