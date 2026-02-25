package teacompat

import (
	tea "charm.land/bubbletea/v2"
	teav1 "github.com/charmbracelet/bubbletea"
)

func Cmd(cmd teav1.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return func() tea.Msg {
		return cmd()
	}
}
