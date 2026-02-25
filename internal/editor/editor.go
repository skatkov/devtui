package editor

import (
	"errors"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/editor"
)

type EditorFinishedMsg struct {
	Err     error
	Content string
}

func OpenEditor(content string, format string) tea.Cmd {
	cb := func(err error, newContent string) tea.Msg {
		return EditorFinishedMsg{Err: err, Content: newContent}
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "editor-*."+format)
	if err != nil {
		return func() tea.Msg { return cb(err, "") }
	}

	// Write the initial content
	if _, err := tmpfile.WriteString(content); err != nil {
		if err := tmpfile.Close(); err != nil {
			// We still need to try to remove the file even if closing failed
			if removeErr := os.Remove(tmpfile.Name()); removeErr != nil {
				// If both close and remove fail, combine the errors
				return func() tea.Msg { return cb(errors.Join(err, removeErr), "") }
			}
			return func() tea.Msg { return cb(err, "") }
		}
	}

	// Open editor with the temp file
	cmd, err := editor.Cmd("", tmpfile.Name())
	if err != nil {
		if removeErr := os.Remove(tmpfile.Name()); removeErr != nil {
			return func() tea.Msg { return cb(errors.Join(err, removeErr), "") }
		}
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		// Read the modified content
		content, readErr := os.ReadFile(tmpfile.Name())
		removeErr := os.Remove(tmpfile.Name())
		if readErr != nil {
			if removeErr != nil {
				return cb(errors.Join(readErr, removeErr), "")
			}
			return cb(readErr, "")
		}
		if removeErr != nil {
			return cb(errors.Join(err, removeErr), string(content))
		}
		return cb(err, string(content))
	})
}
