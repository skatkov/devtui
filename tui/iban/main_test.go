package iban

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/jacoelho/banking/iban"
	"github.com/skatkov/devtui/internal/ui"
)

func TestIBANCountrySelectView(t *testing.T) {
	t.Parallel()

	model := newModel()
	model = updateModel(model, tea.WindowSizeMsg{Width: 80, Height: 40}).(*IBANGenerate)
	view := model.View().Content

	if lipgloss.Height(view) != 40 {
		t.Errorf("expected view to fill 40 lines, got %d", lipgloss.Height(view))
	}
	if !strings.Contains(view, "BH - Bahrain") || strings.Contains(view, "BI - Burundi") {
		t.Errorf("expected nine visible countries, got %q", view)
	}
	if strings.Count(view, "submit") != 1 {
		t.Errorf("expected help to be shown once, got %q", view)
	}
	if !strings.Contains(view, ui.AppTitle+" :: \x1b[1m"+Title) {
		t.Errorf("expected app name and bold title in header, got %q", view)
	}

	lines := strings.Split(view, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Select Country") {
			if strings.TrimSpace(lines[i-1]) != "" || !strings.HasPrefix(line, "  \x1b") {
				t.Errorf("expected form below a blank line without extra indent, got %q and %q", lines[i-1], line)
			}
		}
	}
}

func TestIBANCountrySelectDoesNotGenerateBeforeSubmit(t *testing.T) {
	t.Parallel()

	model := updateModel(newModel(), codeKeypress(tea.KeyDown)).(*IBANGenerate)

	if model.countryCode != "AE" {
		t.Fatalf("expected AE to be highlighted, got %q", model.countryCode)
	}
	if model.generatedIBAN != "" || model.error != "" {
		t.Errorf("expected nothing generated before submit, got %q and error %q", model.generatedIBAN, model.error)
	}
}

func TestIBANAbortedFormDoesNotGenerate(t *testing.T) {
	t.Parallel()

	model := newModel()
	model.form.State = huh.StateAborted
	model = updateModel(model, tea.WindowSizeMsg{Width: 80, Height: 24}).(*IBANGenerate)

	if model.generatedIBAN != "" || model.error != "" {
		t.Errorf("expected nothing generated for an aborted form, got %q and error %q", model.generatedIBAN, model.error)
	}
}

func TestIBANGenerateForSelectedCountry(t *testing.T) {
	t.Parallel()

	model := updateModel(newModel(), codeKeypress(tea.KeyEnter)).(*IBANGenerate)

	if model.form.State != huh.StateCompleted {
		t.Fatalf("expected form to complete, got state %v", model.form.State)
	}
	if !strings.HasPrefix(model.generatedIBAN, "AD") {
		t.Errorf("expected AD IBAN, got %q", model.generatedIBAN)
	}
	if model.formattedIBAN != iban.PaperFormat(model.generatedIBAN) {
		t.Errorf("expected paper format of %q, got %q", model.generatedIBAN, model.formattedIBAN)
	}

	view := model.View().Content
	for _, want := range []string{"Andorra (AD)", model.generatedIBAN, model.formattedIBAN, ":: \x1b[1m" + Title} {
		if !strings.Contains(view, want) {
			t.Errorf("expected view to contain %q, got %q", want, view)
		}
	}

	lines := strings.Split(view, "\n")
	for i, line := range lines {
		if strings.Contains(line, "╭") {
			if strings.TrimSpace(lines[i-1]) != "" || !strings.HasPrefix(line, "  ╭") {
				t.Errorf("expected table below a blank line without extra indent, got %q and %q", lines[i-1], line)
			}
			if width := utf8.RuneCountInString(strings.TrimSpace(line)); width != 100 {
				t.Errorf("expected table width 100, got %d", width)
			}
		}
	}
}

func TestIBANKeepsIBANAfterCompletion(t *testing.T) {
	t.Parallel()

	model := updateModel(newModel(), codeKeypress(tea.KeyEnter)).(*IBANGenerate)
	generated := model.generatedIBAN

	model = updateModel(model, tea.WindowSizeMsg{Width: 100, Height: 30}).(*IBANGenerate)
	model = updateModel(model, codeKeypress(tea.KeyDown)).(*IBANGenerate)
	if model.generatedIBAN != generated {
		t.Fatalf("expected IBAN %q to stay, got %q", generated, model.generatedIBAN)
	}
}

func TestIBANEveryCountryGenerates(t *testing.T) {
	t.Parallel()

	for _, option := range getCountryOptions() {
		if _, err := iban.Generate(option.Value); err != nil {
			t.Errorf("%s: %v", option.Key, err)
		}
	}
}

func TestIBANViewShowsCountryName(t *testing.T) {
	t.Parallel()

	// GF comes after XK in the options, so a loose comparison would pick another name
	model := newModel()
	model.countryCode = "GF"
	model.generateIBAN()
	model.form.State = huh.StateCompleted

	if view := model.View().Content; !strings.Contains(view, "French Guyana (GF)") {
		t.Errorf("expected country name in view, got %q", view)
	}
}

func TestIBANGenerateErrorIsShown(t *testing.T) {
	t.Parallel()

	model := newModel()
	model.countryCode = "ZZ"
	model.generateIBAN()
	model.form.State = huh.StateCompleted

	if model.error != "Failed to generate IBAN: unsupported country code ZZ" {
		t.Errorf("expected unsupported country error, got %q", model.error)
	}

	view := model.View().Content
	if !strings.Contains(view, model.error) || strings.Contains(view, "IBAN (Formatted)") {
		t.Errorf("expected error row instead of IBAN rows, got %q", view)
	}
}

func TestIBANGenerateRequiresCountryCode(t *testing.T) {
	t.Parallel()

	model := newModel()
	model.countryCode = ""
	model.generateIBAN()

	if model.error != "Country code is required" {
		t.Errorf("expected missing country error, got %q", model.error)
	}
}

func newModel() *IBANGenerate {
	common := &ui.CommonModel{Width: 80, Height: 30}
	common.Styles = ui.NewStyle()

	model := NewIBANGenerateModel(common)
	return batchUpdate(model, model.Init()).(*IBANGenerate)
}

func updateModel(model tea.Model, msg tea.Msg) tea.Model {
	nextModel, cmd := model.Update(msg)
	return batchUpdate(nextModel, cmd)
}

func batchUpdate(model tea.Model, cmd tea.Cmd) tea.Model {
	return runCmd(model, cmd, 0)
}

// runCmd feeds the messages of a few command levels back into the model and
// drops commands that wait, such as the cursor blink tick.
func runCmd(model tea.Model, cmd tea.Cmd, depth int) tea.Model {
	if cmd == nil || depth > 3 {
		return model
	}

	done := make(chan tea.Msg, 1)
	go func() { done <- cmd() }()

	var msg tea.Msg
	select {
	case msg = <-done:
	case <-time.After(200 * time.Millisecond):
		return model
	}

	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			model = runCmd(model, c, depth+1)
		}
		return model
	}

	nextModel, nextCmd := model.Update(msg)
	return runCmd(nextModel, nextCmd, depth+1)
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
