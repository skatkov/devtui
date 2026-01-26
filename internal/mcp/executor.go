package mcp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func ExecuteTool(root *cobra.Command, params CallParams) (output string, err error) {
	cmdPath := strings.TrimPrefix(params.Name, "devtui.")
	args := []string{}
	if cmdPath != "" {
		args = append(args, strings.Split(cmdPath, ".")...)
	}
	for name, value := range params.Flags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}
	args = append(args, params.Args...)

	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := stdoutReader.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		if closeErr := stdoutWriter.Close(); closeErr != nil {
			return "", closeErr
		}
		return "", err
	}
	defer func() {
		if closeErr := stderrReader.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	stdoutDone := make(chan error, 1)
	stderrDone := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(stdoutBuf, stdoutReader)
		stdoutDone <- copyErr
	}()
	go func() {
		_, copyErr := io.Copy(stderrBuf, stderrReader)
		stderrDone <- copyErr
	}()

	origStdout := os.Stdout
	origStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	root.SetOut(stdoutWriter)
	root.SetErr(stderrWriter)
	root.SetIn(strings.NewReader(params.Input))
	root.SetArgs(args)

	execErr := root.Execute()

	if closeErr := stdoutWriter.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if closeErr := stderrWriter.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	os.Stdout = origStdout
	os.Stderr = origStderr

	if err := <-stdoutDone; err != nil {
		return "", err
	}
	if err := <-stderrDone; err != nil {
		return "", err
	}

	if execErr != nil {
		if stderrBuf.Len() > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(stderrBuf.String()))
		}
		return "", execErr
	}

	return stdoutBuf.String(), nil
}
