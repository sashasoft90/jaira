# Geteilter Bildschirm: links ADA, rechts BERK, in einer Aufnahme.
#
#   go install github.com/charmbracelet/vhs@latest
#   sudo apt install ttyd ffmpeg tmux
#   vhs scripts/demo-split.tape        → docs/img/demo-split.gif
#
# vhs zeichnet genau ein Terminal auf und kann selbst nicht teilen. Die Teilung
# entsteht deshalb innerhalb der Aufnahme, mit tmux — getippt wird immer ins
# aktive Pane, und "Ctrl+B" gefolgt von "o" wechselt die Seite.

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

# ── ADA legt ein Ticket fuer BERK an ─────────────────────────────────
Type "jaira create 'почини куки на 302' --goal 'куки переживают OAuth' --context 'в Safari выкидывает из сессии' --dod 'сессия переживает редирект' --assignee berk"
Sleep 500ms
Enter
Sleep 3s

Type "ls .jaira/tickets/ | wc -l"
Enter
Sleep 2s

# ── zu BERK ──────────────────────────────────────────────────────────
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

# ── zurueck zu ADA: dieselbe Karte, und sie bekommt sie nicht ────────
Ctrl+B
Type "o"
Sleep 1s

Type "jaira fetch --quiet && jaira pull $(jaira list --json | grep -o '\"id\": \"[^\"]*\"' | head -1 | cut -d'\"' -f4)"
Enter
Sleep 4s

Type "ls .jaira/tickets/ | wc -l"
Enter
Sleep 3s
