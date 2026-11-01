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

// loadUserForUpdate fetches a user by ID and unseals PII for editing.
// Returns nil when an error response was already sent.
func (u *ui) loadUserForUpdate(w http.ResponseWriter, r *http.Request, userID string) userstore.UserInterface {
	user, err := u.UserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController loadUserForUpdate UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading user"))
		return nil
	}
	if user == nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController loadUserForUpdate user not found", slog.String("user_id", userID))
		}
		api.Respond(w, r, api.Error("User not found"))
		return nil
	}

	if u.UserPiiUnseal() != nil {
		unsealed, err := u.UserPiiUnseal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController loadUserForUpdate UserPiiUnseal", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("System error. Saving user failed"))
			return nil
		}
		user = unsealed
	}

	return user
}

// persistUser seals PII and saves the user, then fires OnUserUpdated.
// Returns false when an error response was already sent.
func (u *ui) persistUser(w http.ResponseWriter, r *http.Request, user userstore.UserInterface) bool {
	if u.UserPiiSeal() != nil {
		sealed, err := u.UserPiiSeal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController persistUser UserPiiSeal", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("System error. Saving user failed"))
			return false
		}
		user = sealed
	}

	if err := u.UserStore().UserUpdate(r.Context(), user); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController persistUser UserUpdate", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("System error. Saving user failed"))
		return false
	}

	if u.OnUserUpdated() != nil {
		u.OnUserUpdated()(r.Context(), user.GetID())
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
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return false
	}
	return true
}

// handleStatusUpdateAjax updates only the user's status.
func (u *ui) handleStatusUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID string `json:"user_id"`
		Status string `json:"status"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}
	if payload.Status == "" {
		api.Respond(w, r, api.Error("Status is required"))
		return
	}

	user := u.loadUserForUpdate(w, r, payload.UserID)
	if user == nil {
		return
	}

	user.SetStatus(payload.Status)
	if !u.persistUser(w, r, user) {
		return
	}

	api.Respond(w, r, api.Success("Status saved successfully"))
}

// handleDetailsUpdateAjax updates only the user's profile fields
// (name, email, business name, phone).
func (u *ui) handleDetailsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID       string `json:"user_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		BusinessName string `json:"business_name"`
		Phone        string `json:"phone"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
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

	user := u.loadUserForUpdate(w, r, payload.UserID)
	if user == nil {
		return
	}

	user.SetFirstName(strings.TrimSpace(payload.FirstName))
	user.SetLastName(strings.TrimSpace(payload.LastName))
	user.SetEmail(strings.TrimSpace(payload.Email))
	user.SetBusinessName(strings.TrimSpace(payload.BusinessName))
	user.SetPhone(strings.TrimSpace(payload.Phone))
	if !u.persistUser(w, r, user) {
		return
	}

	api.Respond(w, r, api.Success("User details saved successfully"))
}

// handleRolesUpdateAjax updates only the user's role, role assignments,
// and group memberships.
func (u *ui) handleRolesUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID   string   `json:"user_id"`
		Role     string   `json:"role"`
		RoleIDs  []string `json:"role_ids"`
		GroupIDs []string `json:"group_ids"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}
	if payload.Role != "" && payload.Role != userstore.USER_ROLE_USER && payload.Role != userstore.USER_ROLE_ADMINISTRATOR {
		api.Respond(w, r, api.Error("Invalid role value"))
		return
	}

	user := u.loadUserForUpdate(w, r, payload.UserID)
	if user == nil {
		return
	}

	user.SetRole(payload.Role)
	if !u.persistUser(w, r, user) {
		return
	}
	if !u.syncRoleAssignments(w, r, user.GetID(), payload.RoleIDs) {
		return
	}
	if !u.syncGroupAssignments(w, r, user.GetID(), payload.GroupIDs) {
		return
	}

	api.Respond(w, r, api.Success("Roles saved successfully"))
}

// handleLocationUpdateAjax updates only the user's country and timezone.
func (u *ui) handleLocationUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID   string `json:"user_id"`
		Country  string `json:"country"`
		Timezone string `json:"timezone"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user := u.loadUserForUpdate(w, r, payload.UserID)
	if user == nil {
		return
	}

	user.SetCountry(payload.Country)
	user.SetTimezone(payload.Timezone)
	if !u.persistUser(w, r, user) {
		return
	}

	api.Respond(w, r, api.Success("Location saved successfully"))
}

// handleNotesUpdateAjax updates only the user's admin notes (memo).
func (u *ui) handleNotesUpdateAjax(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID string `json:"user_id"`
		Memo   string `json:"memo"`
	}
	if !u.decodeSectionPayload(w, r, &payload) {
		return
	}
	if payload.UserID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user := u.loadUserForUpdate(w, r, payload.UserID)
	if user == nil {
		return
	}

	user.SetMemo(strings.TrimSpace(payload.Memo))
	if !u.persistUser(w, r, user) {
		return
	}

	api.Respond(w, r, api.Success("Notes saved successfully"))
}
