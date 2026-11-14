package user_impersonate_test

import (
	"errors"
	"net/http"
	"testing"
)

func TestImpersonateSuccess(t *testing.T) {
	var gotUserID string
	env := newTestEnv(t, func(w http.ResponseWriter, r *http.Request, userID string) error {
		gotUserID = userID
		return nil
	})
	defer env.close()
	user := env.seedUser(t, "target@t.com", "active")

	rec := get(t, env, "user_id="+user.GetID())
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect on success, got %d", rec.Code)
	}
	if gotUserID != user.GetID() {
		t.Fatalf("expected impersonate called with %q, got %q", user.GetID(), gotUserID)
	}
}

func TestImpersonateMissingID(t *testing.T) {
	env := newTestEnv(t, nil)
	defer env.close()

	rec := get(t, env, "")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

func TestImpersonateUserNotFound(t *testing.T) {
	env := newTestEnv(t, nil)
	defer env.close()

	rec := get(t, env, "user_id=missing")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

func TestImpersonateInactiveUser(t *testing.T) {
	called := false
	env := newTestEnv(t, func(w http.ResponseWriter, r *http.Request, userID string) error {
		called = true
		return nil
	})
	defer env.close()
	user := env.seedUser(t, "off@t.com", "inactive")

	rec := get(t, env, "user_id="+user.GetID())
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
	if called {
		t.Fatal("expected impersonate NOT called for inactive user")
	}
}

func TestImpersonateNotConfigured(t *testing.T) {
	env := newTestEnv(t, nil)
	defer env.close()
	user := env.seedUser(t, "target@t.com", "active")

	rec := get(t, env, "user_id="+user.GetID())
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

func TestImpersonateCallbackError(t *testing.T) {
	env := newTestEnv(t, func(w http.ResponseWriter, r *http.Request, userID string) error {
		return errors.New("session create failed")
	})
	defer env.close()
	user := env.seedUser(t, "target@t.com", "active")

	rec := get(t, env, "user_id="+user.GetID())
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}
