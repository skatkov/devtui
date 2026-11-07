package uuidgenerate

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
	"github.com/skatkov/devtui/internal/ui"
)

func TestUUIDVersionSelectRespondsToKeyPresses(t *testing.T) {
	t.Parallel()

	common := &ui.CommonModel{Width: 80, Height: 24}
	common.Styles = ui.NewStyle()

	model := NewUUIDGenerateModel(common)
	model = batchUpdate(model, model.Init()).(*UUIDGenerate)
	model = updateModel(model, tea.WindowSizeMsg{Width: 80, Height: 24}).(*UUIDGenerate)

	if model.version != 1 {
		t.Fatalf("expected initial version 1, got %d", model.version)
	}

	model = updateModel(model, codeKeypress(tea.KeyDown)).(*UUIDGenerate)
	if model.version != 2 {
		t.Fatalf("expected version 2 after pressing down, got %d", model.version)
	}

	model = updateModel(model, codeKeypress(tea.KeyEnter)).(*UUIDGenerate)
	if model.form.State != huh.StateCompleted {
		t.Fatalf("expected form to complete after pressing enter, got state %v", model.form.State)
	}
	if model.generatedUUID == uuid.Nil {
		t.Fatal("expected UUID to be generated after form completion")
	}
}

func TestUUIDVersionSelectOptions(t *testing.T) {
	t.Parallel()

	for want := 1; want <= 7; want++ {
		model := newTestModel(24)
		for range want - 1 {
			model = updateModel(model, codeKeypress(tea.KeyDown)).(*UUIDGenerate)
		}
		if model.version != want {
			t.Fatalf("expected version %d after %d presses, got %d", want, want-1, model.version)
		}
	}
}

func TestHideNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		version int
		want    bool
	}{
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{6, true},
		{7, true},
	}

	for _, tt := range tests {
		model := &UUIDGenerate{version: tt.version}
		if got := model.hideNamespace(); got != tt.want {
			t.Errorf("hideNamespace() for version %d = %v, want %v", tt.version, got, tt.want)
		}
	}
}

func TestUUIDGenerateViewFillsWindow(t *testing.T) {
	t.Parallel()

	model := newTestModel(30)
	view := model.View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	if height := lipgloss.Height(view.Content); height != 30 {
		t.Fatalf("expected view height 30, got %d", height)
	}
	if !strings.Contains(view.Content, lipgloss.NewStyle().Bold(true).Render(Title)) {
		t.Fatalf("expected bold title, got %q", view.Content)
	}

	lines := strings.Split(view.Content, "\n")
	if strings.TrimSpace(lines[2]) != "" {
		t.Fatalf("expected blank line above the form, got %q", lines[2])
	}
	if indent := lipgloss.Width(lines[3][:strings.Index(lines[3], "┃")]); indent != 2 {
		t.Fatalf("expected form indent 2, got %d in %q", indent, lines[3])
	}
	if count := strings.Count(view.Content, "submit"); count != 1 {
		t.Fatalf("expected help to be shown once, got %d", count)
	}
}

func TestUUIDGenerateWithNamespace(t *testing.T) {
	t.Parallel()

	model := newTestModel(24)
	for range 4 {
		model = updateModel(model, codeKeypress(tea.KeyDown)).(*UUIDGenerate)
	}
	model = updateModel(model, codeKeypress(tea.KeyEnter)).(*UUIDGenerate)
	if model.form.State != huh.StateNormal {
		t.Fatalf("expected namespace input before completion, got state %v", model.form.State)
	}
	if model.generatedUUID != uuid.Nil {
		t.Fatalf("expected no UUID before completion, got %v", model.generatedUUID)
	}

	for _, r := range "example.com" {
		model.Update(tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)}))
	}
	model = updateModel(model, codeKeypress(tea.KeyEnter)).(*UUIDGenerate)

	want := uuid.NewSHA1(uuid.NameSpaceURL, []byte("example.com"))
	if model.generatedUUID != want {
		t.Fatalf("expected UUID %v, got %v", want, model.generatedUUID)
	}
	view := model.View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	for _, want := range []string{"Version", "5", "Namespace", "example.com", want.String()} {
		if !strings.Contains(view.Content, want) {
			t.Fatalf("expected view to contain %q, got %q", want, view.Content)
		}
	}
	if width := lipgloss.Width(strings.TrimSpace(strings.Split(view.Content, "\n")[1])); width != 100 {
		t.Fatalf("expected table width 100, got %d", width)
	}
}

func TestUUIDGenerateWithoutNamespaceOmitsRow(t *testing.T) {
	t.Parallel()

	model := updateModel(newTestModel(24), codeKeypress(tea.KeyEnter)).(*UUIDGenerate)
	content := model.View().Content
	if strings.Contains(content, "Namespace") {
		t.Fatalf("expected no namespace row, got %q", content)
	}
	if !strings.Contains(content, model.generatedUUID.String()) {
		t.Fatalf("expected generated UUID in view, got %q", content)
	}
}

func TestUUIDGenerateNamespaceAcceptsQ(t *testing.T) {
	t.Parallel()

	model := newTestModel(24)
	for range 4 {
		model = updateModel(model, codeKeypress(tea.KeyDown)).(*UUIDGenerate)
	}
	model = updateModel(model, codeKeypress(tea.KeyEnter)).(*UUIDGenerate)

	for _, r := range "quux" {
		_, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)}))
		if cmd != nil {
			if _, ok := cmd().(tea.QuitMsg); ok {
				t.Fatalf("expected %q to be typed, got quit", r)
			}
		}
	}
	if model.namespace != "quux" {
		t.Fatalf("expected namespace %q, got %q", "quux", model.namespace)
	}
}

func TestUUIDGenerateKeepsUUIDAfterCompletion(t *testing.T) {
	t.Parallel()

	model := updateModel(newTestModel(24), codeKeypress(tea.KeyEnter)).(*UUIDGenerate)
	generated := model.generatedUUID

	model = updateModel(model, tea.WindowSizeMsg{Width: 100, Height: 30}).(*UUIDGenerate)
	model = updateModel(model, codeKeypress(tea.KeyDown)).(*UUIDGenerate)
	if model.generatedUUID != generated {
		t.Fatalf("expected UUID %v to stay, got %v", generated, model.generatedUUID)
	}
}

func TestUUIDGenerateShowsError(t *testing.T) {
	t.Parallel()

	model := newTestModel(24)
	model.version = 0
	model.generateUUID()
	model.form.State = huh.StateCompleted

	content := model.View().Content
	if !strings.Contains(content, "Failed to generate UUID: unsupported uuid version: 0") {
		t.Fatalf("expected error in view, got %q", content)
	}
	if strings.Contains(content, "Generated UUID") {
		t.Fatalf("expected no UUID row, got %q", content)
	}
}

func TestUUIDGenerateAbortedFormGeneratesNothing(t *testing.T) {
	t.Parallel()

	model := newTestModel(24)
	model.form.State = huh.StateAborted
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if model.generatedUUID != uuid.Nil {
		t.Fatalf("expected no UUID for an aborted form, got %v", model.generatedUUID)
	}
}

func TestUUIDGenerateDoesNotWriteToStdout(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = writer
	model := updateModel(newTestModel(24), codeKeypress(tea.KeyEnter)).(*UUIDGenerate)
	os.Stdout = stdout
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 0 {
		t.Fatalf("expected no stdout output, got %q", output)
	}
	if model.generatedUUID == uuid.Nil {
		t.Fatal("expected UUID to be generated")
	}
}

func TestUUIDGenerateKeys(t *testing.T) {
	t.Parallel()

	model := newTestModel(24)

	_, cmd := model.Update(codeKeypress(tea.KeyEscape))
	if msg, ok := cmd().(ui.ReturnToListMsg); !ok || msg.Common != model.common {
		t.Fatalf("expected ReturnToListMsg, got %#v", msg)
	}
	_, cmd = model.Update(tea.KeyPressMsg(tea.Key{Code: 'q', Text: "q"}))
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected quit command")
	}
}

func newTestModel(height int) *UUIDGenerate {
	common := &ui.CommonModel{Width: 80, Height: 24, Styles: ui.NewStyle()}
	model := NewUUIDGenerateModel(common)
	model = batchUpdate(model, model.Init()).(*UUIDGenerate)
	return updateModel(model, tea.WindowSizeMsg{Width: 80, Height: height}).(*UUIDGenerate)
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
