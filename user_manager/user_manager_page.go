package user_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/hb"
)

var (
	//go:embed users.html
	usersHTML string

	//go:embed users.js
	usersJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.UserStore() == nil {
		return shared.ErrorAlert("User store is not configured")
	}

	linksHelper := shared.NewLinksFromRequest(r)

	urlUsersLoad := shared.JSEscapeString(linksHelper.UserManager(map[string]string{"action": actionLoadUsers}))
	urlUserDelete := shared.JSEscapeString(linksHelper.UserManager(map[string]string{"action": actionDeleteUser}))
	urlUserCreate := shared.JSEscapeString(linksHelper.UserManager(map[string]string{"action": actionCreateUser}))
	urlUserUpdate := shared.JSEscapeString(linksHelper.UserUpdate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"}))

	// The impersonate button is controlled by the impersonateEnabled
	// Vue data property. When OnUserImpersonate is nil, the button is
	// hidden via v-if and the URL is set to empty.
	impersonateEnabled := "false"
	urlUserImpersonate := "''"
	if u.OnUserImpersonate() != nil {
		impersonateEnabled = "true"
		urlUserImpersonate = "'" + shared.JSEscapeString(linksHelper.UserImpersonate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"})) + "'"
	}

	html := strings.ReplaceAll(usersHTML, "urlUsersLoad", "'"+urlUsersLoad+"'")
	html = strings.ReplaceAll(html, "urlUserUpdate", "'"+urlUserUpdate+"'")
	html = strings.ReplaceAll(html, "urlUserImpersonate", urlUserImpersonate)
	js := strings.ReplaceAll(usersJS, "urlUsersLoad", "'"+urlUsersLoad+"'")
	js = strings.ReplaceAll(js, "urlUserDelete", "'"+urlUserDelete+"'")
	js = strings.ReplaceAll(js, "urlUserCreate", "'"+urlUserCreate+"'")
	js = strings.ReplaceAll(js, "urlUserUpdate", "'"+urlUserUpdate+"'")
	js = strings.ReplaceAll(js, "__urlUserImpersonate__", urlUserImpersonate)
	js = strings.ReplaceAll(js, "__impersonateEnabled__", impersonateEnabled)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
	})

	nav := hb.Div().Class("mb-3").
		Child(hb.Hyperlink().Class("btn btn-outline-primary btn-sm me-2").
			Child(hb.I().Class("bi bi-shield-lock me-1")).HTML("Roles").
			Href(linksHelper.RoleManager(nil))).
		Child(hb.Hyperlink().Class("btn btn-outline-primary btn-sm").
			Child(hb.I().Class("bi bi-people me-1")).HTML("Groups").
			Href(linksHelper.GroupManager(nil)))

	content := hb.Div().
		Child(nav).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	page := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(content)

	return u.Layout(w, r, "Users | User Manager", page.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		// Vue and SweetAlert2 are already loaded by the default layout
		// (shared.Layout). No need to include them again here.
	})
}
