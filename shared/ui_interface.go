package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/userstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
// This follows the blogadmin/shopadmin pattern.
type UiInterface interface {
	UserStore() userstore.StoreInterface
	GeoResolver() GeoResolverInterface
	Logger() *slog.Logger
	OnUserImpersonate() OnUserImpersonateFunc
	OnUserSearch() OnUserSearchFunc
	OnUserUpdate() OnUserUpdateFunc
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
