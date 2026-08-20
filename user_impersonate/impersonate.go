package user_impersonate

import (
	"errors"
	"net/http"

	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dromara/carbon/v2"
)

// Impersonate creates a new session for the given user ID and sets the
// auth cookie on the response. The secure flag controls whether the
// cookie is marked Secure (use false for HTTP development).
func Impersonate(ss sessionstore.StoreInterface, w http.ResponseWriter, r *http.Request, userID string, secure bool) error {
	if ss == nil {
		return errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(userID).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.GetIP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	err := ss.SessionCreate(r.Context(), session)

	if err != nil {
		return err
	}

	auth.AuthCookieSet(w, r, session.GetKey(), types.WithSecure(secure))

	return nil
}
