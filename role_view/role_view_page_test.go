package role_view_test

import (
	"strings"
	"testing"
)

func TestRoleViewPageContentType(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Team")

	rec := get(t, env, "role_id="+role.GetID())
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
