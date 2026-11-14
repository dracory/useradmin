package useradmin_test

import (
	"errors"
	"testing"

	"github.com/dracory/useradmin"
)

func TestErrorSentinels(t *testing.T) {
	for _, err := range []error{
		useradmin.ErrUserStoreRequired,
		useradmin.ErrLoggerRequired,
		useradmin.ErrGeoResolverRequired,
	} {
		if err == nil || err.Error() == "" {
			t.Fatal("expected non-empty sentinel error")
		}
	}
	if !errors.Is(useradmin.ErrUserStoreRequired, useradmin.ErrUserStoreRequired) {
		t.Fatal("expected errors.Is to match sentinel")
	}
}
