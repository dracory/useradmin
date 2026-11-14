package role_view_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestRoleViewPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Managers")

	rec := get(t, env, "role_id="+role.GetID())
	if rec.Code != http.StatusOK {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Managers") {
		t.Fatal("expected group name on page")
	}
}

func TestRoleViewMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

func TestRoleViewNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "role_id=missing")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}
