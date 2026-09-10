package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/notify"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newFetchCmd is the read half of ticket refs: one round trip, no branch, no
// checkout, no merge.
func newFetchCmd() *cobra.Command {
	var quiet bool
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch the tickets travelling on their own git refs",
		Long: `Fetches refs/jaira/tickets/* from the board's remote and reports what they say.

This is how a ticket assigned to you arrives without anyone sharing a branch:
the ref carries the whole ticket file, so it is readable with no checkout and
nothing to merge. Refs are outside refs/heads, so the default refspec does not
carry them and a teammate who does not use jaira sees none of this.

A ticket newly assigned to you also raises a desktop notification, unless
notifications are turned off in ~/.jaira/settings.json ("notify-off": true).`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := openStore(); err != nil {
				return err
			}
			if err := refs.Usable(); err != nil {
				return fail(ExitError, "no_remote",
					"this board does not carry tickets on refs: %v", err)
			}
			arrivals, err := refs.Incoming()
			if err != nil {
				return err
			}
			// Newly assigned first, then everything else: the reason to run
			// this command is the ticket somebody just handed you.
			sort.SliceStable(arrivals, func(i, j int) bool {
				a, b := arrivals[i], arrivals[j]
				if (a.Mine && a.Changed) != (b.Mine && b.Changed) {
					return a.Mine && a.Changed
				}
				return a.ID < b.ID
			})
			if !quiet {
				announceArrivals(arrivals)
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"arrivals": arrivals})
			}
			printArrivals(cmd.OutOrStdout(), arrivals)
			return nil
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "do not raise a desktop notification")
	return cmd
}

// announceArrivals raises one notification per ticket newly assigned to this
// user. Only Mine and only Changed: a ticket that has been sitting there
// assigned to me since last week is not news, and announcing it on every fetch
// would teach the user to ignore the notification.
func announceArrivals(arrivals []refsync.Arrival) {
	if !settings.Load().NotifyEnabled() {
		return
	}
	for _, a := range arrivals {
		if !a.Mine || !a.Changed {
			continue
		}
		title := a.Title
		if title == "" {
			title = ticket.Handle(a.ID)
		}
		body := fmt.Sprintf("%s — %s", ticket.Handle(a.ID), title)
		if a.Status != "" {
			body += " (" + a.Status + ")"
		}
		notify.Send("jaira: a ticket is yours", body)
	}
}

func printArrivals(w io.Writer, arrivals []refsync.Arrival) {
	if len(arrivals) == 0 {
		fmt.Fprintln(w, "no tickets on refs yet")
		return
	}
	for _, a := range arrivals {
		marks := ""
		if a.Mine {
			marks += " @you"
		}
		if a.Changed {
			marks += " new"
		}
		if !a.Local {
			// The property worth naming: this ticket is in no branch you have.
			marks += " ref-only"
		}
		fmt.Fprintf(w, "%-8s %-12s %s%s\n", ticket.Handle(a.ID), a.Status, a.Title, marks)
	}
}
