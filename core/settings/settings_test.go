package settings_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BeMuCa/jaira/core/settings"
)

// Settings nobody has ever opened must not change how anything behaves: the
// remote is origin and the notification the ticket asked for is on.
func TestDefaultsWithNoFile(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())

	s := settings.Load()
	if got := s.RemoteName(); got != "origin" {
		t.Errorf("remote defaults to %q", got)
	}
	if !s.NotifyEnabled() {
		t.Error("notifications are off by default; the feature was asked for")
	}
	if s.Hook != "" {
		t.Errorf("a hook appeared out of nowhere: %q", s.Hook)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)

	want := settings.Settings{Remote: "board", NotifyOff: true, Hook: "/usr/local/bin/tell-me"}
	if err := settings.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, "settings.json")); err != nil {
		t.Fatalf("settings file: %v", err)
	}
	got := settings.Load()
	if got != want {
		t.Errorf("round trip changed the settings: %+v", got)
	}
	if got.RemoteName() != "board" {
		t.Errorf("remote is %q", got.RemoteName())
	}
	if got.NotifyEnabled() {
		t.Error("notify-off was not honoured")
	}
}

// A broken settings file is one person's preferences, not a reason to refuse to
// open a board.
func TestAMalformedFileFallsBackToDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("JAIRA_HOME", home)
	if err := os.WriteFile(filepath.Join(home, "settings.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	s := settings.Load()
	if s.RemoteName() != "origin" || !s.NotifyEnabled() {
		t.Errorf("a malformed file changed behaviour: %+v", s)
	}
}

func TestBlankRemoteIsStillOrigin(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	if err := settings.Save(settings.Settings{Remote: "   "}); err != nil {
		t.Fatal(err)
	}
	if got := settings.Load().RemoteName(); got != "origin" {
		t.Errorf("whitespace remote resolved to %q", got)
	}
}
