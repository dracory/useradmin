package group_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/useradmin/shared"
)

// UiInterface defines the group manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	GroupManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new group manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionLoadGroups  = "load-groups-ajax"
	actionDeleteGroup = "delete-group-ajax"
	actionCreateGroup = "create-group-ajax"
)

// groupsEnabled reports whether the user store has group tables
// configured (GroupsEnabled).
func (u *ui) groupsEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetGroupTableName() != "" &&
		u.UserStore().GetUserGroupTableName() != ""
}

// GroupManager handles the group manager controller requests
func (u *ui) GroupManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the group manager controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadGroups:
		return u.handleGroupsFetchAjax(w, r)
	case actionDeleteGroup:
		return u.handleGroupDeleteAjax(w, r)
	case actionCreateGroup:
		return u.handleGroupCreateAjax(w, r)
	default:
		return u.renderPage(w, r)
	}
}
