package root

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/ui"
)

func TestRootScreenWithSize(t *testing.T) {
	useTempConfig(t)

	m := RootScreenWithRenderer(nil, 80, 24)
	if m.common.Width != 80 || m.common.Height != 24 || m.common.Styles == nil {
		t.Fatalf("expected styled 80x24 common model, got %+v", m.common)
	}
	if m.currentView != m.listModel || m.Init() != nil {
		t.Fatal("expected root to start on the menu")
	}
}

func TestRootWindowSizeMsgUpdatesCommon(t *testing.T) {
	useTempConfig(t)

	m := RootScreen()
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if m.common.Width != 100 || m.common.Height != 40 || m.listModel.list.Width() != 100 {
		t.Fatalf("expected 100x40, got %dx%d", m.common.Width, m.common.Height)
	}
}

func TestRootReturnToListRestoresSelection(t *testing.T) {
	useTempConfig(t)

	m := RootScreen()
	common := &ui.CommonModel{Width: 90, Height: 30, LastSelectedItem: 2}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	if next.(RootModel).currentView == m.listModel {
		t.Fatal("expected a tool to open")
	}

	next, _ = next.Update(ui.ReturnToListMsg{Common: common})
	got := next.(RootModel)
	if got.currentView != got.listModel || got.listModel.common != common {
		t.Fatal("expected the menu with the returned common model")
	}
	if got.listModel.list.Index() != 2 || got.listModel.list.Width() != 90 {
		t.Fatalf("expected selection 2 at width 90, got %d at %d", got.listModel.list.Index(), got.listModel.list.Width())
	}
	if !got.View().AltScreen {
		t.Fatal("expected alt screen view")
	}
}
