package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestRoleUpdatePrefillsForm verifies the role update page loads the
// seeded role into the Vue form (fields populate only after mount +
// AJAX load completes). v-model is compiled away, so the name input
// is matched by position — it is the first text input in the app.
func TestRoleUpdatePrefillsForm(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("role-update") + "&role_id=" + seedRoleID

	var name string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Poll(
			`(document.querySelectorAll('#app-role-update input[type="text"]')[0]||{}).value !== ''`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Value(`#app-role-update input[type="text"]`, &name, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("role update page failed: %v", err)
	}
	if name != seedRoleName {
		t.Fatalf("expected prefilled role name %q, got %q", seedRoleName, name)
	}
}
