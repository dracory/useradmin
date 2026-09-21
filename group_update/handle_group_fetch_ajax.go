package group_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (u *ui) handleGroupFetchAjax(w http.ResponseWriter, r *http.Request) {
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
	groupID := req.GetStringTrimmed(r, FieldGroupID)
	if groupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return
	}

	group, err := u.UserStore().GroupFindByID(r.Context(), groupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupUpdateController.handleGroupFetchAjax GroupFindByID", slog.String("group_id", groupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading group"))
		return
	}
	if group == nil {
		api.Respond(w, r, api.Error("Group not found"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldID:     group.GetID(),
		FieldName:   group.GetName(),
		FieldHandle: group.GetHandle(),
		FieldStatus: group.GetStatus(),
		FieldMemo:   group.GetMemo(),
	}))
}
