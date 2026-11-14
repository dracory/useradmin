package group_view_test

import (
	"strings"
	"testing"

	"github.com/dracory/useradmin/group_view"
)

func TestFieldGroupIDConstant(t *testing.T) {
	if group_view.FieldGroupID != "group_id" {
		t.Fatalf("expected %q, got %q", "group_id", group_view.FieldGroupID)
	}
}

func TestGroupViewJSEscapes(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	group := env.seedGroup(t, "Team")
	malformed := "');alert(1)//"

	rec := get(t, env, "group_id="+group.GetID()+"&extra="+malformed)
	body := rec.Body.String()
	if strings.Contains(body, "');alert(1)//") {
		t.Fatal("expected user-supplied params escaped in page script")
	}
}
