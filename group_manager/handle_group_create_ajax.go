package group_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleGroupCreateAjax(w http.ResponseWriter, r *http.Request) string {
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

	existing, err := u.UserStore().GroupFindByHandle(r.Context(), handle)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupManagerController.handleGroupCreateAjax GroupFindByHandle", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to check group handle"))
		return ""
	}
	if existing != nil {
		api.Respond(w, r, api.Error("A group with this handle already exists"))
		return ""
	}

	group := userstore.NewGroup().
		SetName(name).
		SetHandle(handle).
		SetStatus(userstore.GROUP_STATUS_ACTIVE)

	if err := u.UserStore().GroupCreate(r.Context(), group); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupManagerController.handleGroupCreateAjax GroupCreate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to create group"))
		return ""
	}

	api.Respond(w, r, api.Success("Group created successfully"))
	return ""
}
