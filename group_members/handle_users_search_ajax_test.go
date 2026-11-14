package group_members_test

import (
	"net/url"
	"testing"
)

func TestUsersSearch(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "alice@t.com")
	env.seedUser(t, "bob@t.com")

	resp := postAction(t, env, "users-search-ajax",
		url.Values{"q": {"alice"}}, nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	users, _ := resp.data()["users"].([]any)
	if len(users) != 1 {
		t.Fatalf("expected 1 matching user, got %d", len(users))
	}
}

func TestUsersSearchEmptyQuery(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	env.seedUser(t, "alice@t.com")

	resp := postAction(t, env, "users-search-ajax", nil, nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
