package role_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleRoleCreateAjax(w http.ResponseWriter, r *http.Request) string {
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
		Name   string `json:"name"`
		Handle string `json:"handle"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	name := strings.TrimSpace(payload.Name)
	handle := strings.TrimSpace(payload.Handle)

	if name == "" {
		api.Respond(w, r, api.Error("Name is required"))
		return ""
	}
	if handle == "" {
		api.Respond(w, r, api.Error("Handle is required"))
		return ""
	}

	existing, err := u.UserStore().RoleFindByHandle(r.Context(), handle)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleManagerController.handleRoleCreateAjax RoleFindByHandle", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to check role handle"))
		return ""
	}
	if existing != nil {
		api.Respond(w, r, api.Error("A role with this handle already exists"))
		return ""
	}

	role := userstore.NewRole().
		SetName(name).
		SetHandle(handle).
		SetStatus(userstore.ROLE_STATUS_ACTIVE)

	if err := u.UserStore().RoleCreate(r.Context(), role); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleManagerController.handleRoleCreateAjax RoleCreate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to create role"))
		return ""
	}

	api.Respond(w, r, api.Success("Role created successfully"))
	return ""
}
