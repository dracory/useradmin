package shared_test

import (
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestErrorAlert(t *testing.T) {
	html := shared.ErrorAlert("oops")
	if !strings.Contains(html, "oops") || !strings.Contains(html, "alert") {
		t.Fatalf("unexpected alert: %q", html)
	}
}
