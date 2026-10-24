package main

import (
	"github.com/skatkov/devtui/cmd"
	"github.com/skatkov/devtui/docs/generator"
)

func main() {
	generator.Generate(cmd.GetRootCmd())
}
