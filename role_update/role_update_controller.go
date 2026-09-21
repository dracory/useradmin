package role_update

import (
	"net/http"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/req"
)

// UiInterface defines the role update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	RoleUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new role update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionRoleFetch  = "role-fetch-ajax"
	actionRoleUpdate = "role-update-ajax"
)

// rolesEnabled reports whether the user store has role tables
// configured (RolesEnabled).
func (u *ui) rolesEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetRoleTableName() != "" &&
		u.UserStore().GetUserRoleTableName() != ""
}

// RoleUpdate handles the role update controller requests
func (u *ui) RoleUpdate(w http.ResponseWriter, r *http.Request) {
	u.Handler(w, r)
}

// Handler processes the role update controller request
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionRoleFetch:
		u.handleRoleFetchAjax(w, r)
	case actionRoleUpdate:
		u.handleRoleUpdateAjax(w, r)
	default:
		html := u.renderPage(w, r)
		if html != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(html))
		}
	}
}
