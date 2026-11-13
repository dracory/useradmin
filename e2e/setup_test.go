// Package e2e contains end-to-end tests that drive the useradmin panel
// with a real headless Chrome via chromedp (CDP). The admin is served
// in-process by httptest on a random port — no binary build needed.
package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dracory/useradmin"
	"github.com/dracory/userstore"
	_ "modernc.org/sqlite"
)

var (
	baseURL        string
	allocCtx       context.Context
	testStore      userstore.StoreInterface
	seedUserID     string
	seedUserMail   string
	seedRoleID     string
	seedRoleName   string
	seedGroupID    string
	seedGroupName  string
	seedMemberMail string
)

const adminPath = "/admin/users"

func TestMain(m *testing.M) {
	// In-process admin server backed by a seeded in-memory SQLite DB.
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		panic("e2e: failed to open db: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		panic("e2e: failed to ping db: " + err.Error())
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
		panic("e2e: failed to init store: " + err.Error())
	}
	testStore = store

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	seedDB(store, logger)

	admin := userstore.NewUser().
		SetFirstName("Admin").
		SetLastName("Root").
		SetEmail("admin@example.com").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_ADMINISTRATOR).
		SetCountry("US")
	if err := store.UserCreate(context.Background(), admin); err != nil {
		panic("e2e: failed to create seed admin user: " + err.Error())
	}
	found, err := store.UserFindByEmail(context.Background(), "admin@example.com")
	if err != nil || found == nil {
		panic("e2e: seeded admin user not found: " + fmt.Sprint(err))
	}
	seedUserID = found.GetID()
	seedUserMail = found.GetEmail()

	// Grab seeded role/group IDs and assign members so the members
	// pages render at least two table rows (the waitTable signal).
	roles, err := store.RoleList(context.Background(), userstore.NewRoleQuery())
	if err != nil || len(roles) == 0 {
		panic("e2e: no seeded roles: " + fmt.Sprint(err))
	}
	groups, err := store.GroupList(context.Background(), userstore.NewGroupQuery())
	if err != nil || len(groups) == 0 {
		panic("e2e: no seeded groups: " + fmt.Sprint(err))
	}
	seedRoleID = roles[0].GetID()
	seedRoleName = roles[0].GetName()
	seedGroupID = groups[0].GetID()
	seedGroupName = groups[0].GetName()

	users, err := store.UserList(context.Background(), userstore.NewUserQuery())
	if err != nil || len(users) < 2 {
		panic("e2e: not enough seeded users: " + fmt.Sprint(err))
	}
	seedMemberMail = users[0].GetEmail()
	for _, u := range users[:2] {
		ur := userstore.NewUserRole().SetUserID(u.GetID()).SetRoleID(seedRoleID)
		if err := store.UserRoleCreate(context.Background(), ur); err != nil {
			panic("e2e: failed to assign user role: " + err.Error())
		}
		ug := userstore.NewUserGroup().SetUserID(u.GetID()).SetGroupID(seedGroupID)
		if err := store.UserGroupCreate(context.Background(), ug); err != nil {
			panic("e2e: failed to assign user group: " + err.Error())
		}
	}

	panel, err := useradmin.New(useradmin.AdminOptions{
		UserStore:         store,
		GeoResolver:       &testGeoResolver{},
		Logger:            logger,
		OnUserImpersonate: testOnUserImpersonate,
		AdminHomeURL:      "/admin",
		UserAdminURL:      adminPath,
	})
	if err != nil {
		panic("e2e: failed to create useradmin: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle(adminPath, http.HandlerFunc(panel.Handle))
	mux.Handle(adminPath+"/", http.HandlerFunc(panel.Handle))
	srv := httptest.NewServer(mux)
	baseURL = srv.URL

	// One shared browser process for the entire test run.
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Headless)
	if bin := chromeBin(); bin != "" {
		opts = append(opts, chromedp.ExecPath(bin))
	}
	var cancel context.CancelFunc
	allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)

	code := m.Run()

	cancel()
	srv.Close()
	os.Exit(code)
}

// seedDB populates the store with users, roles, and groups so the
// manager pages have rows to render. Mirrors example/seed.go but with
// fewer users — e2e only needs a couple of table pages.
func seedDB(store userstore.StoreInterface, logger *slog.Logger) {
	ctx := context.Background()

	firstNames := []string{"Alice", "Bob", "Carol", "David", "Eve"}
	lastNames := []string{"Smith", "Jones", "Brown", "Miller", "Davis"}
	for i := 1; i <= 25; i++ {
		first := firstNames[(i-1)%len(firstNames)]
		last := lastNames[(i-1)%len(lastNames)]
		user := userstore.NewUser().
			SetFirstName(first).
			SetLastName(last).
			SetEmail(fmt.Sprintf("%s.%s.%02d@example.com", first, last, i)).
			SetStatus(userstore.USER_STATUS_ACTIVE).
			SetRole(userstore.USER_ROLE_USER).
			SetCountry("US")
		if err := store.UserCreate(ctx, user); err != nil {
			logger.Error("seedDB: failed to create user", "index", i, "error", err)
		}
	}

	for _, r := range []struct{ name, handle string }{
		{"Editor", "editor"},
		{"Moderator", "moderator"},
		{"Support", "support"},
	} {
		role := userstore.NewRole().SetName(r.name).SetHandle(r.handle).SetStatus(userstore.ROLE_STATUS_ACTIVE)
		if err := store.RoleCreate(ctx, role); err != nil {
			logger.Error("seedDB: failed to create role", "handle", r.handle, "error", err)
		}
	}

	for _, g := range []struct{ name, handle string }{
		{"Customers", "customers"},
		{"Partners", "partners"},
		{"Staff", "staff"},
	} {
		group := userstore.NewGroup().SetName(g.name).SetHandle(g.handle).SetStatus(userstore.GROUP_STATUS_ACTIVE)
		if err := store.GroupCreate(ctx, group); err != nil {
			logger.Error("seedDB: failed to create group", "handle", g.handle, "error", err)
		}
	}
}

// testGeoResolver satisfies useradmin.GeoResolverInterface with a
// static list — no external geo database required.
type testGeoResolver struct{}

var _ useradmin.GeoResolverInterface = (*testGeoResolver)(nil)

func (r *testGeoResolver) Countries(ctx context.Context) ([]useradmin.Country, error) {
	return []useradmin.Country{
		{IsoCode2: "US", Name: "United States"},
		{IsoCode2: "BG", Name: "Bulgaria"},
	}, nil
}

func (r *testGeoResolver) Timezones(ctx context.Context, countryCode ...string) ([]useradmin.Timezone, error) {
	if len(countryCode) == 0 || countryCode[0] == "" {
		return nil, nil
	}
	if countryCode[0] == "US" {
		return []useradmin.Timezone{{Code: "America/New_York"}}, nil
	}
	return nil, nil
}

// testOnUserImpersonate sets a session cookie, like the example does.
// Its presence registers the impersonate route and shows the button.
func testOnUserImpersonate(w http.ResponseWriter, r *http.Request, userID string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "e2e_session",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
	})
	return nil
}

// newTab opens a fresh tab in the shared browser and returns its
// context. Each tab gets a 30s deadline; cleanup closes the tab.
func newTab(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := chromedp.NewContext(allocCtx)
	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	t.Cleanup(func() { cancelTimeout(); cancel() })
	return ctx
}

// chromeBin resolves a headless Chromium binary. Order: CHROME_BIN env,
// Playwright's downloaded headless shell, empty (chromedp finds Chrome
// on PATH / standard install locations).
func chromeBin() string {
	if b := os.Getenv("CHROME_BIN"); b != "" {
		return b
	}
	patterns := []string{
		// Windows
		filepath.Join(os.Getenv("LOCALAPPDATA"),
			"ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe"),
		// WSL: Windows-side Playwright shell via interop
		"/mnt/*/Users/*/AppData/Local/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe",
		// Linux
		filepath.Join(os.Getenv("HOME"),
			".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"),
	}
	for _, pattern := range patterns {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return ""
}
