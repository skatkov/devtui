package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDGenerateCmdDefault(t *testing.T) {
	uuidgenerateVersion = 4
	uuidgenerateNamespace = ""

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"uuidgenerate"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("uuidgenerate command failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	parsed, err := uuid.Parse(output)
	if err != nil {
		t.Fatalf("uuidgenerate output is not a valid uuid: %v", err)
	}
	if parsed.Version() != 4 {
		t.Fatalf("expected version 4 uuid, got version %d", parsed.Version())
	}
}

func TestUUIDGenerateCmdVersion7(t *testing.T) {
	uuidgenerateVersion = 4
	uuidgenerateNamespace = ""

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"uuidgenerate", "--uuid-version", "7"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("uuidgenerate --uuid-version 7 failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	parsed, err := uuid.Parse(output)
	if err != nil {
		t.Fatalf("uuidgenerate output is not a valid uuid: %v", err)
	}
	if parsed.Version() != 7 {
		t.Fatalf("expected version 7 uuid, got version %d", parsed.Version())
	}
}
