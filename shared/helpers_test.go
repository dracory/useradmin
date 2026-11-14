package shared_test

import (
	"context"
	"net/http"

	"github.com/dracory/useradmin/shared"
)

// contextWithEndpoint returns a request context carrying the
// endpoint key used by URLR.
func contextWithEndpoint(r *http.Request, endpoint string) context.Context {
	return context.WithValue(r.Context(), shared.KeyEndpoint, endpoint)
}
