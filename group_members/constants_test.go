package group_members_test

import (
	"testing"

	"github.com/dracory/useradmin/group_members"
)

func TestFieldConstants(t *testing.T) {
	if group_members.FieldGroupID != "group_id" || group_members.FieldUserID != "user_id" {
		t.Fatal("unexpected field constants")
	}
}
