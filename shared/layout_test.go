package shared_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestLayoutFallback(t *testing.T) {
	rec := httptest.NewRecorder()
	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		Styles:     []string{"body{margin:0}"},
		StyleURLs:  []string{"https://x/s.css"},
		Scripts:    []string{"console.log(1)"},
		ScriptURLs: []string{"https://x/j.js"},
	}
	html := shared.Layout(rec, httptest.NewRequest("GET", "/x", nil), "Title", "<p>Body</p>", options)
	for _, want := range []string{"Title", "<p>Body</p>", "body{margin:0}", "https://x/s.css", "console.log(1)", "https://x/j.js"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in layout", want)
		}
	}
}
