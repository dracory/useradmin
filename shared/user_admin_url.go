package shared

import "net/http"

// UserAdminURL returns the user admin base URL from request context
func UserAdminURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyUserAdminURL).(string)
	if !ok {
		return ""
	}
	return value
}
