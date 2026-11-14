package user_manager_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/useradmin/user_manager"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a user manager controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    user_manager.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a user
// manager controller UI.
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

	ui := user_manager.UI(shared.UiConfig{
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
func (e *testEnv) seedUser(t *testing.T, email, status string) userstore.UserInterface {
	t.Helper()

	user := userstore.NewUser().
		SetEmail(email).
		SetFirstName("First").
		SetLastName("Last").
		SetStatus(status)
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// apiResp is the decoded JSON envelope returned by the api package.
type apiResp map[string]any

func (a apiResp) status() string  { s, _ := a["status"].(string); return s }
func (a apiResp) message() string { s, _ := a["message"].(string); return s }
func (a apiResp) data() map[string]any {
	d, _ := a["data"].(map[string]any)
	return d
}

// postAction builds a POST request to the controller with the given
// action and JSON body, executes it, and returns the decoded envelope.
func postAction(t *testing.T, env *testEnv, action string, extraQuery url.Values, payload any) apiResp {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}
	target := "/admin/users?controller=user-manager&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.UserManager(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

// postRaw posts a raw (possibly malformed) body for the given action.
func postRaw(t *testing.T, env *testEnv, action, raw string) apiResp {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/users?controller=user-manager&action="+action,
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.UserManager(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
