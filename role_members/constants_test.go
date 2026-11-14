package role_members_test

import (
	"testing"

	"github.com/dracory/useradmin/role_members"
)

func TestFieldConstants(t *testing.T) {
	if role_members.FieldRoleID != "role_id" || role_members.FieldUserID != "user_id" {
		t.Fatal("unexpected field constants")
	}
}
