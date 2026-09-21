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
		UserID       string   `json:"user_id"`
		Status       string   `json:"status"`
		Role         string   `json:"role"`
		FirstName    string   `json:"first_name"`
		LastName     string   `json:"last_name"`
		Email        string   `json:"email"`
		BusinessName string   `json:"business_name"`
		Phone        string   `json:"phone"`
		Country      string   `json:"country"`
		Timezone     string   `json:"timezone"`
		Memo         string   `json:"memo"`
		RoleIDs      []string `json:"role_ids"`
		GroupIDs     []string `json:"group_ids"`
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
	if strings.TrimSpace(payload.Email) == "" {
		api.Respond(w, r, api.Error("Email is required"))
		return
	}
	if !govalidator.IsEmail(strings.TrimSpace(payload.Email)) {
		api.Respond(w, r, api.Error("Invalid email address"))
		return
	}
	// First name, last name, country, and timezone are optional for admins
	// who often do not have all user details on hand.
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

	if !u.syncRoleAssignments(w, r, user.GetID(), payload.RoleIDs) {
		return
	}
	if !u.syncGroupAssignments(w, r, user.GetID(), payload.GroupIDs) {
		return
	}

	// After a successful update, emit an event so the host can react
	// (e.g. enqueue a blind index rebuild, audit log, notifications).
	if u.OnUserUpdated() != nil {
		u.OnUserUpdated()(r.Context(), user.GetID())
	}

	api.Respond(w, r, api.Success("User saved successfully"))
}

// syncRoleAssignments reconciles the user's role assignments with the
// given role IDs: creates missing assignments and soft-deletes removed
// ones. Returns false if an error response was already sent.
func (u *ui) syncRoleAssignments(w http.ResponseWriter, r *http.Request, userID string, roleIDs []string) bool {
	// Roles are an optional userstore feature — when the tables are
	// not configured, skip syncing entirely.
	if u.UserStore().GetRoleTableName() == "" || u.UserStore().GetUserRoleTableName() == "" {
		return true
	}

	existing, err := u.UserStore().UserRoleList(r.Context(), userstore.NewUserRoleQuery().SetUserID(userID))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.syncRoleAssignments UserRoleList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load user roles"))
		return false
	}

	wanted := map[string]bool{}
	for _, id := range roleIDs {
		wanted[id] = true
	}

	existingByRoleID := map[string]userstore.UserRoleInterface{}
	for _, userRole := range existing {
		existingByRoleID[userRole.GetRoleID()] = userRole
		if !wanted[userRole.GetRoleID()] {
			if err := u.UserStore().UserRoleSoftDelete(r.Context(), userRole); err != nil {
				if u.Logger() != nil {
					u.Logger().Error("userUpdateController.syncRoleAssignments UserRoleSoftDelete", slog.String("error", err.Error()))
				}
				api.Respond(w, r, api.Error("Failed to remove role assignment"))
				return false
			}
		}
	}

	for _, roleID := range roleIDs {
		if _, ok := existingByRoleID[roleID]; ok {
			continue
		}
		if _, err := u.UserStore().UserRoleFindByUserIDAndRoleIDOrCreate(r.Context(), userID, roleID); err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.syncRoleAssignments UserRoleFindByUserIDAndRoleIDOrCreate", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to assign role"))
			return false
		}
	}

	return true
}

// syncGroupAssignments reconciles the user's group memberships with the
// given group IDs: creates missing memberships and soft-deletes removed
// ones. Returns false if an error response was already sent.
func (u *ui) syncGroupAssignments(w http.ResponseWriter, r *http.Request, userID string, groupIDs []string) bool {
	// Groups are an optional userstore feature — when the tables are
	// not configured, skip syncing entirely.
	if u.UserStore().GetGroupTableName() == "" || u.UserStore().GetUserGroupTableName() == "" {
		return true
	}

	existing, err := u.UserStore().UserGroupList(r.Context(), userstore.NewUserGroupQuery().SetUserID(userID))
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.syncGroupAssignments UserGroupList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load user groups"))
		return false
	}

	wanted := map[string]bool{}
	for _, id := range groupIDs {
		wanted[id] = true
	}

	existingByGroupID := map[string]userstore.UserGroupInterface{}
	for _, userGroup := range existing {
		existingByGroupID[userGroup.GetGroupID()] = userGroup
		if !wanted[userGroup.GetGroupID()] {
			if err := u.UserStore().UserGroupSoftDelete(r.Context(), userGroup); err != nil {
				if u.Logger() != nil {
					u.Logger().Error("userUpdateController.syncGroupAssignments UserGroupSoftDelete", slog.String("error", err.Error()))
				}
				api.Respond(w, r, api.Error("Failed to remove group membership"))
				return false
			}
		}
	}

	for _, groupID := range groupIDs {
		if _, ok := existingByGroupID[groupID]; ok {
			continue
		}
		if _, err := u.UserStore().UserGroupFindByUserIDAndGroupIDOrCreate(r.Context(), userID, groupID); err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.syncGroupAssignments UserGroupFindByUserIDAndGroupIDOrCreate", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to assign group"))
			return false
		}
	}

	return true
}
