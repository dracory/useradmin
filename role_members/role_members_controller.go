package role_members

import (
	"net/http"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/req"
)

// UiInterface defines the role members controller's UI interface
type UiInterface interface {
	shared.UiInterface
	RoleMembers(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new role members controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionMembersFetch = "members-fetch-ajax"
	actionMemberAdd    = "member-add-ajax"
	actionMemberRemove = "member-remove-ajax"
	actionUsersSearch  = "users-search-ajax"
)

// rolesEnabled reports whether the user store has role tables
// configured (RolesEnabled).
func (u *ui) rolesEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetRoleTableName() != "" &&
		u.UserStore().GetUserRoleTableName() != ""
}

// RoleMembers handles the role members controller requests
func (u *ui) RoleMembers(w http.ResponseWriter, r *http.Request) {
	u.Handler(w, r)
}

// Handler processes the role members controller request
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionMembersFetch:
		u.handleMembersFetchAjax(w, r)
	case actionMemberAdd:
		u.handleMemberAddAjax(w, r)
	case actionMemberRemove:
		u.handleMemberRemoveAjax(w, r)
	case actionUsersSearch:
		u.handleUsersSearchAjax(w, r)
	default:
		html := u.renderPage(w, r)
		if html != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(html))
		}
	}
}
