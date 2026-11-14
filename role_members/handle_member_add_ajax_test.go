package role_members_test

import "testing"

func TestMemberAddByID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")
	user := env.seedUser(t, "a@t.com")

	resp := postAction(t, env, "member-add-ajax", nil, map[string]any{
		"role_id": role.GetID(), "user_ref": user.GetID(),
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestMemberAddByEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")
	env.seedUser(t, "find@t.com")

	resp := postAction(t, env, "member-add-ajax", nil, map[string]any{
		"role_id": role.GetID(), "user_ref": "find@t.com",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestMemberAddMissingFields(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "member-add-ajax", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestMemberAddUserNotFound(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")

	resp := postAction(t, env, "member-add-ajax", nil, map[string]any{
		"role_id": role.GetID(), "user_ref": "ghost@t.com",
	})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
