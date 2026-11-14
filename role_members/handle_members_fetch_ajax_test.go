package role_members_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/dracory/userstore"
)

func TestMembersFetchSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Editors", "editors")
	user := env.seedUser(t, "m@t.com")

	ur := userstore.NewUserRole()
	ur.SetUserID(user.GetID())
	ur.SetRoleID(role.GetID())
	if err := env.store.UserRoleCreate(context.Background(), ur); err != nil {
		t.Fatalf("seed membership: %v", err)
	}

	resp := postAction(t, env, "members-fetch-ajax",
		url.Values{"role_id": {role.GetID()}}, nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	members, _ := resp.data()["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
}

func TestMembersFetchMissingRoleID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "members-fetch-ajax", nil, nil)
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
