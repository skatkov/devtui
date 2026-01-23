package mcp

import (
	"bufio"
	"encoding/json"
	"io"
)

func ServeStdio(server *Server, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		var req Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}
		resp := server.HandleRequest(req)
		data, err := json.Marshal(resp)
		if err != nil {
			continue
		}
		if _, err := out.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return scanner.Err()
}
