package shared_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestLinksControllers(t *testing.T) {
	l := shared.NewLinks("/admin")
	cases := map[string]string{
		l.Home(nil):            "controller=home",
		l.UserManager(nil):     "controller=user-manager",
		l.UserView(nil):        "controller=user-view",
		l.UserCreate(nil):      "controller=user-create",
		l.UserDelete(nil):      "controller=user-delete",
		l.UserUpdate(nil):      "controller=user-update",
		l.UserImpersonate(nil): "controller=user-impersonate",
		l.RoleManager(nil):     "controller=role-manager",
		l.RoleView(nil):        "controller=role-view",
		l.RoleUpdate(nil):      "controller=role-update",
		l.RoleMembers(nil):     "controller=role-members",
		l.GroupManager(nil):    "controller=group-manager",
		l.GroupView(nil):       "controller=group-view",
		l.GroupUpdate(nil):     "controller=group-update",
		l.GroupMembers(nil):    "controller=group-members",
	}
	for url, want := range cases {
		if !strings.Contains(url, want) {
			t.Fatalf("expected %q in %q", want, url)
		}
		if !strings.HasPrefix(url, "/admin?") {
			t.Fatalf("expected /admin prefix in %q", url)
		}
	}
}

func TestNewLinksFromRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), shared.KeyUserAdminURL, "/admin"))
	l := shared.NewLinksFromRequest(req)
	if !strings.HasPrefix(l.Home(nil), "/admin?") {
		t.Fatalf("expected user admin url from context, got %q", l.Home(nil))
	}
}
