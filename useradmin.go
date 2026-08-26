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

	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

// AdminOptions contains all dependencies and configuration for the
// user admin.
//
// UserStore, GeoResolver, and Logger are required. SessionStore is
// required for the impersonate controller. Blind index stores,
// UserPiiSeal/UserPiiUnseal/UsersPiiUnseal are optional — when nil,
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

	// GeoResolver is required for the user update controller (country
	// and timezone lists).
	GeoResolver shared.GeoResolverInterface

	// Logger is required
	Logger *slog.Logger

	// OnUserImpersonate is optional — when nil, the impersonate
	// button is hidden and the impersonate route is not registered.
	// The host owns the auth mechanism (session+cookie, JWT, etc.).
	OnUserImpersonate shared.OnUserImpersonateFunc

	// OnUserSearch is an optional callback for custom user search
	// (e.g. blind index, Elasticsearch). When nil, useradmin falls
	// back to userstore query-based search.
	OnUserSearch shared.OnUserSearchFunc

	// OnUserUpdate is an optional callback invoked after a user is
	// updated. The host can use it to trigger side effects (blind
	// index rebuild, audit log, notifications, etc.). When nil, the
	// callback is skipped.
	OnUserUpdate shared.OnUserUpdateFunc

	// UserPiiSeal transforms a user from display representation to
	// storage representation (e.g. tokenize, encrypt PII). Optional —
	// when nil, the user is stored as-is (plain text).
	UserPiiSeal shared.UserPiiSealFunc

	// UserPiiUnseal transforms a user from storage representation to
	// display representation (e.g. detokenize, decrypt PII). Optional —
	// when nil, the user is used as-is (plain text).
	UserPiiUnseal shared.UserPiiUnsealFunc

	// UsersPiiUnseal is the batch version of UserPiiUnseal. It allows
	// the host to unseal all users in a single call for efficiency.
	// Optional — when nil, useradmin falls back to UserPiiUnseal per
	// user (or plain text when that is also nil).
	UsersPiiUnseal shared.UsersPiiUnsealFunc

	// FlashRedirect redirects with a flash message. Optional — when
	// nil, plain http.Redirect is used.
	FlashRedirect shared.FlashRedirectFunc

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

// Re-exports of shared types so consumers can import everything from
// the top-level useradmin package without reaching into useradmin/shared.
// Follows the blogadmin/shopadmin convention (e.g. shopadmin.CustomerResolverInterface).
type (
	GeoResolverInterface  = shared.GeoResolverInterface
	Country               = shared.Country
	Timezone              = shared.Timezone
	UserSearchEvent       = shared.UserSearchEvent
	OnUserSearchFunc      = shared.OnUserSearchFunc
	OnUserImpersonateFunc = shared.OnUserImpersonateFunc
	OnUserUpdateFunc      = shared.OnUserUpdateFunc
	UserPiiSealFunc       = shared.UserPiiSealFunc
	UserPiiUnsealFunc     = shared.UserPiiUnsealFunc
	UsersPiiUnsealFunc    = shared.UsersPiiUnsealFunc
	FlashRedirectFunc     = shared.FlashRedirectFunc
)

// admin implements AdminInterface
type admin struct {
	userStore         userstore.StoreInterface
	geoResolver       shared.GeoResolverInterface
	logger            *slog.Logger
	onUserImpersonate shared.OnUserImpersonateFunc
	onUserSearch      shared.OnUserSearchFunc
	onUserUpdate      shared.OnUserUpdateFunc
	userPiiSeal       shared.UserPiiSealFunc
	userPiiUnseal     shared.UserPiiUnsealFunc
	usersPiiUnseal    shared.UsersPiiUnsealFunc
	flashRedirect     shared.FlashRedirectFunc
	funcLayout        func(w http.ResponseWriter, r *http.Request, title string, body string, options struct {
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
// if Logger is nil, ErrGeoResolverRequired if GeoResolver is nil.
//
// OnUserImpersonate is optional — when nil, the impersonate button is
// hidden and the impersonate route returns 404.
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
	if opts.GeoResolver == nil {
		return nil, ErrGeoResolverRequired
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
		userStore:         opts.UserStore,
		geoResolver:       opts.GeoResolver,
		logger:            opts.Logger,
		onUserImpersonate: opts.OnUserImpersonate,
		onUserSearch:      opts.OnUserSearch,
		onUserUpdate:      opts.OnUserUpdate,
		userPiiSeal:       opts.UserPiiSeal,
		userPiiUnseal:     opts.UserPiiUnseal,
		usersPiiUnseal:    opts.UsersPiiUnseal,
		flashRedirect:     opts.FlashRedirect,
		funcLayout:        opts.FuncLayout,
		adminHomeURL:      opts.AdminHomeURL,
		userAdminURL:      opts.UserAdminURL,
		userHomeURL:       opts.UserHomeURL,
	}

	// Build routes once at construction time
	a.routes = a.buildRoutes()

	return a, nil
}

// Handle processes all user admin requests.
// Config values are injected into the request context (following the
// blogadmin/shopadmin pattern). Route lookup is map-based.
//
// Authentication and authorization are the host's responsibility —
// gate the routes with middleware before they reach Handle.
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
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
		UserStore:         a.userStore,
		GeoResolver:       a.geoResolver,
		Logger:            a.logger,
		OnUserImpersonate: a.onUserImpersonate,
		OnUserSearch:      a.onUserSearch,
		OnUserUpdate:      a.onUserUpdate,
		UserPiiSeal:       a.userPiiSeal,
		UserPiiUnseal:     a.userPiiUnseal,
		UsersPiiUnseal:    a.usersPiiUnseal,
		FlashRedirect:     a.flashRedirect,
		Layout:            a.render,
	}

	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_USER_MANAGER: func(w http.ResponseWriter, r *http.Request) { user_manager.UI(uiConfig).UserManager(w, r) },
		shared.CONTROLLER_USER_CREATE:  func(w http.ResponseWriter, r *http.Request) { user_create.UI(uiConfig).UserCreate(w, r) },
		shared.CONTROLLER_USER_DELETE:  func(w http.ResponseWriter, r *http.Request) { user_delete.UI(uiConfig).UserDelete(w, r) },
		shared.CONTROLLER_USER_UPDATE:  func(w http.ResponseWriter, r *http.Request) { user_update.UI(uiConfig).UserUpdate(w, r) },
	}

	// Only register the impersonate route when OnUserImpersonate is
	// configured. When nil, the route is not registered (404).
	if a.onUserImpersonate != nil {
		routes[shared.CONTROLLER_USER_IMPERSONATE] = func(w http.ResponseWriter, r *http.Request) {
			user_impersonate.UI(uiConfig).UserImpersonate(w, r)
		}
	}

	return routes
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
