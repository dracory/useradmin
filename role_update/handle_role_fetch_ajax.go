package role_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (u *ui) handleRoleFetchAjax(w http.ResponseWriter, r *http.Request) {
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
	roleID := req.GetStringTrimmed(r, FieldRoleID)
	if roleID == "" {
		api.Respond(w, r, api.Error("Role ID is required"))
		return
	}

	role, err := u.UserStore().RoleFindByID(r.Context(), roleID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleUpdateController.handleRoleFetchAjax RoleFindByID", slog.String("role_id", roleID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading role"))
		return
	}
	if role == nil {
		api.Respond(w, r, api.Error("Role not found"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldID:     role.GetID(),
		FieldName:   role.GetName(),
		FieldHandle: role.GetHandle(),
		FieldStatus: role.GetStatus(),
		FieldMemo:   role.GetMemo(),
	}))
}
