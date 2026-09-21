package role_members

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
	if !u.rolesEnabled() {
		api.Respond(w, r, api.Error("Roles are not enabled in the user store"))
		return
	}

	var payload struct {
		RoleID  string `json:"role_id"`
		UserRef string `json:"user_ref"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
		return
	}

	userRef := strings.TrimSpace(payload.UserRef)
	if userRef == "" {
		api.Respond(w, r, api.Error("User ID or email is required"))
		return
	}

	role, err := u.UserStore().RoleFindByID(r.Context(), payload.RoleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleMemberAddAjax RoleFindByID", slog.String("role_id", payload.RoleID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading role"))
		return
	}
	if role == nil {
		api.Respond(w, r, api.Error("Role not found"))
		return
	}

	// Resolve the user reference: try ID first, then email.
	user, err := u.UserStore().UserFindByID(r.Context(), userRef)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleMemberAddAjax UserFindByID", slog.String("user_ref", userRef), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error looking up user"))
		return
	}
	if user == nil {
		user, err = u.UserStore().UserFindByEmail(r.Context(), userRef)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("roleUpdateController.handleMemberAddAjax UserFindByEmail", slog.String("user_ref", userRef), slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Error looking up user"))
			return
		}
	}
	if user == nil {
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	userRole, err := u.UserStore().UserRoleFindByUserIDAndRoleIDOrCreate(r.Context(), user.GetID(), role.GetID())
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleMemberAddAjax UserRoleFindByUserIDAndRoleIDOrCreate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to assign role"))
		return
	}
	if userRole == nil {
		api.Respond(w, r, api.Error("Failed to assign role"))
		return
	}

	api.Respond(w, r, api.Success("Member added successfully"))
}
