package shared_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestURLBuildsQuery(t *testing.T) {
	got := shared.URL("/admin", "user-manager", map[string]string{"user_id": "u1"})
	want := "/admin?controller=user-manager&user_id=u1"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestURLRUsesEndpointContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	req = req.WithContext(contextWithEndpoint(req, "/admin"))
	got := shared.URLR(req, "home", nil)
	if got != "/admin?controller=home" {
		t.Fatalf("got %q", got)
	}
}
