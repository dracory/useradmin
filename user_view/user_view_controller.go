package user_view

import (
	"net/http"

	"github.com/dracory/useradmin/shared"
)

// UiInterface defines the user view controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserView(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user view controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// rolesEnabled reports whether the user store has role tables
// configured (RolesEnabled).
func (u *ui) rolesEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetRoleTableName() != "" &&
		u.UserStore().GetUserRoleTableName() != ""
}

// groupsEnabled reports whether the user store has group tables
// configured (GroupsEnabled).
func (u *ui) groupsEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetGroupTableName() != "" &&
		u.UserStore().GetUserGroupTableName() != ""
}

// UserView handles the user view controller requests
func (u *ui) UserView(w http.ResponseWriter, r *http.Request) {
	html := u.renderPage(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}
