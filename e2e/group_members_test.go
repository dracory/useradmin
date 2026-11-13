package e2e

import (
	"strings"
	"testing"
)

// TestGroupMembersListsMembers verifies the group members page renders
// the users seeded into the group. Two members are seeded in setup so
// the second-row wait works.
func TestGroupMembersListsMembers(t *testing.T) {
	url := adminURL("group-members") + "&group_id=" + seedGroupID
	body := waitTable(t, newTab(t), url)
	if !strings.Contains(body, seedMemberMail) {
		t.Fatalf("expected member %q in group members page", seedMemberMail)
	}
}
