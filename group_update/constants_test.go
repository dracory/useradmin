package group_update_test

import (
	"testing"

	"github.com/dracory/useradmin/group_update"
)

func TestFieldConstants(t *testing.T) {
	if group_update.FieldID != "id" || group_update.FieldGroupID != "group_id" ||
		group_update.FieldEmail != "email" {
		t.Fatal("unexpected field constants")
	}
}
