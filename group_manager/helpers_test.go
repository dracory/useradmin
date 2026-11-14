package group_manager_test

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

	"github.com/dracory/useradmin/group_manager"
	"github.com/dracory/useradmin/shared"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a group manager controller UI backed by an
// in-memory SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    group_manager.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a group
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

	ui := group_manager.UI(shared.UiConfig{
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

// seedGroup creates a group with the given name/handle/status.
func (e *testEnv) seedGroup(t *testing.T, name, handle, status string) userstore.GroupInterface {
	t.Helper()

	group := userstore.NewGroup().
		SetName(name).
		SetHandle(handle).
		SetStatus(status)
	if err := e.store.GroupCreate(context.Background(), group); err != nil {
		t.Fatalf("seed group: %v", err)
	}
	return group
}

type apiResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// postAction builds a POST request to the controller with the given
// action and JSON body, executes it, and returns the decoded envelope.
func postAction(t *testing.T, env *testEnv, action string, extraQuery url.Values, payload any) apiResponse {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}
	target := "/admin/users?controller=group-manager&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.GroupManager(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

// postRaw posts a raw (possibly malformed) body for the given action.
func postRaw(t *testing.T, env *testEnv, action, raw string) apiResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/users?controller=group-manager&action="+action,
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.GroupManager(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
