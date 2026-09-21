package role_members

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
	if !u.rolesEnabled() {
		api.Respond(w, r, api.Error("Roles are not enabled in the user store"))
		return
	}

	var payload struct {
		RoleID string `json:"role_id"`
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.RoleID == "" || payload.UserID == "" {
		api.Respond(w, r, api.Error("Role ID and user ID are required"))
		return
	}

	userRole, err := u.UserStore().UserRoleFindByUserIDAndRoleID(r.Context(), payload.UserID, payload.RoleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleMemberRemoveAjax UserRoleFindByUserIDAndRoleID", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load membership"))
		return
	}
	if userRole == nil {
		api.Respond(w, r, api.Error("Membership not found"))
		return
	}

	if err := u.UserStore().UserRoleSoftDelete(r.Context(), userRole); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleMemberRemoveAjax UserRoleSoftDelete", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to remove member"))
		return
	}

	api.Respond(w, r, api.Success("Member removed successfully"))
}
