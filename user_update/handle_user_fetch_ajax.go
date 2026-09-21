package user_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/useradmin/shared"
	"github.com/dracory/userstore"
)

func (u *ui) handleUserFetchAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return
	}
	userID := req.GetStringTrimmed(r, "user_id")
	if userID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, err := u.UserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("handleUserFetchAjax UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading user"))
		return
	}
	if user == nil {
		if u.Logger() != nil {
			u.Logger().Error("handleUserFetchAjax user not found", slog.String("user_id", userID))
		}
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	// Unseal PII for display (e.g. detokenize, decrypt fields).
	// When UserPiiUnseal is nil, the user is used as-is (plain text).
	if u.UserPiiUnseal() != nil {
		unsealed, err := u.UserPiiUnseal()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax UserPiiUnseal", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to unseal user"))
			return
		}
		user = unsealed
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	email := user.GetEmail()
	phone := user.GetPhone()
	business := user.GetBusinessName()
	memo := user.GetMemo()
	status := user.GetStatus()
	role := user.GetRole()
	country := user.GetCountry()
	timezone := user.GetTimezone()

	fieldStatus := map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"email":         true,
		"business_name": true,
		"phone":         true,
		"role":          true,
	}

	if u.GeoResolver() == nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleUserFetchAjax GeoResolver not configured")
		}
		api.Respond(w, r, api.Error("GeoResolver is not configured"))
		return
	}

	countryList, err := u.GeoResolver().Countries(r.Context())
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleUserFetchAjax Countries", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load countries"))
		return
	}
	countries := make([]map[string]string, 0, len(countryList))
	for _, c := range countryList {
		countries = append(countries, map[string]string{
			FieldIsoCode2: c.IsoCode2,
			FieldName:     c.Name,
		})
	}

	var timezoneList []shared.Timezone
	if country != "" {
		timezoneList, err = u.GeoResolver().Timezones(r.Context(), country)
	}
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleUserFetchAjax Timezones", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load timezones"))
		return
	}
	timezones := make([]map[string]string, 0, len(timezoneList))
	for _, tz := range timezoneList {
		timezones = append(timezones, map[string]string{
			FieldTimezone: tz.Code,
		})
	}

	// All available roles and the IDs of the roles assigned to this
	// user. Roles are an optional userstore feature — when the role
	// tables are not configured (empty table names), the lists are
	// returned empty and the checkboxes are hidden in the UI.
	roles := []map[string]string{}
	userRoleIDs := []string{}
	if u.UserStore().GetRoleTableName() != "" && u.UserStore().GetUserRoleTableName() != "" {
		roleList, err := u.UserStore().RoleList(r.Context(), userstore.NewRoleQuery().
			SetOrderBy(userstore.COLUMN_NAME).
			SetSortDirection(userstore.SORT_ORDER_ASC))
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax RoleList", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to load roles"))
			return
		}
		for _, rl := range roleList {
			roles = append(roles, map[string]string{
				FieldID:     rl.GetID(),
				FieldName:   rl.GetName(),
				FieldHandle: rl.GetHandle(),
			})
		}

		userRoleList, err := u.UserStore().UserRoles(r.Context(), userID)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax UserRoles", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to load user roles"))
			return
		}
		for _, rl := range userRoleList {
			userRoleIDs = append(userRoleIDs, rl.GetID())
		}
	}

	// All available groups and the IDs of the groups this user belongs
	// to. Groups are an optional userstore feature — when the group
	// tables are not configured (empty table names), the lists are
	// returned empty and the checkboxes are hidden in the UI.
	groups := []map[string]string{}
	userGroupIDs := []string{}
	if u.UserStore().GetGroupTableName() != "" && u.UserStore().GetUserGroupTableName() != "" {
		groupList, err := u.UserStore().GroupList(r.Context(), userstore.NewGroupQuery().
			SetOrderBy(userstore.COLUMN_NAME).
			SetSortDirection(userstore.SORT_ORDER_ASC))
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax GroupList", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to load groups"))
			return
		}
		for _, g := range groupList {
			groups = append(groups, map[string]string{
				FieldID:     g.GetID(),
				FieldName:   g.GetName(),
				FieldHandle: g.GetHandle(),
			})
		}

		userGroupList, err := u.UserStore().UserGroups(r.Context(), userID)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax UserGroups", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to load user groups"))
			return
		}
		for _, g := range userGroupList {
			userGroupIDs = append(userGroupIDs, g.GetID())
		}
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldStatus:       status,
		FieldRole:         role,
		FieldFirstName:    firstName,
		FieldLastName:     lastName,
		FieldEmail:        email,
		FieldBusinessName: business,
		FieldPhone:        phone,
		FieldCountry:      country,
		FieldTimezone:     timezone,
		FieldMemo:         memo,
		FieldStatusField:  fieldStatus,
		FieldCountries:    countries,
		FieldTimezones:    timezones,
		FieldRoles:        roles,
		FieldUserRoleIDs:  userRoleIDs,
		FieldGroups:       groups,
		FieldUserGroupIDs: userGroupIDs,
	}))
}
