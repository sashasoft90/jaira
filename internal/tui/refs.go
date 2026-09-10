package tui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/notify"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// refFetchEvery is how often the board looks at the remote. Two seconds is
// right for the local rescan; a network round trip on that cadence is not, and
// nobody is handed a ticket twice a minute.
const refFetchEvery = 60 * time.Second

// refFetchedMsg carries the result of a background fetch back into the update
// loop.
type refFetchedMsg struct {
	arrivals []refsync.Arrival
	err      error
}

// newSyncer wires the board's ticket writes to the refs they travel on, with
// the same settings and the same idea of "me" the CLI uses.
func newSyncer(s *ticket.Store, me string) *refsync.Syncer {
	y := refsync.New(s, settings.Load().RemoteName(), me)
	y.IsMine = func(assignee string) bool {
		return assignee != "" && identity.IsMe(s.Root, assignee)
	}
	return y
}

// refMarks is the local, network-free half: what this clone can say about refs
// without asking anyone.
//
// Both answers come from files already on disk — the outbox and the local refs
// — because this runs inside reload, on the two-second timer, and nothing on
// that path may wait for a remote.
func (m *Model) refMarks() {
	m.unsent = map[string]bool{}
	m.refOnly = 0
	if m.refSync == nil || m.refSync.Usable() != nil {
		return
	}
	if entries, err := m.refSync.Box.List(); err == nil {
		for _, e := range entries {
			m.unsent[e.ID] = true
		}
	}
	ids, err := m.refSync.Repo.List()
	if err != nil {
		return
	}
	for _, id := range ids {
		if _, err := m.store.Load(id); err != nil {
			// On a ref, in no branch here. Counted rather than shown as a
			// card: a card the board cannot let you act on needs a read-only
			// passthrough path of its own, and 'jaira fetch' already lists
			// these in full.
			m.refOnly++
		}
	}
}

// fetchRefs asks the remote what the ticket refs now say, off the update loop.
//
// It is a command rather than a call inside Update for the reason every network
// touch in a TUI has to be: View and Update must never block. Bubble Tea runs
// the returned function on its own goroutine and delivers the result as a
// message, which is the documented way to do exactly this.
func fetchRefs(y *refsync.Syncer) tea.Cmd {
	if y == nil {
		return nil
	}
	return func() tea.Msg {
		arrivals, err := y.Incoming()
		return refFetchedMsg{arrivals: arrivals, err: err}
	}
}

// refTick re-arms the background fetch.
func refTick() tea.Cmd {
	return tea.Tick(refFetchEvery, func(time.Time) tea.Msg { return refTickMsg{} })
}

type refTickMsg struct{}

// announceArrivals says what arrived, and raises a desktop notification for a
// ticket that has just become mine.
//
// Only mine, and only changed: a ticket that has been assigned to me since last
// week is not news, and a board that pops up on every fetch teaches the user to
// ignore it.
func (m *Model) announceArrivals(arrivals []refsync.Arrival) {
	notifyOK := settings.Load().NotifyEnabled()
	var mine []refsync.Arrival
	for _, a := range arrivals {
		if a.Mine && a.Changed {
			mine = append(mine, a)
		}
	}
	if len(mine) == 0 {
		return
	}
	for _, a := range mine {
		if !notifyOK {
			break
		}
		title := a.Title
		if title == "" {
			title = ticket.Handle(a.ID)
		}
		notify.Send("jaira: a ticket is yours", fmt.Sprintf("%s — %s", ticket.Handle(a.ID), title))
	}
	if len(mine) == 1 {
		m.notify(fmt.Sprintf("%s is yours: %s", ticket.Handle(mine[0].ID), mine[0].Title), false)
		return
	}
	m.notify(fmt.Sprintf("%d tickets were assigned to you", len(mine)), false)
}
