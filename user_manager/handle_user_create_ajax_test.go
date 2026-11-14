package user_manager_test

import "testing"

func TestCreateUserSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "create-user-ajax", nil, map[string]any{
		"first_name": "New", "last_name": "User", "email": "new@t.com",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestCreateUserMissingEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "create-user-ajax", nil, map[string]any{
		"first_name": "New",
	})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestCreateUserInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postRaw(t, env, "create-user-ajax", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
