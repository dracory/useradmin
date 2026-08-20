package shared

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
)

// ErrorPopup returns a SweetAlert2 error popup tag.
func ErrorPopup(errorMessage string) hb.TagInterface {
	return hb.Swal(hb.SwalOptions{
		Title:            "Error",
		Text:             errorMessage,
		Icon:             "error",
		Timer:            10000,
		TimerProgressBar: true,
	})
}

// SuccessPopup returns a SweetAlert2 success popup tag.
func SuccessPopup(successMessage string) hb.TagInterface {
	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             successMessage,
		Icon:             "success",
		Timer:            10000,
		TimerProgressBar: true,
	})
}

// SuccessPopupWithRedirect returns a SweetAlert2 success popup with an
// optional redirect. If redirectUrl is empty, no redirect is configured.
func SuccessPopupWithRedirect(successMessage string, redirectUrl string, redirectSeconds int) hb.TagInterface {
	if redirectUrl != "" {
		return hb.Swal(hb.SwalOptions{
			Title:            "Success",
			Text:             successMessage,
			Icon:             "success",
			Timer:            redirectSeconds * 1000,
			TimerProgressBar: true,
			RedirectURL:      redirectUrl,
			RedirectSeconds:  redirectSeconds,
		})
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             successMessage,
		Icon:             "success",
		Timer:            redirectSeconds * 1000,
		TimerProgressBar: true,
	})
}

// FlashRedirectFunc redirects the user with a flash message. Host
// projects that have a flash-message system (cache store + /flash
// route) provide this callback so useradmin can surface messages
// across redirects. If nil, useradmin falls back to a plain
// http.Redirect and the message is dropped.
//
// messageType is one of "error", "success", "info", "warning".
type FlashRedirectFunc func(w http.ResponseWriter, r *http.Request, messageType, message, redirectURL string, seconds int) string

// FlashError performs a flash redirect for an error message. If
// flashRedirect is nil, it falls back to a plain http.Redirect.
func FlashError(flashRedirect FlashRedirectFunc, w http.ResponseWriter, r *http.Request, message, redirectURL string, seconds int) string {
	if flashRedirect != nil {
		return flashRedirect(w, r, "error", message, redirectURL, seconds)
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return `<a href="` + redirectURL + `">See Other</a>`
}

// FlashSuccess performs a flash redirect for a success message. If
// flashRedirect is nil, it falls back to a plain http.Redirect.
func FlashSuccess(flashRedirect FlashRedirectFunc, w http.ResponseWriter, r *http.Request, message, redirectURL string, seconds int) string {
	if flashRedirect != nil {
		return flashRedirect(w, r, "success", message, redirectURL, seconds)
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return `<a href="` + redirectURL + `">See Other</a>`
}

// JSEscapeString escapes a string for safe embedding inside a JavaScript
// single-quoted string literal. It escapes backslash, single quote, and
// newlines so that user-supplied values (e.g. user IDs from URL params)
// cannot break out of the string context.
func JSEscapeString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
