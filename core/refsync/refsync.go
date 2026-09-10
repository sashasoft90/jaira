// Package refsync connects the ticket store to the ref a ticket travels on.
//
// It is the only place that knows all three: core/ticket writes the file,
// core/gitref carries it on refs/jaira/tickets/<id>, core/outbox holds what has
// not been sent. Keeping that knowledge here is what lets core/ticket stay
// unaware of git and core/gitref stay unaware of tickets.
//
// Nothing here pushes on the write path. A ticket write is recorded into the
// outbox and returns immediately, because the write path runs inside the TUI
// and inside the git merge driver as well as in the CLI, and none of those may
// wait for a network round trip. Sending is a separate act (Flush), run where
// waiting is acceptable.
package refsync

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/merge"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Syncer records ticket writes for the remote and sends them.
type Syncer struct {
	Repo *gitref.Repo
	Box  *outbox.Box

	// Actor is who the queued write is from, so a flush can say whose write it
	// is sending and a rejection can name both sides.
	Actor string

	// Store is the working tree's tickets, used only to tell a ticket that
	// arrived on a ref alone from one this clone already has a file for.
	Store *ticket.Store

	// IsMine decides whether an assignee is this user. It is a function rather
	// than a name because "me" includes the aliases a person has recorded, and
	// core/identity already knows how to answer it — this package should not
	// hold a second opinion.
	IsMine func(assignee string) bool

	// SeenPath records which ref SHA was last reported for each ticket, so
	// being handed a ticket is announced once instead of on every fetch.
	SeenPath string

	// usable caches whether this checkout can carry refs at all. It is asked
	// once per process rather than per write: the answer is a property of the
	// checkout, and every ticket write would otherwise pay two git invocations
	// to learn the same thing.
	once   sync.Once
	usable error

	// dirty records that this process queued something. It is what lets a read
	// command stay entirely off the network: only a command that actually
	// wrote a ticket has a reason to send, and jaira list must never wait for
	// a remote.
	dirty bool
}

// Dirty reports whether this process queued a write. A caller uses it to decide
// whether flushing is worth a round trip at all.
func (y *Syncer) Dirty() bool { return y != nil && y.dirty }

// Usable reports whether this board carries tickets on refs. A board with no
// repository or no remote is not an error — it is jaira used the way it also
// works, and the answer is simply no.
func (y *Syncer) Usable() error {
	if y == nil {
		return gitref.ErrNoRepo
	}
	y.once.Do(func() { y.usable = y.Repo.Usable() })
	return y.usable
}

// New builds a syncer for a store, or returns nil when this board has no
// business talking to a remote: no git, no repository, or no remote configured.
// A nil *Syncer is a valid recorder that records nothing, so callers do not
// need to branch.
func New(s *ticket.Store, remote, actor string) *Syncer {
	repo := &gitref.Repo{Dir: s.Root, Remote: remote, AuthorName: actor}
	return &Syncer{
		Repo:     repo,
		Box:      outbox.At(s),
		Actor:    actor,
		Store:    s,
		SeenPath: filepath.Join(s.StateDir(), "refs-seen.json"),
	}
}

// Record queues the ticket's current bytes for its ref. It satisfies
// ticket.WriteRecorder.
//
// The lease is the local ref's SHA, which is only ever moved after the remote
// has accepted a push — so it is exactly "the state the remote last confirmed
// to us", which is what a compare-and-swap has to be checked against. A ticket
// with no ref yet leases the empty string, meaning "I expect this ref not to
// exist".
func (y *Syncer) Record(id string, content []byte) error {
	if y == nil || y.Usable() != nil {
		// A board without a repository or without a remote does not carry
		// tickets on refs. Saying so on every single write would be noise, and
		// queueing a write that can never be sent would be worse: the card
		// would show "unsent" forever.
		return nil
	}
	lease, err := y.Repo.SHA(id)
	if err != nil && !errors.Is(err, gitref.ErrNoRef) {
		return err
	}
	if err := y.Box.Queue(id, outbox.OpWrite, content, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// RecordDelete queues the removal of a ticket's ref, for a ticket that has been
// logged or archived. It goes through the same queue as a write so that taking
// a ticket off the board works offline too — otherwise the last act of a
// finished ticket would be the one act that needs a network.
func (y *Syncer) RecordDelete(id string) error {
	if y == nil || y.Usable() != nil {
		return nil
	}
	lease, err := y.Repo.SHA(id)
	if err != nil {
		if errors.Is(err, gitref.ErrNoRef) {
			// Nothing to delete: the ticket never reached a ref.
			return nil
		}
		return err
	}
	if err := y.Box.Queue(id, outbox.OpDelete, nil, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// Winner describes the state that beat a rejected write, read from the ref
// itself. A rejection that only says "you lost" leaves the user to go and find
// out who and what — which they cannot do without knowing the ref exists.
type Winner struct {
	ID        string `json:"id"`
	Assignee  string `json:"assignee,omitempty"`
	Status    string `json:"status,omitempty"`
	UpdatedBy string `json:"updated-by,omitempty"`
	UpdatedAt string `json:"updated-at,omitempty"`
}

// Describe renders the winner as the one line a command can print.
func (w Winner) Describe() string {
	who := w.UpdatedBy
	if who == "" {
		who = w.Assignee
	}
	if who == "" {
		who = "someone else"
	}
	parts := []string{who + " wrote it first"}
	if w.Status != "" {
		parts = append(parts, "it is now in "+w.Status)
	}
	if w.Assignee != "" && !strings.EqualFold(w.Assignee, who) {
		parts = append(parts, "assigned to "+w.Assignee)
	}
	if w.UpdatedAt != "" {
		parts = append(parts, "at "+w.UpdatedAt)
	}
	return strings.Join(parts, ", ")
}

// Holds renders the same state as a statement of ownership rather than of a
// race. The two refusals are different facts and must not borrow each other's
// wording: "berk wrote it first" is about who won a push, "berk has it" is
// about who the ticket belongs to.
func (w Winner) Holds() string {
	who := w.Assignee
	if who == "" {
		who = w.UpdatedBy
	}
	if who == "" {
		who = "somebody else"
	}
	out := who + " has it"
	if w.Status != "" {
		out += ", it is in " + w.Status
	}
	return out
}

// Report is one queued write's fate, with the other side named when we lost.
type Report struct {
	outbox.Result
	Winner *Winner `json:"winner,omitempty"`
}

// Flush sends what the outbox holds and reports each entry, naming who was
// quicker for every write the remote refused.
//
// The fetch comes first, and deliberately: a queued write leases a SHA this
// clone has seen, so without fetching, the very first push would be refused
// for a race that has in fact already been resolved on the remote. Fetching
// also makes the winner readable, which is the difference between "you lost"
// and "Berk moved it to review".
func (y *Syncer) Flush() ([]Report, error) {
	if y == nil || y.Usable() != nil {
		return nil, nil
	}
	entries, err := y.Box.List()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		// An unreachable remote is not an error here: it is the state the
		// outbox exists for. The queue stays as it is and the next command
		// tries again.
		if errors.Is(err, gitref.ErrOffline) || errors.Is(err, gitref.ErrNoGit) {
			return nil, nil
		}
		// A remote that exists but is not a repository is the opposite case:
		// no amount of waiting fixes it, and staying quiet would leave every
		// write queued with nobody told why.
		return nil, err
	}

	results, err := y.Box.Flush(y.Repo)
	if err != nil {
		return nil, err
	}
	reports := make([]Report, 0, len(results))
	for _, r := range results {
		rep := Report{Result: r}
		if r.Outcome == outbox.Rejected {
			if w, err := y.winner(r.ID); err == nil {
				rep.Winner = w
			}
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

// Pending reports an unsent write for one ticket, for the marker on a board
// card.
func (y *Syncer) Pending(id string) (outbox.Entry, bool) {
	if y == nil {
		return outbox.Entry{}, false
	}
	return y.Box.Pending(id)
}

// PendingAge is how long this ticket's write has been waiting, or zero.
func (y *Syncer) PendingAge(id string) time.Duration {
	e, ok := y.Pending(id)
	if !ok {
		return 0
	}
	return e.Age()
}

// winner reads the state now on the ref.
func (y *Syncer) winner(id string) (*Winner, error) {
	content, _, err := y.Repo.Read(id)
	if err != nil {
		return nil, err
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, fmt.Errorf("refsync: the ref for %s is not a ticket: %w", id, err)
	}
	w := &Winner{ID: id}
	for _, f := range []struct {
		key string
		dst *string
	}{
		{ticket.FieldAssignee, &w.Assignee},
		{ticket.FieldStatus, &w.Status},
		{ticket.FieldUpdatedBy, &w.UpdatedBy},
		{ticket.FieldUpdatedAt, &w.UpdatedAt},
	} {
		if v, ok, err := d.Scalar(f.key); err == nil && ok {
			*f.dst = v
		}
	}
	return w, nil
}

// Arrival is a ticket as the ref now has it, seen from this clone.
type Arrival struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
	SHA      string `json:"sha"`

	// Mine says the ticket is assigned to this user. It is the whole reason
	// this exists: being handed a ticket is the event nobody currently learns
	// about.
	Mine bool `json:"mine"`

	// Changed says the ref moved since this clone last looked. A ticket that
	// has been sitting there assigned to me for a week is not news, and
	// announcing it on every fetch would train the user to ignore the thing.
	Changed bool `json:"changed"`

	// Local says a ticket file for this id exists in the working tree. False
	// means the ticket reached this clone only on its ref — the case a
	// teammate's unmerged branch cannot deliver.
	Local bool `json:"local"`
}

// seen is the per-clone record of which ref SHA was last reported, so an
// arrival is announced once rather than on every fetch.
type seen map[string]string

// Incoming fetches the ticket refs and reports what they now say.
//
// It is the read half of the feature: a fetch of refs/jaira/tickets/* costs one
// round trip and needs no branch, no checkout and no merge, so a clone can be
// told "this is yours now" without anyone having pushed a branch.
func (y *Syncer) Incoming() ([]Arrival, error) {
	if y == nil || y.Usable() != nil {
		return nil, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		if errors.Is(err, gitref.ErrOffline) || errors.Is(err, gitref.ErrNoGit) {
			return nil, nil
		}
		return nil, err
	}
	ids, err := y.Repo.List()
	if err != nil {
		return nil, err
	}
	known := y.loadSeen()
	next := make(seen, len(ids))
	var out []Arrival
	for _, id := range ids {
		content, sha, err := y.Repo.Read(id)
		if err != nil {
			// A ref whose tree is not a ticket is not worth failing a fetch
			// over; it is also not something this tool wrote.
			continue
		}
		next[id] = sha
		a := Arrival{ID: id, SHA: sha, Changed: known[id] != sha}
		d, err := ticket.ParseDoc(content)
		if err != nil {
			continue
		}
		for _, f := range []struct {
			key string
			dst *string
		}{
			{ticket.FieldTitle, &a.Title},
			{ticket.FieldStatus, &a.Status},
			{ticket.FieldAssignee, &a.Assignee},
		} {
			if v, ok, err := d.Scalar(f.key); err == nil && ok {
				*f.dst = v
			}
		}
		if y.IsMine != nil {
			a.Mine = y.IsMine(a.Assignee)
		}
		a.Local = y.hasLocal(id)
		out = append(out, a)
	}
	y.saveSeen(next)
	return out, nil
}

// isMe answers for the caller's own name, falling back to a plain comparison
// when nobody told this syncer how identity works.
func (y *Syncer) isMe(who string) bool {
	if y.IsMine != nil {
		return y.IsMine(who)
	}
	return strings.EqualFold(strings.TrimSpace(who), strings.TrimSpace(y.Actor))
}

// hasLocal reports whether the working tree holds a file for this ticket.
func (y *Syncer) hasLocal(id string) bool {
	if y.Store == nil {
		return false
	}
	_, err := y.Store.Load(id)
	return err == nil
}

func (y *Syncer) loadSeen() seen {
	s := seen{}
	if y.SeenPath == "" {
		return s
	}
	b, err := os.ReadFile(y.SeenPath)
	if err != nil {
		return s
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return seen{}
	}
	return s
}

// saveSeen records what was just reported. A failure here is silent on purpose:
// the worst it costs is one repeated notification, and refusing a fetch over it
// would cost the feature.
func (y *Syncer) saveSeen(s seen) {
	if y.SeenPath == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(y.SeenPath), 0o755); err != nil {
		return
	}
	if b, err := json.MarshalIndent(s, "", "  "); err == nil {
		_ = os.WriteFile(y.SeenPath, append(b, '\n'), 0o644)
	}
}

// Reconciled is one ticket as the board should show it: the local file and the
// ref brought together, with the two state markers the user needs.
type Reconciled struct {
	ID string `json:"id"`

	// Content is what to show. On a clean merge it is the two sides merged; a
	// ticket that exists on only one side is that side unchanged.
	Content []byte `json:"-"`

	// Conflicts are the prose fields the merge could not settle. They are the
	// existing conflict path, not a new one: the same fields 'jaira resolve'
	// handles after a git merge.
	Conflicts []merge.Conflict `json:"conflicts,omitempty"`

	// RefOnly says the ticket reached this clone on its ref alone — it is in
	// no branch here. This is the case the whole feature exists for.
	RefOnly bool `json:"ref-only"`

	// Unsent says this clone holds a write for the ticket that has not reached
	// the remote yet.
	Unsent bool `json:"unsent"`
}

// Reconcile brings the local ticket file and the ref together for one ticket.
//
// It is deliberately not "show whichever is newer". The two sides are merged
// field by field by core/merge, which is the same resolver the git merge driver
// uses: status by progress along the lane chain (never backwards), lists by
// union, other scalars by updated-at, prose as a real conflict. A newer local
// edit must not drag a ticket back out of review because someone touched it
// after the reviewer did.
//
// The base for that merge is the ref's parent blob — the state both sides
// started from — which is available precisely because a ref commit keeps its
// leased SHA as parent.
func (y *Syncer) Reconcile(id string, lanes *lane.Set) (*Reconciled, error) {
	out := &Reconciled{ID: id}
	if y == nil {
		return out, nil
	}
	_, out.Unsent = y.Pending(id)

	var local []byte
	if y.Store != nil {
		if t, err := y.Store.Load(id); err == nil {
			if b, err := os.ReadFile(t.Path); err == nil {
				local = b
			}
		}
	}
	if y.Usable() != nil {
		out.Content = local
		return out, nil
	}
	theirs, _, err := y.Repo.Read(id)
	if err != nil {
		// No ref: the local file is the whole story.
		out.Content = local
		return out, nil
	}
	if local == nil {
		out.RefOnly = true
		out.Content = theirs
		return out, nil
	}
	base, err := y.Repo.ReadParent(id)
	if err != nil {
		// The ticket's first write has no earlier state; the two sides are all
		// there is, and core/merge treats an empty base as "both added it".
		base = local
	}
	res, err := merge.Merge(base, local, theirs, lanes)
	if err != nil {
		// A merge that cannot be performed must not blank the card: showing
		// the local file is always defensible, since it is what this clone
		// wrote.
		out.Content = local
		return out, nil
	}
	out.Content = res.Merged
	out.Conflicts = res.Conflicts
	return out, nil
}

// ErrTaken means the ticket on the ref already belongs to somebody else.
//
// It is a separate refusal from a lost race, and both are needed. The
// compare-and-swap only rules out two writes landing in the same instant; it
// cannot say "this ticket is not yours", because a pull re-reads the ref
// immediately before writing and would therefore always hold a current lease.
// Ownership is a fact recorded in the ticket, so it is checked as one.
var ErrTaken = errors.New("refsync: the ticket already belongs to someone else")

// Pulled reports what a pull did.
type Pulled struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path,omitempty"`

	// AlreadyHere means the file was already on this machine and nothing
	// needed doing. A pull is safe to repeat, so this is a success.
	AlreadyHere bool `json:"already-here"`

	// Winner is set when the ref refused the take-over: somebody else got
	// there first, and no file was written.
	Winner *Winner `json:"winner,omitempty"`
}

// Pull takes a ticket that lives on its ref and makes it this clone's to work
// on: it claims the ticket on the ref and only then writes the file here.
//
// Two things gate it, and neither is enough alone:
//
//   - The ref's assignee. A ticket somebody else has already pulled is
//     refused with ErrTaken, naming them. This is the check that makes "one
//     clone has the file" true, because a pull re-reads the ref right before
//     writing and so would otherwise always hold a current lease.
//   - The compare-and-swap on the push. That covers the remaining case the
//     first check cannot see: two clones pulling in the same instant, each
//     having read an unassigned ticket.
//
// The order is the other half of the mechanism. The push goes first. Writing
// the file first would leave the loser holding a ticket file that belongs to
// somebody else — the duplicate this design exists to make impossible rather
// than to resolve.
//
// This is also the one write in jaira that is synchronous and does not go
// through the outbox. That is not an exception to the offline rule but a
// consequence of it: with no route to the remote there is no ref to read, so
// there is nothing to take over in the first place.
func (y *Syncer) Pull(id string, steal bool) (*Pulled, error) {
	if y == nil || y.Store == nil {
		return nil, errors.New("refsync: no board here")
	}
	if err := y.Usable(); err != nil {
		return nil, err
	}
	// Already here: a repeated pull is a successful no-op, the same as moving a
	// ticket to the lane it is already in.
	if t, err := y.Store.Load(id); err == nil {
		return &Pulled{ID: t.ID, Title: t.Title, Path: t.Path, AlreadyHere: true}, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		return nil, err
	}
	content, lease, err := y.Repo.Read(id)
	if err != nil {
		return nil, err
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, fmt.Errorf("refsync: the ref for %s does not carry a ticket: %w", id, err)
	}
	if holder, _, _ := d.Scalar(ticket.FieldAssignee); strings.TrimSpace(holder) != "" && !y.isMe(holder) && !steal {
		out := &Pulled{ID: id}
		if w, wErr := y.winner(id); wErr == nil {
			out.Winner = w
		}
		return out, fmt.Errorf("%w: %s has it", ErrTaken, holder)
	}
	if err := d.SetScalar(ticket.FieldAssignee, y.Actor); err != nil {
		return nil, err
	}
	if err := ticket.Touch(d, time.Now()); err != nil {
		return nil, err
	}
	if err := ticket.TouchBy(d, y.Actor); err != nil {
		return nil, err
	}
	taken := d.Bytes()

	if _, err := y.Repo.Write(id, taken, lease); err != nil {
		if errors.Is(err, gitref.ErrRaceLost) {
			out := &Pulled{ID: id}
			if w, wErr := y.winner(id); wErr == nil {
				out.Winner = w
			}
			return out, err
		}
		return nil, err
	}

	title, _, _ := d.Scalar(ticket.FieldTitle)
	path := filepath.Join(y.Store.TicketsDir(), ticket.Filename(id, title))
	if err := os.MkdirAll(y.Store.TicketsDir(), 0o755); err != nil {
		return nil, err
	}
	// Written with the bytes the remote accepted, not with a freshly built
	// document: what is on the board here and what the team sees on the ref are
	// then the same file, byte for byte.
	if err := ticket.WriteAtomic(path, taken); err != nil {
		return nil, err
	}
	return &Pulled{ID: id, Title: title, Path: path}, nil
}
