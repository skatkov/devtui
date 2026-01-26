package mcp

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func BuildTools(root *cobra.Command) []ToolSchema {
	var tools []ToolSchema
	var walk func(cmd *cobra.Command, path []string)
	walk = func(cmd *cobra.Command, path []string) {
		if cmd.IsAvailableCommand() && !cmd.HasSubCommands() {
			name := "devtui." + strings.Join(path, ".")
			tools = append(tools, ToolSchema{
				Name:        name,
				Description: cmd.Short,
				InputSchema: buildSchema(cmd),
			})
		}

		for _, child := range cmd.Commands() {
			if child.IsAvailableCommand() {
				walk(child, append(path, child.Name()))
			}
		}
	}

	for _, child := range root.Commands() {
		if child.IsAvailableCommand() {
			walk(child, []string{child.Name()})
		}
	}

	return tools
}

func buildSchema(cmd *cobra.Command) JSONSchema {
	schema := JSONSchema{
		Type:       "object",
		Properties: map[string]JSONSchema{},
	}
	schema.Properties["input"] = JSONSchema{Type: "string"}
	schema.Properties["args"] = JSONSchema{Type: "array"}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		schema.Properties[flag.Name] = JSONSchema{
			Type:        flagType(flag.Value.Type()),
			Description: flag.Usage,
			Default:     flag.DefValue,
		}
	})

	return schema
}

func flagType(valueType string) string {
	switch valueType {
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "stringSlice", "stringArray", "intSlice", "intArray":
		return "array"
	default:
		return "string"
	}
}
