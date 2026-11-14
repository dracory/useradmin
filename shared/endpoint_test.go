package shared_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestEndpointFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(contextWithEndpoint(req, "/admin"))
	if got := shared.Endpoint(req); got != "/admin" {
		t.Fatalf("expected /admin, got %q", got)
	}
}

func TestEndpointEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	if got := shared.Endpoint(req); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
