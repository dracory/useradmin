package user_update_test

import (
	"net/url"
	"testing"
)

// TestUserFetchReturnsUser verifies the fetch endpoint returns the user
// fields needed to populate the form.
func TestUserFetchReturnsUser(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "user-fetch-ajax", url.Values{"user_id": {user.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}
	if resp.Data["email"] != "target@test.com" {
		t.Errorf("expected email in response, got %v", resp.Data["email"])
	}
}

// TestUserFetchRejectsMissingUser verifies an unknown user errors out.
func TestUserFetchRejectsMissingUser(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	resp := postAction(t, env, "user-fetch-ajax", nil, map[string]any{
		"user_id": "nonexistent",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}
