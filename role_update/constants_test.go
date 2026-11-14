package role_update_test

import (
	"testing"

	"github.com/dracory/useradmin/role_update"
)

func TestFieldConstants(t *testing.T) {
	if role_update.FieldID != "id" || role_update.FieldRoleID != "role_id" ||
		role_update.FieldEmail != "email" {
		t.Fatal("unexpected field constants")
	}
}
