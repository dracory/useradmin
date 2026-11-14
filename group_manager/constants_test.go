package group_manager_test

import (
	"testing"

	"github.com/dracory/useradmin/group_manager"
)

// constants.go only declares field names; verify they are non-empty.
func TestFieldConstants(t *testing.T) {
	if group_manager.FieldID == "" || group_manager.FieldName == "" || group_manager.FieldHandle == "" {
		t.Fatal("expected non-empty field constants")
	}
}
