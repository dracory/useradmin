package group_members_test

import (
	"context"
	"testing"

	"github.com/dracory/userstore"
)

func TestMemberRemoveSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	g := env.seedGroup(t, "Editors", "editors")
	user := env.seedUser(t, "m@t.com")

	ug := userstore.NewUserGroup()
	ug.SetUserID(user.GetID())
	ug.SetGroupID(g.GetID())
	if err := env.store.UserGroupCreate(context.Background(), ug); err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	resp := postAction(t, env, "member-remove-ajax", nil, map[string]any{
		"group_id": g.GetID(), "user_id": user.GetID(),
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestMemberRemoveNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "member-remove-ajax", nil, map[string]any{
		"group_id": "g", "user_id": "u",
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
