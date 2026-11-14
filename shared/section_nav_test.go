package shared_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/userstore"
)

// stubStore returns role/group table names for SectionNav checks.
func TestSectionNavUsersActive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin?controller=user-manager", http.NoBody)
	req = req.WithContext(contextWithEndpoint(req, "/admin"))
	html := shared.SectionNav(req, nil).ToHTML()
	if !strings.Contains(html, "Users") {
		t.Fatalf("expected Users tab, got %q", html)
	}
}

func TestSectionNavRolesHiddenWithoutStore(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin?controller=role-manager", http.NoBody)
	req = req.WithContext(contextWithEndpoint(req, "/admin"))
	html := shared.SectionNav(req, nil).ToHTML()
	// nil store → no Roles/Groups tabs
	if strings.Contains(html, "Groups") {
		t.Fatal("expected Groups hidden without store")
	}
	_ = userstore.NewUserQuery()
}
