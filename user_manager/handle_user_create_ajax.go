package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/api"
	"github.com/dracory/useradmin/shared"
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

	// First and last name are optional for admins who often do not
	// have all user details on hand. Email is the only required field.
	email := strings.TrimSpace(reqBody.Email)
	if email == "" {
		api.Respond(w, r, api.Error("Email is required"))
		return ""
	}
	if !govalidator.IsEmail(email) {
		api.Respond(w, r, api.Error("Invalid email address"))
		return ""
	}

	// Check email uniqueness before creating. When OnUserSearch is
	// provided, use it (e.g. blind index when vault tokenization is
	// enabled). Otherwise, query the userstore directly.
	if u.OnUserSearch() != nil {
		ids, err := u.OnUserSearch()(r.Context(), []shared.SearchCondition{
			{Field: shared.SearchFieldEmail, Op: shared.SearchOpEquals, Value: email},
		})
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUserCreateAjax OnUserSearch", slog.String("error", err.Error()))
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
