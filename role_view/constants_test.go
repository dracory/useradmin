package role_view_test

import (
	"strings"
	"testing"

	"github.com/dracory/useradmin/role_view"
)

func TestFieldRoleIDConstant(t *testing.T) {
	if role_view.FieldRoleID != "role_id" {
		t.Fatalf("expected %q, got %q", "role_id", role_view.FieldRoleID)
	}
}

func TestRoleViewJSEscapes(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	role := env.seedRole(t, "Team")
	malformed := "');alert(1)//"

	rec := get(t, env, "role_id="+role.GetID()+"&extra="+malformed)
	body := rec.Body.String()
	if strings.Contains(body, "');alert(1)//") {
		t.Fatal("expected user-supplied params escaped in page script")
	}
}
