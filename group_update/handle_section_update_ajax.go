package group_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

// loadGroupForUpdate fetches a group by ID.
// Returns nil when an error response was already sent.
func (u *ui) loadGroupForUpdate(w http.ResponseWriter, r *http.Request, groupID string) userstore.GroupInterface {
	group, err := u.UserStore().GroupFindByID(r.Context(), groupID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupUpdateController loadGroupForUpdate GroupFindByID", slog.String("group_id", groupID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading group"))
		return nil
	}
	if group == nil {
		api.Respond(w, r, api.Error("Group not found"))
		return nil
	}

	return group
}

// persistGroup saves the group.
// Returns false when an error response was already sent.
func (u *ui) persistGroup(w http.ResponseWriter, r *http.Request, group userstore.GroupInterface) bool {
	if err := u.UserStore().GroupUpdate(r.Context(), group); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("groupUpdateController persistGroup GroupUpdate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving group failed"))
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
	if !u.groupsEnabled() {
		api.Respond(w, r, api.Error("Groups are not enabled in the user store"))
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return false
	}
	return true
}

// handleStatusUpdateAjax updates only the group's status.
func (u *ui) handleStatusUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		GroupID string `json:"group_id"`
		Status  string `json:"status"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
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

	group := u.loadGroupForUpdate(w, r, payload.GroupID)
	if group == nil {
		return
	}

	group.SetStatus(payload.Status)
	if !u.persistGroup(w, r, group) {
		return
	}

	api.Respond(w, r, api.Success("Status saved successfully"))
}

// handleDetailsUpdateAjax updates only the group's name and handle.
func (u *ui) handleDetailsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		GroupID string `json:"group_id"`
		Name    string `json:"name"`
		Handle  string `json:"handle"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
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

	group := u.loadGroupForUpdate(w, r, payload.GroupID)
	if group == nil {
		return
	}

	if handle != group.GetHandle() {
		existing, err := u.UserStore().GroupFindByHandle(r.Context(), handle)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("groupUpdateController.handleDetailsUpdateAjax GroupFindByHandle", slog.String("error", err.Error()))
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
	if !u.persistGroup(w, r, group) {
		return
	}

	api.Respond(w, r, api.Success("Group details saved successfully"))
}

// handleNotesUpdateAjax updates only the group's admin notes (memo).
func (u *ui) handleNotesUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		GroupID string `json:"group_id"`
		Memo    string `json:"memo"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.GroupID == "" {
		api.Respond(w, r, api.Error("Group ID is required"))
		return
	}

	group := u.loadGroupForUpdate(w, r, payload.GroupID)
	if group == nil {
		return
	}

	group.SetMemo(strings.TrimSpace(payload.Memo))
	if !u.persistGroup(w, r, group) {
		return
	}

	api.Respond(w, r, api.Success("Notes saved successfully"))
}
