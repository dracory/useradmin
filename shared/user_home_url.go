package shared

import "net/http"

// UserHomeURL returns the user home URL from request context.
// Used by the impersonate controller to redirect after a successful
// impersonation.
func UserHomeURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyUserHomeURL).(string)
	if !ok {
		return ""
	}
	return value
}
