package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleUserDeleteAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}

	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if reqBody.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return ""
	}

	user, err := u.UserStore().UserFindByID(r.Context(), reqBody.UserID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUserDeleteAjax", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("User not found"))
		return ""
	}
	if user == nil {
		api.Respond(w, r, api.Error("User not found"))
		return ""
	}

	if err := u.UserStore().UserSoftDelete(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUserDeleteAjax", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to delete user"))
		return ""
	}

	api.Respond(w, r, api.Success("User deleted successfully"))
	return ""
}
