package user_update_test

import (
	"net/url"
	"testing"

	"github.com/dracory/useradmin/shared"
)

// TestTimezonesFetchReturnsList verifies the timezones endpoint proxies
// the geo resolver for the given country.
func TestTimezonesFetchReturnsList(t *testing.T) {
	env := newTestEnv(t, &stubGeo{timezones: []shared.Timezone{{Code: "Europe/London"}}})
	defer env.close()

	resp := postAction(t, env, "get-timezones-ajax", url.Values{"country_code": {"GB"}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}
	list, ok := resp.Data["timezones"].([]interface{})
	if !ok || len(list) != 1 {
		t.Fatalf("expected 1 timezone, got %v", resp.Data["timezones"])
	}
	if list[0].(map[string]interface{})["timezone"] != "Europe/London" {
		t.Errorf("expected Europe/London, got %v", list[0])
	}
}

// TestTimezonesFetchRequiresCountryCode verifies the country code is
// required.
func TestTimezonesFetchRequiresCountryCode(t *testing.T) {
	env := newTestEnv(t, &stubGeo{})
	defer env.close()

	resp := postAction(t, env, "get-timezones-ajax", nil, map[string]any{})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}
