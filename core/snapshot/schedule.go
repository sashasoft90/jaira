package snapshot

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// stampName is the per-working-tree record of when this board was last
// snapshotted. Per tree rather than per machine, unlike the release check: how
// old a board's backup is says nothing about any other board's.
const stampName = "snapshot.json"

// Stamp is what gets recorded.
type Stamp struct {
	RanAt  time.Time `json:"ran_at"`
	Commit string    `json:"commit,omitempty"`
}

// StampPath is where the record lives, under the state directory the caller
// passes — never inside the repository, because when the backup last ran is
// this machine's business and committing it would be one more thing for two
// clones to disagree about.
func StampPath(stateDir string) string { return filepath.Join(stateDir, stampName) }

// ReadStamp returns the record, or a zero Stamp. Every error is swallowed: a
// missing or damaged record is indistinguishable from "never ran" for every
// purpose this serves, and the cost of getting it wrong is one extra snapshot.
func ReadStamp(stateDir string) Stamp {
	b, err := os.ReadFile(StampPath(stateDir))
	if err != nil {
		return Stamp{}
	}
	var s Stamp
	if err := json.Unmarshal(b, &s); err != nil {
		return Stamp{}
	}
	return s
}

// WriteStamp records that a run happened. Written before the run rather than
// after: a run that fails must not make the next command try again
// immediately, or an unreachable remote turns into a retry on every command.
func WriteStamp(stateDir string, s Stamp) error {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := StampPath(stateDir) + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, StampPath(stateDir))
}

// Due reports whether the last snapshot is older than every.
func Due(stateDir string, every time.Duration) bool {
	if Disabled() {
		return false
	}
	if every <= 0 {
		every = DefaultEvery
	}
	s := ReadStamp(stateDir)
	if s.RanAt.IsZero() {
		return true
	}
	return time.Since(s.RanAt) >= every
}

// Disabled reports whether snapshots are switched off for this process. It is
// also the recursion guard handed to a spawned child, so there is one
// mechanism here rather than two.
func Disabled() bool { return os.Getenv("JAIRA_NO_SNAPSHOT") == "1" }

// SpawnRun starts a detached child that takes the snapshot, and does not wait
// for it.
//
// Everything here is copied deliberately from the release check's own
// background refresh, including the reasons:
//
//   - The stamp is written by the caller before spawning, so a child that
//     dies or hangs cannot make the next command spawn another.
//   - The child gets JAIRA_NO_SNAPSHOT=1 so it can never decide a snapshot is
//     due and spawn a grandchild.
//   - No Wait: the parent is about to exit and the child is reparented to
//     init. A snapshot must never be something a command waits for — it
//     touches the network, and the whole point of the three-day cadence is
//     that nobody is waiting for it.
//   - Never from a test binary, which has no such subcommand and would
//     re-run its own suite, detached, reaching the network on the way.
//   - Stdio to os.DevNull, because a child writing to the terminal after the
//     parent has exited is the failure this design exists to avoid.
func SpawnRun() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if isTestBinary(exe) {
		return nil
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer devNull.Close()

	cmd := exec.Command(exe, "snapshot", "--json")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = devNull, devNull, devNull
	cmd.Env = append(os.Environ(), "JAIRA_NO_SNAPSHOT=1")
	return cmd.Start()
}

func isTestBinary(exe string) bool {
	base := strings.TrimSuffix(filepath.Base(exe), ".exe")
	return strings.HasSuffix(base, ".test")
}
