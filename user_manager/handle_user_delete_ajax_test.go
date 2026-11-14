package user_manager_test

import (
	"testing"
)

func TestDeleteUserSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	user := env.seedUser(t, "doomed@t.com", "active")

	resp := postAction(t, env, "delete-user-ajax", nil, map[string]any{
		"user_id": user.GetID(),
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestDeleteUserMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "delete-user-ajax", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "delete-user-ajax", nil, map[string]any{
		"user_id": "missing",
	})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
