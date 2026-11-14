package shared_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/dracory/useradmin/shared"
)

// UiBase implements UiInterface; verify it satisfies the interface.
func TestUiBaseImplementsUiInterface(t *testing.T) {
	var _ shared.UiInterface = shared.NewUiBase(shared.UiConfig{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
}
