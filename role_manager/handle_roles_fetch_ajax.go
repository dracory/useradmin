package role_manager

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleRolesFetchAjax(w http.ResponseWriter, r *http.Request) string {
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

	roleList, err := u.UserStore().RoleList(r.Context(), userstore.NewRoleQuery().
		SetOrderBy(userstore.COLUMN_NAME).
		SetSortDirection(userstore.SORT_ORDER_ASC))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("roleManagerController.handleRolesFetchAjax RoleList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load roles"))
		return ""
	}

	roles := make([]map[string]interface{}, 0, len(roleList))
	for _, role := range roleList {
		roles = append(roles, map[string]interface{}{
			FieldID:        role.GetID(),
			FieldName:      role.GetName(),
			FieldHandle:    role.GetHandle(),
			FieldStatus:    role.GetStatus(),
			FieldMemo:      role.GetMemo(),
			FieldCreatedAt: role.GetCreatedAtCarbon().Format("d M Y"),
			FieldUpdatedAt: role.GetUpdatedAtCarbon().Format("d M Y"),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{
		FieldRoles: roles,
	}))
	return ""
}
