package group_update_test

import (
	"context"
	"testing"

	"github.com/dracory/userstore"
)

// TestStatusUpdateOnlyChangesStatus verifies the status endpoint touches
// only the status field.
func TestStatusUpdateOnlyChangesStatus(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"status":   userstore.GROUP_STATUS_INACTIVE,
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.GroupFindByID(context.Background(), group.GetID())
	if saved.GetStatus() != userstore.GROUP_STATUS_INACTIVE {
		t.Errorf("expected status inactive, got %q", saved.GetStatus())
	}
	if saved.GetName() != "Members" || saved.GetHandle() != "members" || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestStatusUpdateRejectsEmpty verifies status is required.
func TestStatusUpdateRejectsEmpty(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"status":   "",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestStatusUpdateRejectsInvalid verifies invalid status values are rejected.
func TestStatusUpdateRejectsInvalid(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"status":   "bogus",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestDetailsUpdateOnlyChangesDetails verifies the details endpoint
// touches only name and handle.
func TestDetailsUpdateOnlyChangesDetails(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"name":     "Subscribers",
		"handle":   "subscribers",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.GroupFindByID(context.Background(), group.GetID())
	if saved.GetName() != "Subscribers" || saved.GetHandle() != "subscribers" {
		t.Error("expected details updated")
	}
	if saved.GetStatus() != userstore.GROUP_STATUS_ACTIVE || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestDetailsUpdateRejectsEmptyName verifies name is required.
func TestDetailsUpdateRejectsEmptyName(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"name":     "",
		"handle":   "members",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestDetailsUpdateRejectsDuplicateHandle verifies a handle used by
// another group is rejected.
func TestDetailsUpdateRejectsDuplicateHandle(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	other := userstore.NewGroup().
		SetName("Editors").
		SetHandle("editors").
		SetStatus(userstore.GROUP_STATUS_ACTIVE)
	if err := env.store.GroupCreate(context.Background(), other); err != nil {
		t.Fatalf("seed second group: %v", err)
	}

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"name":     "Members",
		"handle":   "editors",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestNotesUpdateOnlyChangesMemo verifies the notes endpoint touches only
// the memo field.
func TestNotesUpdateOnlyChangesMemo(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	group := env.seedGroup(t)

	resp := postAction(t, env, "notes-update-ajax", nil, map[string]any{
		"group_id": group.GetID(),
		"memo":     "new memo",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.GroupFindByID(context.Background(), group.GetID())
	if saved.GetMemo() != "new memo" {
		t.Errorf("expected memo updated, got %q", saved.GetMemo())
	}
	if saved.GetName() != "Members" || saved.GetStatus() != userstore.GROUP_STATUS_ACTIVE {
		t.Error("expected other fields to remain unchanged")
	}
}
