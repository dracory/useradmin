package role_update_test

import (
	"net/url"
	"testing"
)

// TestRoleFetchReturnsRole verifies the fetch endpoint returns the role
// fields needed to populate the form.
func TestRoleFetchReturnsRole(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	resp := postAction(t, env, "role-fetch-ajax", url.Values{"role_id": {role.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}
	if resp.Data["handle"] != "editor" {
		t.Errorf("expected handle in response, got %v", resp.Data["handle"])
	}
}

// TestRoleFetchRejectsMissingRole verifies an unknown role errors out.
func TestRoleFetchRejectsMissingRole(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "role-fetch-ajax", url.Values{"role_id": {"nonexistent"}}, map[string]any{})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}
