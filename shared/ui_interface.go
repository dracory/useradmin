package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/geostore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
// This follows the blogadmin/shopadmin pattern.
type UiInterface interface {
	UserStore() userstore.StoreInterface
	GeoStore() geostore.StoreInterface
	Logger() *slog.Logger
	SessionStore() sessionstore.StoreInterface
	BlindIndexFirstName() blindindexstore.StoreInterface
	BlindIndexLastName() blindindexstore.StoreInterface
	BlindIndexEmail() blindindexstore.StoreInterface
	TaskStore() taskstore.StoreInterface
	BlindIndexRebuildTaskAlias() string
	VaultTokenizer() VaultTokenizer
	AuthUser(r *http.Request) userstore.UserInterface
	FlashRedirect() FlashRedirectFunc
	SecureCookie() bool

	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
