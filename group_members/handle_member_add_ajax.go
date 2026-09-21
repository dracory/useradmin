package group_members

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
)

func (u *ui) handleMemberAddAjax(w http.ResponseWriter, r *http.Request) {
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
		UserRef string `json:"user_ref"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return
	}

	userRef := strings.TrimSpace(payload.UserRef)
	if userRef == "" {
		api.Respond(w, r, api.Error("User ID or email is required"))
		return
	}

	group, err := u.UserStore().GroupFindByID(r.Context(), payload.GroupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMemberAddAjax GroupFindByID", slog.String("group_id", payload.GroupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading group"))
		return
	}
	if group == nil {
		api.Respond(w, r, api.Error("Group not found"))
		return
	}

	// Resolve the user reference: try ID first, then email.
	user, err := u.UserStore().UserFindByID(r.Context(), userRef)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMemberAddAjax UserFindByID", slog.String("user_ref", userRef), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error looking up user"))
		return
	}
	if user == nil {
		user, err = u.UserStore().UserFindByEmail(r.Context(), userRef)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupMembersController.handleMemberAddAjax UserFindByEmail", slog.String("user_ref", userRef), slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Error looking up user"))
			return
		}
	}
	if user == nil {
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	userGroup, err := u.UserStore().UserGroupFindByUserIDAndGroupIDOrCreate(r.Context(), user.GetID(), group.GetID())
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupMembersController.handleMemberAddAjax UserGroupFindByUserIDAndGroupIDOrCreate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to assign group"))
		return
	}
	if userGroup == nil {
		api.Respond(w, r, api.Error("Failed to assign group"))
		return
	}

	api.Respond(w, r, api.Success("Member added successfully"))
}
