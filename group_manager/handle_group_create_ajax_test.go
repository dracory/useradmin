package group_manager_test

import "testing"

func TestCreateGroupSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "create-group-ajax", nil, map[string]any{
		"name": "Editors", "handle": "editors",
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestCreateGroupMissingFields(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	for _, payload := range []map[string]any{
		{"name": "", "handle": "x"},
		{"name": "x", "handle": ""},
	} {
		resp := postAction(t, env, "create-group-ajax", nil, payload)
		if resp.Status != "error" {
			t.Fatalf("expected error for %+v, got %+v", payload, resp)
		}
	}
}

func TestCreateGroupDuplicateHandle(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedGroup(t, "Editors", "editors", "active")

	resp := postAction(t, env, "create-group-ajax", nil, map[string]any{
		"name": "Editors 2", "handle": "editors",
	})
	if resp.Status != "error" {
		t.Fatalf("expected duplicate-handle error, got %+v", resp)
	}
}

func TestCreateGroupInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postRaw(t, env, "create-group-ajax", "{oops")
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
