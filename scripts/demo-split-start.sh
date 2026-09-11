#!/usr/bin/env bash
# Zwei Panes nebeneinander, links ADA, rechts BERK, beide in einem eigenen Klon
# desselben Boards — und haengt sich dran. Fuer die Aufnahme gedacht:
#
#   scripts/demo-split-start.sh        # oder ueber scripts/demo-split.tape
#
# Warum tmux und nicht zwei Fenster: eine Aufnahme zeichnet genau ein Terminal
# auf. Der geteilte Bildschirm muss deshalb INNERHALB dieses einen Terminals
# entstehen, sonst ist die zweite Seite im Bild nicht zu sehen.
set -eu
J=${J:-/tmp/jaira}
S=${SESSION:-jdemo}
R=$(J="$J" "$(dirname "$0")/demo-split-setup.sh")

tmux kill-session -t "$S" 2>/dev/null || true
tmux new-session -d -s "$S" -c "$R/ada" "bash --norc"
A=$(tmux list-panes -t "$S" -F '#{pane_id}' | head -1)
B=$(tmux split-window -h -t "$A" -c "$R/berk" -P -F '#{pane_id}' "bash --norc")

# Namen ueber den Panes: ohne sie ist auf der Aufnahme nicht zu sehen, wer wer
# ist, und genau das ist die Aussage des ganzen Bildes.
tmux set -t "$S" -g pane-border-status top
tmux set -t "$S" -g pane-border-format ' #{pane_title} '
tmux set -t "$S" -g status off
tmux select-pane -t "$A" -T " ADA "
tmux select-pane -t "$B" -T " BERK "

for p in "$A" "$B"; do
  tmux send-keys -t "$p" "export JAIRA_HOME=$R/home PATH=$(dirname "$J"):\$PATH PS1='\$ '; clear" Enter
done
tmux select-pane -t "$A"
echo "$R" > /tmp/jaira-demo-sandbox
exec tmux attach -t "$S"
