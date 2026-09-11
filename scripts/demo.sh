#!/bin/bash
# Демонстрация: как тикет переезжает от одного человека к другому через git-реф.
# Рассчитана на запись: печатает команду, ждёт, показывает результат.
#   bash scripts/demo.sh              — обычный темп
#   SPEED=0 bash scripts/demo.sh      — без пауз (для проверки)
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
echo -e "\033[1mjaira: тикет едет к человеку по git-рефу, без общей ветки\033[0m"
pause 1.5

who "ADA заводит задачу для BERK"
say "ни сервера, ни аккаунтов — только git-remote, который и так есть"
run ada "$J create 'почини куки на 302' --goal 'куки переживают OAuth' --context 'в Safari выкидывает из сессии' --dod 'сессия переживает редирект' --assignee berk"
ID=$(cd $R/ada && $J list --json | grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4)
H=${ID: -6}

say "у ADA на диске файла НЕТ — тикет живёт на своём рефе"
run ada "ls .jaira/tickets/ | wc -l"
run ada "git ls-remote origin 'refs/jaira/*'"

who "BERK, другой человек, другой клон"
say "ни одной ветки ADA у него нет"
run berk "git branch -a | wc -l"
run berk "$J list"
say "один заход в remote — и задача у него на доске"
run berk "$J fetch"

who "читать можно, писать — нет"
run berk "$J note $H 'начну завтра'; echo exit=\$?"
say "потому что это ещё не его: файла здесь нет"

who "BERK забирает задачу"
run berk "$J pull $H"
say "compare-and-swap на рефе: выиграть может ровно один"
run berk "$J note $H 'смотрю на редирект'"

who "ADA пробует ту же задачу"
run ada "$J fetch --quiet >/dev/null; $J pull $H; echo exit=\$?"
run ada "ls .jaira/tickets/ | wc -l"
say "у проигравшего не появляется ничего — дублировать нечего"

who "BERK возвращает задачу на доску"
run berk "$J release $H"
say "и теперь её может взять кто угодно"
run ada "$J fetch --quiet >/dev/null; $J pull $H"

who "доска целиком — в отдельной ветке, как файлы"
run ada "$J snapshot"
run ada "git ls-tree -r --name-only jaira/board"
run ada "git log --oneline jaira/board"

echo; echo -e "\033[1mвсё это — один git-remote. ни сервера, ни демона.\033[0m"
echo -e "\033[2mпесочница: $R\033[0m"
