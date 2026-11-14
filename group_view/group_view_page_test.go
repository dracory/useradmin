package group_view_test

import (
	"strings"
	"testing"
)

func TestGroupViewPageContentType(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	group := env.seedGroup(t, "Team")

	rec := get(t, env, "group_id="+group.GetID())
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
