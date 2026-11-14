package user_view_test

import (
	"testing"

	"github.com/dracory/useradmin/user_view"
)

func TestFieldUserIDConstant(t *testing.T) {
	if user_view.FieldUserID != "user_id" {
		t.Fatalf("expected %q, got %q", "user_id", user_view.FieldUserID)
	}
}
