package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/userstore"
)

// UiConfig holds the dependencies passed to subcontroller UI factories.
// This follows the blogadmin/shopadmin pattern.
//
// UserStore, GeoResolver, and Logger are required for core controllers.
// OnUserImpersonate is optional — when nil, the impersonate button is
// hidden and the impersonate route is not registered. OnUserSearch is
// optional — when nil, useradmin falls back to userstore query-based
// search. OnUserUpdated is optional — when nil, the callback is skipped.
// UserPiiSeal/UserPiiUnseal/UsersPiiUnseal are optional — when nil,
// user fields are treated as plain text.
//
// Authentication and authorization are the host's responsibility —
// gate the routes with middleware before they reach useradmin.
type UiConfig struct {
	UserStore   userstore.StoreInterface
	GeoResolver GeoResolverInterface
	Logger      *slog.Logger

	// OnUserImpersonate is optional — when nil, the impersonate
	// button is hidden and the impersonate route is not registered.
	OnUserImpersonate OnUserImpersonateFunc

	// OnUserSearch is an optional callback for custom user search
	// (e.g. blind index, Elasticsearch). When nil, useradmin falls
	// back to userstore query-based search.
	OnUserSearch OnUserSearchFunc

	// OnUserUpdated is an optional callback invoked after a user is
	// updated. The host can use it to trigger side effects (blind
	// index rebuild, audit log, notifications, etc.). When nil, the
	// callback is skipped.
	OnUserUpdated OnUserUpdatedFunc

	// UserPiiSeal transforms a user from display representation to
	// storage representation (e.g. tokenize, encrypt PII). Optional —
	// when nil, the user is stored as-is (plain text).
	UserPiiSeal UserPiiSealFunc

	// UserPiiUnseal transforms a user from storage representation to
	// display representation (e.g. detokenize, decrypt PII). Optional —
	// when nil, the user is used as-is (plain text).
	UserPiiUnseal UserPiiUnsealFunc

	// UsersPiiUnseal is the batch version of UserPiiUnseal. It allows
	// the host to unseal all users in a single call for efficiency.
	// Optional — when nil, useradmin falls back to UserPiiUnseal per
	// user (or plain text when that is also nil).
	UsersPiiUnseal UsersPiiUnsealFunc

	// FlashRedirect redirects with a flash message. Optional — when
	// nil, plain http.Redirect is used.
	FlashRedirect FlashRedirectFunc

	// Layout is the layout renderer callback.
	Layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
