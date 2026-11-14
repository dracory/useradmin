package shared_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/dracory/useradmin/shared"
)

// UiConfig is a plain options struct; verify fields are settable.
func TestUiConfigFields(t *testing.T) {
	cfg := shared.UiConfig{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	if cfg.Logger == nil {
		t.Fatal("expected logger")
	}
}
