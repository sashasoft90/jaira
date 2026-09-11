#!/bin/bash
# A walkthrough of how a ticket travels from one person to another on a git ref.
# Written to be recorded: it prints the command, waits, then shows the result.
#
#   bash scripts/demo.sh              normal pace
#   SPEED=0 bash scripts/demo.sh      no pauses, for checking it still works
set -u
J=${J:-/tmp/jaira}
P=${SPEED:-1}
pause() { sleep $(echo "$P * ${1:-1}" | bc 2>/dev/null || echo 1); }

R=$(mktemp -d); export JAIRA_HOME=$R/home
git init -q --bare $R/board.git
for n in ada berk; do
  git clone -q $R/board.git $R/$n 2>/dev/null
  git -C $R/$n config user.name $n; git -C $R/$n config user.email $n@example.test
  $J -C $R/$n init >/dev/null
done

who()  { echo; echo -e "\033[1;36m─── $1 ───\033[0m"; pause 0.6; }
say()  { echo -e "\033[2m# $1\033[0m"; pause 0.8; }
run()  { echo -e "\033[1;32m$\033[0m $2"; pause 0.5; ( cd $R/$1 && eval "$2" ); pause 1.2; }

clear
echo -e "\033[1mjaira: a ticket reaches a person over a git ref, with no shared branch\033[0m"
pause 1.5

who "ADA files a ticket for BERK"
say "no server and no accounts — only the git remote that is already there"
run ada "$J create 'session cookie dropped on 302' --goal 'the cookie survives the OAuth round-trip' --context 'reported in chat: Safari logs people out mid-flow' --dod 'session survives the redirect' --assignee berk"
ID=$(cd $R/ada && $J list --json | grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4)
H=${ID: -6}

say "nothing on ADA's disk: the ticket lives on its own ref"
run ada "ls .jaira/tickets/ | wc -l"
run ada "git ls-remote origin 'refs/jaira/*'"

who "BERK — another person, another clone"
say "he has none of ADA's branches"
run berk "git branch -a | wc -l"
run berk "$J list"
say "one round trip to the remote, and the ticket is on his board"
run berk "$J fetch"

who "reading works, writing does not"
run berk "$J note $H 'will start tomorrow'; echo exit=\$?"
say "because it is not his yet — there is no file here"

who "BERK takes the ticket"
run berk "$J pull $H"
say "a compare-and-swap on the ref: exactly one clone can win"
run berk "$J note $H 'looking at the redirect'"

who "ADA tries the same ticket"
run ada "$J fetch --quiet >/dev/null; $J pull $H; echo exit=\$?"
run ada "ls .jaira/tickets/ | wc -l"
say "the loser gets nothing at all — there is no second copy to duplicate"

who "BERK hands it back"
run berk "$J release $H"
say "and now anybody can take it"
run ada "$J fetch --quiet >/dev/null; $J pull $H"

who "the whole board, as files, on a branch of its own"
run ada "$J snapshot"
run ada "git ls-tree -r --name-only jaira/board"
run ada "git log --oneline jaira/board"

echo; echo -e "\033[1mall of it over one git remote. no server, no daemon.\033[0m"
echo -e "\033[2msandbox: $R\033[0m"
