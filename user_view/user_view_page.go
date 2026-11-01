package user_view

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/hb"
	"github.com/dracory/req"
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	linksHelper := shared.NewLinksFromRequest(r)
	userManagerURL := linksHelper.UserManager(nil)

	if u.UserStore() == nil {
		return shared.FlashError(u.FlashRedirect(), w, r, "User store is not configured", userManagerURL, 10)
	}

	userID := req.GetStringTrimmed(r, FieldUserID)

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
		return shared.FlashError(u.FlashRedirect(), w, r, "User not found", userManagerURL, 10)
	}

	// Unseal PII for display (e.g. detokenize, decrypt fields).
	if u.UserPiiUnseal() != nil {
		unsealed, err := u.UserPiiUnseal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("renderPage UserPiiUnseal", slog.String("error", err.Error()))
			}
		} else {
			user = unsealed
		}
	}

	displayName := strings.TrimSpace(user.GetFirstName() + " " + user.GetLastName())
	if displayName == "" {
		displayName = user.GetID()
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: userManagerURL},
		{Name: displayName, URL: linksHelper.UserView(map[string]string{FieldUserID: userID})},
	})

	buttonBack := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(userManagerURL)

	heading := hb.Heading1().HTML("View User").Child(buttonBack)

	userTitle := hb.Heading2().Class("mb-3").Text("User: ").Text(displayName)

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-pencil-square me-2")).
		HTML("Edit").
		Href(linksHelper.UserUpdate(map[string]string{FieldUserID: userID}))

	detailsTable := hb.Table().Class("table table-bordered").
		Child(hb.Tr().Child(hb.Th().Style("width:200px;").Text("ID")).Child(hb.Td().Text(user.GetID()))).
		Child(hb.Tr().Child(hb.Th().Text("First Name")).Child(hb.Td().Text(user.GetFirstName()))).
		Child(hb.Tr().Child(hb.Th().Text("Middle Names")).Child(hb.Td().Text(user.GetMiddleNames()))).
		Child(hb.Tr().Child(hb.Th().Text("Last Name")).Child(hb.Td().Text(user.GetLastName()))).
		Child(hb.Tr().Child(hb.Th().Text("Email")).Child(hb.Td().Text(user.GetEmail()))).
		Child(hb.Tr().Child(hb.Th().Text("Phone")).Child(hb.Td().Text(user.GetPhone()))).
		Child(hb.Tr().Child(hb.Th().Text("Business Name")).Child(hb.Td().Text(user.GetBusinessName()))).
		Child(hb.Tr().Child(hb.Th().Text("Role")).Child(hb.Td().Text(user.GetRole()))).
		Child(hb.Tr().Child(hb.Th().Text("Status")).Child(hb.Td().Text(user.GetStatus()))).
		Child(hb.Tr().Child(hb.Th().Text("Country")).Child(hb.Td().Text(user.GetCountry()))).
		Child(hb.Tr().Child(hb.Th().Text("Timezone")).Child(hb.Td().Text(user.GetTimezone()))).
		Child(hb.Tr().Child(hb.Th().Text("Memo")).Child(hb.Td().Text(user.GetMemo()))).
		Child(hb.Tr().Child(hb.Th().Text("Created")).Child(hb.Td().Text(user.GetCreatedAtCarbon().Format("d M Y H:i")))).
		Child(hb.Tr().Child(hb.Th().Text("Modified")).Child(hb.Td().Text(user.GetUpdatedAtCarbon().Format("d M Y H:i"))))

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("User Details").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEdit),
	).Child(
		hb.Div().Class("card-body").Child(detailsTable),
	)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(heading).
		Child(shared.SectionNav(r, u.UserStore())).
		Child(userTitle).
		Child(card).
		Child(u.rolesCard(w, r, userID)).
		Child(u.groupsCard(w, r, userID))

	return u.Layout(w, r, "View User | Users", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})
}

// rolesCard renders the read-only list of roles assigned to the user.
// Role assignment happens on the dedicated user update page.
func (u *ui) rolesCard(w http.ResponseWriter, r *http.Request, userID string) hb.TagInterface {
	if !u.rolesEnabled() {
		return hb.Span()
	}

	linksHelper := shared.NewLinksFromRequest(r)

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-pencil-square me-2")).
		HTML("Edit Roles").
		Href(linksHelper.UserUpdate(map[string]string{FieldUserID: userID}))

	roleList, err := u.UserStore().UserRoles(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userViewController.rolesCard UserRoles", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		return hb.Div().Class("card mt-3").Child(
			hb.Div().Class("card-body").Child(hb.Raw(shared.ErrorAlert("Failed to load roles"))),
		)
	}

	var rows []hb.TagInterface
	for _, role := range roleList {
		rows = append(rows, hb.Tr().
			Child(hb.Td().Child(
				hb.Hyperlink().
					Class("text-decoration-none").
					Text(role.GetName()).
					Href(linksHelper.RoleView(map[string]string{"role_id": role.GetID()})),
			)).
			Child(hb.Td().Text(role.GetHandle())).
			Child(hb.Td().Text(role.GetStatus())))
	}

	var body hb.TagInterface
	if len(rows) == 0 {
		body = hb.Div().Class("text-muted text-center py-3").Text("This user has no roles.")
	} else {
		body = hb.Div().Class("table-responsive").Child(
			hb.Table().Class("table table-striped table-hover table-bordered align-middle").
				Child(hb.Thead().Child(hb.Tr().
					Child(hb.Th().Text("Name")).
					Child(hb.Th().Text("Handle")).
					Child(hb.Th().Text("Status")))).
				Child(hb.Tbody().Children(rows)),
		)
	}

	return hb.Div().Class("card mt-3").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Roles").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEdit),
	).Child(
		hb.Div().Class("card-body").Child(body),
	)
}

// groupsCard renders the read-only list of groups the user belongs to.
// Group assignment happens on the dedicated user update page.
func (u *ui) groupsCard(w http.ResponseWriter, r *http.Request, userID string) hb.TagInterface {
	if !u.groupsEnabled() {
		return hb.Span()
	}

	linksHelper := shared.NewLinksFromRequest(r)

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-pencil-square me-2")).
		HTML("Edit Groups").
		Href(linksHelper.UserUpdate(map[string]string{FieldUserID: userID}))

	groupList, err := u.UserStore().UserGroups(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userViewController.groupsCard UserGroups", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		return hb.Div().Class("card mt-3").Child(
			hb.Div().Class("card-body").Child(hb.Raw(shared.ErrorAlert("Failed to load groups"))),
		)
	}

	var rows []hb.TagInterface
	for _, group := range groupList {
		rows = append(rows, hb.Tr().
			Child(hb.Td().Child(
				hb.Hyperlink().
					Class("text-decoration-none").
					Text(group.GetName()).
					Href(linksHelper.GroupView(map[string]string{"group_id": group.GetID()})),
			)).
			Child(hb.Td().Text(group.GetHandle())).
			Child(hb.Td().Text(group.GetStatus())))
	}

	var body hb.TagInterface
	if len(rows) == 0 {
		body = hb.Div().Class("text-muted text-center py-3").Text("This user has no groups.")
	} else {
		body = hb.Div().Class("table-responsive").Child(
			hb.Table().Class("table table-striped table-hover table-bordered align-middle").
				Child(hb.Thead().Child(hb.Tr().
					Child(hb.Th().Text("Name")).
					Child(hb.Th().Text("Handle")).
					Child(hb.Th().Text("Status")))).
				Child(hb.Tbody().Children(rows)),
		)
	}

	return hb.Div().Class("card mt-3").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Groups").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEdit),
	).Child(
		hb.Div().Class("card-body").Child(body),
	)
}
