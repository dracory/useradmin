package user_view_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/useradmin/user_view"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a user view controller UI backed by an in-memory
// SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    user_view.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore (roles and groups
// enabled) and a user view controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	return newTestEnvWithOptions(t, true, true)
}

// newTestEnvWithOptions boots a user view controller UI with roles and
// groups optionally enabled — the page must hide the corresponding
// cards when the tables are not configured.
func newTestEnvWithOptions(t *testing.T, rolesEnabled, groupsEnabled bool) *testEnv {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping db: %v", err)
	}

	opts := userstore.NewStoreOptions{
		DB:                 db,
		UserTableName:      "user",
		AutomigrateEnabled: true,
	}
	if rolesEnabled {
		opts.RolesEnabled = true
		opts.RoleTableName = "role"
		opts.UserRoleTableName = "user_role"
	}
	if groupsEnabled {
		opts.GroupsEnabled = true
		opts.GroupTableName = "group"
		opts.UserGroupTableName = "user_group"
	}

	store, err := userstore.NewStore(opts)
	if err != nil {
		_ = db.Close()
		t.Fatalf("new userstore: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := user_view.UI(shared.UiConfig{
		UserStore: store,
		Logger:    logger,
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

func (e *testEnv) close() {
	if e.db != nil {
		_ = e.db.Close()
	}
}

// seedUser creates a user with the common fields populated and returns it.
func (e *testEnv) seedUser(t *testing.T) userstore.UserInterface {
	t.Helper()

	user := userstore.NewUser().
		SetEmail("target@test.com").
		SetFirstName("First").
		SetLastName("Last").
		SetBusinessName("Biz").
		SetPhone("123").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_USER).
		SetCountry("GB").
		SetTimezone("Europe/London").
		SetMemo("some memo")
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// seedRole creates an active role and returns it.
func (e *testEnv) seedRole(t *testing.T) userstore.RoleInterface {
	t.Helper()

	role := userstore.NewRole().
		SetName("Guardian").
		SetHandle("guardian").
		SetStatus(userstore.ROLE_STATUS_ACTIVE)
	if err := e.store.RoleCreate(context.Background(), role); err != nil {
		t.Fatalf("seed role: %v", err)
	}
	return role
}

// seedGroup creates an active group and returns it.
func (e *testEnv) seedGroup(t *testing.T) userstore.GroupInterface {
	t.Helper()

	group := userstore.NewGroup().
		SetName("Members").
		SetHandle("members").
		SetStatus(userstore.GROUP_STATUS_ACTIVE)
	if err := e.store.GroupCreate(context.Background(), group); err != nil {
		t.Fatalf("seed group: %v", err)
	}
	return group
}

// get issues a GET request against the user view controller and
// returns the recorder.
func get(t *testing.T, env *testEnv, target string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	env.ui.UserView(rec, req)
	return rec
}
