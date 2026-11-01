package role_update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRenderPageContainsSections verifies the edit page renders all
// section cards: Status, Role Details, and Admin Notes.
func TestRenderPageContainsSections(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=role-update&role_id="+role.GetID(), nil)
	rec := httptest.NewRecorder()
	env.ui.RoleUpdate(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		"Status",
		"Role Details",
		"Admin Notes",
		"app-role-update",
		"status-update-ajax",
		"details-update-ajax",
		"notes-update-ajax",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected page to contain %q", want)
		}
	}
}

// TestRenderPageMissingRoleRedirects verifies that an unknown role_id
// redirects back to the role manager.
func TestRenderPageMissingRoleRedirects(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=role-update&role_id=nonexistent", nil)
	rec := httptest.NewRecorder()
	env.ui.RoleUpdate(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect %d, got %d", http.StatusSeeOther, rec.Code)
	}
}
