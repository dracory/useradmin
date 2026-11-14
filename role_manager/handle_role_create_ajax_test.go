package role_manager_test

import "testing"

func TestCreateRoleSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "create-role-ajax", nil, map[string]any{
		"name": "Editors", "handle": "editors",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestCreateRoleMissingFields(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	for _, payload := range []map[string]any{
		{"name": "", "handle": "x"},
		{"name": "x", "handle": ""},
	} {
		resp := postAction(t, env, "create-role-ajax", nil, payload)
		if resp.status() != "error" {
			t.Fatalf("expected error for %+v, got %+v", payload, resp)
		}
	}
}

func TestCreateRoleDuplicateHandle(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedRole(t, "Editors", "editors", "active")

	resp := postAction(t, env, "create-role-ajax", nil, map[string]any{
		"name": "Editors 2", "handle": "editors",
	})
	if resp.status() != "error" {
		t.Fatalf("expected duplicate-handle error, got %+v", resp)
	}
}

func TestCreateRoleInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postRaw(t, env, "create-role-ajax", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
