package user_impersonate_test

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/useradmin/user_impersonate"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a user impersonate controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    user_impersonate.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a user
// impersonate controller UI wired to the given impersonate func.
func newTestEnv(t *testing.T, onImpersonate shared.OnUserImpersonateFunc) *testEnv {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	store, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:                 db,
		UserTableName:      "user",
		RolesEnabled:       true,
		RoleTableName:      "role",
		UserRoleTableName:  "user_role",
		GroupsEnabled:      true,
		GroupTableName:     "group",
		UserGroupTableName: "user_group",
		AutomigrateEnabled: true,
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("new userstore: %v", err)
	}

	ui := user_impersonate.UI(shared.UiConfig{
		UserStore:         store,
		Logger:            slog.New(slog.NewTextHandler(os.Stderr, nil)),
		OnUserImpersonate: onImpersonate,
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})

	return &testEnv{ui: ui, store: store, db: db}
}

func (e *testEnv) close() { _ = e.db.Close() }

// seedUser creates a user with the given email/status.
func (e *testEnv) seedUser(t *testing.T, email, status string) userstore.UserInterface {
	t.Helper()

	user := userstore.NewUser().
		SetEmail(email).
		SetStatus(status)
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// get performs a GET request against the controller.
func get(t *testing.T, env *testEnv, query string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet,
		"/admin/users?controller=user-impersonate&"+query, http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.UserImpersonate(rec, req)
	return rec
}

var errImpersonate = errors.New("impersonate failed")
