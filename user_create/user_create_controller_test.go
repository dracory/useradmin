package user_create_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUserCreateModalRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	req := httptest.NewRequest(http.MethodGet, "/admin/users?controller=user-create", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"ModaluserCreate", "user_first_name", "user_last_name", "user_email"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected %q in modal", want)
		}
	}
}

func TestUserCreateSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := postForm(t, env, url.Values{
		"user_first_name": {"New"},
		"user_last_name":  {"User"},
		"user_email":      {"new@t.com"},
	})
	if !strings.Contains(rec.Body.String(), "user created successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}

	found, err := env.store.UserFindByEmail(context.Background(), "new@t.com")
	if err != nil || found == nil {
		t.Fatal("expected user persisted")
	}
}

func TestUserCreateMissingEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := postForm(t, env, url.Values{"user_first_name": {"New"}})
	if !strings.Contains(rec.Body.String(), "user email is required") {
		t.Fatalf("expected email-required error, got %q", rec.Body.String())
	}
}

func TestUserCreateInvalidEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := postForm(t, env, url.Values{"user_email": {"not-an-email"}})
	if !strings.Contains(rec.Body.String(), "invalid") {
		t.Fatalf("expected invalid-email error, got %q", rec.Body.String())
	}
}

func TestUserCreateDuplicateEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "dup@t.com")

	rec := postForm(t, env, url.Values{"user_email": {"dup@t.com"}})
	if !strings.Contains(rec.Body.String(), "already exists") {
		t.Fatalf("expected duplicate-email error, got %q", rec.Body.String())
	}
}

func TestUserCreateUnsealSkips(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
}
