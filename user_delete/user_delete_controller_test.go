package user_delete_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUserDeleteModalRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	user := env.seedUser(t, "d@t.com")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/users?controller=user-delete&user_id="+user.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserDelete(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"ModalUserDelete", "user_id", "Delete"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected %q in modal", want)
		}
	}
}

func TestUserDeleteSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	user := env.seedUser(t, "doomed@t.com")

	rec := postForm(t, env, url.Values{"user_id": {user.GetID()}})
	if !strings.Contains(rec.Body.String(), "user deleted successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}

	found, err := env.store.UserFindByID(context.Background(), user.GetID())
	if err != nil || found != nil {
		t.Fatal("expected user soft-deleted")
	}
}

func TestUserDeleteMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := postForm(t, env, url.Values{})
	if !strings.Contains(rec.Body.String(), "user id is required") {
		t.Fatalf("expected id-required error, got %q", rec.Body.String())
	}
}

func TestUserDeleteNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := postForm(t, env, url.Values{"user_id": {"missing"}})
	if !strings.Contains(rec.Body.String(), "User not found") {
		t.Fatalf("expected not-found error, got %q", rec.Body.String())
	}
}
