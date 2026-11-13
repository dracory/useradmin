package e2e

import (
	"strings"
	"testing"
)

// TestRoleMembersListsMembers verifies the role members page renders
// the users seeded into the role. Two members are seeded in setup so
// the second-row wait works.
func TestRoleMembersListsMembers(t *testing.T) {
	url := adminURL("role-members") + "&role_id=" + seedRoleID
	body := waitTable(t, newTab(t), url)
	if !strings.Contains(body, seedMemberMail) {
		t.Fatalf("expected member %q in role members page", seedMemberMail)
	}
}
