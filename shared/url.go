package shared

import (
	"net/http"
	"net/url"
)

// URL builds a URL for the given endpoint and controller with params.
// The controller is placed in the params map under the "controller" key.
// The params map is copied before mutation (does not modify caller's map).
func URL(endpoint string, controller string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	// Copy the map to avoid mutating the caller's map
	p := make(map[string]string, len(params)+1)
	for k, v := range params {
		p[k] = v
	}
	p["controller"] = controller
	return endpoint + query(p)
}

// URLR builds a URL using the endpoint from the request context.
func URLR(r *http.Request, controller string, params map[string]string) string {
	return URL(Endpoint(r), controller, params)
}

func query(queryData map[string]string) string {
	if len(queryData) == 0 {
		return ""
	}
	v := url.Values{}
	for key, value := range queryData {
		v.Set(key, value)
	}
	return "?" + v.Encode()
}
