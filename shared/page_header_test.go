package shared_test

import (
	"strings"
	"testing"

	"github.com/dracory/hb"
	"github.com/dracory/useradmin/shared"
)

func TestPageHeader(t *testing.T) {
	html := shared.PageHeader("Users", hb.Button().HTML("New")).ToHTML()
	if !strings.Contains(html, "Users") || !strings.Contains(html, "New") {
		t.Fatalf("unexpected header: %q", html)
	}
}
