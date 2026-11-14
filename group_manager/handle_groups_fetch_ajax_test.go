package group_manager_test

import "testing"

func TestLoadGroupsSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedGroup(t, "Admins", "admins", "active")
	env.seedGroup(t, "Members", "members", "active")

	resp := postAction(t, env, "load-groups-ajax", nil, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	groups, _ := resp.Data["groups"].([]any)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestLoadGroupsEmpty(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "load-groups-ajax", nil, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
