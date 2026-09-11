#!/usr/bin/env bash
# Builds the sandbox for the split-screen walkthrough and prints its path.
# Called during the hidden phase of a recording, so the picture carries only
# the commands that are the point.
set -eu
J=${J:-/tmp/jaira}
R=$(mktemp -d)
git init -q --bare "$R/board.git"
for n in ada grace; do
  git clone -q "$R/board.git" "$R/$n" 2>/dev/null
  git -C "$R/$n" config user.name "$n"
  git -C "$R/$n" config user.email "$n@example.test"
  JAIRA_HOME="$R/home" "$J" -C "$R/$n" init >/dev/null
done
echo "$R"
