package uuidgenerate

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
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
