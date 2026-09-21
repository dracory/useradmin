package group_manager

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleGroupsFetchAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}
	if !u.groupsEnabled() {
		api.Respond(w, r, api.Error("Groups are not enabled in the user store"))
		return ""
	}

	groupList, err := u.UserStore().GroupList(r.Context(), userstore.NewGroupQuery().
		SetOrderBy(userstore.COLUMN_NAME).
		SetSortDirection(userstore.SORT_ORDER_ASC))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupManagerController.handleGroupsFetchAjax GroupList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load groups"))
		return ""
	}

	groups := make([]map[string]interface{}, 0, len(groupList))
	for _, group := range groupList {
		groups = append(groups, map[string]interface{}{
			FieldID:        group.GetID(),
			FieldName:      group.GetName(),
			FieldHandle:    group.GetHandle(),
			FieldStatus:    group.GetStatus(),
			FieldMemo:      group.GetMemo(),
			FieldCreatedAt: group.GetCreatedAtCarbon().Format("d M Y"),
			FieldUpdatedAt: group.GetUpdatedAtCarbon().Format("d M Y"),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{
		FieldGroups: groups,
	}))
	return ""
}
