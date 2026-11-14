package group_members_test

import (
	"strings"
	"testing"
)

func TestGroupMembersPageRenders(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	g := env.seedGroup(t, "Editors", "editors")

	rec := get(t, env, "group_id="+g.GetID())
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Editor") && !strings.Contains(rec.Body.String(), "Members") {
		t.Fatal("expected members page")
	}
}
