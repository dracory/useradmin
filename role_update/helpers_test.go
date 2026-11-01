package role_update_test

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

	"github.com/dracory/useradmin/role_update"
	"github.com/dracory/useradmin/shared"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a role update controller UI backed by an in-memory
// SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    role_update.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a role update
// controller UI wired to it.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping db: %v", err)
	}

	store, err := userstore.NewStore(userstore.NewStoreOptions{
		DB:                 db,
		UserTableName:      "user",
		RolesEnabled:       true,
		RoleTableName:      "role",
		UserRoleTableName:  "user_role",
		AutomigrateEnabled: true,
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("new userstore: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := role_update.UI(shared.UiConfig{
		UserStore: store,
		Logger:    logger,
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml +
				"<script>" + strings.Join(options.Scripts, " ") + "</script></body></html>"
		},
	})

	return &testEnv{ui: ui, store: store, db: db}
}

func (e *testEnv) close() {
	if e.db != nil {
		_ = e.db.Close()
	}
}

// seedRole creates a role with the common fields populated and returns it.
func (e *testEnv) seedRole(t *testing.T) userstore.RoleInterface {
	t.Helper()

	role := userstore.NewRole().
		SetName("Editor").
		SetHandle("editor").
		SetStatus(userstore.ROLE_STATUS_ACTIVE).
		SetMemo("old memo")
	if err := e.store.RoleCreate(context.Background(), role); err != nil {
		t.Fatalf("seed role: %v", err)
	}
	return role
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
	target := "/admin/users?controller=role-update&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.RoleUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
