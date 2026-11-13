package e2e

import (
	"testing"

	"github.com/chromedp/chromedp"
)

// TestUserDeleteShowsModal verifies the user delete controller renders
// its confirmation modal for a seeded user. The delete itself is not
// submitted — mutating state is covered by the create-flow tests.
func TestUserDeleteShowsModal(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("user-delete") + "&user_id=" + seedUserID

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#ModalUserDelete", chromedp.ByID),
	)
	if err != nil {
		t.Fatalf("user delete modal did not render: %v", err)
	}
}
