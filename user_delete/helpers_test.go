package user_delete_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/useradmin/user_delete"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a user delete controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    user_delete.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a user
// delete controller UI.
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

	ui := user_delete.UI(shared.UiConfig{
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

// seedUser creates a user with the given email/status.
func (e *testEnv) seedUser(t *testing.T, email string) userstore.UserInterface {
	t.Helper()

	user := userstore.NewUser().
		SetEmail(email).
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// postForm posts a form-encoded body to the controller.
func postForm(t *testing.T, env *testEnv, form url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/users?controller=user-delete",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	env.ui.UserDelete(rec, req)
	return rec
}
