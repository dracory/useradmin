package shared_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestUserHomeURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), shared.KeyUserHomeURL, "/"))
	if got := shared.UserHomeURL(req); got != "/" {
		t.Fatalf("got %q", got)
	}
}
