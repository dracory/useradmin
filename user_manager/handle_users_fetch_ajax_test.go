package user_manager_test

import "testing"

func TestLoadUsersSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "a@t.com", "active")
	env.seedUser(t, "b@t.com", "unverified")

	resp := postAction(t, env, "load-users-ajax", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	if total, _ := resp.data()["total"].(float64); total != 2 {
		t.Fatalf("expected total 2, got %v", resp.data()["total"])
	}
}

func TestLoadUsersStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "a@t.com", "active")
	env.seedUser(t, "b@t.com", "unverified")

	resp := postAction(t, env, "load-users-ajax", nil, map[string]any{"status": "active"})
	users, _ := resp.data()["users"].([]any)
	if len(users) != 1 {
		t.Fatalf("expected 1 active user, got %d", len(users))
	}
}

func TestLoadUsersEmailFilter(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "alice@t.com", "active")
	env.seedUser(t, "bob@t.com", "active")

	resp := postAction(t, env, "load-users-ajax", nil, map[string]any{"email": "alice"})
	users, _ := resp.data()["users"].([]any)
	if len(users) != 1 {
		t.Fatalf("expected 1 matching user, got %d", len(users))
	}
}

func TestLoadUsersInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postRaw(t, env, "load-users-ajax", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
