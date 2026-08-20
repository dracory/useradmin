package shared

// Context keys for config values injected by Handle()
const KeyEndpoint = "endpoint"
const KeyAdminHomeURL = "admin_home_url"
const KeyUserAdminURL = "user_admin_url"
const KeyUserHomeURL = "user_home_url"

// Controller names used in the ?controller= query parameter
const (
	CONTROLLER_USER_MANAGER     = "user-manager"
	CONTROLLER_USER_CREATE      = "user-create"
	CONTROLLER_USER_DELETE      = "user-delete"
	CONTROLLER_USER_UPDATE      = "user-update"
	CONTROLLER_USER_IMPERSONATE = "user-impersonate"
)

// CatchAll is the catch-all route suffix
const CatchAll = "/*"

// Error messages
const ERROR_USER_STORE_IS_NIL = "user store cannot be nil"
const ERROR_LOGGER_IS_NIL = "logger cannot be nil"
