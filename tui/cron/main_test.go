package cron

import (
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/skatkov/devtui/internal/ui"
)

func TestValidateCronExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"every five minutes", "*/5 * * * *", ""},
		{"fixed time", "30 4 1 * *", ""},
		{"not a cron expression", "hello", "invalid cron expression format"},
		{"out of range minute", "99 * * * *", "invalid cron expression: failed to parse CRON expression"},
		{"sunday as zero", "0 0 * * 0", ""},
		{"range with step", "0-30/5 * * * *", ""},
		{"trailing garbage", "*/5 * * * *.", "invalid cron expression format"},
		{"unsupported macro", "@daily", "invalid cron expression format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateCronExpression(tt.input)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateCronExpression(%q) = %v, want nil", tt.input, err)
				}
				return
			}
			if err == nil || !strings.HasPrefix(err.Error(), tt.wantErr) {
				t.Fatalf("validateCronExpression(%q) = %v, want prefix %q", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestDescribeCron(t *testing.T) {
	t.Parallel()

	desc, err := describeCron("*/5 * * * *")
	if err != nil {
		t.Fatalf("describeCron() error = %v", err)
	}
	if desc != "Every 5 minutes" {
		t.Errorf("expected %q, got %q", "Every 5 minutes", desc)
	}

	if _, err := describeCron("99 * * * *"); err == nil {
		t.Error("expected error for out of range minute")
	}
}

func TestDescriptionUsesTwentyFourHourTime(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m.cronExpression = "30 14 * * *"
	want := strings.Repeat(" ", 20) + "At 14:30"
	if got := stripANSI(m.description()); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDescriptionEmptyForInvalidExpression(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m.cronExpression = "99 * * * *"
	if got := m.description(); got != "" {
		t.Errorf("expected empty description, got %q", got)
	}
}

func TestDescriptionNamesWeekdaysFromSunday(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"0 0 * * 0", "At 00:00, only on Sunday"},
		{"0 0 * * 1-5", "At 00:00, Monday through Friday"},
		{"0 0 * * 6", "At 00:00, only on Saturday"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			m := newTestModel(24)
			m.cronExpression = tt.input
			if got := stripANSI(m.description()); !strings.Contains(got, tt.want) {
				t.Errorf("description() for %q = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestViewFillsWindowHeight(t *testing.T) {
	t.Parallel()

	for _, height := range []int{24, 40} {
		m := newTestModel(height)
		view := m.View()
		if got := lipgloss.Height(view.Content); got != height {
			t.Errorf("expected view height %d, got %d", height, got)
		}
		if !view.AltScreen {
			t.Error("expected alt screen view")
		}
		if !strings.Contains(view.Content, Title) {
			t.Errorf("expected view to contain title %q", Title)
		}
	}
}

func TestViewLayout(t *testing.T) {
	t.Parallel()

	m := newTestModel(14)
	m = batchUpdate(m, m.Init()).(*CronModel)
	content := stripANSI(m.View().Content)
	lines := strings.Split(content, "\n")
	if len(lines) != 14 {
		t.Fatalf("expected 14 lines, got %d: %q", len(lines), content)
	}
	if !strings.HasPrefix(lines[3], "  ┃ Cron Expression") {
		t.Errorf("expected form below header, got %q", lines[3])
	}
	if got := strings.Count(content, "enter submit"); got != 1 {
		t.Errorf("expected help once, got %d times", got)
	}
	if !strings.HasPrefix(lines[12], "  enter submit") {
		t.Errorf("expected help at the bottom, got %q", lines[12])
	}
}

func TestViewUsesBoldTitles(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	if got := m.View().Content; !strings.Contains(got, "\x1b[1m"+Title) {
		t.Errorf("expected bold module title, got %q", got)
	}

	m.form.State = huh.StateCompleted
	if got := m.View().Content; !strings.Contains(got, "\x1b[1;38;2;255;105;180m*/5 * * * *") {
		t.Errorf("expected bold expression, got %q", got)
	}
}

func TestViewShowsCompletedDescription(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m.cronExpression = "30 14 * * *"
	m.form.State = huh.StateCompleted
	got := stripANSI(m.View().Content)
	if !strings.Contains(got, "30 14 * * *") || !strings.Contains(got, "At 14:30") {
		t.Errorf("expected expression and description, got %q", got)
	}
}

func TestViewShowsCompletedParseError(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m.cronExpression = "99 * * * *"
	m.form.State = huh.StateCompleted
	got := stripANSI(m.View().Content)
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected padded error on 5 lines, got %q", got)
	}
	if !strings.HasPrefix(lines[2], "  Error parsing cron expression: ") {
		t.Errorf("expected padded parse error, got %q", lines[2])
	}
}

func TestEscReturnsToList(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	_, cmd := m.Update(codeKeypress(tea.KeyEscape))
	msg, ok := cmd().(ui.ReturnToListMsg)
	if !ok {
		t.Fatalf("expected ReturnToListMsg, got %T", cmd())
	}
	if msg.Common != m.common {
		t.Error("expected common model to be passed back")
	}
}

func TestQuitKeys(t *testing.T) {
	t.Parallel()

	for _, key := range []tea.KeyPressMsg{
		tea.KeyPressMsg(tea.Key{Code: 'q', Text: "q"}),
		tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}),
	} {
		m := newTestModel(24)
		_, cmd := m.Update(key)
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Errorf("expected quit for %q, got %T", key.String(), cmd())
		}
	}
}

func TestWindowSizeMsgUpdatesCommon(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	if m.common.Width != 100 || m.common.Height != 30 {
		t.Errorf("expected size 100x30, got %dx%d", m.common.Width, m.common.Height)
	}
}

func TestTypingUpdatesExpression(t *testing.T) {
	t.Parallel()

	m := newTestModel(24)
	m = batchUpdate(m, m.Init()).(*CronModel)
	m = updateModel(m, codeKeypress(tea.KeyBackspace)).(*CronModel)
	if m.cronExpression != "*/5 * * * " {
		t.Errorf("expected %q, got %q", "*/5 * * * ", m.cronExpression)
	}
}

var ansiPattern = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func newTestModel(height int) *CronModel {
	common := &ui.CommonModel{Width: 80, Height: height}
	common.Styles = ui.NewStyle()
	return NewCronModel(common)
}

func updateModel(model tea.Model, msg tea.Msg) tea.Model {
	nextModel, cmd := model.Update(msg)
	return batchUpdate(nextModel, cmd)
}

func batchUpdate(model tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return model
	}

	msg := cmd()
	nextModel, nextCmd := model.Update(msg)
	if nextCmd == nil {
		return nextModel
	}

	msg = nextCmd()
	nextModel, _ = nextModel.Update(msg)
	return nextModel
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
