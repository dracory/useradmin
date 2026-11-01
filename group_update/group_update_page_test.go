package group_update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRenderPageContainsSections verifies the edit page renders all
// section cards: Status, Group Details, and Admin Notes.
func TestRenderPageContainsSections(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=group-update&group_id="+group.GetID(), nil)
	rec := httptest.NewRecorder()
	env.ui.GroupUpdate(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		"Status",
		"Group Details",
		"Admin Notes",
		"app-group-update",
		"status-update-ajax",
		"details-update-ajax",
		"notes-update-ajax",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected page to contain %q", want)
		}
	}
}

// TestRenderPageMissingGroupRedirects verifies that an unknown group_id
// redirects back to the group manager.
func TestRenderPageMissingGroupRedirects(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=group-update&group_id=nonexistent", nil)
	rec := httptest.NewRecorder()
	env.ui.GroupUpdate(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect %d, got %d", http.StatusSeeOther, rec.Code)
	}
}
