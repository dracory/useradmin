// Package useradmin provides a standalone user admin interface following
// the folder-per-controller pattern. Each controller is in its own
// subfolder and handles its own views and AJAX data.
//
// This module is modeled on github.com/dracory/blogadmin.
package useradmin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/useradmin/user_create"
	"github.com/dracory/useradmin/user_delete"
	"github.com/dracory/useradmin/user_impersonate"
	"github.com/dracory/useradmin/user_manager"
	"github.com/dracory/useradmin/user_update"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/geostore"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

// AdminOptions contains all dependencies and configuration for the
// user admin.
//
// UserStore, GeoStore, and Logger are required. SessionStore is
// required for the impersonate controller. Blind index stores,
// TaskStore, and VaultTokenizer are optional — when nil, the
// corresponding features degrade gracefully (filtered search disabled,
// email-change rebuild skipped, user fields treated as plain text).
//
// FuncLayout is an optional function to render the admin interface
// inside your own layout (branding, menus, etc.). If nil, a default
// bare-bones HTML page is used (Bootstrap + Vue CDN). Uses an
// anonymous struct to match blogadmin/shopadmin exactly, so consumers
// can reuse their existing layout function for useradmin.
type AdminOptions struct {
	// UserStore is required
	UserStore userstore.StoreInterface

	// GeoStore is required for the user update controller (country and
	// timezone lists).
	GeoStore geostore.StoreInterface

	// Logger is required
	Logger *slog.Logger

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
	VaultTokenizer shared.VaultTokenizer

	// AuthUser returns the authenticated user from the request, or
	// nil if unauthenticated. Used by the create/delete/impersonate
	// controllers for authorization checks.
	AuthUser func(r *http.Request) userstore.UserInterface

	// AuthUserID returns the authenticated user ID from the request.
	// If it returns "", the user is treated as unauthenticated and
	// redirected to AdminHomeURL. If nil, the auth check is skipped
	// (the host project is expected to gate the route).
	AuthUserID func(r *http.Request) string

	// FlashRedirect redirects with a flash message. Optional — when
	// nil, plain http.Redirect is used.
	FlashRedirect shared.FlashRedirectFunc

	// SecureCookie controls whether the impersonation cookie is marked
	// Secure. Set to false for HTTP (development), true for HTTPS
	// (production). Defaults to true.
	SecureCookie bool

	// FuncLayout is an optional function to render the admin interface
	// inside your own layout (branding, menus, etc.). It receives the
	// request and response writer so the host project can access
	// request context (auth user, locale, etc.) when rendering the
	// layout.
	FuncLayout func(w http.ResponseWriter, r *http.Request, title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string

	// AdminHomeURL is the URL for the admin home page (default: "/admin")
	AdminHomeURL string

	// UserAdminURL is the base URL for the user admin (default: "/admin/users")
	UserAdminURL string

	// UserHomeURL is the URL the impersonate controller redirects to
	// after a successful impersonation (default: "/"). This is where
	// the impersonated user lands.
	UserHomeURL string
}

// AdminInterface defines the interface for the user admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	userStore              userstore.StoreInterface
	geoStore               geostore.StoreInterface
	logger                 *slog.Logger
	sessionStore           sessionstore.StoreInterface
	blindIndexFirstName    blindindexstore.StoreInterface
	blindIndexLastName     blindindexstore.StoreInterface
	blindIndexEmail        blindindexstore.StoreInterface
	taskStore              taskstore.StoreInterface
	blindIndexRebuildAlias string
	vaultTokenizer         shared.VaultTokenizer
	authUser               func(r *http.Request) userstore.UserInterface
	authUserID             func(r *http.Request) string
	flashRedirect          shared.FlashRedirectFunc
	secureCookie           bool
	funcLayout             func(w http.ResponseWriter, r *http.Request, title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	adminHomeURL string
	userAdminURL string
	userHomeURL  string
	routes       map[string]func(w http.ResponseWriter, r *http.Request)
}

// New creates a new user admin instance.
// Returns ErrUserStoreRequired if UserStore is nil, ErrLoggerRequired
// if Logger is nil, ErrGeoStoreRequired if GeoStore is nil, and
// ErrSessionStoreRequired if SessionStore is nil.
//
// This makes misconfiguration fail fast at construction instead of
// surfacing as runtime errors inside individual controllers.
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.UserStore == nil {
		return nil, ErrUserStoreRequired
	}
	if opts.Logger == nil {
		return nil, ErrLoggerRequired
	}
	if opts.GeoStore == nil {
		return nil, ErrGeoStoreRequired
	}
	if opts.SessionStore == nil {
		return nil, ErrSessionStoreRequired
	}

	// Set defaults
	if opts.AdminHomeURL == "" {
		opts.AdminHomeURL = "/admin"
	}
	if opts.UserAdminURL == "" {
		opts.UserAdminURL = "/admin/users"
	}
	if opts.UserHomeURL == "" {
		opts.UserHomeURL = "/"
	}

	a := &admin{
		userStore:              opts.UserStore,
		geoStore:               opts.GeoStore,
		logger:                 opts.Logger,
		sessionStore:           opts.SessionStore,
		blindIndexFirstName:    opts.BlindIndexFirstName,
		blindIndexLastName:     opts.BlindIndexLastName,
		blindIndexEmail:        opts.BlindIndexEmail,
		taskStore:              opts.TaskStore,
		blindIndexRebuildAlias: opts.BlindIndexRebuildTaskAlias,
		vaultTokenizer:         opts.VaultTokenizer,
		authUser:               opts.AuthUser,
		authUserID:             opts.AuthUserID,
		flashRedirect:          opts.FlashRedirect,
		secureCookie:           opts.SecureCookie,
		funcLayout:             opts.FuncLayout,
		adminHomeURL:           opts.AdminHomeURL,
		userAdminURL:           opts.UserAdminURL,
		userHomeURL:            opts.UserHomeURL,
	}

	// Build routes once at construction time
	a.routes = a.buildRoutes()

	return a, nil
}

// Handle processes all user admin requests.
// Config values are injected into the request context (following the
// blogadmin/shopadmin pattern). Route lookup is map-based.
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if a.authUserID != nil && a.authUserID(r) == "" {
		http.Redirect(w, r, a.adminHomeURL, http.StatusSeeOther)
		return
	}

	// Inject config into request context (like blogadmin/shopadmin)
	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)
	ctx = context.WithValue(ctx, shared.KeyUserAdminURL, a.userAdminURL)
	ctx = context.WithValue(ctx, shared.KeyUserHomeURL, a.userHomeURL)
	r = r.WithContext(ctx)

	// Map-based route lookup
	controller := req.GetStringTrimmed(r, "controller")
	if controller == "" {
		controller = shared.CONTROLLER_USER_MANAGER
	}

	handler, ok := a.routes[controller]
	if !ok {
		handler = a.routes[shared.CONTROLLER_USER_MANAGER]
	}

	handler(w, r)
}

// buildRoutes creates the handler dispatch map once at construction time.
func (a *admin) buildRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	uiConfig := shared.UiConfig{
		UserStore:                  a.userStore,
		GeoStore:                   a.geoStore,
		Logger:                     a.logger,
		SessionStore:               a.sessionStore,
		BlindIndexFirstName:        a.blindIndexFirstName,
		BlindIndexLastName:         a.blindIndexLastName,
		BlindIndexEmail:            a.blindIndexEmail,
		TaskStore:                  a.taskStore,
		BlindIndexRebuildTaskAlias: a.blindIndexRebuildAlias,
		VaultTokenizer:             a.vaultTokenizer,
		AuthUser:                   a.authUser,
		FlashRedirect:              a.flashRedirect,
		SecureCookie:               a.secureCookie,
		Layout:                     a.render,
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_USER_MANAGER:     func(w http.ResponseWriter, r *http.Request) { user_manager.UI(uiConfig).UserManager(w, r) },
		shared.CONTROLLER_USER_CREATE:      func(w http.ResponseWriter, r *http.Request) { user_create.UI(uiConfig).UserCreate(w, r) },
		shared.CONTROLLER_USER_DELETE:      func(w http.ResponseWriter, r *http.Request) { user_delete.UI(uiConfig).UserDelete(w, r) },
		shared.CONTROLLER_USER_UPDATE:      func(w http.ResponseWriter, r *http.Request) { user_update.UI(uiConfig).UserUpdate(w, r) },
		shared.CONTROLLER_USER_IMPERSONATE: func(w http.ResponseWriter, r *http.Request) { user_impersonate.UI(uiConfig).UserImpersonate(w, r) },
	}
}

// render wraps content in the layout. If FuncLayout is provided and
// returns non-empty, it is used; otherwise the default shared.Layout
// is used (following the blogadmin/shopadmin pattern).
//
// When FuncLayout is set, the default shared.Layout is NOT computed
// (avoids wasted work).
func (a *admin) render(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	// If a custom layout is provided, try it first
	if a.funcLayout != nil {
		custom := a.funcLayout(w, r, webpageTitle, webpageHtml, options)
		if custom != "" {
			if w != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(custom))
				return ""
			}
			return custom
		}
	}

	webpage := shared.Layout(w, r, webpageTitle, webpageHtml, options)

	if w != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(webpage))
		return ""
	}

	return webpage
}
