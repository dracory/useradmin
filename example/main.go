// Example useradmin server.
//
// Boots an HTTP server with an in-memory SQLite database and mounts the
// useradmin panel at /admin/users. No external services required.
//
// Run from the useradmin module root:
//
//	go run ./example
//
// Then open http://localhost:8080/ in your browser.
//
// The example is open (no authentication) so you can click around
// immediately. To wire authentication, provide AuthUserID and AuthUser
// callbacks that read your session/JWT from the request context.
package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/dracory/useradmin"

	"github.com/dracory/geostore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

const (
	addr     = ":8080"
	dbFile   = ":memory:"
	adminURL = "/admin/users"
	homeURL  = "/admin"
	dbDriver = "sqlite"
	dsnSuffix = "?parseTime=true"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(dbFile)
	if err != nil {
		logger.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	userStore, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:                 db,
		UserTableName:      "user",
		AutomigrateEnabled: true,
	})
	if err != nil {
		logger.Error("failed to create userstore", "err", err)
		os.Exit(1)
	}

	geoStore, err := geostore.NewStore(geostore.NewStoreOptions{
		DB:                 db,
		CountryTableName:   "geo_country",
		StateTableName:     "geo_state",
		TimezoneTableName:  "geo_timezone",
		AutomigrateEnabled: true,
		AutoseedEnabled:    true,
	})
	if err != nil {
		logger.Error("failed to create geostore", "err", err)
		os.Exit(1)
	}

	sessionStore, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session",
		AutomigrateEnabled: true,
	})
	if err != nil {
		logger.Error("failed to create sessionstore", "err", err)
		os.Exit(1)
	}

	// Seed the in-memory DB with sample users so the list page has data.
	if dbFile == ":memory:" {
		seedUsers(userStore, logger)
	}

	admin, err := useradmin.New(useradmin.AdminOptions{
		UserStore:    userStore,
		GeoStore:     geoStore,
		Logger:       logger,
		SessionStore: sessionStore,
		AdminHomeURL: homeURL,
		UserAdminURL: adminURL,
		// AuthUserID and AuthUser are intentionally nil so the example
		// is open. Provide them in a real integration.
		SecureCookie: false, // example runs on HTTP
	})
	if err != nil {
		logger.Error("failed to create useradmin", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// useradmin.AdminInterface exposes Handle(w, r) but not ServeHTTP,
	// so wrap it in an http.HandlerFunc.
	mux.Handle(adminURL, http.HandlerFunc(admin.Handle))
	mux.Handle(adminURL+"/", http.HandlerFunc(admin.Handle))

	// Landing page that links into the admin.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, landingHTML, adminURL)
	})

	logger.Info("useradmin example server starting",
		"addr", addr,
		"admin", adminURL,
		"db", dbFile,
	)
	fmt.Printf("\n  useradmin example running:\n    Landing:  http://localhost%s/\n    Admin:    http://localhost%s\n\n", portFromAddr(addr), adminURL)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}

// openDB opens a SQLite database. For ":memory:" a fresh database is
// created on every run; pass a file path to persist data across restarts.
func openDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" {
		// Start clean so re-runs don't collide with stale schema.
		if err := os.Remove(filepath); err != nil && !strings.Contains(err.Error(), "no such file") {
			return nil, err
		}
	}

	db, err := sql.Open(dbDriver, filepath+dsnSuffix)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// portFromAddr extracts the port from a listen address like ":8080" or
// "127.0.0.1:8080". Falls back to the original addr if parsing fails.
func portFromAddr(addr string) string {
	i := strings.LastIndex(addr, ":")
	if i < 0 {
		return addr
	}
	return addr[i:]
}

const landingHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>useradmin example</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="bg-light">
  <div class="container py-5">
    <div class="row justify-content-center">
      <div class="col-md-7">
        <div class="card shadow-sm">
          <div class="card-body p-5">
            <h1 class="h3 mb-3">useradmin example</h1>
            <p class="text-muted mb-4">
              Standalone user admin panel running on an in-memory SQLite
              database. Data is reset on every restart.
            </p>
            <a href="%s" class="btn btn-primary">Open User Admin &rarr;</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</body>
</html>`
