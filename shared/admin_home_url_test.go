package shared_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestAdminHomeURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), shared.KeyAdminHomeURL, "/admin"))
	if got := shared.AdminHomeURL(req); got != "/admin" {
		t.Fatalf("got %q", got)
	}
}
