package role_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleRoleDeleteAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}
	if !u.rolesEnabled() {
		api.Respond(w, r, api.Error("Roles are not enabled in the user store"))
		return ""
	}

	var payload struct {
		RoleID string `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
		return ""
	}

	if err := u.UserStore().RoleSoftDeleteByID(r.Context(), payload.RoleID); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleManagerController.handleRoleDeleteAjax RoleSoftDeleteByID", slog.String("role_id", payload.RoleID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to delete role"))
		return ""
	}

	api.Respond(w, r, api.Success("Role deleted successfully"))
	return ""
}
