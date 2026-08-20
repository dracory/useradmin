package user_update

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/api"
	"github.com/dracory/taskstore"
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

	originalEmail := user.GetEmail()
	if u.VaultTokenizer() != nil {
		_, _, em, _, _, err := u.VaultTokenizer().Untokenize(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserUpdateAjax Untokenize originalEmail", slog.String("error", err.Error()))
			}
			// On untokenize failure, skip the email-change detection
			// rather than enqueuing a spurious blind index rebuild.
			originalEmail = strings.TrimSpace(payload.Email)
		} else {
			originalEmail = em
		}
	}

	user.SetMemo(strings.TrimSpace(payload.Memo))
	user.SetStatus(payload.Status)
	user.SetRole(payload.Role)
	user.SetCountry(payload.Country)
	user.SetTimezone(payload.Timezone)

	if u.VaultTokenizer() != nil {
		firstToken, lastToken, emailToken, phoneToken, businessToken, err := u.VaultTokenizer().Tokenize(
			r.Context(),
			user,
			strings.TrimSpace(payload.FirstName),
			strings.TrimSpace(payload.LastName),
			strings.TrimSpace(payload.Email),
			strings.TrimSpace(payload.Phone),
			strings.TrimSpace(payload.BusinessName),
		)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("Error tokenizing user", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("System error. Saving user failed"))
			return
		}
		user.SetFirstName(firstToken)
		user.SetLastName(lastToken)
		user.SetEmail(emailToken)
		user.SetPhone(phoneToken)
		user.SetBusinessName(businessToken)
	} else {
		user.SetFirstName(strings.TrimSpace(payload.FirstName))
		user.SetLastName(strings.TrimSpace(payload.LastName))
		user.SetEmail(strings.TrimSpace(payload.Email))
		user.SetPhone(strings.TrimSpace(payload.Phone))
		user.SetBusinessName(strings.TrimSpace(payload.BusinessName))
	}

	if err := u.UserStore().UserUpdate(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("Error updating user", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving user failed"))
		return
	}

	// When vault tokenization is enabled and the email changed, enqueue
	// a blind index rebuild task so search stays consistent.
	if u.VaultTokenizer() != nil && u.TaskStore() != nil && u.BlindIndexRebuildTaskAlias() != "" {
		if originalEmail != strings.TrimSpace(payload.Email) {
			_, err := u.TaskStore().TaskDefinitionEnqueueByAlias(
				r.Context(),
				taskstore.DefaultQueueName,
				u.BlindIndexRebuildTaskAlias(),
				map[string]any{
					"index":    "email",
					"truncate": "no",
				},
			)
			if err != nil {
				if u.Logger() != nil {
					u.Logger().Error("Error enqueuing blind index rebuild", slog.String("error", err.Error()))
				}
			}
		}
	}

	api.Respond(w, r, api.Success("User saved successfully"))
}
