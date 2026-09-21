package group_view

import (
	"net/http"

	"github.com/dracory/useradmin/shared"
)

// UiInterface defines the group view controller's UI interface
type UiInterface interface {
	shared.UiInterface
	GroupView(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new group view controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// groupsEnabled reports whether the user store has group tables
// configured (GroupsEnabled).
func (u *ui) groupsEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetGroupTableName() != "" &&
		u.UserStore().GetUserGroupTableName() != ""
}

// GroupView handles the group view controller requests
func (u *ui) GroupView(w http.ResponseWriter, r *http.Request) {
	html := u.renderPage(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}
