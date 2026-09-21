package role_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleRoleUpdateAjax(w http.ResponseWriter, r *http.Request) {
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
		Name   string `json:"name"`
		Handle string `json:"handle"`
		Status string `json:"status"`
		Memo   string `json:"memo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
		return
	}

	role, err := u.UserStore().RoleFindByID(r.Context(), payload.RoleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleRoleUpdateAjax RoleFindByID", slog.String("role_id", payload.RoleID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading role"))
		return
	}
	if role == nil {
		api.Respond(w, r, api.Error("Role not found"))
		return
	}

	name := strings.TrimSpace(payload.Name)
	handle := strings.TrimSpace(payload.Handle)

	if name == "" {
		api.Respond(w, r, api.Error("Name is required"))
		return
	}
	if handle == "" {
		api.Respond(w, r, api.Error("Handle is required"))
		return
	}
	if payload.Status == "" {
		api.Respond(w, r, api.Error("Status is required"))
		return
	}
	if payload.Status != userstore.ROLE_STATUS_ACTIVE && payload.Status != userstore.ROLE_STATUS_INACTIVE && payload.Status != userstore.ROLE_STATUS_DELETED {
		api.Respond(w, r, api.Error("Invalid status value"))
		return
	}

	if handle != role.GetHandle() {
		existing, err := u.UserStore().RoleFindByHandle(r.Context(), handle)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("roleUpdateController.handleRoleUpdateAjax RoleFindByHandle", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to check role handle"))
			return
		}
		if existing != nil && existing.GetID() != role.GetID() {
			api.Respond(w, r, api.Error("A role with this handle already exists"))
			return
		}
	}

	role.SetName(name)
	role.SetHandle(handle)
	role.SetStatus(payload.Status)
	role.SetMemo(strings.TrimSpace(payload.Memo))

	if err := u.UserStore().RoleUpdate(r.Context(), role); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleRoleUpdateAjax RoleUpdate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving role failed"))
		return
	}

	api.Respond(w, r, api.Success("Role saved successfully"))
}
