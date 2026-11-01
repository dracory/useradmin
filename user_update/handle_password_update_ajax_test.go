package user_update_test

import (
	"context"
	"testing"
)

// TestPasswordUpdateChangesPassword verifies that the password endpoint
// only changes the password hash and leaves all other fields untouched.
func TestPasswordUpdateChangesPassword(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "password-update-ajax", nil, map[string]any{
		"user_id":  user.GetID(),
		"password": "newpassword456",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, err := env.store.UserFindByID(context.Background(), user.GetID())
	if err != nil || saved == nil {
		t.Fatalf("UserFindByID: %v", err)
	}
	if !saved.PasswordCompare("newpassword456") {
		t.Error("expected new password to verify")
	}
	if saved.PasswordCompare("oldpassword123") {
		t.Error("expected old password to no longer verify")
	}
	if saved.GetEmail() != "target@test.com" || saved.GetFirstName() != "First" {
		t.Error("expected profile fields to remain unchanged")
	}
}

// TestPasswordUpdateRejectsShortPassword verifies the minimum length rule.
func TestPasswordUpdateRejectsShortPassword(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "password-update-ajax", nil, map[string]any{
		"user_id":  user.GetID(),
		"password": "short",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}
