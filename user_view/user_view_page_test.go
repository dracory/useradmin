package user_view_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestUserViewRendersUserDetails(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	user := env.seedUser(t)

	rec := get(t, env, "/admin/users?controller=user-view&user_id="+user.GetID())

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		"View User", "User Details",
		"target@test.com", "First", "Last", "Biz", "123",
		"GB", "Europe/London", "some memo",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %q", want)
		}
	}
	// The Edit button must point at the user update controller.
	if !strings.Contains(body, "controller=user-update&amp;user_id="+user.GetID()) &&
		!strings.Contains(body, "controller=user-update&user_id="+user.GetID()) {
		t.Error("expected Edit link to user-update controller")
	}
}

func TestUserViewRendersRolesAndGroups(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	user := env.seedUser(t)
	role := env.seedRole(t)
	group := env.seedGroup(t)

	if _, err := env.store.UserRoleFindByUserIDAndRoleIDOrCreate(context.Background(), user.GetID(), role.GetID()); err != nil {
		t.Fatalf("assign role: %v", err)
	}
	if _, err := env.store.UserGroupFindByUserIDAndGroupIDOrCreate(context.Background(), user.GetID(), group.GetID()); err != nil {
		t.Fatalf("assign group: %v", err)
	}

	rec := get(t, env, "/admin/users?controller=user-view&user_id="+user.GetID())

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{"Roles", "Guardian", "guardian", "Groups", "Members", "members"} {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %q", want)
		}
	}
	if !strings.Contains(body, "controller=role-view") {
		t.Error("expected role name to link to role-view controller")
	}
	if !strings.Contains(body, "controller=group-view") {
		t.Error("expected group name to link to group-view controller")
	}
}

func TestUserViewHidesCardsWhenTablesNotConfigured(t *testing.T) {
	env := newTestEnvWithOptions(t, false, false)
	defer env.close()

	user := env.seedUser(t)

	rec := get(t, env, "/admin/users?controller=user-view&user_id="+user.GetID())

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "User Details") {
		t.Error("expected User Details card to render")
	}
	if strings.Contains(body, ">Roles<") {
		t.Error("expected Roles card to be hidden when role tables are not configured")
	}
	if strings.Contains(body, ">Groups<") {
		t.Error("expected Groups card to be hidden when group tables are not configured")
	}
}

func TestUserViewMissingUserIDRedirects(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "/admin/users?controller=user-view")

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.Contains(loc, "controller=user-manager") {
		t.Errorf("expected redirect to user-manager, got %q", loc)
	}
}

func TestUserViewUnknownUserRedirects(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	rec := get(t, env, "/admin/users?controller=user-view&user_id=nonexistent")

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.Contains(loc, "controller=user-manager") {
		t.Errorf("expected redirect to user-manager, got %q", loc)
	}
}
