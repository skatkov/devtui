package numbers

import (
	"regexp"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/skatkov/devtui/internal/ui"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestNumbersConvertsDecimalByDefault(t *testing.T) {
	t.Parallel()

	model := newTestModel(30)
	model = updateModel(model, codeKeypress(tea.KeyEnter))
	model = typeText(model, "255")
	model = updateModel(model, codeKeypress(tea.KeyEnter))

	m := model.(NumbersModel)
	if m.form.State != huh.StateCompleted {
		t.Fatalf("expected completed form, got state %v", m.form.State)
	}
	if m.result.Base != 10 || m.result.Value != 255 {
		t.Fatalf("expected decimal 255, got base %d value %d", m.result.Base, m.result.Value)
	}

	view := m.View().Content
	if !strings.Contains(view, "11111111") || !strings.Contains(view, "ff") {
		t.Fatalf("expected conversions in view, got %q", view)
	}
	if got := lipgloss.Width(view); got != 106 {
		t.Fatalf("expected 100 column table inside padding, got width %d", got)
	}
}

func TestNumbersBaseSelectionWrapsToHexadecimal(t *testing.T) {
	t.Parallel()

	model := newTestModel(30)
	for range 3 {
		model = updateModel(model, codeKeypress(tea.KeyUp))
	}
	model = updateModel(model, codeKeypress(tea.KeyEnter))
	model = typeText(model, "ff")
	model = updateModel(model, codeKeypress(tea.KeyEnter))

	m := model.(NumbersModel)
	if m.result.Base != 16 || m.result.Value != 255 {
		t.Fatalf("expected hexadecimal 255, got base %d value %d", m.result.Base, m.result.Value)
	}
}

func TestNumbersPlaceholderNamesSelectedBase(t *testing.T) {
	t.Parallel()

	model := updateModel(newTestModel(30), tea.WindowSizeMsg{Width: 80, Height: 30})
	model = updateModel(model, codeKeypress(tea.KeyUp))
	model = updateModel(model, codeKeypress(tea.KeyEnter))

	view := ansiPattern.ReplaceAllString(model.View().Content, "")
	if !strings.Contains(view, "Enter a Base 8 (octal) number") {
		t.Fatalf("expected octal placeholder, got %q", view)
	}
}

func TestNumbersInputValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantMsg string
	}{
		{"empty input", "", "number cannot be empty"},
		{"invalid digits", "12a", "please enter a valid Base 10 (decimal) number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			model := newTestModel(30)
			model = updateModel(model, codeKeypress(tea.KeyEnter))
			model = typeText(model, tt.input)
			model = updateModel(model, codeKeypress(tea.KeyEnter))

			m := model.(NumbersModel)
			if m.form.State == huh.StateCompleted {
				t.Fatal("expected form to stay open")
			}
			if !strings.Contains(m.View().Content, tt.wantMsg) {
				t.Fatalf("expected %q in view, got %q", tt.wantMsg, m.View().Content)
			}
		})
	}
}

func TestNumbersFormViewLayout(t *testing.T) {
	t.Parallel()

	view := newTestModel(30).View().Content
	if got := lipgloss.Height(view); got != 30 {
		t.Fatalf("expected view to fill 30 lines, got %d", got)
	}
	if !strings.Contains(view, "\x1b[1mNumber Base Converter") {
		t.Fatalf("expected bold title, got %q", view)
	}

	lines := strings.Split(ansiPattern.ReplaceAllString(view, ""), "\n")
	if strings.TrimSpace(lines[2]) != "" {
		t.Fatalf("expected blank line above the form, got %q", lines[2])
	}
	if !strings.HasPrefix(lines[3], "  ┃ Select Base") {
		t.Fatalf("expected form at the left edge, got %q", lines[3])
	}
	if got := strings.Count(view, "enter"); got != 1 {
		t.Fatalf("expected help to be shown once, got %d times", got)
	}
}

func TestNumbersWindowSizeMsgUpdatesCommon(t *testing.T) {
	t.Parallel()

	common := &ui.CommonModel{Styles: ui.NewStyle()}
	NewNumberModel(common).Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	if common.Width != 100 || common.Height != 40 {
		t.Fatalf("expected 100x40, got %dx%d", common.Width, common.Height)
	}
}

func TestNumbersEscReturnsToList(t *testing.T) {
	t.Parallel()

	common := &ui.CommonModel{Width: 80, Height: 24, Styles: ui.NewStyle()}
	_, cmd := NewNumberModel(common).Update(codeKeypress(tea.KeyEscape))

	msg, ok := cmd().(ui.ReturnToListMsg)
	if !ok || msg.Common != common {
		t.Fatalf("expected return to list with the shared common model, got %#v", msg)
	}
}

func TestNumbersQuitKeys(t *testing.T) {
	t.Parallel()

	for _, key := range []tea.KeyPressMsg{
		tea.KeyPressMsg(tea.Key{Code: 'q', Text: "q"}),
		tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}),
	} {
		_, cmd := newTestModel(24).Update(key)
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("expected quit for %q", key.String())
		}
	}
}

func newTestModel(height int) tea.Model {
	model := NewNumberModel(&ui.CommonModel{Width: 80, Height: height, Styles: ui.NewStyle()})
	return batchUpdate(model, model.Init())
}

func updateModel(model tea.Model, msg tea.Msg) tea.Model {
	nextModel, cmd := model.Update(msg)
	return batchUpdate(nextModel, cmd)
}

func batchUpdate(model tea.Model, cmd tea.Cmd) tea.Model {
	return runCmd(model, cmd, 0)
}

// runCmd feeds the messages of a few command levels back into the model and
// drops commands that wait, such as the cursor blink tick.
func runCmd(model tea.Model, cmd tea.Cmd, depth int) tea.Model {
	if cmd == nil || depth > 3 {
		return model
	}

	done := make(chan tea.Msg, 1)
	go func() { done <- cmd() }()

	var msg tea.Msg
	select {
	case msg = <-done:
	case <-time.After(200 * time.Millisecond):
		return model
	}

	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			model = runCmd(model, c, depth+1)
		}
		return model
	}

	nextModel, nextCmd := model.Update(msg)
	return runCmd(nextModel, nextCmd, depth+1)
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func typeText(model tea.Model, text string) tea.Model {
	for _, r := range text {
		model = updateModel(model, tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)}))
	}
	return model
}
