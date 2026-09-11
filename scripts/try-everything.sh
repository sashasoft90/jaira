#!/bin/bash
# Песочница для проверки всего, что сделано. Ничего не трогает вне $R.
# Запуск:  bash try-everything.sh        (J=/путь/к/jaira чтобы взять другой бинарь)
set -u
J=${J:-/tmp/jaira}
R=$(mktemp -d); export JAIRA_HOME=$R/home
say() { echo; echo "=== $*"; }
refs() { git -C "$1" ls-remote origin 'refs/jaira/tickets/*' 2>/dev/null | wc -l; }
files() { ls "$1"/.jaira/tickets 2>/dev/null | wc -l; }
id_of() { grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4; }

say "0. три клона одного bare-репо: ada, berk, carol"
git init -q --bare $R/board.git
for n in ada berk carol; do
  git clone -q $R/board.git $R/$n 2>/dev/null
  git -C $R/$n config user.name $n; git -C $R/$n config user.email $n@x.test
  $J -C $R/$n init >/dev/null
done
echo "   песочница: $R"

say "1. ada заводит тикет для berk — файла у неё НЕ появляется"
A=$($J -C $R/ada create "почини куки на 302" --goal "куки переживают OAuth" \
      --context "репорт из чата про Safari" --dod "переживает редирект" --assignee berk --json | id_of)
echo "   файлов у ada: $(files $R/ada)   рефов на remote: $(refs $R/ada)"

say "2. berk ничего не делал — просто читает"
$J -C $R/berk fetch --quiet
echo "   ^ '@you new ref-only' = тикет назначен тебе и лежит в чужом бранче"

say "3. список и показ работают без файла, запись — нет"
$J -C $R/berk list | sed -n '2,3p'
$J -C $R/berk note $A "попробую написать"; echo "   ^ exit=$? (ждём 3)"

say "4. carol пытается забрать чужой тикет"
$J -C $R/carol fetch --quiet >/dev/null; $J -C $R/carol pull $A; echo "   exit=$? файлов у carol: $(files $R/carol)"

say "5. berk забирает свой — файл появляется ТОЛЬКО у него"
$J -C $R/berk pull $A | head -1
echo "   berk: $(files $R/berk)   ada: $(files $R/ada)   carol: $(files $R/carol)"

say "6. теперь berk может писать"
$J -C $R/berk note $A "начал"; echo "   exit=$?"

say "7. berk отдаёт тикет назад"
$J -C $R/berk release $A | head -2
echo "   файлов у berk: $(files $R/berk)"

say "8. теперь его берёт carol"
$J -C $R/carol fetch --quiet >/dev/null; $J -C $R/carol pull $A | head -1

say "9. ada отбирает силой — громко и со следом на тикете"
$J -C $R/ada fetch --quiet >/dev/null; $J -C $R/ada pull $A --steal | head -2
$J -C $R/ada show $A | grep -A2 'Progress' | tail -2

say "10. офлайн: запись остаётся локально и ждёт"
# считаем только очередь ada: JAIRA_HOME в песочнице общий на все три клона
outbox_ada() { ls $R/home/state/ada-*/outbox/ 2>/dev/null | wc -l; }
git -C $R/ada remote set-url origin https://127.0.0.1:1/nope.git
$J -C $R/ada note $A "написано без сети" >/dev/null 2>&1
echo "   в outbox у ada: $(outbox_ada) запись(ей)   (ждём 1)"
git -C $R/ada remote set-url origin $R/board.git
$J -C $R/ada note $A "сеть вернулась" >/dev/null 2>&1
echo "   после возврата сети осталось: $(outbox_ada)   (ждём 0)"

say "11. фоновый fetch: carol только читает, и всё равно узнаёт новое"
B=$($J -C $R/ada create "второй тикет" --goal g --context "для проверки фонового fetch" --dod d --json | id_of)
rm -f $R/home/state/*/fetch.json
/usr/bin/time -f "   list занял %e s" $J -C $R/carol list >/dev/null
sleep 1; echo "   тикетов у carol теперь: $($J -C $R/carol list --json 2>/dev/null | grep -c '"id"')"

say "12. снапшот: доска как файлы в отдельной ветке"
$J -C $R/ada snapshot
git -C $R/ada ls-tree -r --name-only jaira/board
echo "   повторный прогон:"; $J -C $R/ada snapshot

say "13. тикет доводим до конца — реф ПЕРЕЖИВАЕТ архив"
$J -C $R/ada share >/dev/null 2>&1; git -C $R/ada add -A .jaira .gitignore >/dev/null 2>&1
git -C $R/ada commit -qm "share" >/dev/null 2>&1; git -C $R/ada push -q -u origin master 2>/dev/null
$J -C $R/ada archive $A >/dev/null
echo "   после archive рефов: $(refs $R/ada)  (ждём: реф ещё жив)"
$J -C $R/ada snapshot | tail -1
echo "   после снапшота без приземления: $(refs $R/ada)  (ждём: всё ещё жив)"

say "14. приземляем в master — снапшот жнёт реф"
git -C $R/ada add -A .jaira >/dev/null; git -C $R/ada commit -qm "archive" >/dev/null; git -C $R/ada push -q origin master
$J -C $R/ada snapshot | tail -1
echo "   рефов: $(refs $R/ada)  (ждём: на один меньше)"
echo "   тикет остался в снапшоте:"; git -C $R/ada ls-tree -r --name-only jaira/board | head -3

say "15. история доски — обычный git"
git -C $R/ada log --oneline jaira/board | head -4

say "16. мерж двух веток на одном тикете (add/add)"
cd $R/ada
git checkout -q -b work; $J pull $B >/dev/null 2>&1; $J set $B assignee=ada >/dev/null
git add -A .jaira >/dev/null; git commit -qm "work on $B" >/dev/null
git checkout -q master; git checkout -q -b other; git checkout -q work -- .jaira/tickets 2>/dev/null
F=$(ls .jaira/tickets | head -1)
[ -n "$F" ] && sed -i 's/^status: .*/status: todo/' .jaira/tickets/$F && git add -A .jaira >/dev/null && git commit -qm "other" >/dev/null
git checkout -q work; git merge --no-edit other 2>&1 | tail -1
echo "   маркеров конфликта: $(grep -c '<<<<<<<' .jaira/tickets/$F 2>/dev/null | head -1)  (ждём 0)"
cd - >/dev/null

say "17. борд без remote ведёт себя как раньше"
mkdir -p $R/plain && $J -C $R/plain init >/dev/null
$J -C $R/plain create "без remote" --goal g --context "каталог без git" --dod d | tail -1
echo "   файлов: $(files $R/plain)  (ждём 1 — файл остаётся)"

echo; echo "=== песочница осталась: $R"
echo "    убрать:  rm -rf $R"
