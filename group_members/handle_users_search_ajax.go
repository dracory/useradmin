package group_members

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

// usersSearchLimit caps each search query and the merged result set.
const usersSearchLimit = 50

func (u *ui) handleUsersSearchAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return
	}
	if !u.groupsEnabled() {
		api.Respond(w, r, api.Error("Groups are not enabled in the user store"))
		return
	}

	q := strings.TrimSpace(req.GetStringTrimmed(r, "q"))

	userList, err := u.searchUsers(r, q)
	if err != nil {
		api.Respond(w, r, api.Error("Failed to search users"))
		return
	}

	users := make([]map[string]any, 0, len(userList))
	for _, user := range userList {
		name := strings.TrimSpace(user.GetFirstName() + " " + user.GetLastName())
		users = append(users, map[string]any{
			FieldUserID:    user.GetID(),
			FieldName:      name,
			FieldFirstName: user.GetFirstName(),
			FieldLastName:  user.GetLastName(),
			FieldEmail:     user.GetEmail(),
			FieldStatus:    user.GetStatus(),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		"users": users,
	}))
}

// searchUsers finds users matching q against first name, last name, and
// email (any match), merged and deduplicated. An empty q returns the
// first page of users.
func (u *ui) searchUsers(r *http.Request, q string) ([]userstore.UserInterface, error) {
	ctx := r.Context()

	if q == "" {
		return u.UserStore().UserList(ctx, userstore.NewUserQuery().
			SetOrderBy(userstore.COLUMN_CREATED_AT).
			SetSortDirection(userstore.SORT_ORDER_DESC).
			SetLimit(usersSearchLimit))
	}

	seen := map[string]bool{}
	var merged []userstore.UserInterface

	queries := []*userstore.UserQueryInterface{}
	firstNameQuery := userstore.NewUserQuery().SetFirstNameLike(q).SetLimit(usersSearchLimit)
	lastNameQuery := userstore.NewUserQuery().SetLastNameLike(q).SetLimit(usersSearchLimit)
	emailQuery := userstore.NewUserQuery().SetEmailLike(q).SetLimit(usersSearchLimit)
	queries = append(queries, &firstNameQuery, &lastNameQuery, &emailQuery)

	for _, query := range queries {
		list, err := u.UserStore().UserList(ctx, *query)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupMembersController.handleUsersSearchAjax UserList", slog.String("error", err.Error()))
			}
			return nil, err
		}
		for _, user := range list {
			if seen[user.GetID()] {
				continue
			}
			seen[user.GetID()] = true
			merged = append(merged, user)
			if len(merged) >= usersSearchLimit {
				return merged, nil
			}
		}
	}

	return merged, nil
}
