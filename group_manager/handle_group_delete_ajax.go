package group_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleGroupDeleteAjax(w http.ResponseWriter, r *http.Request) string {
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

	var payload struct {
		GroupID string `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return ""
	}

	if err := u.UserStore().GroupSoftDeleteByID(r.Context(), payload.GroupID); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupManagerController.handleGroupDeleteAjax GroupSoftDeleteByID", slog.String("group_id", payload.GroupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to delete group"))
		return ""
	}

	api.Respond(w, r, api.Success("Group deleted successfully"))
	return ""
}
