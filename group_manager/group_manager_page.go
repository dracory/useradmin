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
	linksHelper := shared.NewLinksFromRequest(r)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Groups", URL: ""},
	})

	pageShell := func(children ...hb.TagInterface) string {
		content := hb.Div().Children(children)
		page := hb.Div().
			Class("container").
			Class("py-4").
			Child(breadcrumbs).
			Child(shared.PageHeader("Users. Groups")).
			Child(shared.SectionNav(r, u.UserStore())).
			Child(content)
		return u.Layout(w, r, "Groups | User Manager", page.ToHTML(), struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}{})
	}

	if u.UserStore() == nil {
		return pageShell(hb.Raw(shared.ErrorAlert("User store is not configured")))
	}
	if !u.groupsEnabled() {
		return pageShell(hb.Raw(shared.ErrorAlert("Groups are not enabled in the user store")))
	}

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

	// Prepend the loadVueIfNeeded guard so the script is self-contained
	// and does not depend on the layout defining it.
	js = shared.VueLoaderJS + "\n" + js

	return pageShell(hb.Raw(html), hb.Script(js))
}
