package group_manager_test

import (
	"context"
	"testing"
)

func TestDeleteGroupSuccess(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	g := env.seedGroup(t, "Doomed", "doomed", "active")

	resp := postAction(t, env, "delete-group-ajax", nil, map[string]any{
		"group_id": g.GetID(),
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	found, err := env.store.GroupFindByID(context.Background(), g.GetID())
	if err != nil || found != nil {
		t.Fatal("expected group deleted")
	}
}

func TestDeleteGroupMissingID(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postAction(t, env, "delete-group-ajax", nil, map[string]any{})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
