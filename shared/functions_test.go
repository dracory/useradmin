package shared_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestErrorPopup(t *testing.T) {
	if !strings.Contains(shared.ErrorPopup("bad").ToHTML(), "bad") {
		t.Fatal("expected message")
	}
}

func TestSuccessPopup(t *testing.T) {
	if !strings.Contains(shared.SuccessPopup("ok").ToHTML(), "ok") {
		t.Fatal("expected message")
	}
}

func TestSuccessPopupWithRedirect(t *testing.T) {
	html := shared.SuccessPopupWithRedirect("ok", "/next", 2).ToHTML()
	if !strings.Contains(html, "ok") || !strings.Contains(html, "/next") {
		t.Fatalf("unexpected: %q", html)
	}
}

func TestFlashErrorFallback(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	html := shared.FlashError(nil, rec, req, "bad", "/err", 3)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}
	if !strings.Contains(html, "/err") {
		t.Fatalf("expected link, got %q", html)
	}
}

func TestFlashSuccessFallback(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	html := shared.FlashSuccess(nil, rec, req, "ok", "/ok", 3)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}
	if !strings.Contains(html, "/ok") {
		t.Fatalf("expected link, got %q", html)
	}
}

func TestFlashRedirectCustom(t *testing.T) {
	called := false
	fn := func(w http.ResponseWriter, r *http.Request, messageType, message, redirectURL string, seconds int) string {
		called = true
		if messageType != "error" || message != "bad" {
			t.Fatal("wrong args")
		}
		return "custom"
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", http.NoBody)
	if got := shared.FlashError(fn, rec, req, "bad", "/err", 1); got != "custom" || !called {
		t.Fatal("expected custom flash redirect")
	}
}

func TestJSEscapeString(t *testing.T) {
	cases := map[string]string{
		"plain":  "plain",
		"it's":   `it\'s`,
		`a\b`:    `a\\b`,
		"a\nb":   `a\nb`,
		"a\rb":   `a\rb`,
		"a\tb":   `a\tb`,
		"');x()": `\');x()`,
	}
	for in, want := range cases {
		if got := shared.JSEscapeString(in); got != want {
			t.Fatalf("input %q: expected %q, got %q", in, want, got)
		}
	}
}
