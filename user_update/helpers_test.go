package user_update_test

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
	"github.com/dracory/useradmin/user_update"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv bundles a user update controller UI backed by an in-memory
// SQLite userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    user_update.UiInterface
	store userstore.StoreInterface
	db    *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a user update
// controller UI wired to the given geo resolver (may be nil).
func newTestEnv(t *testing.T, geo shared.GeoResolverInterface) *testEnv {
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
		GroupsEnabled:      true,
		GroupTableName:     "group",
		UserGroupTableName: "user_group",
		AutomigrateEnabled: true,
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("new userstore: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := user_update.UI(shared.UiConfig{
		UserStore:   store,
		GeoResolver: geo,
		Logger:      logger,
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
		SetMemo("old memo")
	if err := user.SetPasswordAndHash("oldpassword123"); err != nil {
		t.Fatalf("SetPasswordAndHash: %v", err)
	}
	if err := e.store.UserCreate(context.Background(), user); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

// stubGeo satisfies shared.GeoResolverInterface with canned data.
type stubGeo struct {
	timezones []shared.Timezone
}

var _ shared.GeoResolverInterface = (*stubGeo)(nil)

func (g *stubGeo) Countries(ctx context.Context) ([]shared.Country, error) {
	return []shared.Country{{IsoCode2: "GB", Name: "United Kingdom"}}, nil
}

func (g *stubGeo) Timezones(ctx context.Context, countryCode ...string) ([]shared.Timezone, error) {
	return g.timezones, nil
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
	target := "/admin/users?controller=user-update&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.UserUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
