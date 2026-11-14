package shared_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestNewUiBaseFields(t *testing.T) {
	base := shared.NewUiBase(shared.UiConfig{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if base.Logger() == nil {
		t.Fatal("expected logger")
	}
	if base.UserStore() != nil || base.OnUserImpersonate() != nil ||
		base.OnUserSearch() != nil || base.UserPiiSeal() != nil {
		t.Fatal("expected nil optionals")
	}
}
