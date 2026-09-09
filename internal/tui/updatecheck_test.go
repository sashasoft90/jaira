package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/selfupdate"
)

// setReleaseCurrent points release.Current at v for the duration of the
// test, mirroring internal/cli/update_test.go's setCurrent.
func setReleaseCurrent(t *testing.T, v string) {
	t.Helper()
	orig := release.Current
	release.Current = v
	t.Cleanup(func() { release.Current = orig })
}

func TestVersionLineNeverCheckedShowsVersionAlone(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "") // re-enable; TestMain disables by default

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") {
		t.Errorf("versionLine() = %q, want it to name the running version", line)
	}
	if strings.Contains(line, "up to date") || strings.Contains(line, "available") {
		t.Errorf("versionLine() = %q, must not claim a checked state when the cache has never been written", line)
	}
}

func TestVersionLineUpToDate(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.0.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") || !strings.Contains(line, "up to date") {
		t.Errorf("versionLine() = %q, want it to say up to date", line)
	}
}

func TestVersionLineUpdateAvailable(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "1.3.0") || !strings.Contains(line, "jaira self upgrade") {
		t.Errorf("versionLine() = %q, want it to name the available release and 'jaira self upgrade'", line)
	}
}

// TestVersionLineDisabledShowsVersionAlone asserts JAIRA_NO_UPDATE_CHECK=1
// suppresses the "available" half even when a stale cache would otherwise
// have something to say, and reports the version alone rather than nothing.
func TestVersionLineDisabledShowsVersionAlone(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "1")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	line := versionLine()
	if !strings.Contains(line, "jaira 1.0.0") {
		t.Errorf("versionLine() = %q, want it to name the running version", line)
	}
	if strings.Contains(line, "1.3.0") || strings.Contains(line, "available") {
		t.Errorf("versionLine() = %q, JAIRA_NO_UPDATE_CHECK=1 must suppress the available half", line)
	}
}

// TestHomeHeadCarriesTheVersionIndicator asserts the launcher's top left
// corner — not its footer, where the line lived until P1AE82 — carries the
// same indicator versionLine() produces.
func TestHomeHeadCarriesTheVersionIndicator(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.0.0"}); err != nil {
		t.Fatal(err)
	}

	h, err := NewHome(nil)
	if err != nil {
		t.Fatal(err)
	}
	h.width, h.height = 100, 30
	out := h.render()
	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[0], "up to date") {
		t.Errorf("home first line = %q, want the version indicator", lines[0])
	}
	if strings.Count(out, "up to date") != 1 {
		t.Errorf("home = %q, want the version indicator exactly once", out)
	}
}

// TestBoardHeadCarriesTheVersionIndicator asserts the board's top left corner
// carries the indicator and its status bar no longer does. Until P1AE82 the
// line sat in the status bar; it says which binary is running, which is asked
// before the board is read rather than after.
//
// The cache is written only after newTestStore has run, because that helper
// sets its own isolated JAIRA_HOME — writing it first would target a
// directory the test's actual Model never reads from.
func TestBoardHeadCarriesTheVersionIndicator(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	s := newTestStore(t)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.3.0"}); err != nil {
		t.Fatal(err)
	}

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	first := strings.Split(m.renderBoard(), "\n")[0]
	if !strings.Contains(first, "1.3.0") || !strings.Contains(first, "jaira self upgrade") {
		t.Errorf("board first line = %q, want the version indicator naming the available release", first)
	}
	if sb := m.statusBar(); strings.Contains(sb, "1.3.0") || strings.Contains(sb, "jaira 1.0.0") {
		t.Errorf("status bar = %q, want no version indicator", sb)
	}
}

// TestBoardVersionSitsInTheTopLeftCorner is the placement the ticket asks for,
// stated as geometry rather than as "somewhere in the output": the version is
// the board's very first row, flush against the left edge, above the head line
// that carries the ticket count.
//
// It also pins the room left for a second row underneath it: the follow-up
// puts an "^ <version>" pill there, and that only stays possible while the
// version owns a row of its own instead of sharing the head line.
func TestBoardVersionSitsInTheTopLeftCorner(t *testing.T) {
	setReleaseCurrent(t, "1.0.0")
	s := newTestStore(t)
	t.Setenv("JAIRA_NO_UPDATE_CHECK", "")
	if err := selfupdate.Write(selfupdate.Check{CheckedAt: time.Now().UTC(), Latest: "1.0.0"}); err != nil {
		t.Fatal(err)
	}

	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 150, 32
	lines := strings.Split(m.renderBoard(), "\n")
	if got := stripANSI(lines[0]); !strings.HasPrefix(got, "jaira 1.0.0") {
		t.Errorf("board first line = %q, want it to start flush left with the version", got)
	}
	// The head line, with the ticket count, is below it rather than sharing it.
	if got := stripANSI(lines[1]); !strings.Contains(got, "tickets") {
		t.Errorf("board second line = %q, want the head line with the ticket count", got)
	}
}
