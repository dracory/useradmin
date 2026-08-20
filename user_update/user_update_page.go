package user_update

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
	userManagerURL := linksHelper.UserManager(nil)

	if u.UserStore() == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "User store is not configured", userManagerURL, 10)
	}

	userID := req.GetStringTrimmed(r, "user_id")

	if userID == "" {
		return shared.FlashError(u.FlashRedirect(), w, r, "User ID is required", userManagerURL, 10)
	}

	user, err := u.UserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("renderPage UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		return shared.FlashError(u.FlashRedirect(), w, r, "Error loading user", userManagerURL, 10)
	}
	if user == nil {
		if u.Logger() != nil {
			u.Logger().Error("renderPage user not found", slog.String("user_id", userID))
		}
		return shared.FlashError(u.FlashRedirect(), w, r, "User not found", userManagerURL, 10)
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	if u.VaultTokenizer() != nil {
		fn, ln, _, _, _, err := u.VaultTokenizer().Untokenize(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("At userUpdateController > renderPage", slog.String("error", err.Error()))
			}
		} else {
			firstName = fn
			lastName = ln
		}
	}

	displayName := strings.TrimSpace(firstName + " " + lastName)
	if displayName == "" {
		displayName = user.GetID()
	}

	returnURL := shared.JSEscapeString(userManagerURL)
	urlGetUser := shared.JSEscapeString(linksHelper.UserUpdate(map[string]string{"action": actionUserFetch, "user_id": userID}))
	urlGetTimezones := shared.JSEscapeString(linksHelper.UserUpdate(map[string]string{"action": actionGetTimezones}))
	urlUpdateUser := shared.JSEscapeString(linksHelper.UserUpdate(map[string]string{"action": actionUserUpdate}))
	escapedUserID := shared.JSEscapeString(userID)

	html := strings.ReplaceAll(formHTML, "USER_ID_PLACEHOLDER", "'"+escapedUserID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(formJS, "USER_ID_PLACEHOLDER", "'"+escapedUserID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetUser", "'"+urlGetUser+"'")
	js = strings.ReplaceAll(js, "urlGetTimezones", "'"+urlGetTimezones+"'")
	js = strings.ReplaceAll(js, "urlUpdateUser", "'"+urlUpdateUser+"'")

	// form.html already contains <div id="app-user-update" class="mt-3">
	// so we inject it as raw HTML — no wrapper div needed (avoids
	// duplicate id in the DOM). Vue is loaded by the default layout
	// (shared.Layout); no need to include it again here.
	appHTML := hb.Raw(html)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: userManagerURL},
		{Name: "Edit User", URL: linksHelper.UserUpdate(map[string]string{"user_id": userID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(userManagerURL)

	heading := hb.Heading1().HTML("Edit User").Child(buttonCancel)

	userTitle := hb.Heading2().Class("mb-3").Text("User: ").Text(displayName)

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("User Details").Style("margin-bottom:0;display:inline-block;")),
	).Child(
		hb.Div().Class("card-body").Child(appHTML),
	)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(heading).
		Child(userTitle).
		Child(card)

	return u.Layout(w, r, "Edit User | Users", content.ToHTML(), struct {
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
