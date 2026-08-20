package user_update

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/geostore"
	"github.com/dracory/neat"
	"github.com/dracory/req"
)

func (u *ui) handleTimezonesFetchAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	countryCode := req.GetStringTrimmed(r, "country_code")
	if countryCode == "" {
		api.Respond(w, r, api.Error("Country code is required"))
		return
	}

	if u.GeoStore() == nil {
		if u.Logger() != nil {
			u.Logger().Error("userUpdateController.handleTimezonesFetchAjax GeoStore not configured")
		}
		api.Respond(w, r, api.Error("GeoStore is not configured"))
		return
	}

	timezoneList, err := u.GeoStore().TimezoneList(r.Context(), geostore.TimezoneQueryOptions{
		SortOrder:   neat.SortAsc,
		OrderBy:     geostore.COLUMN_TIMEZONE,
		CountryCode: countryCode,
	})
	if err != nil {
		api.Respond(w, r, api.Error("Failed to load timezones"))
		return
	}

	timezones := make([]map[string]string, 0, len(timezoneList))
	for _, tz := range timezoneList {
		timezones = append(timezones, map[string]string{
			FieldTimezone: tz.Timezone(),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldTimezones: timezones,
	}))
}
