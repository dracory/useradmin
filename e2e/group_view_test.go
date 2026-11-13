package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestGroupViewShowsGroup verifies the group view page loads the seeded
// group's details. The page is server-rendered — no Vue mount delay.
func TestGroupViewShowsGroup(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("group-view") + "&group_id=" + seedGroupID

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("group view page failed: %v", err)
	}
	if !strings.Contains(body, seedGroupName) {
		t.Fatalf("expected seeded group name %q in group view", seedGroupName)
	}
}
