package role_members_test

import (
	"strings"
	"testing"
)

func TestRoleMembersPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")

	rec := get(t, env, "role_id="+role.GetID())
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Editor") && !strings.Contains(rec.Body.String(), "Members") {
		t.Fatal("expected members page")
	}
}
