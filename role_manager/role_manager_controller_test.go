package role_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoleManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=role-manager", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.RoleManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Role") {
		t.Fatal("expected role manager page")
	}
}
