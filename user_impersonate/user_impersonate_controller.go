package user_impersonate

import (
	"log/slog"
	"net/http"

	"github.com/dracory/useradmin/shared"
	"github.com/dracory/userstore"

	"github.com/dracory/req"
)

// UiInterface defines the user impersonate controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserImpersonate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user impersonate controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// UserImpersonate handles the user impersonate controller requests
func (u *ui) UserImpersonate(w http.ResponseWriter, r *http.Request) {
	u.Handler(w, r)
}

// Handler processes the user impersonate controller request
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) {
	linksHelper := shared.NewLinksFromRequest(r)
	usersURL := linksHelper.UserManager(nil)
	userHomeURL := shared.UserHomeURL(r)
	if userHomeURL == "" {
		userHomeURL = shared.AdminHomeURL(r)
	}

	authUser := u.AuthUser(r)

	if authUser == nil {
		shared.FlashError(u.FlashRedirect(), w, r, "User not found", usersURL, 15)
		return
	}

	if !authUser.IsAdministrator() {
		shared.FlashError(u.FlashRedirect(), w, r, "Not authorized", usersURL, 15)
		return
	}

	userID := req.GetStringTrimmed(r, "user_id")

	if userID == "" {
		shared.FlashError(u.FlashRedirect(), w, r, "User ID not found", usersURL, 15)
		return
	}

	// Verify the target user exists and is active before creating a
	// session. Prevents impersonating deleted/non-existent users.
	if u.UserStore() == nil {
		shared.FlashError(u.FlashRedirect(), w, r, "User store not configured", usersURL, 15)
		return
	}

	targetUser, err := u.UserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userImpersonateController UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		shared.FlashError(u.FlashRedirect(), w, r, "Error loading user", usersURL, 15)
		return
	}
	if targetUser == nil {
		shared.FlashError(u.FlashRedirect(), w, r, "User not found", usersURL, 15)
		return
	}
	if targetUser.GetStatus() != userstore.USER_STATUS_ACTIVE {
		shared.FlashError(u.FlashRedirect(), w, r, "Cannot impersonate an inactive user", usersURL, 15)
		return
	}

	// Prevent self-impersonation (no-op that confuses the session).
	if authUser != nil && authUser.GetID() == userID {
		shared.FlashError(u.FlashRedirect(), w, r, "You cannot impersonate yourself", usersURL, 15)
		return
	}

	err = Impersonate(u.SessionStore(), w, r, userID, u.SecureCookie())

	if err != nil {
		shared.FlashError(u.FlashRedirect(), w, r, err.Error(), usersURL, 15)
		return
	}

	shared.FlashSuccess(u.FlashRedirect(), w, r, "Impersonation is successful", userHomeURL, 15)
}
