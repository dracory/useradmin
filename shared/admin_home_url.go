package shared

import "net/http"

// AdminHomeURL returns the admin home URL from request context
func AdminHomeURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyAdminHomeURL).(string)
	if !ok {
		return ""
	}
	return value
}
