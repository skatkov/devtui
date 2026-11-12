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

// createTemp creates the file handed to the editor; tests replace it.
var createTemp = os.CreateTemp

func OpenEditor(content string, format string) tea.Cmd {
	path, err := writeTempFile(content, format)
	if err != nil {
		return func() tea.Msg { return EditorFinishedMsg{Err: err} }
	}

	cmd, err := editor.Cmd("", path)
	if err != nil {
		return func() tea.Msg { return EditorFinishedMsg{Err: errors.Join(err, os.Remove(path))} }
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return finishEdit(path, err)
	})
}

// writeTempFile stores content in a new temporary file and returns its path.
func writeTempFile(content, format string) (string, error) {
	tmpfile, err := createTemp("", "editor-*."+format)
	if err != nil {
		return "", err
	}

	_, err = tmpfile.WriteString(content)
	if err = errors.Join(err, tmpfile.Close()); err != nil {
		return "", errors.Join(err, os.Remove(tmpfile.Name()))
	}
	return tmpfile.Name(), nil
}

// finishEdit reads the edited file at path, removes it and reports the result.
func finishEdit(path string, err error) tea.Msg {
	content, readErr := os.ReadFile(path)
	removeErr := os.Remove(path)
	if readErr != nil {
		return EditorFinishedMsg{Err: errors.Join(readErr, removeErr)}
	}
	return EditorFinishedMsg{Err: errors.Join(err, removeErr), Content: string(content)}
}
