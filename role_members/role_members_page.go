package role_members

import (
	_ "embed"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

var (
	//go:embed members.html
	membersHTML string

	//go:embed members.js
	membersJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	linksHelper := shared.NewLinksFromRequest(r)
	roleManagerURL := linksHelper.RoleManager(nil)

	if u.UserStore() == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "User store is not configured", roleManagerURL, 10)
	}
	if !u.rolesEnabled() {
		return shared.FlashError(u.FlashRedirect(), w, r, "Roles are not enabled in the user store", roleManagerURL, 10)
	}

	roleID := req.GetStringTrimmed(r, FieldRoleID)

	if roleID == "" {
		return shared.FlashError(u.FlashRedirect(), w, r, "Role ID is required", roleManagerURL, 10)
	}

	role, err := u.UserStore().RoleFindByID(r.Context(), roleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("renderPage RoleFindByID", slog.String("role_id", roleID), slog.String("error", err.Error()))
		}
		return shared.FlashError(u.FlashRedirect(), w, r, "Error loading role", roleManagerURL, 10)
	}
	if role == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "Role not found", roleManagerURL, 10)
	}

	displayName := role.GetName()
	if displayName == "" {
		displayName = role.GetID()
	}

	roleViewURL := linksHelper.RoleView(map[string]string{FieldRoleID: roleID})

	returnURL := shared.JSEscapeString(roleViewURL)
	urlGetMembers := shared.JSEscapeString(linksHelper.RoleMembers(map[string]string{"action": actionMembersFetch}))
	urlMemberAdd := shared.JSEscapeString(linksHelper.RoleMembers(map[string]string{"action": actionMemberAdd}))
	urlMemberRemove := shared.JSEscapeString(linksHelper.RoleMembers(map[string]string{"action": actionMemberRemove}))
	urlUsersSearch := shared.JSEscapeString(linksHelper.RoleMembers(map[string]string{"action": actionUsersSearch}))
	escapedRoleID := shared.JSEscapeString(roleID)

	html := strings.ReplaceAll(membersHTML, "ROLE_ID_PLACEHOLDER", "'"+escapedRoleID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(membersJS, "ROLE_ID_PLACEHOLDER", "'"+escapedRoleID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetMembers", "'"+urlGetMembers+"'")
	js = strings.ReplaceAll(js, "urlMemberAdd", "'"+urlMemberAdd+"'")
	js = strings.ReplaceAll(js, "urlMemberRemove", "'"+urlMemberRemove+"'")
	js = strings.ReplaceAll(js, "urlUsersSearch", "'"+urlUsersSearch+"'")

	appHTML := hb.Raw(html)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Roles", URL: roleManagerURL},
		{Name: displayName, URL: roleViewURL},
		{Name: "Members", URL: linksHelper.RoleMembers(map[string]string{FieldRoleID: roleID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(roleViewURL)

	heading := hb.Heading1().HTML("Edit Members").Child(buttonCancel)

	roleTitle := hb.Heading2().Class("mb-3").Text("Role: ").Text(displayName)

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Members").Style("margin-bottom:0;display:inline-block;")),
	).Child(
		hb.Div().Class("card-body").Child(appHTML),
	)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(heading).
		Child(roleTitle).
		Child(card)

	return u.Layout(w, r, "Edit Members | Roles", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
		Scripts: []string{js},
	})
}
