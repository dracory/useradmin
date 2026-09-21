package role_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/hb"
)

var (
	//go:embed roles.html
	rolesHTML string

	//go:embed roles.js
	rolesJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.UserStore() == nil {
		return shared.ErrorAlert("User store is not configured")
	}
	if !u.rolesEnabled() {
		return shared.ErrorAlert("Roles are not enabled in the user store")
	}

	linksHelper := shared.NewLinksFromRequest(r)

	urlRolesLoad := shared.JSEscapeString(linksHelper.RoleManager(map[string]string{"action": actionLoadRoles}))
	urlRoleDelete := shared.JSEscapeString(linksHelper.RoleManager(map[string]string{"action": actionDeleteRole}))
	urlRoleCreate := shared.JSEscapeString(linksHelper.RoleManager(map[string]string{"action": actionCreateRole}))
	urlRoleUpdate := shared.JSEscapeString(linksHelper.RoleUpdate(map[string]string{"role_id": "ROLE_ID_PLACEHOLDER"}))
	urlRoleView := shared.JSEscapeString(linksHelper.RoleView(map[string]string{"role_id": "ROLE_ID_PLACEHOLDER"}))

	html := strings.ReplaceAll(rolesHTML, "urlRoleUpdate", "'"+urlRoleUpdate+"'")
	html = strings.ReplaceAll(html, "urlRoleView", "'"+urlRoleView+"'")
	js := strings.ReplaceAll(rolesJS, "urlRolesLoad", "'"+urlRolesLoad+"'")
	js = strings.ReplaceAll(js, "urlRoleDelete", "'"+urlRoleDelete+"'")
	js = strings.ReplaceAll(js, "urlRoleCreate", "'"+urlRoleCreate+"'")
	js = strings.ReplaceAll(js, "urlRoleUpdate", "'"+urlRoleUpdate+"'")
	js = strings.ReplaceAll(js, "urlRoleView", "'"+urlRoleView+"'")

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Roles", URL: linksHelper.RoleManager(nil)},
	})

	content := hb.Div().
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	page := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.Heading1().HTML("Roles")).
		Child(content)

	return u.Layout(w, r, "Roles | User Manager", page.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		// Vue and SweetAlert2 are already loaded by the default layout
		// (shared.Layout). No need to include them again here.
	})
}
