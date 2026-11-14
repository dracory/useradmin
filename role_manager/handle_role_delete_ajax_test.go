package role_manager_test

import (
	"context"
	"testing"
)

func TestDeleteRoleSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Doomed", "doomed", "active")

	resp := postAction(t, env, "delete-role-ajax", nil, map[string]any{
		"role_id": role.GetID(),
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	found, err := env.store.RoleFindByID(context.Background(), role.GetID())
	if err != nil || found != nil {
		t.Fatal("expected role deleted")
	}
}

func TestDeleteRoleMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "delete-role-ajax", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
