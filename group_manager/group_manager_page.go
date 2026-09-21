package group_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/hb"
)

var (
	//go:embed groups.html
	groupsHTML string

	//go:embed groups.js
	groupsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.UserStore() == nil {
		return shared.ErrorAlert("User store is not configured")
	}
	if !u.groupsEnabled() {
		return shared.ErrorAlert("Groups are not enabled in the user store")
	}

	linksHelper := shared.NewLinksFromRequest(r)

	urlGroupsLoad := shared.JSEscapeString(linksHelper.GroupManager(map[string]string{"action": actionLoadGroups}))
	urlGroupDelete := shared.JSEscapeString(linksHelper.GroupManager(map[string]string{"action": actionDeleteGroup}))
	urlGroupCreate := shared.JSEscapeString(linksHelper.GroupManager(map[string]string{"action": actionCreateGroup}))
	urlGroupUpdate := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"group_id": "GROUP_ID_PLACEHOLDER"}))
	urlGroupView := shared.JSEscapeString(linksHelper.GroupView(map[string]string{"group_id": "GROUP_ID_PLACEHOLDER"}))

	html := strings.ReplaceAll(groupsHTML, "urlGroupUpdate", "'"+urlGroupUpdate+"'")
	html = strings.ReplaceAll(html, "urlGroupView", "'"+urlGroupView+"'")
	js := strings.ReplaceAll(groupsJS, "urlGroupsLoad", "'"+urlGroupsLoad+"'")
	js = strings.ReplaceAll(js, "urlGroupDelete", "'"+urlGroupDelete+"'")
	js = strings.ReplaceAll(js, "urlGroupCreate", "'"+urlGroupCreate+"'")
	js = strings.ReplaceAll(js, "urlGroupUpdate", "'"+urlGroupUpdate+"'")
	js = strings.ReplaceAll(js, "urlGroupView", "'"+urlGroupView+"'")

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Groups", URL: linksHelper.GroupManager(nil)},
	})

	content := hb.Div().
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	page := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.Heading1().HTML("Groups")).
		Child(content)

	return u.Layout(w, r, "Groups | User Manager", page.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		// Vue and SweetAlert2 are already loaded by the default layout
		// (shared.Layout). No need to include them again here.
	})
}
