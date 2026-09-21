package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dracory/userstore"
)

// seedUsers populates the user store with sample users so the admin
// panel has rows to exercise pagination, filtering, and sorting.
//
// Creates 40 users with mixed statuses and varied names/emails.
func seedUsers(store userstore.StoreInterface, logger *slog.Logger) {
	ctx := context.Background()

	statuses := []string{
		userstore.USER_STATUS_ACTIVE,
		userstore.USER_STATUS_ACTIVE,
		userstore.USER_STATUS_ACTIVE,
		userstore.USER_STATUS_INACTIVE,
		userstore.USER_STATUS_UNVERIFIED,
	}

	firstNames := []string{"Alice", "Bob", "Carol", "David", "Eve", "Frank", "Grace", "Heidi"}
	lastNames := []string{"Smith", "Jones", "Brown", "Miller", "Davis", "Wilson", "Taylor", "Clark"}

	for i := 1; i <= 40; i++ {
		status := statuses[(i-1)%len(statuses)]
		first := firstNames[(i-1)%len(firstNames)]
		last := lastNames[(i-1)%len(lastNames)]

		user := userstore.NewUser().
			SetFirstName(first).
			SetLastName(last).
			SetEmail(fmt.Sprintf("%s.%s.%02d@example.com", first, last, i)).
			SetStatus(status).
			SetRole(userstore.USER_ROLE_USER).
			SetCountry("US")

		if err := store.UserCreate(ctx, user); err != nil {
			logger.Error("seedUsers: failed to create user", "index", i, "error", err)
		}
	}

	// Add one administrator so the impersonate controller has a
	// privileged user to act as.
	admin := userstore.NewUser().
		SetFirstName("Admin").
		SetLastName("Root").
		SetEmail("admin@example.com").
		SetStatus(userstore.USER_STATUS_ACTIVE).
		SetRole(userstore.USER_ROLE_ADMINISTRATOR).
		SetCountry("US")
	if err := store.UserCreate(ctx, admin); err != nil {
		logger.Error("seedUsers: failed to create admin user", "error", err)
	}

	// Seed a few roles and groups so the manager pages have data.
	for _, r := range []struct{ name, handle string }{
		{"Editor", "editor"},
		{"Moderator", "moderator"},
		{"Support", "support"},
	} {
		role := userstore.NewRole().SetName(r.name).SetHandle(r.handle).SetStatus(userstore.ROLE_STATUS_ACTIVE)
		if err := store.RoleCreate(ctx, role); err != nil {
			logger.Error("seedUsers: failed to create role", "handle", r.handle, "error", err)
		}
	}

	for _, g := range []struct{ name, handle string }{
		{"Customers", "customers"},
		{"Partners", "partners"},
		{"Staff", "staff"},
	} {
		group := userstore.NewGroup().SetName(g.name).SetHandle(g.handle).SetStatus(userstore.GROUP_STATUS_ACTIVE)
		if err := store.GroupCreate(ctx, group); err != nil {
			logger.Error("seedUsers: failed to create group", "handle", g.handle, "error", err)
		}
	}

	logger.Info("seedUsers complete", "users", 41)
}
