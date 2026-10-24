package editor

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenEditorWritesTempFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("EDITOR", "true")
	t.Setenv("SNAP_REVISION", "")

	cmd := OpenEditor("hello", "json")
	if msg, ok := cmd().(EditorFinishedMsg); ok {
		t.Fatalf("expected the editor to be started, got %+v", msg)
	}

	files, err := filepath.Glob(filepath.Join(dir, "editor-*.json"))
	if err != nil || len(files) != 1 {
		t.Fatalf("expected one temp file, got %v (%v)", files, err)
	}
	content, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read temp file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("expected temp file content %q, got %q", "hello", content)
	}
}

func TestOpenEditorCreateTempFailureReturnsError(t *testing.T) {
	t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "missing"))

	msg, ok := OpenEditor("hello", "json")().(EditorFinishedMsg)
	if !ok {
		t.Fatal("expected EditorFinishedMsg when the temp file cannot be created")
	}
	if msg.Err == nil {
		t.Fatal("expected error when the temp file cannot be created")
	}
}

func TestFinishEdit(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "edit.json")
	if err := os.WriteFile(path, []byte("edited"), 0o600); err != nil {
		t.Fatal(err)
	}

	msg := finishEdit(path, nil).(EditorFinishedMsg)
	if msg.Err != nil {
		t.Fatalf("unexpected error: %v", msg.Err)
	}
	if msg.Content != "edited" {
		t.Fatalf("expected content %q, got %q", "edited", msg.Content)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected temp file to be removed, got %v", err)
	}
}

func TestFinishEditKeepsEditorError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "edit.json")
	if err := os.WriteFile(path, []byte("edited"), 0o600); err != nil {
		t.Fatal(err)
	}
	editorErr := errors.New("exit status 1")

	msg := finishEdit(path, editorErr).(EditorFinishedMsg)
	if !errors.Is(msg.Err, editorErr) {
		t.Fatalf("expected editor error, got %v", msg.Err)
	}
	if msg.Content != "edited" {
		t.Fatalf("expected content %q, got %q", "edited", msg.Content)
	}
}

func TestFinishEditMissingFileReturnsError(t *testing.T) {
	t.Parallel()

	msg := finishEdit(filepath.Join(t.TempDir(), "missing.json"), nil).(EditorFinishedMsg)
	if !errors.Is(msg.Err, os.ErrNotExist) {
		t.Fatalf("expected not-exist error, got %v", msg.Err)
	}
	if msg.Content != "" {
		t.Fatalf("expected empty content, got %q", msg.Content)
	}
}

func TestOpenEditorWriteFailureReturnsError(t *testing.T) {
	dir := t.TempDir()
	orig := createTemp
	createTemp = func(dir, pattern string) (*os.File, error) {
		f, err := os.CreateTemp(dir, pattern)
		if err != nil {
			return nil, err
		}
		f.Close()
		return os.Open(f.Name())
	}
	t.Cleanup(func() { createTemp = orig })
	t.Setenv("TMPDIR", dir)

	msg, ok := OpenEditor("hello", "json")().(EditorFinishedMsg)
	if !ok || msg.Err == nil {
		t.Fatalf("expected EditorFinishedMsg with error, got %+v", msg)
	}
	if files, _ := filepath.Glob(filepath.Join(dir, "editor-*")); len(files) != 0 {
		t.Fatalf("expected temp file to be removed, got %v", files)
	}
}

func TestOpenEditorCommandFailureReturnsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("SNAP_REVISION", "1")

	msg, ok := OpenEditor("hello", "json")().(EditorFinishedMsg)
	if !ok || msg.Err == nil {
		t.Fatalf("expected EditorFinishedMsg with error, got %+v", msg)
	}
	if files, _ := filepath.Glob(filepath.Join(dir, "editor-*")); len(files) != 0 {
		t.Fatalf("expected temp file to be removed, got %v", files)
	}
}
