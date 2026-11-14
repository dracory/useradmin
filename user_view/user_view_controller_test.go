package user_view_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserViewPageMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin?controller=user-view", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserView(rec, req)

	if rec.Code == 0 {
		t.Fatal("expected response written")
	}
}
