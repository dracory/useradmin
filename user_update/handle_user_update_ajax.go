package user_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (u *ui) handleUserUpdateAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return
	}

	var payload struct {
		UserID       string `json:"user_id"`
		Status       string `json:"status"`
		Role         string `json:"role"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		BusinessName string `json:"business_name"`
		Phone        string `json:"phone"`
		Country      string `json:"country"`
		Timezone     string `json:"timezone"`
		Memo         string `json:"memo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return
	}

	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, err := u.UserStore().UserFindByID(r.Context(), payload.UserID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleUserUpdateAjax UserFindByID", slog.String("user_id", payload.UserID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading user"))
		return
	}
	if user == nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleUserUpdateAjax user not found", slog.String("user_id", payload.UserID))
		}
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	if payload.Status == "" {
		api.Respond(w, r, api.Error("Status is required"))
		return
	}
	if strings.TrimSpace(payload.FirstName) == "" {
		api.Respond(w, r, api.Error("First name is required"))
		return
	}
	if strings.TrimSpace(payload.LastName) == "" {
		api.Respond(w, r, api.Error("Last name is required"))
		return
	}
	if strings.TrimSpace(payload.Email) == "" {
		api.Respond(w, r, api.Error("Email is required"))
		return
	}
	if !govalidator.IsEmail(strings.TrimSpace(payload.Email)) {
		api.Respond(w, r, api.Error("Invalid email address"))
		return
	}
	if payload.Country == "" {
		api.Respond(w, r, api.Error("Country is required"))
		return
	}
	if payload.Timezone == "" {
		api.Respond(w, r, api.Error("Timezone is required"))
		return
	}
	if payload.Role != "" && payload.Role != userstore.USER_ROLE_USER && payload.Role != userstore.USER_ROLE_ADMINISTRATOR {
		api.Respond(w, r, api.Error("Invalid role value"))
		return
	}

	// Unseal PII for display (e.g. detokenize, decrypt fields).
	if u.UserPiiUnseal() != nil {
		unsealed, err := u.UserPiiUnseal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserUpdateAjax UserPiiUnseal", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("System error. Saving user failed"))
			return
		}
		user = unsealed
	}

	user.SetMemo(strings.TrimSpace(payload.Memo))
	user.SetStatus(payload.Status)
	user.SetRole(payload.Role)
	user.SetCountry(payload.Country)
	user.SetTimezone(payload.Timezone)
	user.SetFirstName(strings.TrimSpace(payload.FirstName))
	user.SetLastName(strings.TrimSpace(payload.LastName))
	user.SetEmail(strings.TrimSpace(payload.Email))
	user.SetPhone(strings.TrimSpace(payload.Phone))
	user.SetBusinessName(strings.TrimSpace(payload.BusinessName))

	// Seal PII for storage (e.g. tokenize, encrypt fields).
	if u.UserPiiSeal() != nil {
		sealed, err := u.UserPiiSeal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("Error sealing user PII", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("System error. Saving user failed"))
			return
		}
		user = sealed
	}

	if err := u.UserStore().UserUpdate(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("Error updating user", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving user failed"))
		return
	}

	// After a successful update, emit an event so the host can react
	// (e.g. enqueue a blind index rebuild, audit log, notifications).
	if u.OnUserUpdated() != nil {
		u.OnUserUpdated()(r.Context(), user.GetID())
	}

	api.Respond(w, r, api.Success("User saved successfully"))
}
