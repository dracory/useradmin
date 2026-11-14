package shared_test

import (
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestBreadcrumbStruct(t *testing.T) {
	b := shared.Breadcrumb{Name: "Home", URL: "/admin"}
	if b.Name != "Home" || b.URL != "/admin" {
		t.Fatal("unexpected fields")
	}
}
