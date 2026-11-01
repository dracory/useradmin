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
	urlUpdateStatus := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"action": actionStatusUpdate}))
	urlUpdateDetails := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"action": actionDetailsUpdate}))
	urlUpdateNotes := shared.JSEscapeString(linksHelper.GroupUpdate(map[string]string{"action": actionNotesUpdate}))
	escapedGroupID := shared.JSEscapeString(groupID)

	html := strings.ReplaceAll(formHTML, "GROUP_ID_PLACEHOLDER", "'"+escapedGroupID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(formJS, "GROUP_ID_PLACEHOLDER", "'"+escapedGroupID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetGroup", "'"+urlGetGroup+"'")
	js = strings.ReplaceAll(js, "urlUpdateStatus", "'"+urlUpdateStatus+"'")
	js = strings.ReplaceAll(js, "urlUpdateDetails", "'"+urlUpdateDetails+"'")
	js = strings.ReplaceAll(js, "urlUpdateNotes", "'"+urlUpdateNotes+"'")

	// Prepend the loadVueIfNeeded guard so the script is self-contained
	// and does not depend on the layout defining it.
	js = shared.VueLoaderJS + "\n" + js

	appHTML := hb.Raw(html)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Groups", URL: groupManagerURL},
		{Name: "Edit Group", URL: ""},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-outline-secondary btn-sm").
		Child(hb.I().Class("bi bi-chevron-left me-1")).
		HTML("Back").
		Href(groupManagerURL)

	heading := shared.PageHeader("Users. Edit Group", buttonCancel)

	groupTitle := hb.Heading2().Class("mb-3").Text("Group: ").Text(displayName)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(heading).
		Child(shared.SectionNav(r)).
		Child(groupTitle).
		Child(appHTML)

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
