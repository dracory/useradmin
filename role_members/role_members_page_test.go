package role_members_test

import (
	"net/http"
	"testing"
)

func TestRoleMembersPageRequiresRoleID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "")
	if rec.Code != http.StatusSeeOther && rec.Code != http.StatusOK {
		t.Fatalf("expected redirect or alert, got %d", rec.Code)
	}
}
