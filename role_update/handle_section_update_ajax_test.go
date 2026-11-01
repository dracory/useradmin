package role_update_test

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

	role := env.seedRole(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"status":  userstore.ROLE_STATUS_INACTIVE,
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.RoleFindByID(context.Background(), role.GetID())
	if saved.GetStatus() != userstore.ROLE_STATUS_INACTIVE {
		t.Errorf("expected status inactive, got %q", saved.GetStatus())
	}
	if saved.GetName() != "Editor" || saved.GetHandle() != "editor" || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestStatusUpdateRejectsEmpty verifies status is required.
func TestStatusUpdateRejectsEmpty(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"status":  "",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestStatusUpdateRejectsInvalid verifies invalid status values are rejected.
func TestStatusUpdateRejectsInvalid(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"status":  "bogus",
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

	role := env.seedRole(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"name":    "Moderator",
		"handle":  "moderator",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.RoleFindByID(context.Background(), role.GetID())
	if saved.GetName() != "Moderator" || saved.GetHandle() != "moderator" {
		t.Error("expected details updated")
	}
	if saved.GetStatus() != userstore.ROLE_STATUS_ACTIVE || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestDetailsUpdateRejectsEmptyName verifies name is required.
func TestDetailsUpdateRejectsEmptyName(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"name":    "",
		"handle":  "editor",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestDetailsUpdateRejectsDuplicateHandle verifies a handle used by
// another role is rejected.
func TestDetailsUpdateRejectsDuplicateHandle(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	role := env.seedRole(t)

	other := userstore.NewRole().
		SetName("Admin").
		SetHandle("admin").
		SetStatus(userstore.ROLE_STATUS_ACTIVE)
	if err := env.store.RoleCreate(context.Background(), other); err != nil {
		t.Fatalf("seed second role: %v", err)
	}

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"name":    "Editor",
		"handle":  "admin",
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

	role := env.seedRole(t)

	resp := postAction(t, env, "notes-update-ajax", nil, map[string]any{
		"role_id": role.GetID(),
		"memo":    "new memo",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.RoleFindByID(context.Background(), role.GetID())
	if saved.GetMemo() != "new memo" {
		t.Errorf("expected memo updated, got %q", saved.GetMemo())
	}
	if saved.GetName() != "Editor" || saved.GetStatus() != userstore.ROLE_STATUS_ACTIVE {
		t.Error("expected other fields to remain unchanged")
	}
}
