package shared_test

import (
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestBreadcrumbsRender(t *testing.T) {
	html := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: "/admin"},
		{Name: "Users", URL: "/admin/users"},
		{Name: "View"},
	}).ToHTML()

	for _, want := range []string{"Home", "Users", "View", "/admin"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in breadcrumbs", want)
		}
	}
}
