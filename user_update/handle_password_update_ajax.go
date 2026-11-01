package user_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
)

// handlePasswordUpdateAjax changes a user's password. It is kept
// separate from the profile update so an admin can reset a password
// without touching the (possibly tokenized) PII fields.
func (u *ui) handlePasswordUpdateAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return
	}

	var payload struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}
	if len(strings.TrimSpace(payload.Password)) < 8 {
		api.Respond(w, r, api.Error("Password must be at least 8 characters"))
		return
	}

	user, err := u.UserStore().UserFindByID(r.Context(), payload.UserID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handlePasswordUpdateAjax UserFindByID", slog.String("user_id", payload.UserID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading user"))
		return
	}
	if user == nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handlePasswordUpdateAjax user not found", slog.String("user_id", payload.UserID))
		}
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	if err := user.SetPasswordAndHash(payload.Password); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handlePasswordUpdateAjax SetPasswordAndHash", slog.String("user_id", payload.UserID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Changing password failed"))
		return
	}

	if err := u.UserStore().UserUpdate(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handlePasswordUpdateAjax UserUpdate", slog.String("user_id", payload.UserID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Changing password failed"))
		return
	}

	if u.OnUserUpdated() != nil {
		u.OnUserUpdated()(r.Context(), user.GetID())
	}

	api.Respond(w, r, api.Success("Password changed successfully"))
}
