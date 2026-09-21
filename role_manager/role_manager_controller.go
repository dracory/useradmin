package role_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/useradmin/shared"
)

// UiInterface defines the role manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	RoleManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new role manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionLoadRoles  = "load-roles-ajax"
	actionDeleteRole = "delete-role-ajax"
	actionCreateRole = "create-role-ajax"
)

// rolesEnabled reports whether the user store has role tables
// configured (RolesEnabled).
func (u *ui) rolesEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetRoleTableName() != "" &&
		u.UserStore().GetUserRoleTableName() != ""
}

// RoleManager handles the role manager controller requests
func (u *ui) RoleManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the role manager controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadRoles:
		return u.handleRolesFetchAjax(w, r)
	case actionDeleteRole:
		return u.handleRoleDeleteAjax(w, r)
	case actionCreateRole:
		return u.handleRoleCreateAjax(w, r)
	default:
		return u.renderPage(w, r)
	}
}
