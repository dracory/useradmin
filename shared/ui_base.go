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

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs can embed this to get all the accessor
// methods for free, following the blogadmin/shopadmin pattern.
type UiBase struct {
	UserStoreField              userstore.StoreInterface
	GeoStoreField               geostore.StoreInterface
	LoggerField                 *slog.Logger
	SessionStoreField           sessionstore.StoreInterface
	BlindIndexFirstNameField    blindindexstore.StoreInterface
	BlindIndexLastNameField     blindindexstore.StoreInterface
	BlindIndexEmailField        blindindexstore.StoreInterface
	TaskStoreField              taskstore.StoreInterface
	BlindIndexRebuildTaskAliasField string
	VaultTokenizerField         VaultTokenizer
	AuthUserField               func(r *http.Request) userstore.UserInterface
	FlashRedirectField          FlashRedirectFunc
	SecureCookieField           bool
	LayoutField                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

func (u UiBase) UserStore() userstore.StoreInterface      { return u.UserStoreField }
func (u UiBase) GeoStore() geostore.StoreInterface        { return u.GeoStoreField }
func (u UiBase) Logger() *slog.Logger                     { return u.LoggerField }
func (u UiBase) SessionStore() sessionstore.StoreInterface {
	return u.SessionStoreField
}
func (u UiBase) BlindIndexFirstName() blindindexstore.StoreInterface {
	return u.BlindIndexFirstNameField
}
func (u UiBase) BlindIndexLastName() blindindexstore.StoreInterface {
	return u.BlindIndexLastNameField
}
func (u UiBase) BlindIndexEmail() blindindexstore.StoreInterface {
	return u.BlindIndexEmailField
}
func (u UiBase) TaskStore() taskstore.StoreInterface { return u.TaskStoreField }
func (u UiBase) BlindIndexRebuildTaskAlias() string  { return u.BlindIndexRebuildTaskAliasField }
func (u UiBase) VaultTokenizer() VaultTokenizer      { return u.VaultTokenizerField }
func (u UiBase) AuthUser(r *http.Request) userstore.UserInterface {
	if u.AuthUserField == nil {
		return nil
	}
	return u.AuthUserField(r)
}
func (u UiBase) FlashRedirect() FlashRedirectFunc { return u.FlashRedirectField }
func (u UiBase) SecureCookie() bool               { return u.SecureCookieField }

func (u UiBase) Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	return u.LayoutField(w, r, webpageTitle, webpageHtml, options)
}

// NewUiBase creates a UiBase from a UiConfig
func NewUiBase(config UiConfig) UiBase {
	return UiBase{
		UserStoreField:                 config.UserStore,
		GeoStoreField:                  config.GeoStore,
		LoggerField:                    config.Logger,
		SessionStoreField:              config.SessionStore,
		BlindIndexFirstNameField:       config.BlindIndexFirstName,
		BlindIndexLastNameField:        config.BlindIndexLastName,
		BlindIndexEmailField:           config.BlindIndexEmail,
		TaskStoreField:                 config.TaskStore,
		BlindIndexRebuildTaskAliasField: config.BlindIndexRebuildTaskAlias,
		VaultTokenizerField:            config.VaultTokenizer,
		AuthUserField:                  config.AuthUser,
		FlashRedirectField:             config.FlashRedirect,
		SecureCookieField:              config.SecureCookie,
		LayoutField:                    config.Layout,
	}
}
