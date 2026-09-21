package group_view

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

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "User Manager", URL: linksHelper.UserManager(nil)},
		{Name: "Groups", URL: groupManagerURL},
		{Name: displayName, URL: linksHelper.GroupView(map[string]string{FieldGroupID: groupID})},
	})

	buttonBack := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(groupManagerURL)

	heading := hb.Heading1().HTML("View Group").Child(buttonBack)

	groupTitle := hb.Heading2().Class("mb-3").Text("Group: ").Text(displayName)

	buttonEdit := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-pencil-square me-2")).
		HTML("Edit").
		Href(linksHelper.GroupUpdate(map[string]string{FieldGroupID: groupID}))

	detailsTable := hb.Table().Class("table table-bordered").
		Child(hb.Tr().Child(hb.Th().Style("width:200px;").Text("ID")).Child(hb.Td().Text(group.GetID()))).
		Child(hb.Tr().Child(hb.Th().Text("Name")).Child(hb.Td().Text(group.GetName()))).
		Child(hb.Tr().Child(hb.Th().Text("Handle")).Child(hb.Td().Text(group.GetHandle()))).
		Child(hb.Tr().Child(hb.Th().Text("Status")).Child(hb.Td().Text(group.GetStatus()))).
		Child(hb.Tr().Child(hb.Th().Text("Memo")).Child(hb.Td().Text(group.GetMemo()))).
		Child(hb.Tr().Child(hb.Th().Text("Created")).Child(hb.Td().Text(group.GetCreatedAtCarbon().Format("d M Y H:i")))).
		Child(hb.Tr().Child(hb.Th().Text("Modified")).Child(hb.Td().Text(group.GetUpdatedAtCarbon().Format("d M Y H:i"))))

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("Group Details").Style("margin-bottom:0;display:inline-block;")).
			Child(buttonEdit),
	).Child(
		hb.Div().Class("card-body").Child(detailsTable),
	)

	membersCard := u.membersCard(w, r, groupID)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(heading).
		Child(groupTitle).
		Child(card).
		Child(membersCard)

	return u.Layout(w, r, "View Group | Groups", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})
}

// membersCard renders the read-only member list for the group. Member
// management happens on the dedicated group members page.
func (u *ui) membersCard(w http.ResponseWriter, r *http.Request, groupID string) hb.TagInterface {
	linksHelper := shared.NewLinksFromRequest(r)

	buttonEditMembers := hb.Hyperlink().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-people me-2")).
		HTML("Edit Members").
		Href(linksHelper.GroupMembers(map[string]string{FieldGroupID: groupID}))

	userGroupList, err := u.UserStore().UserGroupList(r.Context(), userstore.NewUserGroupQuery().SetGroupID(groupID))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupViewController.membersCard UserGroupList", slog.String("group_id", groupID), slog.String("error", err.Error()))
		}
		return hb.Div().Class("card mt-3").Child(
			hb.Div().Class("card-body").Child(hb.Raw(shared.ErrorAlert("Failed to load members"))),
		)
	}

	userIDs := make([]string, 0, len(userGroupList))
	for _, userGroup := range userGroupList {
		userIDs = append(userIDs, userGroup.GetUserID())
	}

	var rows []hb.TagInterface
	if len(userIDs) > 0 {
		userList, err := u.UserStore().UserList(r.Context(), userstore.NewUserQuery().SetIDIn(userIDs))
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupViewController.membersCard UserList", slog.String("error", err.Error()))
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
		body = hb.Div().Class("text-muted text-center py-3").Text("No users have this group.")
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
