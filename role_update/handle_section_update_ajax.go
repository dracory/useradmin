package role_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

// loadRoleForUpdate fetches a role by ID.
// Returns nil when an error response was already sent.
func (u *ui) loadRoleForUpdate(w http.ResponseWriter, r *http.Request, roleID string) userstore.RoleInterface {
	role, err := u.UserStore().RoleFindByID(r.Context(), roleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController loadRoleForUpdate RoleFindByID", slog.String("role_id", roleID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading role"))
		return nil
	}
	if role == nil {
		api.Respond(w, r, api.Error("Role not found"))
		return nil
	}

	return role
}

// persistRole saves the role.
// Returns false when an error response was already sent.
func (u *ui) persistRole(w http.ResponseWriter, r *http.Request, role userstore.RoleInterface) bool {
	if err := u.UserStore().RoleUpdate(r.Context(), role); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController persistRole RoleUpdate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving role failed"))
		return false
	}

	return true
}

// decodeSectionPayload validates the request method, store availability,
// and decodes the JSON body. Returns false when an error response was sent.
func (u *ui) decodeSectionPayload(w http.ResponseWriter, r *http.Request, payload any) bool {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return false
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return false
	}
	if !u.rolesEnabled() {
		api.Respond(w, r, api.Error("Roles are not enabled in the user store"))
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return false
	}
	return true
}

// handleStatusUpdateAjax updates only the role's status.
func (u *ui) handleStatusUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RoleID string `json:"role_id"`
		Status string `json:"status"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
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

	role := u.loadRoleForUpdate(w, r, payload.RoleID)
	if role == nil {
		return
	}

	role.SetStatus(payload.Status)
	if !u.persistRole(w, r, role) {
		return
	}

	api.Respond(w, r, api.Success("Status saved successfully"))
}

// handleDetailsUpdateAjax updates only the role's name and handle.
func (u *ui) handleDetailsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RoleID string `json:"role_id"`
		Name   string `json:"name"`
		Handle string `json:"handle"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
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

	role := u.loadRoleForUpdate(w, r, payload.RoleID)
	if role == nil {
		return
	}

	if handle != role.GetHandle() {
		existing, err := u.UserStore().RoleFindByHandle(r.Context(), handle)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("roleUpdateController.handleDetailsUpdateAjax RoleFindByHandle", slog.String("error", err.Error()))
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
	if !u.persistRole(w, r, role) {
		return
	}

	api.Respond(w, r, api.Success("Role details saved successfully"))
}

// handleNotesUpdateAjax updates only the role's admin notes (memo).
func (u *ui) handleNotesUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RoleID string `json:"role_id"`
		Memo   string `json:"memo"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.RoleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
		return
	}

	role := u.loadRoleForUpdate(w, r, payload.RoleID)
	if role == nil {
		return
	}

	role.SetMemo(strings.TrimSpace(payload.Memo))
	if !u.persistRole(w, r, role) {
		return
	}

	api.Respond(w, r, api.Success("Notes saved successfully"))
}
