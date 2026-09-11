#!/usr/bin/env bash
# Baut die Sandbox fuer die geteilte Vorfuehrung und druckt die zwei Pfade.
# Wird in der versteckten Phase der Aufnahme aufgerufen, damit im Bild nur die
# eigentlichen Kommandos stehen.
set -eu
J=${J:-/tmp/jaira}
R=$(mktemp -d)
git init -q --bare "$R/board.git"
for n in ada berk; do
  git clone -q "$R/board.git" "$R/$n" 2>/dev/null
  git -C "$R/$n" config user.name "$n"
  git -C "$R/$n" config user.email "$n@example.test"
  JAIRA_HOME="$R/home" "$J" -C "$R/$n" init >/dev/null
done
echo "$R"
