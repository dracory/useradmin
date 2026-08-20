package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/geostore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

// UiConfig holds the dependencies passed to subcontroller UI factories.
// This follows the blogadmin/shopadmin pattern.
//
// UserStore, GeoStore, and Logger are required for core controllers.
// SessionStore is required for the impersonate controller. The blind
// index stores are optional — when nil, the corresponding search filter
// is disabled. VaultTokenizer is optional — when nil, user fields are
// treated as plain text. TaskStore is optional — when nil, blind index
// rebuild enqueue on email change is skipped.
type UiConfig struct {
	UserStore userstore.StoreInterface
	GeoStore  geostore.StoreInterface
	Logger    *slog.Logger

	// SessionStore is required for the impersonate controller.
	SessionStore sessionstore.StoreInterface

	// BlindIndexFirstName/LastName/Email enable filtered search by
	// the corresponding field. Optional.
	BlindIndexFirstName blindindexstore.StoreInterface
	BlindIndexLastName  blindindexstore.StoreInterface
	BlindIndexEmail     blindindexstore.StoreInterface

	// TaskStore is used to enqueue a blind index rebuild when a user's
	// email changes and vault tokenization is enabled. Optional.
	TaskStore taskstore.StoreInterface

	// BlindIndexRebuildTaskAlias is the task alias enqueued on email
	// change. If empty, the enqueue is skipped.
	BlindIndexRebuildTaskAlias string

	// VaultTokenizer abstracts vault tokenization. Optional — when
	// nil, user fields are treated as plain text.
	VaultTokenizer VaultTokenizer

	// AuthUser returns the authenticated user from the request, or
	// nil if unauthenticated. Used by the create/delete/impersonate
	// controllers for authorization checks.
	AuthUser func(r *http.Request) userstore.UserInterface

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
