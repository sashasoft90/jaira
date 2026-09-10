// Package settings holds the few choices that belong to a person rather than
// to a board.
//
// It exists because two of them arrived with the ticket refs: which remote the
// refs go to, and whether an assignment pops up a desktop notification. Neither
// belongs in the repository — one teammate wanting to be notified says nothing
// about another, and a remote name is a property of a clone, not of the project.
//
// The file is ~/.jaira/settings.json, beside projects.json, and every field has
// a working default: a missing file, an unreadable one and an empty one all mean
// "the defaults", because settings a user never opened must never be the reason
// a board does not start.
package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/gitref"
)

// Settings is the whole file.
type Settings struct {
	// Remote is where ticket refs are pushed and fetched. Empty means origin.
	Remote string `json:"remote,omitempty"`

	// Notify turns the desktop notification for an assignment off. It is
	// phrased as the negative so that the zero value — a file that does not
	// mention it, or no file at all — leaves notifications on: the feature was
	// asked for, and a default that silently withholds it would answer the
	// request with nothing.
	NotifyOff bool `json:"notify-off,omitempty"`

	// Hook is a script called when a ticket moves or is claimed, for immediate
	// delivery through whatever the user already uses (a Slack webhook,
	// ntfy.sh, Telegram). jaira only calls it and brings no dependency of its
	// own; an empty value calls nothing.
	Hook string `json:"hook,omitempty"`
}

// Path is the settings file, honouring JAIRA_HOME so tests and a sandboxed run
// do not touch the real one.
func Path() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "settings.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "settings.json")
}

// Load returns the settings, or the defaults for anything missing.
//
// An unreadable or malformed file is deliberately not an error: it is one
// person's preferences, and refusing to open a board over it would trade a
// small annoyance for a total one.
func Load() Settings {
	var s Settings
	path := Path()
	if path == "" {
		return s
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return Settings{}
	}
	return s
}

// Save writes the settings, creating ~/.jaira if it is not there yet.
func Save(s Settings) error {
	path := Path()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// RemoteOr returns the configured remote, or the default.
func (s Settings) RemoteName() string {
	if r := strings.TrimSpace(s.Remote); r != "" {
		return r
	}
	return gitref.DefaultRemote
}

// NotifyEnabled reports whether an assignment should raise a desktop
// notification.
func (s Settings) NotifyEnabled() bool { return !s.NotifyOff }
