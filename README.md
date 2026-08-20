# useradmin

[![Tests Status](https://github.com/dracory/useradmin/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/useradmin/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/useradmin)](https://goreportcard.com/report/github.com/dracory/useradmin)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/useradmin)](https://pkg.go.dev/github.com/dracory/useradmin)

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Introduction

Admin interface for [`github.com/dracory/userstore`](https://github.com/dracory/userstore).
Provides a ready-to-use admin panel for managing users: list with AJAX,
create, update, delete, and impersonate.

Modeled after [`github.com/dracory/blogadmin`](https://github.com/dracory/blogadmin)
— same folder-per-controller pattern, same `UiConfig`/`UiBase` conventions.

## Features

- **User management** — list with AJAX pagination/sorting/filtering,
  create, update, delete
- **User impersonation** — sign in as another user for support/debugging
- **Geo integration** — country and timezone pickers via `geostore`
- **Vault tokenization** — pluggable vault tokenizer for encrypted user
  fields (first name, last name, email, phone, business name)
- **Blind index search** — pluggable blind index stores for filtered
  search on tokenized fields
- **Custom layouts** — bring your own layout via `FuncLayout`
- **Bootstrap + Vue CDN** — default UI works out of the box

## Installation

```bash
go get github.com/dracory/useradmin
```

## Quick Start

```go
package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/dracory/useradmin"

    "github.com/dracory/geostore"
    "github.com/dracory/sessionstore"
    "github.com/dracory/userstore"
)

func main() {
    userStore, _ := userstore.NewStore(userstore.NewStoreOptions{
        DB:                 yourDB,
        UserTableName:      "user",
        AutomigrateEnabled: true,
    })

    geoStore, _ := geostore.NewStore(geostore.NewStoreOptions{
        DB:                 yourDB,
        CountryTableName:   "geo_country",
        TimezoneTableName:  "geo_timezone",
        AutomigrateEnabled: true,
        AutoseedEnabled:    true,
    })

    sessionStore, _ := sessionstore.NewStore(sessionstore.NewStoreOptions{
        DB:                 yourDB,
        SessionTableName:   "session",
        AutomigrateEnabled: true,
    })

    admin, err := useradmin.New(useradmin.AdminOptions{
        UserStore:    userStore,
        GeoStore:     geoStore,
        Logger:       slog.New(slog.NewTextHandler(os.Stderr, nil)),
        SessionStore: sessionStore,
        AdminHomeURL: "/admin",
        UserAdminURL: "/admin/users",
    })
    if err != nil {
        log.Fatal(err)
    }

    http.Handle("/admin/users", http.HandlerFunc(admin.Handle))
    http.ListenAndServe(":8080", nil)
}
```

See [`example/`](example/) for a complete runnable server with
in-memory SQLite and seed data.

## Integration with a Router

`useradmin.AdminInterface` exposes `Handle(w, r)`, which is an
`http.HandlerFunc`-compatible method. Wire it into any router that
accepts standard `http.Handler`:

```go
// stdlib
mux.Handle("/admin/users", http.HandlerFunc(admin.Handle))

// github.com/dracory/rtr
route := rtr.NewRoute().
    SetName("Admin > Users").
    SetPath("/admin/users").
    SetHTMLHandler(admin.Handle)
```

## Optional Dependencies

The following `AdminOptions` fields are optional. When nil/empty, the
corresponding features degrade gracefully:

- **BlindIndexFirstName / BlindIndexLastName / BlindIndexEmail** — when
  nil, filtered search by that field is disabled (the filter returns no
  matches instead of erroring).
- **TaskStore + BlindIndexRebuildTaskAlias** — when nil/empty, the
  blind index rebuild enqueue on email change is skipped.
- **VaultTokenizer** — when nil, user fields are treated as plain text
  (no tokenization/untokenization).
- **AuthUser** — when nil, the create/delete/impersonate controllers
  treat the request as unauthenticated.
- **FlashRedirect** — when nil, plain `http.Redirect` is used instead of
  a flash-message redirect.
- **FuncLayout** — when nil, a default bare-bones HTML page is used
  (Bootstrap + Vue CDN).

## Custom Layout

By default, useradmin renders a bare-bones HTML page with Bootstrap and
Vue from CDN. To embed the admin inside your own layout (branding, menus,
etc.), provide `FuncLayout`:

```go
admin, _ := useradmin.New(useradmin.AdminOptions{
    UserStore:    userStore,
    GeoStore:     geoStore,
    Logger:       logger,
    SessionStore: sessionStore,
    FuncLayout: func(w http.ResponseWriter, r *http.Request, title, body string, opts struct {
        Styles     []string
        StyleURLs  []string
        Scripts    []string
        ScriptURLs []string
    }) string {
        return myLayout(w, r, title, body, opts)
    },
})
```

`FuncLayout` receives the request and response writer so the host
project can access request context (auth user, locale, etc.) when
rendering the layout.

## Testing

```bash
go test ./...
```
