package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestHomePageLoads(t *testing.T) {
	ctx := newTab(t)
	var body string
	err := chromedp.Run(ctx, chromedp.Navigate(baseURL+adminPath))
	if err != nil {
		t.Fatalf("navigate: %v", err)
	}
	err = chromedp.Run(ctx,
		chromedp.WaitReady("body"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if strings.TrimSpace(body) == "" {
		t.Fatal("expected non-empty home page")
	}
}
