package ui

import "testing"

func TestFormatHelpColumnsHandlesShortColumnList(t *testing.T) {
	t.Parallel()

	m := BasePagerModel{Common: &CommonModel{Width: 80, Height: 24}}

	if m.FormatHelpColumns([]string{"copy", "edit"}) == "" {
		t.Fatal("expected non-empty help output")
	}
}

func TestHelpColumnValueBounds(t *testing.T) {
	t.Parallel()

	columns := []string{"a", "b"}

	if got := helpColumnValue(columns, 0); got != "a" {
		t.Fatalf("expected first value to be 'a', got %q", got)
	}

	if got := helpColumnValue(columns, 5); got != "" {
		t.Fatalf("expected out-of-range value to be empty, got %q", got)
	}
}
