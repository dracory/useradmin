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
	OnUserUpdateField      OnUserUpdateFunc
	OnUserDecodeField      OnUserDecodeFunc
	OnUserEncodeField      OnUserEncodeFunc
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
func (u UiBase) OnUserUpdate() OnUserUpdateFunc           { return u.OnUserUpdateField }
func (u UiBase) OnUserDecode() OnUserDecodeFunc           { return u.OnUserDecodeField }
func (u UiBase) OnUserEncode() OnUserEncodeFunc           { return u.OnUserEncodeField }
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
		OnUserUpdateField:      config.OnUserUpdate,
		OnUserDecodeField:      config.OnUserDecode,
		OnUserEncodeField:      config.OnUserEncode,
		FlashRedirectField:     config.FlashRedirect,
		LayoutField:            config.Layout,
	}
}
