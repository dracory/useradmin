package role_members_test

import (
	"context"
	"testing"

	"github.com/dracory/userstore"
)

func TestMemberRemoveSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")
	user := env.seedUser(t, "m@t.com")

	ur := userstore.NewUserRole()
	ur.SetUserID(user.GetID())
	ur.SetRoleID(role.GetID())
	if err := env.store.UserRoleCreate(context.Background(), ur); err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	resp := postAction(t, env, "member-remove-ajax", nil, map[string]any{
		"role_id": role.GetID(), "user_id": user.GetID(),
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestMemberRemoveNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "member-remove-ajax", nil, map[string]any{
		"role_id": "r", "user_id": "u",
	})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestMemberRemoveMissingFields(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "member-remove-ajax", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
