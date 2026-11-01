package user_update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRenderPageContainsSections verifies the edit page renders all
// section cards: Status, User Details, Roles, Location, Change Password,
// and Admin Notes.
func TestRenderPageContainsSections(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-update&user_id="+user.GetID(), nil)
	rec := httptest.NewRecorder()
	env.ui.UserUpdate(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		"User Details",
		"Change Password",
		"Admin Notes",
		"Location",
		"app-user-update",
		"status-update-ajax",
		"password-update-ajax",
		"details-update-ajax",
		"roles-update-ajax",
		"location-update-ajax",
		"notes-update-ajax",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected page to contain %q", want)
		}
	}
}

// TestRenderPageMissingUserRedirects verifies that an unknown user_id
// redirects back to the user manager.
func TestRenderPageMissingUserRedirects(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-update&user_id=nonexistent", nil)
	rec := httptest.NewRecorder()
	env.ui.UserUpdate(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect %d, got %d", http.StatusSeeOther, rec.Code)
	}
}
