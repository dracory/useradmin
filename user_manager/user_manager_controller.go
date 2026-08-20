package user_manager

import (
	"net/http"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/req"
)

// UiInterface defines the user manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionLoadUsers  = "load-users-ajax"
	actionDeleteUser = "delete-user-ajax"
	actionCreateUser = "create-user-ajax"
)

// UserManager handles the user manager controller requests
func (u *ui) UserManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the user manager controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadUsers:
		return u.handleUsersFetchAjax(w, r)
	case actionDeleteUser:
		return u.handleUserDeleteAjax(w, r)
	case actionCreateUser:
		return u.handleUserCreateAjax(w, r)
	default:
		return u.renderPage(w, r)
	}
}
