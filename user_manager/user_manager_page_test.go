package user_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserManagerPageContentType(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-manager", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserManager(rec, req)

	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
