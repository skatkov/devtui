package base64decoder

import (
	"testing"

	"github.com/skatkov/devtui/internal/ui"
)

func TestHelpViewDoesNotPanic(t *testing.T) {
	t.Parallel()

	m := NewBase64Model(&ui.CommonModel{Width: 80, Height: 24})
	if m.helpView() == "" {
		t.Fatal("expected non-empty help view")
	}
}
