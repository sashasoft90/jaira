# Split screen: ADA on the left, BERK on the right, in one recording.
#
#   go install github.com/charmbracelet/vhs@latest
#   sudo apt install ttyd ffmpeg tmux
#   vhs scripts/demo-split.tape        → docs/img/demo-split.gif
#
# vhs records exactly one terminal and cannot split by itself. The split
# therefore happens inside the recording, with tmux: typing always goes to the
# active pane, and "Ctrl+B" followed by "o" moves to the other side.

Output docs/img/demo-split.gif

Set Shell bash
Set FontSize 14
Set Width 1500
Set Height 700
Set Padding 15
Set Theme "Catppuccin Mocha"
Set TypingSpeed 30ms

Hide
Type "go build -o /tmp/jaira ./cmd/jaira && bash scripts/demo-split-start.sh" Enter
Sleep 8s
Show
Sleep 1s

# ── ADA files a ticket for BERK ──────────────────────────────────────
Type "jaira create 'session cookie dropped on 302' --goal 'the cookie survives the OAuth round-trip' --context 'reported in chat: Safari logs people out mid-flow' --dod 'session survives the redirect' --assignee berk"
Sleep 500ms
Enter
Sleep 3s

Type "ls .jaira/tickets/ | wc -l"
Enter
Sleep 2s

# ── over to BERK ─────────────────────────────────────────────────────
Ctrl+B
Type "o"
Sleep 1s

Type "git branch -a | wc -l"
Enter
Sleep 2s

Type "jaira fetch"
Enter
Sleep 3s

Type "jaira pull $(jaira list --json | grep -o '\"id\": \"[^\"]*\"' | head -1 | cut -d'\"' -f4)"
Enter
Sleep 3s

# ── back to ADA: the same ticket, and she does not get it ────────────
Ctrl+B
Type "o"
Sleep 1s

Type "jaira fetch --quiet && jaira pull $(jaira list --json | grep -o '\"id\": \"[^\"]*\"' | head -1 | cut -d'\"' -f4)"
Enter
Sleep 4s

Type "ls .jaira/tickets/ | wc -l"
Enter
Sleep 3s
