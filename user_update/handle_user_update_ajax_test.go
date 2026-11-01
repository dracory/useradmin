package user_update_test

import (
	"context"
	"testing"

	"github.com/dracory/userstore"
)

// TestUserUpdateLegacyEndpointStillWorks verifies the legacy full-payload
// endpoint still saves a user — it is kept for backward compatibility.
func TestUserUpdateLegacyEndpointStillWorks(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	user := env.seedUser(t)

	resp := postAction(t, env, "user-update-ajax", nil, map[string]any{
		"user_id":    user.GetID(),
		"status":     userstore.USER_STATUS_ACTIVE,
		"email":      "legacy@test.com",
		"first_name": "Legacy",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.UserFindByID(context.Background(), user.GetID())
	if saved.GetEmail() != "legacy@test.com" || saved.GetFirstName() != "Legacy" {
		t.Error("expected user updated via legacy endpoint")
	}
}
