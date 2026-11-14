package user_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-manager", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "User") {
		t.Fatal("expected user manager page")
	}
}
