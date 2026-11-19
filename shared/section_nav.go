package shared

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/userstore"
)

// SectionNav renders the useradmin section navigation as pills, matching the
// nav-pills card pattern used across the admin packages (shopadmin, taskadmin).
// The active tab is derived from the ?controller= query parameter. The Roles
// and Groups tabs are only shown when the user store has the corresponding
// tables enabled.
func SectionNav(r *http.Request, store userstore.StoreInterface) hb.TagInterface {
	linksHelper := NewLinksFromRequest(r)

	active := "users"
	controller := ""
	if r != nil && r.URL != nil {
		controller = r.URL.Query().Get("controller")
	}
	if strings.HasPrefix(controller, "role-") {
		active = "roles"
	}
	if strings.HasPrefix(controller, "group-") {
		active = "groups"
	}

	tab := func(label, key string, url string) hb.TagInterface {
		class := "nav-link"
		if key == active {
			class += " active"
		}
		return hb.LI().Class("nav-item").Child(
			hb.Hyperlink().HTML(label).Href(url).Class(class))
	}

	tabs := []hb.TagInterface{
		tab("Users", "users", linksHelper.UserManager(nil)),
	}
	if store != nil && store.GetRoleTableName() != "" && store.GetUserRoleTableName() != "" {
		tabs = append(tabs, tab("Roles", "roles", linksHelper.RoleManager(nil)))
	}
	if store != nil && store.GetGroupTableName() != "" && store.GetUserGroupTableName() != "" {
		tabs = append(tabs, tab("Groups", "groups", linksHelper.GroupManager(nil)))
	}

	ulNav := hb.UL().
		Class("nav nav-pills p-1 justify-content-center flex-grow-1").
		Children(tabs)

	label := hb.Span().
		Class("fw-bolder text-uppercase ms-3 me-3").
		Text("Users").
		Child(hb.I().Class("bi bi-chevron-right ms-2"))

	return hb.Div().
		Class("card card-default mt-3 mb-4 bg-primary-subtle").
		Child(hb.Div().
			Class("card-body d-flex align-items-center").
			Style("padding: 2px;").
			Child(label).
			Child(ulNav))
}
