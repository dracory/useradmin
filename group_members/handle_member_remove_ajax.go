package group_members

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleMemberRemoveAjax(w http.ResponseWriter, r *http.Request) {
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

	var payload struct {
		GroupID string `json:"group_id"`
		UserID  string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.GroupID == "" || payload.UserID == "" {
		api.Respond(w, r, api.Error("Group ID and user ID are required"))
		return
	}

	userGroup, err := u.UserStore().UserGroupFindByUserIDAndGroupID(r.Context(), payload.UserID, payload.GroupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMemberRemoveAjax UserGroupFindByUserIDAndGroupID", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load membership"))
		return
	}
	if userGroup == nil {
		api.Respond(w, r, api.Error("Membership not found"))
		return
	}

	if err := u.UserStore().UserGroupSoftDelete(r.Context(), userGroup); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMemberRemoveAjax UserGroupSoftDelete", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to remove member"))
		return
	}

	api.Respond(w, r, api.Success("Member removed successfully"))
}
