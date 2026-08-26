package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
// This follows the blogadmin/shopadmin pattern.
type UiInterface interface {
	UserStore() userstore.StoreInterface
	GeoResolver() GeoResolverInterface
	Logger() *slog.Logger
	SessionResolver() SessionResolverInterface
	BlindIndexFirstName() BlindIndexResolverInterface
	BlindIndexLastName() BlindIndexResolverInterface
	BlindIndexEmail() BlindIndexResolverInterface
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
