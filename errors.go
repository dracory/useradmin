package useradmin

import "errors"

// Common errors
var (
	// ErrUserStoreRequired is returned when UserStore is not provided
	ErrUserStoreRequired = errors.New("user store is required")

	// ErrLoggerRequired is returned when Logger is not provided
	ErrLoggerRequired = errors.New("logger is required")

	// ErrSessionStoreRequired is returned when SessionStore is not
	// provided. The impersonate controller needs it to create a new
	// session for the impersonated user.
	ErrSessionStoreRequired = errors.New("session store is required for the impersonate controller")

	// ErrGeoStoreRequired is returned when GeoStore is not provided.
	// The user update controller needs it to list countries and
	// timezones.
	ErrGeoStoreRequired = errors.New("geo store is required for the user update controller")
)
