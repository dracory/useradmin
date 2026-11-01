package shared

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
)

// SectionNav renders the useradmin section navigation as tabs, matching the
// nav-tabs pattern used across the admin packages. The active tab is derived
// from the ?controller= query parameter.
func SectionNav(r *http.Request) hb.TagInterface {
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

	return hb.UL().Class("nav nav-tabs mb-3").Children([]hb.TagInterface{
		tab("Users", "users", linksHelper.UserManager(nil)),
		tab("Roles", "roles", linksHelper.RoleManager(nil)),
		tab("Groups", "groups", linksHelper.GroupManager(nil)),
	})
}
