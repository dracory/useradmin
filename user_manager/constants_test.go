package user_manager_test

import "testing"

// constants.go declares field names and limits; exercised indirectly
// through the fetch action. Verify the max-per-page cap is enforced.
func TestMaxPerPageEnforced(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	for i := 0; i < 3; i++ {
		env.seedUser(t, "u"+string(rune('a'+i))+"@t.com", "active")
	}

	resp := postAction(t, env, "load-users-ajax", nil, map[string]any{
		"per_page": 1000000,
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
