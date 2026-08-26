# Refactor Plan: Decouple `useradmin` from non-`userstore` store packages

## Principle

**`useradmin`'s public API must only expose types from `useradmin` itself,
`userstore`, `net/http`, and `log/slog`. The host must not be forced to import
any other package to wire up `useradmin`.**

If the host has to import `geostore`, `sessionstore`, `blindindexstore`,
`taskstore`, `auth`, or `carbon` to construct a `useradmin.AdminOptions`,
then `useradmin` is dictating the host's infrastructure choices. The host
may have completely different implementations for geo data, sessions, blind
index, and tasks.

## Problem

The `useradmin` module was extracted from `blueprint/pkg/useradmin` into a
standalone module (`github.com/dracory/useradmin`). The extraction is a
positive change, but it is over-eager: `useradmin` leaks 6 infrastructure
packages through its public API (`AdminOptions`, `UiConfig`, `UiInterface`),
forcing the host to import them.

`VaultTokenizer` already shows the right pattern — an interface in `shared/`
with the host providing an adapter. The other infrastructure dependencies
should follow the same pattern.

## Dependency audit

### Packages that leak through the public API (must be abstracted)

These appear in `AdminOptions` / `UiConfig` / `UiInterface` field types, so
the host must import them to construct or implement the config:

| Package | Used for | Also couples to |
|---|---|---|
| `geostore` | `CountryList`, `TimezoneList` | `CountryQueryOptions`, `TimezoneQueryOptions`, `COLUMN_NAME`, `COLUMN_TIMEZONE` |
| `sessionstore` | `SessionCreate` | `NewSession()` builder, `Session` type |
| `blindindexstore` | `Search` | `SEARCH_TYPE_EQUALS`, `SEARCH_TYPE_CONTAINS` constants |
| `taskstore` | `TaskDefinitionEnqueueByAlias` | `DefaultQueueName` constant |
| `auth` | `AuthCookieSet` | `auth/types.WithSecure` |
| `carbon/v2` | `Now(UTC).AddHours(2)...` | session expiry policy |

### Packages that are internal implementation details (fine as transitive deps)

These are used inside `useradmin`'s controllers to render HTML, parse
requests, and format JSON responses. They do **not** appear in any exported
field type, so the host never imports them directly — Go resolves them as
transitive dependencies automatically:

| Package | Used for | Usage count |
|---|---|---|
| `hb` | HTML builder (nodes, children, attributes) | ~60 calls across 9 files |
| `bs` | Bootstrap HTML components (FormGroup, Modal, etc.) | ~21 calls across 2 files |
| `api` | JSON response helpers (`Respond`, `Error`, `Success`) | ~62 calls across 6 files |
| `req` | Request parsing (`GetStringTrimmed`, `GetIP`) | ~13 calls across 8 files |
| `cdn` | CDN URL constants for Bootstrap/Vue/jQuery | ~10 calls in layout/page files |
| `govalidator` | `IsEmail` validation | 3 calls |
| `neat` | `SortAsc`/`SortDesc` sort-order constants | 4 calls (all go away with `GeoLister` abstraction) |

After the refactor, `neat` is dropped entirely (its 4 uses disappear when
`GeoLister` replaces `geostore` query options). `govalidator` could
optionally be replaced with stdlib `net/mail.ParseAddress`, but it does not
leak through the API so it is not required for this refactor.

## New abstractions in `shared/` (new file `shared/resolvers.go`)

Naming follows the existing convention used by `blogadmin`/`shopadmin`:
`XxxResolverInterface` (interface) + `XxxResolver` (host's struct). Func
types stay as `XxxFunc` (matching `FlashRedirectFunc`).

```go
// GeoResolverInterface provides countries and timezones. The host
// implements this against whatever geo data source it uses
// (dracory/geostore, a static list, an external API, etc.).
type GeoResolverInterface interface {
	Countries(ctx context.Context) ([]Country, error)
	Timezones(ctx context.Context, countryCode string) ([]Timezone, error)
}

type Country struct {
	IsoCode2 string
	Name     string
}

type Timezone struct {
	Code string
}

// BlindIndexSearchType selects the match mode for
// BlindIndexResolverInterface.Search. Constants live in useradmin (not
// blindindexstore) so the host is free to map them to its own search
// backend.
type BlindIndexSearchType string

const (
	BlindIndexSearchEquals   BlindIndexSearchType = "equals"
	BlindIndexSearchContains BlindIndexSearchType = "contains"
)

// BlindIndexResolverInterface searches a blind index for user IDs matching
// a value. One instance per indexed field (first name, last name, email).
type BlindIndexResolverInterface interface {
	Search(ctx context.Context, value string, searchType BlindIndexSearchType) ([]string, error)
}

// TaskEnqueuerFunc enqueues a background task by alias with the given
// payload. The host routes the alias to its own task system.
type TaskEnqueuerFunc func(ctx context.Context, alias string, payload map[string]any) error

// SessionCreatorFunc creates a new session for the given user ID and sets
// the auth cookie on the response. The host owns the session store, cookie
// format, and expiry policy. `secure` controls the Secure cookie flag
// (false for HTTP development, true for HTTPS production).
type SessionCreatorFunc func(w http.ResponseWriter, r *http.Request, userID string, secure bool) error
```

`TaskEnqueuerFunc` and `SessionCreatorFunc` remain func types (matching
`FlashRedirectFunc`) since they are single-operation callbacks, not
multi-method resolvers. If you prefer them as `TaskResolverInterface` /
`SessionResolverInterface` instead, say so.

## Files changed in `useradmin`

| File | Change |
|---|---|
| `shared/resolvers.go` | **NEW** — the 4 abstractions above |
| `shared/ui_interface.go` | Replace `GeoStore()/SessionStore()/BlindIndex*/TaskStore()` return types with `GeoResolverInterface/SessionCreatorFunc/BlindIndexResolverInterface/TaskEnqueuerFunc`. Drop 4 store imports. |
| `shared/ui_base.go` | Same field-type swaps. Drop 4 store imports. |
| `shared/ui_config.go` | Same field-type swaps. Drop 4 store imports. |
| `useradmin.go` | `AdminOptions` + `admin` struct field-type swaps. Drop 4 store imports. Update `New()` validation (rename `ErrGeoStoreRequired`→`ErrGeoResolverRequired`, `ErrSessionStoreRequired`→`ErrSessionCreatorRequired`). |
| `user_update/handle_user_fetch_ajax.go` | `GeoStore().CountryList(ctx, opts)` → `GeoResolver().Countries(ctx)`; `TimezoneList(ctx, opts)` → `Timezones(ctx, country)`. Field access `c.IsoCode2()`→`c.IsoCode2`, `tz.Timezone()`→`tz.Code`. Drop `geostore`, `neat` imports. |
| `user_update/handle_timezones_fetch_ajax.go` | Same `TimezoneList`→`Timezones` swap. Drop `geostore`, `neat` imports. |
| `user_update/handle_user_update_ajax.go` | `TaskStore().TaskDefinitionEnqueueByAlias(ctx, queueName, alias, payload)` → `TaskEnqueuer()(ctx, alias, payload)`. Drop `taskstore` import. |
| `user_manager/handle_users_fetch_ajax.go` | `blindindexstore.SEARCH_TYPE_CONTAINS` → `shared.BlindIndexSearchContains`. Drop `blindindexstore` import. |
| `user_manager/handle_user_create_ajax.go` | `blindindexstore.SEARCH_TYPE_EQUALS` → `shared.BlindIndexSearchEquals`. Drop `blindindexstore` import. |
| `user_create/user_create_controller.go` | `"equals"` string literal → `shared.BlindIndexSearchEquals`. (No import change — already doesn't import blindindexstore.) |
| `user_impersonate/impersonate.go` | **DELETE** — session+cookie logic moves to host's `SessionCreatorFunc` adapter. |
| `user_impersonate/user_impersonate_controller.go` | `Impersonate(u.SessionStore(), w, r, userID, secure)` → `u.SessionCreator()(w, r, userID, secure)` with nil check. |
| `errors.go` | Rename `ErrGeoStoreRequired`→`ErrGeoResolverRequired`, `ErrSessionStoreRequired`→`ErrSessionCreatorRequired` (if they exist there). |
| `go.mod` | Remove direct deps: `auth`, `blindindexstore`, `geostore`, `sessionstore`, `taskstore`, `carbon/v2`, `neat`. |

## Files changed in `blueprint`

| File | Change |
|---|---|
| `internal/controllers/admin/adapters/adapters.go` | Add 4 adapters: `GeoResolver` (wraps `geostore.StoreInterface`), `BlindIndexResolver` (wraps `blindindexstore.StoreInterface`), `NewTaskEnqueuerFunc` (wraps `taskstore.StoreInterface`), `NewSessionCreatorFunc` (wraps `sessionstore.StoreInterface` + `auth.AuthCookieSet` + `carbon` expiry — absorbs the deleted `impersonate.go` logic). |
| `internal/controllers/admin/users/users_controller.go` | Wire adapters: `GeoResolver: adapters.NewGeoResolver(app.GetGeoStore())`, `SessionCreator: adapters.NewSessionCreatorFunc(app)`, `BlindIndexFirstName/LastName/Email: adapters.NewBlindIndexResolver(...)`, `TaskEnqueuer: adapters.NewTaskEnqueuerFunc(app.GetTaskStore())`. |

## `example/` changes

`example/main.go` creates real `geostore`/`sessionstore` instances — needs
thin adapter wrappers to satisfy the new interfaces. Small change, demo-only.

## Public API after refactor

The host only imports these to wire `useradmin`:

| Import | Why |
|---|---|
| `github.com/dracory/useradmin` | `New()`, `AdminOptions`, `AdminInterface` |
| `github.com/dracory/useradmin/shared` | `VaultTokenizer`, `FlashRedirectFunc`, `GeoResolverInterface`, `BlindIndexResolverInterface`, `TaskEnqueuerFunc`, `SessionCreatorFunc` |
| `github.com/dracory/userstore` | `StoreInterface`, `UserInterface`, `UserQuery` |
| `net/http` | `http.ResponseWriter`, `*http.Request` |
| `log/slog` | `*slog.Logger` |

No other package is required. The 6 infrastructure packages (`geostore`,
`sessionstore`, `blindindexstore`, `taskstore`, `auth`, `carbon`) are only
imported by the host's adapter implementations — that's the host's choice,
not `useradmin`'s requirement.

## What `useradmin` imports after (direct deps)

`govalidator`, `api`, `bs`, `cdn`, `hb`, `req`, `userstore` +
`modernc.org/sqlite` (example only). All are internal implementation details
that do not leak through the public API. The 6 infrastructure packages move
to the host side as adapter implementations.

## Tradeoff

This **diverges from blogadmin/shopadmin**, which import their concrete store
packages directly. The divergence is justified: blogadmin's stores
(`blogstore`, `customstore`, `settingstore`) are blog-specific, while
useradmin's were generic infrastructure (geo/session/index/tasks) that hosts
reasonably implement differently.

## Verification

- `go build ./...` in `useradmin`
- `go build ./...` in `blueprint`
- `go test ./...` in `useradmin` (existing tests use the old signatures — will
  need updating; per `AGENTS.md` tests are "difficult to maintain, use only
  test functions", so existing test functions will have their signatures
  updated rather than adding new ones)
