#!/usr/bin/env bash
# Two panes side by side — ADA on the left, BERK on the right, each in its own
# clone of the same board — and attaches to them. Meant for recording:
#
#   scripts/demo-split-start.sh        # or through scripts/demo-split.tape
#
# Why tmux rather than two windows: a recording captures exactly one terminal.
# The split therefore has to happen INSIDE that one terminal, or the second
# side is not in the picture at all.
set -eu
J=${J:-/tmp/jaira}
S=${SESSION:-jdemo}
R=$(J="$J" "$(dirname "$0")/demo-split-setup.sh")

tmux kill-session -t "$S" 2>/dev/null || true
tmux new-session -d -s "$S" -c "$R/ada" "bash --norc"
A=$(tmux list-panes -t "$S" -F '#{pane_id}' | head -1)
B=$(tmux split-window -h -t "$A" -c "$R/berk" -P -F '#{pane_id}' "bash --norc")

# Names above the panes: without them a viewer cannot tell who is who, and who
# is who is the whole point of the picture.
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
