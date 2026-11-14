package shared_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestUserAdminURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), shared.KeyUserAdminURL, "/admin/users"))
	if got := shared.UserAdminURL(req); got != "/admin/users" {
		t.Fatalf("got %q", got)
	}
}
