package role_members_test

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

	"github.com/dracory/useradmin/role_members"
	"github.com/dracory/useradmin/shared"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a role members controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    role_members.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a role
// members controller UI.
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

	ui := role_members.UI(shared.UiConfig{
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

// seedUser creates a user with the given email.
func (e *testEnv) seedUser(t *testing.T, email string) userstore.UserInterface {
	t.Helper()

	user := userstore.NewUser().
		SetEmail(email).
		SetFirstName("First").
		SetLastName("Last").
		SetStatus(userstore.USER_STATUS_ACTIVE)
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// seedRole creates a role with the given name/handle.
func (e *testEnv) seedRole(t *testing.T, name, handle string) userstore.RoleInterface {
	t.Helper()

	role := userstore.NewRole().SetName(name).SetHandle(handle).SetStatus("active")
	if err := e.store.RoleCreate(context.Background(), role); err != nil {
		t.Fatalf("seed role: %v", err)
	}
	return role
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
// action/query and JSON body, executes it, and returns the decoded
// envelope.
func postAction(t *testing.T, env *testEnv, action string, extraQuery url.Values, payload any) apiResp {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}
	target := "/admin/users?controller=role-members&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.RoleMembers(rec, req)

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
		"/admin/users?controller=role-members&action="+action,
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.RoleMembers(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

// get performs a GET request against the controller.
func get(t *testing.T, env *testEnv, query string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet,
		"/admin/users?controller=role-members&"+query, http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.RoleMembers(rec, req)
	return rec
}
