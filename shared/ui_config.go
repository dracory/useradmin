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
// search. OnUserUpdate is optional — when nil, the callback is skipped.
// VaultTokenizer is optional — when nil, user fields are treated as
// plain text.
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

	// OnUserUpdate is an optional callback invoked after a user is
	// updated. The host can use it to trigger side effects (blind
	// index rebuild, audit log, notifications, etc.). When nil, the
	// callback is skipped.
	OnUserUpdate OnUserUpdateFunc

	// VaultTokenizer abstracts vault tokenization. Optional — when
	// nil, user fields are treated as plain text.
	VaultTokenizer VaultTokenizer

	// FlashRedirect redirects with a flash message. Optional — when
	// nil, plain http.Redirect is used.
	FlashRedirect FlashRedirectFunc

	// SecureCookie controls whether the impersonation cookie is marked
	// Secure. Set to false for HTTP (development), true for HTTPS
	// (production). Defaults to true.
	SecureCookie bool

	// Layout is the layout renderer callback.
	Layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
