package group_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGroupManagerPageContentType(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=group-manager", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.GroupManager(rec, req)

	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
