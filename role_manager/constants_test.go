package role_manager_test

import (
	"testing"

	"github.com/dracory/useradmin/role_manager"
)

// constants.go only declares field names; verify they are non-empty.
func TestFieldConstants(t *testing.T) {
	if role_manager.FieldID == "" || role_manager.FieldName == "" || role_manager.FieldHandle == "" {
		t.Fatal("expected non-empty field constants")
	}
}
