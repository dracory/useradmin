package group_members_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/dracory/userstore"
)

func TestMembersFetchSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	g := env.seedGroup(t, "Editors", "editors")
	user := env.seedUser(t, "m@t.com")

	ug := userstore.NewUserGroup()
	ug.SetUserID(user.GetID())
	ug.SetGroupID(g.GetID())
	if err := env.store.UserGroupCreate(context.Background(), ug); err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	resp := postAction(t, env, "members-fetch-ajax",
		url.Values{"group_id": {g.GetID()}}, nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	members, _ := resp.data()["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
}

func TestMembersFetchMissingGroupID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "members-fetch-ajax", nil, nil)
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
