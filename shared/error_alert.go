package shared

import "github.com/dracory/hb"

// ErrorAlert returns an inline HTML error alert for the given message.
// This replaces the former ToFlashError flash-message pattern, which
// required a cache store and a /flash route handler that did not exist
// in the standalone module.
func ErrorAlert(message string) string {
	return hb.Div().
		Class("alert alert-danger").
		Style("margin: 20px;").
		HTML(message).
		ToHTML()
}
