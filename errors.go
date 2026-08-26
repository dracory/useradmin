package useradmin

import "errors"

// Common errors
var (
	// ErrUserStoreRequired is returned when UserStore is not provided
	ErrUserStoreRequired = errors.New("user store is required")

	// ErrLoggerRequired is returned when Logger is not provided
	ErrLoggerRequired = errors.New("logger is required")

	// ErrSessionResolverRequired is returned when SessionResolver is
	// not provided. The impersonate controller needs it to create a
	// new session for the impersonated user.
	ErrSessionResolverRequired = errors.New("session resolver is required for the impersonate controller")

	// ErrGeoResolverRequired is returned when GeoResolver is not
	// provided. The user update controller needs it to list countries
	// and timezones.
	ErrGeoResolverRequired = errors.New("geo resolver is required for the user update controller")
)
