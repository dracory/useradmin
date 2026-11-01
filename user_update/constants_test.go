package user_update

import "testing"

// TestActionConstants verifies the AJAX action names stay stable — the
// frontend JS posts these strings, so renaming one silently breaks a
// section save.
func TestActionConstants(t *testing.T) {
	actions := map[string]string{
		actionUserFetch:      "user-fetch-ajax",
		actionGetTimezones:   "get-timezones-ajax",
		actionUserUpdate:     "user-update-ajax",
		actionPasswordUpdate: "password-update-ajax",
		actionStatusUpdate:   "status-update-ajax",
		actionDetailsUpdate:  "details-update-ajax",
		actionRolesUpdate:    "roles-update-ajax",
		actionLocationUpdate: "location-update-ajax",
		actionNotesUpdate:    "notes-update-ajax",
	}
	for got, want := range actions {
		if got != want {
			t.Errorf("expected action %q, got %q", want, got)
		}
	}
}

// TestFieldConstants verifies a sample of the JSON field names used in
// handler payloads.
func TestFieldConstants(t *testing.T) {
	fields := map[string]string{
		FieldStatus:    "status",
		FieldEmail:     "email",
		FieldTimezone:  "timezone",
		FieldTimezones: "timezones",
	}
	for got, want := range fields {
		if got != want {
			t.Errorf("expected field %q, got %q", want, got)
		}
	}
}
