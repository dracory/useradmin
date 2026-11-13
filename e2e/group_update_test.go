package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestGroupUpdatePrefillsForm verifies the group update page loads the
// seeded group into the Vue form (fields populate only after mount +
// AJAX load completes). v-model is compiled away, so the name input
// is matched by position — it is the first text input in the app.
func TestGroupUpdatePrefillsForm(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("group-update") + "&group_id=" + seedGroupID

	var name string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Poll(
			`(document.querySelectorAll('#app-group-update input[type="text"]')[0]||{}).value !== ''`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Value(`#app-group-update input[type="text"]`, &name, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("group update page failed: %v", err)
	}
	if name != seedGroupName {
		t.Fatalf("expected prefilled group name %q, got %q", seedGroupName, name)
	}
}
