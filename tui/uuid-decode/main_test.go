package uuiddecode

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/uuidutil"
)

const sampleUUID = "4326ff5f-774d-4506-a18c-4bc50c761863"

func TestExtractUUIDData(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		uuid     uuid.UUID
		expected [][]string
	}{
		{
			name: "Valid UUID v4",
			uuid: uuid.MustParse("4326ff5f-774d-4506-a18c-4bc50c761863"),
			expected: [][]string{
				{"Standard String Format", "4326ff5f-774d-4506-a18c-4bc50c761863"},
				{"Single Integer Value", "89260762576260186387968344354060638307"},
				{"Version", "4"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents", "43:26:FF:5F:77:4D:45:06:A1:8C:4B:C5:0C:76:18:63"},
			},
		},
		{
			name: "Valid UUID v1",
			uuid: uuid.MustParse("550e8400-e29b-11d4-a716-446655440000"),
			expected: [][]string{
				{"Standard String Format", "550e8400-e29b-11d4-a716-446655440000"},
				{"Single Integer Value", "113059749145936098728763079434011148288"},
				{"Version", "1"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents - Time", "2001-01-04 23:43:07.540992 UTC"},
				{"Contents - Clock", "10006"},
				{"Contents - Node", "44:66:55:44:00:00"},
			},
		},
		{
			name: "UUID v2",
			uuid: uuid.MustParse("000003e8-e64c-21ef-ac00-325096b39f47"),
			expected: [][]string{
				{"Standard String Format", "000003e8-e64c-21ef-ac00-325096b39f47"},
				{"Single Integer Value", "79299436105144797400979033268039"},
				{"Version", "2"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents", "00:00:03:E8:E6:4C:21:EF:AC:00:32:50:96:B3:9F:47"},
			},
		},
		{
			name: "UUID v3",
			uuid: uuid.MustParse("bd278ea4-0b40-366b-b7b1-f067a7817f35"),
			expected: [][]string{
				{"Standard String Format", "bd278ea4-0b40-366b-b7b1-f067a7817f35"},
				{"Single Integer Value", "251429483879896292488202237579433705269"},
				{"Version", "3"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents", "BD:27:8E:A4:0B:40:36:6B:B7:B1:F0:67:A7:81:7F:35"},
			},
		},
		{
			name:     "Nil UUID",
			uuid:     uuid.Nil,
			expected: [][]string{},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uuidutil.FieldsToRows(uuidutil.Decode(tt.uuid))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("\nTest: %s\nGot:\n", tt.name)
				for _, row := range got {
					t.Errorf("  %v\n", row)
				}
				t.Errorf("Want:\n")
				for _, row := range tt.expected {
					t.Errorf("  %v\n", row)
				}
			}
		})
	}
}

func TestUUIDDecodeViewFillsWindow(t *testing.T) {
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

func TestUUIDDecodeShowsDecodedFields(t *testing.T) {
	t.Parallel()

	model := submitUUID(newTestModel(24), sampleUUID)
	if model.form.State != huh.StateCompleted {
		t.Fatalf("expected form to complete, got state %v", model.form.State)
	}

	view := model.View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	content := view.Content
	for _, want := range []string{"Standard String Format", sampleUUID, "Version", "4"} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected view to contain %q, got %q", want, content)
		}
	}
	if width := lipgloss.Width(strings.TrimSpace(strings.Split(content, "\n")[1])); width != 100 {
		t.Fatalf("expected table width 100, got %d", width)
	}
}

func TestUUIDDecodeRejectsInvalidUUID(t *testing.T) {
	t.Parallel()

	model := submitUUID(newTestModel(24), "not-a-uuid")
	if model.form.State != huh.StateNormal {
		t.Fatalf("expected form to stay open, got state %v", model.form.State)
	}
}

func TestUUIDDecodeKeys(t *testing.T) {
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

func newTestModel(height int) *UUIDDecode {
	common := &ui.CommonModel{Width: 80, Height: 24, Styles: ui.NewStyle()}
	model := NewUUIDDecodeModel(common)
	model.Update(model.Init()())
	model.Update(tea.WindowSizeMsg{Width: 80, Height: height})
	return model
}

func submitUUID(model *UUIDDecode, value string) *UUIDDecode {
	for _, r := range value {
		model.Update(tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)}))
	}
	return updateModel(model, codeKeypress(tea.KeyEnter)).(*UUIDDecode)
}

func updateModel(model tea.Model, msg tea.Msg) tea.Model {
	nextModel, cmd := model.Update(msg)
	return batchUpdate(nextModel, cmd)
}

func batchUpdate(model tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return model
	}

	msg := cmd()
	nextModel, nextCmd := model.Update(msg)
	if nextCmd == nil {
		return nextModel
	}

	msg = nextCmd()
	nextModel, _ = nextModel.Update(msg)
	return nextModel
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
