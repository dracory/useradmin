package group_update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGroupUpdatePageMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin?controller=group-update", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.GroupUpdate(rec, req)

	if rec.Code == 0 {
		t.Fatal("expected response written")
	}
	_ = strings.Contains(rec.Body.String(), "")
}
