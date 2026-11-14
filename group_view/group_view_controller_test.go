package group_view_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestGroupViewPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	group := env.seedGroup(t, "Admins")

	rec := get(t, env, "group_id="+group.GetID())
	if rec.Code != http.StatusOK {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Admins") {
		t.Fatal("expected group name on page")
	}
}

func TestGroupViewMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}

func TestGroupViewNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "group_id=missing")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rec.Code)
	}
}
