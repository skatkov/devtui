package cmd

import (
	"bytes"
	"testing"
)

func TestMCPCommandListsTools(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(bytes.NewBufferString("{\"id\":1,\"method\":\"tools/list\"}\n"))
	cmd.SetArgs([]string{"mcp"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
}
