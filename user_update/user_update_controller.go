package user_update

import (
	"net/http"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/req"
)

// UiInterface defines the user update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionUserFetch      = "user-fetch-ajax"
	actionGetTimezones   = "get-timezones-ajax"
	actionUserUpdate     = "user-update-ajax"
	actionPasswordUpdate = "password-update-ajax"
	actionStatusUpdate   = "status-update-ajax"
	actionDetailsUpdate  = "details-update-ajax"
	actionRolesUpdate    = "roles-update-ajax"
	actionLocationUpdate = "location-update-ajax"
	actionNotesUpdate    = "notes-update-ajax"
)

// UserUpdate handles the user update controller requests
func (u *ui) UserUpdate(w http.ResponseWriter, r *http.Request) {
	u.Handler(w, r)
}

// Handler processes the user update controller request
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionUserFetch:
		u.handleUserFetchAjax(w, r)
	case actionGetTimezones:
		u.handleTimezonesFetchAjax(w, r)
	case actionUserUpdate:
		u.handleUserUpdateAjax(w, r)
	case actionPasswordUpdate:
		u.handlePasswordUpdateAjax(w, r)
	case actionStatusUpdate:
		u.handleStatusUpdateAjax(w, r)
	case actionDetailsUpdate:
		u.handleDetailsUpdateAjax(w, r)
	case actionRolesUpdate:
		u.handleRolesUpdateAjax(w, r)
	case actionLocationUpdate:
		u.handleLocationUpdateAjax(w, r)
	case actionNotesUpdate:
		u.handleNotesUpdateAjax(w, r)
	default:
		html := u.renderPage(w, r)
		if html != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(html))
		}
	}
}
