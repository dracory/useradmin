package user_update_test

import (
	"context"
	"testing"

	"github.com/dracory/userstore"
)

// TestStatusUpdateOnlyChangesStatus verifies the status endpoint touches
// only the status field.
func TestStatusUpdateOnlyChangesStatus(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"user_id": user.GetID(),
		"status":  userstore.USER_STATUS_INACTIVE,
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetStatus() != userstore.USER_STATUS_INACTIVE {
		t.Errorf("expected status inactive, got %q", saved.GetStatus())
	}
	if saved.GetEmail() != "target@test.com" || saved.GetCountry() != "GB" || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestStatusUpdateRejectsEmpty verifies status is required.
func TestStatusUpdateRejectsEmpty(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "status-update-ajax", nil, map[string]any{
		"user_id": user.GetID(),
		"status":  "",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestDetailsUpdateOnlyChangesDetails verifies the details endpoint
// touches only profile fields.
func TestDetailsUpdateOnlyChangesDetails(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"user_id":       user.GetID(),
		"first_name":    "New",
		"last_name":     "Name",
		"email":         "new@test.com",
		"business_name": "NewBiz",
		"phone":         "999",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetFirstName() != "New" || saved.GetEmail() != "new@test.com" {
		t.Error("expected details updated")
	}
	if saved.GetStatus() != userstore.USER_STATUS_ACTIVE || saved.GetCountry() != "GB" || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestDetailsUpdateRejectsBadEmail verifies email validation still applies.
func TestDetailsUpdateRejectsBadEmail(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "details-update-ajax", nil, map[string]any{
		"user_id": user.GetID(),
		"email":   "not-an-email",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

// TestRolesUpdateOnlyChangesRoles verifies the roles endpoint touches only
// role-related fields.
func TestRolesUpdateOnlyChangesRoles(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "roles-update-ajax", nil, map[string]any{
		"user_id":   user.GetID(),
		"role":      userstore.USER_ROLE_ADMINISTRATOR,
		"role_ids":  []string{},
		"group_ids": []string{},
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetRole() != userstore.USER_ROLE_ADMINISTRATOR {
		t.Errorf("expected role administrator, got %q", saved.GetRole())
	}
	if saved.GetEmail() != "target@test.com" || saved.GetStatus() != userstore.USER_STATUS_ACTIVE {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestLocationUpdateOnlyChangesLocation verifies the location endpoint
// touches only country and timezone.
func TestLocationUpdateOnlyChangesLocation(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "location-update-ajax", nil, map[string]any{
		"user_id":  user.GetID(),
		"country":  "US",
		"timezone": "America/New_York",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetCountry() != "US" || saved.GetTimezone() != "America/New_York" {
		t.Error("expected location updated")
	}
	if saved.GetEmail() != "target@test.com" || saved.GetMemo() != "old memo" {
		t.Error("expected other fields to remain unchanged")
	}
}

// TestNotesUpdateOnlyChangesMemo verifies the notes endpoint touches only
// the memo field.
func TestNotesUpdateOnlyChangesMemo(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "notes-update-ajax", nil, map[string]any{
		"user_id": user.GetID(),
		"memo":    "new memo",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetMemo() != "new memo" {
		t.Errorf("expected memo updated, got %q", saved.GetMemo())
	}
	if saved.GetEmail() != "target@test.com" || saved.GetStatus() != userstore.USER_STATUS_ACTIVE {
		t.Error("expected other fields to remain unchanged")
	}
}
