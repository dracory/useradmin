package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/api"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/userstore"
)

func (u *ui) handleUserCreateAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}

	var reqBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if strings.TrimSpace(reqBody.FirstName) == "" {
		api.Respond(w, r, api.Error("First name is required"))
		return ""
	}
	if strings.TrimSpace(reqBody.LastName) == "" {
		api.Respond(w, r, api.Error("Last name is required"))
		return ""
	}
	email := strings.TrimSpace(reqBody.Email)
	if email == "" {
		api.Respond(w, r, api.Error("Email is required"))
		return ""
	}
	if !govalidator.IsEmail(email) {
		api.Respond(w, r, api.Error("Invalid email address"))
		return ""
	}

	// Check email uniqueness before creating. When vault tokenization
	// is enabled, the blind index email store is the only way to detect
	// duplicates (the userstore holds tokens, not plaintext). When vault
	// is disabled, query the userstore directly.
	if u.VaultTokenizer() != nil && u.BlindIndexEmail() != nil {
		ids, err := u.BlindIndexEmail().Search(r.Context(), email, blindindexstore.SEARCH_TYPE_EQUALS)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUserCreateAjax blind index email", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to verify email uniqueness"))
			return ""
		}
		if len(ids) > 0 {
			api.Respond(w, r, api.Error("A user with this email already exists"))
			return ""
		}
	} else {
		query := userstore.NewUserQuery().SetEmail(email)
		existing, err := u.UserStore().UserList(r.Context(), query)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUserCreateAjax UserList", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to verify email uniqueness"))
			return ""
		}
		if len(existing) > 0 {
			api.Respond(w, r, api.Error("A user with this email already exists"))
			return ""
		}
	}

	user := userstore.NewUser()
	user.SetFirstName(strings.TrimSpace(reqBody.FirstName))
	user.SetLastName(strings.TrimSpace(reqBody.LastName))
	user.SetEmail(email)

	if err := u.UserStore().UserCreate(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUserCreateAjax", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to create user"))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData("User created successfully", map[string]interface{}{FieldUserID: user.GetID()}))
	return ""
}
