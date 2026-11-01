package group_update_test

import (
	"net/url"
	"testing"
)

// TestGroupFetchReturnsGroup verifies the fetch endpoint returns the group
// fields needed to populate the form.
func TestGroupFetchReturnsGroup(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "group-fetch-ajax", url.Values{"group_id": {group.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}
	if resp.Data["handle"] != "members" {
		t.Errorf("expected handle in response, got %v", resp.Data["handle"])
	}
}

// TestGroupFetchRejectsMissingGroup verifies an unknown group errors out.
func TestGroupFetchRejectsMissingGroup(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "group-fetch-ajax", url.Values{"group_id": {"nonexistent"}}, map[string]any{})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}
