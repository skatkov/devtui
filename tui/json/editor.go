package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/editor"
)

type editorFinishedMsg struct {
	err     error
	content string
}

func openEditor(content string, format string) tea.Cmd {
	cb := func(err error, newContent string) tea.Msg {
		return editorFinishedMsg{err: err, content: newContent}
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "editor-*."+format)
	if err != nil {
		return func() tea.Msg { return cb(err, "") }
	}

	// Write the initial content
	if _, err := tmpfile.WriteString(content); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return func() tea.Msg { return cb(err, "") }
	}
	tmpfile.Close()

	// Open editor with the temp file
	cmd, err := editor.Cmd("", tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		return func() tea.Msg { return cb(err, "") }
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		// Read the modified content
		content, readErr := os.ReadFile(tmpfile.Name())
		os.Remove(tmpfile.Name()) // Clean up
		if readErr != nil {
			return cb(readErr, "")
		}
		return cb(err, string(content))
	})
}
