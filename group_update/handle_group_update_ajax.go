package group_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleGroupUpdateAjax(w http.ResponseWriter, r *http.Request) {
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

	var payload struct {
		GroupID string `json:"group_id"`
		Name    string `json:"name"`
		Handle  string `json:"handle"`
		Status  string `json:"status"`
		Memo    string `json:"memo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return
	}

	group, err := u.UserStore().GroupFindByID(r.Context(), payload.GroupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupUpdateController.handleGroupUpdateAjax GroupFindByID", slog.String("group_id", payload.GroupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading group"))
		return
	}
	if group == nil {
		api.Respond(w, r, api.Error("Group not found"))
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
	if payload.Status != userstore.GROUP_STATUS_ACTIVE && payload.Status != userstore.GROUP_STATUS_INACTIVE && payload.Status != userstore.GROUP_STATUS_DELETED {
		api.Respond(w, r, api.Error("Invalid status value"))
		return
	}

	if handle != group.GetHandle() {
		existing, err := u.UserStore().GroupFindByHandle(r.Context(), handle)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupUpdateController.handleGroupUpdateAjax GroupFindByHandle", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to check group handle"))
			return
		}
		if existing != nil && existing.GetID() != group.GetID() {
			api.Respond(w, r, api.Error("A group with this handle already exists"))
			return
		}
	}

	group.SetName(name)
	group.SetHandle(handle)
	group.SetStatus(payload.Status)
	group.SetMemo(strings.TrimSpace(payload.Memo))

	if err := u.UserStore().GroupUpdate(r.Context(), group); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupUpdateController.handleGroupUpdateAjax GroupUpdate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving group failed"))
		return
	}

	api.Respond(w, r, api.Success("Group saved successfully"))
}
