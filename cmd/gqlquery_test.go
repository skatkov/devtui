package cmd

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	graphqlquery "github.com/skatkov/devtui/tui/graphql-query"
)

const gqlqueryInput = "# note\nquery { a { b } }"

func TestGqlqueryCmd(t *testing.T) {
	captureProgram(t, nil)

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default two spaces without comments", nil, "query {\n  a {\n    b\n  }\n}\n\n"},
		{"custom indent", []string{"--indent", "    "}, "query {\n    a {\n        b\n    }\n}\n\n"},
		{"empty indent uses tabs", []string{"--indent", ""}, "query {\n\ta {\n\t\tb\n\t}\n}\n\n"},
		{"comments", []string{"--with-comments"}, "# note\nquery {\n  a {\n    b\n  }\n}\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runDevtui(strings.NewReader(gqlqueryInput), append([]string{"gqlquery"}, tt.args...)...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := stdout.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGqlqueryCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "gqlquery", "{ a }")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "query {\n  a\n}\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestGqlqueryCmdInvalidQuery(t *testing.T) {
	_, _, err := runDevtui(strings.NewReader("query {"), "gqlquery")
	if err == nil {
		t.Fatal("expected error for an invalid query")
	}
}

func TestGqlqueryCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	_, _, err := runDevtui(iotest.ErrReader(want), "gqlquery")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestGqlqueryCmdWithDescriptionsDefault(t *testing.T) {
	if got := gqlfmtCmd.Flags().Lookup("with-descriptions").DefValue; got != "false" {
		t.Fatalf("expected default false, got %q", got)
	}
}

func TestGqlqueryCmdTUI(t *testing.T) {
	model := captureProgram(t, nil)

	_, _, err := runDevtui(strings.NewReader("{ a }"), "gqlquery", "--tui")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := (*model).(graphqlquery.GraphQLQueryModel)
	if !ok {
		t.Fatalf("expected graphqlquery.GraphQLQueryModel, got %T", *model)
	}
	if m.Content != "{ a }" {
		t.Fatalf("expected content %q, got %q", "{ a }", m.Content)
	}
}

func TestGqlqueryCmdTUIInvalidQuery(t *testing.T) {
	model := captureProgram(t, nil)

	_, _, err := runDevtui(strings.NewReader("query {"), "gqlquery", "--tui")
	if err == nil {
		t.Fatal("expected error for an invalid query")
	}
	if *model != nil {
		t.Fatalf("expected no program, got %T", *model)
	}
}

func TestGqlqueryCmdTUIError(t *testing.T) {
	want := errors.New("program failed")
	captureProgram(t, want)

	_, _, err := runDevtui(strings.NewReader("{ a }"), "gqlquery", "--tui")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
