package shared_test

import (
	"strings"
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestVueLoaderJS(t *testing.T) {
	for _, want := range []string{"loadVueIfNeeded", "cdn.jsdelivr.net/npm/vue@3", "callback"} {
		if !strings.Contains(shared.VueLoaderJS, want) {
			t.Fatalf("expected %q in VueLoaderJS", want)
		}
	}
}
