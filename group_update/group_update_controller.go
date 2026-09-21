package group_update

import (
	"net/http"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/req"
)

// UiInterface defines the group update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	GroupUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new group update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

const (
	actionGroupFetch  = "group-fetch-ajax"
	actionGroupUpdate = "group-update-ajax"
)

// groupsEnabled reports whether the user store has group tables
// configured (GroupsEnabled).
func (u *ui) groupsEnabled() bool {
	return u.UserStore() != nil &&
		u.UserStore().GetGroupTableName() != "" &&
		u.UserStore().GetUserGroupTableName() != ""
}

// GroupUpdate handles the group update controller requests
func (u *ui) GroupUpdate(w http.ResponseWriter, r *http.Request) {
	u.Handler(w, r)
}

// Handler processes the group update controller request
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionGroupFetch:
		u.handleGroupFetchAjax(w, r)
	case actionGroupUpdate:
		u.handleGroupUpdateAjax(w, r)
	default:
		html := u.renderPage(w, r)
		if html != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(html))
		}
	}
}
