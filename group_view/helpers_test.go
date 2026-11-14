package group_view_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracory/useradmin/group_view"
	"github.com/dracory/useradmin/shared"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a group view controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    group_view.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a group
// view controller UI.
func newTestEnv(t *testing.T) *testEnv {
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

	ui := group_view.UI(shared.UiConfig{
		UserStore: store,
		Logger:    slog.New(slog.NewTextHandler(os.Stderr, nil)),
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

// seedGroup creates a group with the given name.
func (e *testEnv) seedGroup(t *testing.T, name string) userstore.GroupInterface {
	t.Helper()

	group := userstore.NewGroup().SetName(name)
	if err := e.store.GroupCreate(context.Background(), group); err != nil {
		t.Fatalf("seed group: %v", err)
	}
	return group
}

// get performs a GET request against the controller.
func get(t *testing.T, env *testEnv, query string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet,
		"/admin/groups/view?controller=group-view&"+query, http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.GroupView(rec, req)
	return rec
}
