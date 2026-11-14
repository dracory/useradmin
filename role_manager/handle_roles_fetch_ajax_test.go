package role_manager_test

import "testing"

func TestLoadRolesSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedRole(t, "Admins", "admins", "active")
	env.seedRole(t, "Members", "members", "active")

	resp := postAction(t, env, "load-roles-ajax", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	roles, _ := resp.data()["roles"].([]any)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
}
