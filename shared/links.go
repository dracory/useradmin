package shared

import "net/http"

// Links provides URL helpers for useradmin controllers.
// The base URL is read from request context (injected by Handle()),
// not hardcoded. This follows the blogadmin/shopadmin pattern.
type Links struct {
	baseURL string
}

// NewLinks creates a Links helper with the given base URL.
// If baseURL is empty, defaults to "/admin/users".
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/users"
	}
	return &Links{baseURL: baseURL}
}

// NewLinksFromRequest creates a Links helper using the user admin URL
// from the request context.
func NewLinksFromRequest(r *http.Request) *Links {
	return NewLinks(UserAdminURL(r))
}

// Home builds the URL for the user manager controller (default)
func (l *Links) Home(params map[string]string) string {
	return l.url(CONTROLLER_USER_MANAGER, params)
}

// UserManager builds the URL for the user manager controller
func (l *Links) UserManager(params map[string]string) string {
	return l.url(CONTROLLER_USER_MANAGER, params)
}

// UserCreate builds the URL for the user create controller
func (l *Links) UserCreate(params map[string]string) string {
	return l.url(CONTROLLER_USER_CREATE, params)
}

// UserDelete builds the URL for the user delete controller
func (l *Links) UserDelete(params map[string]string) string {
	return l.url(CONTROLLER_USER_DELETE, params)
}

// UserUpdate builds the URL for the user update controller
func (l *Links) UserUpdate(params map[string]string) string {
	return l.url(CONTROLLER_USER_UPDATE, params)
}

// UserImpersonate builds the URL for the user impersonate controller
func (l *Links) UserImpersonate(params map[string]string) string {
	return l.url(CONTROLLER_USER_IMPERSONATE, params)
}

// url builds a URL for the given controller. The params map is copied
// before mutation (does not modify caller's map).
func (l *Links) url(controller string, params map[string]string) string {
	return URL(l.baseURL, controller, params)
}
