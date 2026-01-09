package input

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func ReadFromArgsOrStdin(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	data, err := readBytes(cmd.InOrStdin())
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadBytesFromArgsOrStdin(cmd *cobra.Command, args []string) ([]byte, error) {
	if len(args) > 0 {
		return []byte(args[0]), nil
	}

	return readBytes(cmd.InOrStdin())
}

func readBytes(r io.Reader) ([]byte, error) {
	if f, ok := r.(*os.File); ok {
		if isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd()) {
			return nil, nil
		}
	}

	return io.ReadAll(r)
}
