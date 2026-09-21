package group_members

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

func (u *ui) handleMembersFetchAjax(w http.ResponseWriter, r *http.Request) {
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
	groupID := req.GetStringTrimmed(r, FieldGroupID)
	if groupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return
	}

	userGroupList, err := u.UserStore().UserGroupList(r.Context(), userstore.NewUserGroupQuery().SetGroupID(groupID))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMembersFetchAjax UserGroupList", slog.String("group_id", groupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load members"))
		return
	}

	userIDs := make([]string, 0, len(userGroupList))
	for _, userGroup := range userGroupList {
		userIDs = append(userIDs, userGroup.GetUserID())
	}

	members := make([]map[string]any, 0, len(userIDs))
	if len(userIDs) > 0 {
		userList, err := u.UserStore().UserList(r.Context(), userstore.NewUserQuery().SetIDIn(userIDs))
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupMembersController.handleMembersFetchAjax UserList", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to load members"))
			return
		}

		for _, user := range userList {
			name := strings.TrimSpace(user.GetFirstName() + " " + user.GetLastName())
			members = append(members, map[string]any{
				FieldUserID:    user.GetID(),
				FieldName:      name,
				FieldFirstName: user.GetFirstName(),
				FieldLastName:  user.GetLastName(),
				FieldEmail:     user.GetEmail(),
				FieldStatus:    user.GetStatus(),
			})
		}
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldMembers: members,
	}))
}
