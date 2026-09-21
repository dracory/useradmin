package useradmin_test

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
	"testing"

	"github.com/dracory/useradmin"

	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

// testEnv holds a useradmin instance backed by an in-memory SQLite
// userstore, ready to serve HTTP requests in a test.
type testEnv struct {
	admin     useradmin.AdminInterface
	userStore userstore.StoreInterface
	db        *sql.DB
}

// newTestEnv boots an in-memory SQLite userstore and a useradmin
// instance wired to a stub geo resolver. Modeled on example/main.go.
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

	userStore, err := userstore.NewStore(userstore.NewStoreOptions{
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

	admin, err := useradmin.New(useradmin.AdminOptions{
		UserStore:    userStore,
		GeoResolver:  &stubGeoResolver{},
		Logger:       logger,
		AdminHomeURL: "/admin",
		UserAdminURL: "/admin/users",
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("new useradmin: %v", err)
	}

	return &testEnv{admin: admin, userStore: userStore, db: db}
}

func (e *testEnv) close() {
	if e.db != nil {
		_ = e.db.Close()
	}
}

// stubGeoResolver satisfies useradmin.GeoResolverInterface with empty
// lists. The validation tests do not exercise geo data.
type stubGeoResolver struct{}

var _ useradmin.GeoResolverInterface = (*stubGeoResolver)(nil)

func (r *stubGeoResolver) Countries(ctx context.Context) ([]useradmin.Country, error) {
	return nil, nil
}

func (r *stubGeoResolver) Timezones(ctx context.Context, countryCode ...string) ([]useradmin.Timezone, error) {
	return nil, nil
}

// postJSON builds a POST request to the useradmin handler with the given
// query params and JSON body, executes it, and returns the decoded
// response envelope.
func postJSON(t *testing.T, env *testEnv, queryParams url.Values, payload any) apiResponse {
	t.Helper()

	target := "/admin/users?" + queryParams.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.admin.Handle(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

type apiResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// TestUserUpdateAcceptsMinimalPayload verifies that an admin can save a
// user with only email + status — first name, last name, country, and
// timezone are optional. Guards the validation relaxation in
// user_update/handle_user_update_ajax.go against regression.
func TestUserUpdateAcceptsMinimalPayload(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	ctx := context.Background()

	target := userstore.NewUser().
		SetEmail("target-update@test.com").
		SetFirstName("Old").
		SetLastName("Name").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_USER).
		SetCountry("GB").
		SetTimezone("Europe/London").
		SetPassword("password123")
	if err := env.userStore.UserCreate(ctx, target); err != nil {
		t.Fatalf("seed target: %v", err)
	}

	resp := postJSON(t, env, url.Values{
		"controller": {"user-update"},
		"action":     {"user-update-ajax"},
	}, map[string]string{
		"user_id": target.GetID(),
		"status":  userstore.USER_STATUS_ACTIVE,
		"email":   "target-update@test.com",
		// first_name, last_name, country, timezone intentionally omitted
	})

	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	// Verify the optional fields were cleared in storage
	saved, err := env.userStore.UserFindByID(ctx, target.GetID())
	if err != nil {
		t.Fatalf("UserFindByID: %v", err)
	}
	if saved == nil {
		t.Fatal("saved user not found")
	}
	if saved.GetFirstName() != "" {
		t.Errorf("expected first name cleared, got %q", saved.GetFirstName())
	}
	if saved.GetLastName() != "" {
		t.Errorf("expected last name cleared, got %q", saved.GetLastName())
	}
	if saved.GetCountry() != "" {
		t.Errorf("expected country cleared, got %q", saved.GetCountry())
	}
	if saved.GetTimezone() != "" {
		t.Errorf("expected timezone cleared, got %q", saved.GetTimezone())
	}
}

// TestUserUpdateRejectsMissingEmail verifies that email remains required
// — the relaxation must not extend to the login identifier.
func TestUserUpdateRejectsMissingEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	ctx := context.Background()

	target := userstore.NewUser().
		SetEmail("target-noemail@test.com").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_USER).
		SetPassword("password123")
	if err := env.userStore.UserCreate(ctx, target); err != nil {
		t.Fatalf("seed target: %v", err)
	}

	resp := postJSON(t, env, url.Values{
		"controller": {"user-update"},
		"action":     {"user-update-ajax"},
	}, map[string]string{
		"user_id": target.GetID(),
		"status":  userstore.USER_STATUS_ACTIVE,
		"email":   "", // intentionally empty
	})

	if resp.Status != "error" {
		t.Fatalf("expected status=error for missing email, got %q", resp.Status)
	}
	if resp.Message != "Email is required" {
		t.Errorf("expected message 'Email is required', got %q", resp.Message)
	}
}

// TestUserUpdateRejectsMissingStatus verifies that status remains
// required — the relaxation must not extend to the account state field.
func TestUserUpdateRejectsMissingStatus(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	ctx := context.Background()

	target := userstore.NewUser().
		SetEmail("target-nostatus@test.com").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_USER).
		SetPassword("password123")
	if err := env.userStore.UserCreate(ctx, target); err != nil {
		t.Fatalf("seed target: %v", err)
	}

	resp := postJSON(t, env, url.Values{
		"controller": {"user-update"},
		"action":     {"user-update-ajax"},
	}, map[string]string{
		"user_id": target.GetID(),
		"status":  "", // intentionally empty
		"email":   "target-nostatus@test.com",
	})

	if resp.Status != "error" {
		t.Fatalf("expected status=error for missing status, got %q", resp.Status)
	}
	if resp.Message != "Status is required" {
		t.Errorf("expected message 'Status is required', got %q", resp.Message)
	}
}

// TestUserCreateAcceptsEmailOnly verifies that an admin can create a
// user with only an email — first/last name are optional. Guards the
// relaxation in user_manager/handle_user_create_ajax.go.
func TestUserCreateAcceptsEmailOnly(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	ctx := context.Background()

	resp := postJSON(t, env, url.Values{
		"controller": {"user-manager"},
		"action":     {"create-user-ajax"},
	}, map[string]string{
		"email": "new-email-only@test.com",
		// first_name, last_name intentionally omitted
	})

	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	userID, _ := resp.Data["user_id"].(string)
	if userID == "" {
		t.Fatal("response data.user_id is empty")
	}

	created, err := env.userStore.UserFindByID(ctx, userID)
	if err != nil {
		t.Fatalf("UserFindByID: %v", err)
	}
	if created == nil {
		t.Fatal("created user not found")
	}
	if created.GetEmail() != "new-email-only@test.com" {
		t.Errorf("expected email new-email-only@test.com, got %q", created.GetEmail())
	}
	if created.GetFirstName() != "" {
		t.Errorf("expected empty first name, got %q", created.GetFirstName())
	}
	if created.GetLastName() != "" {
		t.Errorf("expected empty last name, got %q", created.GetLastName())
	}
}

// TestUserCreateRejectsMissingEmail verifies that email remains required
// for user creation — the relaxation must not extend to the login
// identifier.
func TestUserCreateRejectsMissingEmail(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	resp := postJSON(t, env, url.Values{
		"controller": {"user-manager"},
		"action":     {"create-user-ajax"},
	}, map[string]string{
		"email": "", // intentionally empty
	})

	if resp.Status != "error" {
		t.Fatalf("expected status=error for missing email, got %q", resp.Status)
	}
	if resp.Message != "Email is required" {
		t.Errorf("expected message 'Email is required', got %q", resp.Message)
	}
}
