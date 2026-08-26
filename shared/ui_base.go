package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/userstore"
)

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs can embed this to get all the accessor
// methods for free, following the blogadmin/shopadmin pattern.
type UiBase struct {
	UserStoreField         userstore.StoreInterface
	GeoResolverField       GeoResolverInterface
	LoggerField            *slog.Logger
	OnUserImpersonateField OnUserImpersonateFunc
	OnUserSearchField      OnUserSearchFunc
	OnUserUpdatedField      OnUserUpdatedFunc
	UserPiiSealField       UserPiiSealFunc
	UserPiiUnsealField     UserPiiUnsealFunc
	UsersPiiUnsealField    UsersPiiUnsealFunc
	FlashRedirectField     FlashRedirectFunc
	LayoutField            func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

func (u UiBase) UserStore() userstore.StoreInterface      { return u.UserStoreField }
func (u UiBase) GeoResolver() GeoResolverInterface        { return u.GeoResolverField }
func (u UiBase) Logger() *slog.Logger                     { return u.LoggerField }
func (u UiBase) OnUserImpersonate() OnUserImpersonateFunc { return u.OnUserImpersonateField }
func (u UiBase) OnUserSearch() OnUserSearchFunc           { return u.OnUserSearchField }
func (u UiBase) OnUserUpdated() OnUserUpdatedFunc           { return u.OnUserUpdatedField }
func (u UiBase) UserPiiSeal() UserPiiSealFunc             { return u.UserPiiSealField }
func (u UiBase) UserPiiUnseal() UserPiiUnsealFunc         { return u.UserPiiUnsealField }
func (u UiBase) UsersPiiUnseal() UsersPiiUnsealFunc       { return u.UsersPiiUnsealField }
func (u UiBase) FlashRedirect() FlashRedirectFunc         { return u.FlashRedirectField }

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
		UserStoreField:         config.UserStore,
		GeoResolverField:       config.GeoResolver,
		LoggerField:            config.Logger,
		OnUserImpersonateField: config.OnUserImpersonate,
		OnUserSearchField:      config.OnUserSearch,
		OnUserUpdatedField:      config.OnUserUpdated,
		UserPiiSealField:       config.UserPiiSeal,
		UserPiiUnsealField:     config.UserPiiUnseal,
		UsersPiiUnsealField:    config.UsersPiiUnseal,
		FlashRedirectField:     config.FlashRedirect,
		LayoutField:            config.Layout,
	}
}
