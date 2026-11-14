package group_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGroupManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=group-manager", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.GroupManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Group") {
		t.Fatal("expected group manager page")
	}
}
