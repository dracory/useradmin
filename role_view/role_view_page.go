package role_view

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/userstore"

	"github.com/dracory/hb"
	"github.com/dracory/req"
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

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Roles", URL: roleManagerURL},
		{Name: displayName, URL: linksHelper.RoleView(map[string]string{FieldRoleID: roleID})},
	})

	buttonBack := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(roleManagerURL)

	heading := hb.Heading1().HTML("View Role").Child(buttonBack)

	roleTitle := hb.Heading2().Class("mb-3").Text("Role: ").Text(displayName)

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-pencil-square me-2")).
		HTML("Edit").
		Href(linksHelper.RoleUpdate(map[string]string{FieldRoleID: roleID}))

	detailsTable := hb.Table().Class("table table-bordered").
		Child(hb.Tr().Child(hb.Th().Style("width:200px;").Text("ID")).Child(hb.Td().Text(role.GetID()))).
		Child(hb.Tr().Child(hb.Th().Text("Name")).Child(hb.Td().Text(role.GetName()))).
		Child(hb.Tr().Child(hb.Th().Text("Handle")).Child(hb.Td().Text(role.GetHandle()))).
		Child(hb.Tr().Child(hb.Th().Text("Status")).Child(hb.Td().Text(role.GetStatus()))).
		Child(hb.Tr().Child(hb.Th().Text("Memo")).Child(hb.Td().Text(role.GetMemo()))).
		Child(hb.Tr().Child(hb.Th().Text("Created")).Child(hb.Td().Text(role.GetCreatedAtCarbon().Format("d M Y H:i")))).
		Child(hb.Tr().Child(hb.Th().Text("Modified")).Child(hb.Td().Text(role.GetUpdatedAtCarbon().Format("d M Y H:i"))))

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Role Details").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEdit),
	).Child(
		hb.Div().Class("card-body").Child(detailsTable),
	)

	membersCard := u.membersCard(w, r, roleID)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(heading).
		Child(roleTitle).
		Child(card).
		Child(membersCard)

	return u.Layout(w, r, "View Role | Roles", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})
}

// membersCard renders the read-only member list for the role. Member
// management happens on the dedicated role members page.
func (u *ui) membersCard(w http.ResponseWriter, r *http.Request, roleID string) hb.TagInterface {
	linksHelper := shared.NewLinksFromRequest(r)

	buttonEditMembers := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-people me-2")).
		HTML("Edit Members").
		Href(linksHelper.RoleMembers(map[string]string{FieldRoleID: roleID}))

	userRoleList, err := u.UserStore().UserRoleList(r.Context(), userstore.NewUserRoleQuery().SetRoleID(roleID))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleViewController.membersCard UserRoleList", slog.String("role_id", roleID), slog.String("error", err.Error()))
		}
		return hb.Div().Class("card mt-3").Child(
			hb.Div().Class("card-body").Child(hb.Raw(shared.ErrorAlert("Failed to load members"))),
		)
	}

	userIDs := make([]string, 0, len(userRoleList))
	for _, userRole := range userRoleList {
		userIDs = append(userIDs, userRole.GetUserID())
	}

	var rows []hb.TagInterface
	if len(userIDs) > 0 {
		userList, err := u.UserStore().UserList(r.Context(), userstore.NewUserQuery().SetIDIn(userIDs))
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("roleViewController.membersCard UserList", slog.String("error", err.Error()))
			}
			return hb.Div().Class("card mt-3").Child(
				hb.Div().Class("card-body").Child(hb.Raw(shared.ErrorAlert("Failed to load members"))),
			)
		}
		for _, user := range userList {
			name := strings.TrimSpace(user.GetFirstName() + " " + user.GetLastName())
			if name == "" {
				name = user.GetID()
			}
			rows = append(rows, hb.Tr().
				Child(hb.Td().Text(name)).
				Child(hb.Td().Text(user.GetEmail())).
				Child(hb.Td().Text(user.GetStatus())))
		}
	}

	var body hb.TagInterface
	if len(rows) == 0 {
		body = hb.Div().Class("text-muted text-center py-3").Text("No users have this role.")
	} else {
		body = hb.Div().Class("table-responsive").Child(
			hb.Table().Class("table table-striped table-hover table-bordered align-middle").
				Child(hb.Thead().Child(hb.Tr().
					Child(hb.Th().Text("Name")).
					Child(hb.Th().Text("Email")).
					Child(hb.Th().Text("Status")))).
				Child(hb.Tbody().Children(rows)),
		)
	}

	return hb.Div().Class("card mt-3").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Members").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEditMembers),
	).Child(
		hb.Div().Class("card-body").Child(body),
	)
}
