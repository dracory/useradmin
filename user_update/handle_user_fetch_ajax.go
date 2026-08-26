package user_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/useradmin/shared"
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

	// Decode the user for display (e.g. decrypt tokenized fields).
	// When OnUserDecode is nil, the user is used as-is (plain text).
	if u.OnUserDecode() != nil {
		decoded, err := u.OnUserDecode()(r.Context(), user)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userUpdateController.handleUserFetchAjax OnUserDecode", slog.String("error", err.Error()))
			}
			api.Respond(w, r, api.Error("Failed to decode user"))
			return
		}
		user = decoded
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
	}))
}
