package refsync_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Each side is a real clone with a real store, wired the way the CLI wires it:
// the store's recorder is the syncer, so an ordinary ticket write is what
// queues the ref write. Testing through Mutate rather than through Record
// directly is the point — the claim being made is that the existing write path
// now carries tickets to the remote, not that a new function can.
type side struct {
	store  *ticket.Store
	syncer *refsync.Syncer
	dir    string
}

func twoSides(t *testing.T) (ada, berk side, remote string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))

	bare := filepath.Join(root, "board.git")
	run(t, root, "git", "init", "--bare", "--quiet", bare)

	mk := func(name string) side {
		dir := filepath.Join(root, name)
		run(t, root, "git", "clone", "--quiet", bare, dir)
		run(t, dir, "git", "config", "user.name", name)
		run(t, dir, "git", "config", "user.email", name+"@example.test")
		s, err := ticket.At(dir)
		if err != nil {
			t.Fatalf("store at %s: %v", dir, err)
		}
		if _, err := s.Init(); err != nil {
			t.Fatalf("init %s: %v", dir, err)
		}
		s.Actor = name
		y := refsync.New(s, "origin", name)
		s.Recorder = y
		return side{store: s, syncer: y, dir: dir}
	}
	return mk("ada"), mk("berk"), bare
}

func create(t *testing.T, s side, title string) string {
	t.Helper()
	id := ticket.NewID(time.Now())
	_, err := s.store.Create(map[string]string{
		ticket.FieldID:        id,
		ticket.FieldTitle:     title,
		ticket.FieldStatus:    "backlog",
		ticket.FieldAssignee:  s.store.Actor,
		ticket.FieldCreator:   s.store.Actor,
		ticket.FieldUpdatedBy: s.store.Actor,
	}, nil, "# "+title+"\n")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	return id
}

// An ordinary ticket write now also travels to the remote, and arrives at
// someone who has none of the writer's branches.
func TestAWriteOnTheStoreReachesTheOtherClone(t *testing.T) {
	ada, berk, _ := twoSides(t)

	id := create(t, ada, "session cookie dropped on 302")
	if _, ok := ada.syncer.Pending(id); !ok {
		t.Fatal("creating a ticket queued nothing for the remote")
	}
	reports, err := ada.syncer.Flush()
	if err != nil {
		t.Fatalf("flush: %v", err)
	}
	if len(reports) != 1 || reports[0].Outcome != outbox.Sent {
		t.Fatalf("want one sent, got %+v", reports)
	}
	if _, ok := ada.syncer.Pending(id); ok {
		t.Error("a sent write is still queued")
	}

	if err := berk.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("berk fetch: %v", err)
	}
	got, _, err := berk.syncer.Repo.Read(id)
	if err != nil {
		t.Fatalf("berk read: %v", err)
	}
	if !strings.Contains(string(got), "session cookie dropped on 302") {
		t.Errorf("the ticket did not arrive:\n%s", got)
	}
	// And it arrived without any branch of Ada's being involved.
	if out := run(t, berk.dir, "git", "branch", "--list", "--all"); strings.Contains(out, "jaira") {
		t.Errorf("the ref showed up as a branch: %q", out)
	}
}

// A move made while the other side has already moved the same ticket is
// refused, and the report says who was quicker and where the ticket now is.
// "You lost" on its own would leave the user with nothing to act on.
func TestARejectedWriteNamesWhoWasQuicker(t *testing.T) {
	ada, berk, _ := twoSides(t)

	id := create(t, ada, "shared ticket")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("ada flush: %v", err)
	}

	// Berk takes the ticket over from the ref, without any branch.
	if err := berk.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("berk fetch: %v", err)
	}
	content, lease, err := berk.syncer.Repo.Read(id)
	if err != nil {
		t.Fatalf("berk read: %v", err)
	}
	moved := strings.ReplaceAll(string(content), "status: backlog", "status: review")
	moved = strings.ReplaceAll(moved, "assignee: ada", "assignee: berk")
	moved = strings.ReplaceAll(moved, "updated-by: ada", "updated-by: berk")
	if _, err := berk.syncer.Repo.Write(id, []byte(moved), lease); err != nil {
		t.Fatalf("berk write: %v", err)
	}

	// Ada, who knows nothing of that, writes the ticket again.
	if _, err := ada.store.Mutate(id, func(tk *ticket.Ticket) error {
		return tk.Doc().SetScalar(ticket.FieldStatus, "in-progress")
	}); err != nil {
		t.Fatalf("ada mutate: %v", err)
	}
	reports, err := ada.syncer.Flush()
	if err != nil {
		t.Fatalf("ada flush: %v", err)
	}
	if len(reports) != 1 || reports[0].Outcome != outbox.Rejected {
		t.Fatalf("want one rejected, got %+v", reports)
	}
	w := reports[0].Winner
	if w == nil {
		t.Fatal("a rejection with no winner tells the user nothing")
	}
	if w.UpdatedBy != "berk" && w.Assignee != "berk" {
		t.Errorf("winner is %+v, want berk", w)
	}
	if w.Status != "review" {
		t.Errorf("winner status is %q, want review", w.Status)
	}
	if line := w.Describe(); !strings.Contains(line, "berk") || !strings.Contains(line, "review") {
		t.Errorf("the line shown to the user is %q", line)
	}

	// Ada's local file is untouched by the rejection: the board reconciles it
	// with the ref, and nothing about that belongs to the queue.
	tk, err := ada.store.Load(id)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if tk.Status != "in-progress" {
		t.Errorf("the local file was changed by a rejection: %q", tk.Status)
	}
}

// Taking a ticket off the board removes its ref, and does so through the same
// queue — otherwise the last act of a finished ticket would be the only one
// needing a network.
func TestRecordDeleteTakesTheRefDown(t *testing.T) {
	ada, berk, _ := twoSides(t)

	id := create(t, ada, "finished work")
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if err := berk.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("berk fetch: %v", err)
	}
	if _, err := berk.syncer.Repo.SHA(id); err != nil {
		t.Fatalf("berk should see the ref: %v", err)
	}

	if err := ada.syncer.RecordDelete(id); err != nil {
		t.Fatalf("record delete: %v", err)
	}
	e, ok := ada.syncer.Pending(id)
	if !ok || e.Op != outbox.OpDelete {
		t.Fatalf("delete was not queued: %+v", e)
	}
	if _, err := ada.syncer.Flush(); err != nil {
		t.Fatalf("flush delete: %v", err)
	}
	if err := berk.syncer.Repo.Fetch(); err != nil {
		t.Fatalf("berk re-fetch: %v", err)
	}
	if _, err := berk.syncer.Repo.SHA(id); err == nil {
		t.Error("the ref survived on the other clone")
	}
}

// A board with no remote must keep working exactly as before: writes succeed,
// nothing is queued, and no command waits for git.
func TestABoardWithoutARemoteRecordsNothing(t *testing.T) {
	root := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(root, "home"))
	dir := filepath.Join(root, "plain")
	if err := exec.Command("mkdir", "-p", dir).Run(); err != nil {
		t.Fatal(err)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	s.Actor = "ada"
	y := refsync.New(s, "origin", "ada")
	s.Recorder = y

	id := ticket.NewID(time.Now())
	if _, err := s.Create(map[string]string{
		ticket.FieldID:     id,
		ticket.FieldTitle:  "no remote here",
		ticket.FieldStatus: "backlog",
	}, nil, "# no remote here\n"); err != nil {
		t.Fatalf("create without a repository failed: %v", err)
	}
	if _, ok := y.Pending(id); ok {
		t.Error("a directory that is not a repository queued a ref write")
	}
	reports, err := y.Flush()
	if err != nil {
		t.Fatalf("flush on a board with no remote: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("want nothing to report, got %+v", reports)
	}
}

func TestANilSyncerIsAWorkingRecorder(t *testing.T) {
	var y *refsync.Syncer
	if err := y.Record("01AAA", []byte("x")); err != nil {
		t.Errorf("nil syncer refused a write: %v", err)
	}
	if err := y.RecordDelete("01AAA"); err != nil {
		t.Errorf("nil syncer refused a delete: %v", err)
	}
	if reports, err := y.Flush(); err != nil || reports != nil {
		t.Errorf("nil syncer flush: %v %v", reports, err)
	}
	if _, ok := y.Pending("01AAA"); ok {
		t.Error("nil syncer reported a pending write")
	}
}

func TestWinnerDescribeFallsBackToSomeoneElse(t *testing.T) {
	if got := (refsync.Winner{}).Describe(); !strings.Contains(got, "someone else") {
		t.Errorf("empty winner describes as %q", got)
	}
	w := refsync.Winner{Assignee: "berk", Status: "review"}
	if got := w.Describe(); !strings.Contains(got, "berk") {
		t.Errorf("assignee-only winner describes as %q", got)
	}
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.test",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}

var _ ticket.WriteRecorder = (*refsync.Syncer)(nil)

var _ = gitref.Prefix
