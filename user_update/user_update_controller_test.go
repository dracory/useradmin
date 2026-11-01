package user_update_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerDispatchesSectionActions verifies the controller routes each
// section action to its handler — a missing user_id must produce the
// handler's validation error, not the page render fallback.
func TestHandlerDispatchesSectionActions(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	actions := []string{
		"status-update-ajax",
		"details-update-ajax",
		"roles-update-ajax",
		"location-update-ajax",
		"notes-update-ajax",
		"password-update-ajax",
	}
	for _, action := range actions {
		resp := postAction(t, env, action, nil, map[string]any{})
		if resp.Status != "error" {
			t.Errorf("action %s: expected status=error, got %q", action, resp.Status)
		}
		if resp.Message != "User ID is required" {
			t.Errorf("action %s: expected 'User ID is required', got %q", action, resp.Message)
		}
	}
}

// TestHandlerUnknownActionRendersPage verifies that a request without a
// recognized action falls through to the page renderer — which redirects
// to the user manager when no user_id is given.
func TestHandlerUnknownActionRedirects(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-update", nil)
	rec := httptest.NewRecorder()
	env.ui.UserUpdate(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect %d, got %d", http.StatusSeeOther, rec.Code)
	}
}
