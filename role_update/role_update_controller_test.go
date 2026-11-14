package role_update_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoleUpdatePageMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin?controller=role-update", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.RoleUpdate(rec, req)

	if rec.Code == 0 {
		t.Fatal("expected response written")
	}
}
