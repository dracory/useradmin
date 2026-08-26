package shared

import (
	"context"
	"net/http"

	"github.com/dracory/userstore"
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

// SearchOp defines the comparison operator for a SearchCondition.
type SearchOp string

const (
	SearchOpEquals      SearchOp = "eq"
	SearchOpContains    SearchOp = "contains"
	SearchOpNotContains SearchOp = "not_contains"
	SearchOpStartsWith  SearchOp = "starts_with"
)

// SearchField constants identify the user fields that can be filtered.
type SearchField string

const (
	SearchFieldFirstName    SearchField = "first_name"
	SearchFieldLastName     SearchField = "last_name"
	SearchFieldEmail        SearchField = "email"
	SearchFieldPhone        SearchField = "phone"
	SearchFieldBusinessName SearchField = "business_name"
)

// SearchCombine defines how a condition combines with the previous one.
// The first condition's CombineWith is ignored.
type SearchCombine string

const (
	SearchAnd SearchCombine = "AND"
	SearchOr  SearchCombine = "OR"
)

// SearchCondition is a single filter criterion passed to OnUserSearch.
type SearchCondition struct {
	// Field is the user field to filter on.
	Field SearchField
	// Op is the comparison operator.
	Op SearchOp
	// Value is the value to compare against.
	Value string
	// CombineWith specifies how this condition combines with the
	// previous one. Defaults to AND. Ignored for the first condition.
	CombineWith SearchCombine
}

// OnUserSearchFunc is an optional callback for custom user search.
// The host receives a list of conditions and applies whatever logic
// it wants (blind index, Elasticsearch, etc.). When nil, useradmin
// falls back to userstore query-based search.
type OnUserSearchFunc func(ctx context.Context, conditions []SearchCondition) ([]string, error)

// OnUserImpersonateFunc is an optional callback invoked when an admin
// impersonates a user. The host owns the auth mechanism — it can
// create a session record and set a cookie, issue a JWT, or anything
// else. When nil, impersonation is disabled.
type OnUserImpersonateFunc func(w http.ResponseWriter, r *http.Request, userID string) error

// OnUserUpdatedFunc is an optional callback invoked after a user is
// updated. The host can load the user by ID and react to whatever
// changed (blind index rebuild, audit log, notifications, etc.).
// When nil, the callback is skipped.
type OnUserUpdatedFunc func(ctx context.Context, userID string)

// UserPiiSealFunc transforms a user from display representation to
// storage representation (e.g. tokenize, encrypt, mask PII fields).
// The host owns the mechanism. When nil, the user is stored as-is
// (plain text).
type UserPiiSealFunc func(ctx context.Context, user userstore.UserInterface) (userstore.UserInterface, error)

// UserPiiUnsealFunc transforms a user from storage representation to
// display representation (e.g. detokenize, decrypt, reveal PII fields).
// The host owns the mechanism. When nil, the user is used as-is
// (plain text).
type UserPiiUnsealFunc func(ctx context.Context, user userstore.UserInterface) (userstore.UserInterface, error)

// UsersPiiUnsealFunc is the batch version of UserPiiUnsealFunc. It
// allows the host to unseal all users in a single vault batch call
// for efficiency. When nil, useradmin falls back to calling
// UserPiiUnsealFunc per user (or plain text when that is also nil).
type UsersPiiUnsealFunc func(ctx context.Context, users []userstore.UserInterface) ([]userstore.UserInterface, error)
