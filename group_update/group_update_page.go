package group_update

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
	//go:embed form.html
	formHTML string

	//go:embed form.js
	formJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	linksHelper := shared.NewLinksFromRequest(r)
	groupManagerURL := linksHelper.GroupManager(nil)

	if u.UserStore() == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "User store is not configured", groupManagerURL, 10)
	}
	if !u.groupsEnabled() {
		return shared.FlashError(u.FlashRedirect(), w, r, "Groups are not enabled in the user store", groupManagerURL, 10)
	}

	groupID := req.GetStringTrimmed(r, FieldGroupID)

	if groupID == "" {
		return shared.FlashError(u.FlashRedirect(), w, r, "Group ID is required", groupManagerURL, 10)
	}

	group, err := u.UserStore().GroupFindByID(r.Context(), groupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("renderPage GroupFindByID", slog.String("group_id", groupID), slog.String("error", err.Error()))
		}
		return shared.FlashError(u.FlashRedirect(), w, r, "Error loading group", groupManagerURL, 10)
	}
	if group == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "Group not found", groupManagerURL, 10)
	}

	displayName := group.GetName()
	if displayName == "" {
		displayName = group.GetID()
	}

	returnURL := shared.JSEscapeString(groupManagerURL)
	urlGetGroup := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"action": actionGroupFetch, FieldGroupID: groupID}))
	urlUpdateGroup := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"action": actionGroupUpdate}))
	escapedGroupID := shared.JSEscapeString(groupID)

	html := strings.ReplaceAll(formHTML, "GROUP_ID_PLACEHOLDER", "'"+escapedGroupID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(formJS, "GROUP_ID_PLACEHOLDER", "'"+escapedGroupID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetGroup", "'"+urlGetGroup+"'")
	js = strings.ReplaceAll(js, "urlUpdateGroup", "'"+urlUpdateGroup+"'")

	appHTML := hb.Raw(html)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Groups", URL: groupManagerURL},
		{Name: "Edit Group", URL: linksHelper.GroupUpdate(map[string]string{FieldGroupID: groupID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(groupManagerURL)

	heading := hb.Heading1().HTML("Edit Group").Child(buttonCancel)

	groupTitle := hb.Heading2().Class("mb-3").Text("Group: ").Text(displayName)

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Group Details").Style("margin-bottom:0;display:inline-block;")),
	).Child(
		hb.Div().Class("card-body").Child(appHTML),
	)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(heading).
		Child(groupTitle).
		Child(card)

	return u.Layout(w, r, "Edit Group | Groups", content.ToHTML(), struct {
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
