package shared_test

import (
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestConsts(t *testing.T) {
	if shared.KeyEndpoint != "endpoint" ||
		shared.KeyAdminHomeURL != "admin_home_url" ||
		shared.KeyUserAdminURL != "user_admin_url" ||
		shared.KeyUserHomeURL != "user_home_url" {
		t.Fatal("unexpected context keys")
	}
	if shared.CatchAll != "/*" {
		t.Fatal("unexpected CatchAll")
	}
	if shared.ERROR_USER_STORE_IS_NIL == "" || shared.ERROR_LOGGER_IS_NIL == "" {
		t.Fatal("expected error constants")
	}
}
