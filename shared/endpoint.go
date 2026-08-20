package shared

import "net/http"

// Endpoint returns the request endpoint URL from context
func Endpoint(r *http.Request) string {
	value, ok := r.Context().Value(KeyEndpoint).(string)
	if !ok {
		return ""
	}
	return value
}
