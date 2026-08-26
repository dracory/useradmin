package shared

import (
	"context"
	"net/http"
)

// GeoResolverInterface provides countries and timezones. The host
// implements this against whatever geo data source it uses
// (dracory/geostore, a static list, an external API, etc.).
type GeoResolverInterface interface {
	// Countries returns all countries.
	Countries(ctx context.Context) ([]Country, error)

	// Timezones returns timezones for the given country code. The
	// country code is optional — pass no argument or an empty string
	// when no country is selected; implementations should return an
	// empty list in that case.
	Timezones(ctx context.Context, countryCode ...string) ([]Timezone, error)
}

// Country is a single country entry returned by GeoResolverInterface.
type Country struct {
	IsoCode2 string
	Name     string
}

// Timezone is a single timezone entry returned by GeoResolverInterface.
type Timezone struct {
	Code string
}

// UserSearchEvent is passed to OnUserSearch callbacks when the user
// list is filtered or when an email uniqueness check is performed.
type UserSearchEvent struct {
	// FirstName filters by first name. Empty means no filter.
	FirstName string
	// LastName filters by last name. Empty means no filter.
	LastName string
	// Email filters by email. Empty means no filter.
	Email string
	// ExactMatch when true means exact match (used for uniqueness
	// checks). When false means substring match (used for list
	// filtering).
	ExactMatch bool
}

// OnUserSearchFunc is an optional callback for custom user search.
// The host can use it to search via blind index, Elasticsearch, or any
// other backend. When nil, useradmin falls back to userstore
// query-based search (SetFirstNameLike, SetLastNameLike, SetEmailLike
// for substring; SetEmail for exact match).
type OnUserSearchFunc func(ctx context.Context, event UserSearchEvent) ([]string, error)

// SessionResolverInterface creates sessions for impersonation. The host
// owns the session store, cookie format, and expiry policy.
type SessionResolverInterface interface {
	// Create creates a new session for the given user ID and sets the
	// auth cookie on the response. The secure flag controls whether
	// the cookie is marked Secure (false for HTTP development, true
	// for HTTPS production).
	Create(w http.ResponseWriter, r *http.Request, userID string, secure bool) error
}

// UserUpdateEvent is passed to OnUserUpdate callbacks after a user is
// updated. It carries the information the host may need to react —
// e.g. enqueuing a blind index rebuild when the email changed.
type UserUpdateEvent struct {
	// UserID is the ID of the updated user.
	UserID string
	// OriginalEmail is the user's email before the update.
	OriginalEmail string
	// NewEmail is the user's email after the update.
	NewEmail string
}

// OnUserUpdateFunc is an optional callback invoked after a user is
// updated. The host can use it to trigger side effects (blind index
// rebuild, audit log, notifications, etc.) without useradmin dictating
// how. When nil, the callback is skipped.
type OnUserUpdateFunc func(ctx context.Context, event UserUpdateEvent)
