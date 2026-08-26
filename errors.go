package useradmin

import "errors"

// Common errors
var (
	// ErrUserStoreRequired is returned when UserStore is not provided
	ErrUserStoreRequired = errors.New("user store is required")

	// ErrLoggerRequired is returned when Logger is not provided
	ErrLoggerRequired = errors.New("logger is required")

	// ErrGeoResolverRequired is returned when GeoResolver is not
	// provided. The user update controller needs it to list countries
	// and timezones.
	ErrGeoResolverRequired = errors.New("geo resolver is required for the user update controller")
)
