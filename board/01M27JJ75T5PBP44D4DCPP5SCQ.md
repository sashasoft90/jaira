---
id: 01M27JJ75T5PBP44D4DCPP5SCQ
title: "Auch wer nur die CLI benutzt, bekommt die Refs ohne danach zu fragen"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket, das mir jemand zuweist, erreicht mich auch dann, wenn ich das Board nie oeffne und nur lesende Kommandos benutze - ohne dass ein einziges Kommando auf das Netz wartet"
context: |-
  Gefunden beim Abnehmen von 8566KF (Frage von Alexander): 'wenn wir jaira nicht als Tool starten, kommt der fetch dann ueberhaupt nie?'

  Heute holt die Refs: das TUI alle 60 Sekunden (internal/tui/refs.go) und jede SCHREIBENDE Kommandozeile, weil refsync.Flush vor dem Senden fetcht. Wer nur in der CLI lebt und nur liest - jaira list, next, show - fetcht nie. Fuer den existiert eine Zuweisung nicht, bis er von Hand 'jaira fetch' tippt. Das ist genau die Person, fuer die das Feature gebaut wurde: ein Agent an einer bash-Zeile.

  Ausdruecklich NICHT die Loesung: den Fetch in 'jaira list' haengen. Ein Lesekommando darf nie auf ein Remote warten - das ist die Regel, die das Board schnell haelt, und Alexander hat sie beim Beantworten ausdruecklich bestaetigt ('надо подправить но не в list').

  Das Muster dafuer liegt schon zweimal im Haus und ist beide Male dasselbe: die Update-Pruefung (core/selfupdate/cache.go SpawnRefresh) und der Snapshot (core/snapshot/schedule.go). Stempel VOR dem Lauf, abgekoppeltes Kind, Rekursionsschutz per Umgebungsvariable, kein Wait, stdio auf DevNull, nie aus einem Testbinary. Ein drittes handgeschriebenes Exemplar davon waere die Stelle, an der eine dieser Feinheiten irgendwann fehlt.
definition-of-done: "nach jedem Kommando - auch einem rein lesenden - wird ein abgekoppelter 'jaira fetch' gestartet, wenn der letzte laenger her ist als das eingestellte Intervall (Default zehn Minuten, 'fetch-every-minutes' in settings.json), und kein Kommando wartet darauf; das Muster aus Stempel-vor-dem-Lauf, Rekursionsschutz, kein Wait, DevNull und Testbinary-Schutz steht an EINER Stelle, die der Snapshot-Lauf und dieser Fetch beide benutzen, statt ein drittes Mal abgeschrieben zu werden; JAIRA_NO_FETCH=1 schaltet es ab und ist zugleich der Rekursionsschutz des Kindes; Tests decken ab, dass faellig/nicht faellig richtig entschieden wird und dass ein abgeschaltetes Board nie faellig ist"
tags:
  - cli
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T06:30:31Z
updated-at: 2026-09-11T06:31:04Z
claimed-by: DESKTOP-RFTCH11-13032
claimed-at: 2026-09-11T06:30:57Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
---

# Auch wer nur die CLI benutzt, bekommt die Refs ohne danach zu fragen

## Definition of Done

- [ ] nach jedem Kommando - auch einem rein lesenden - wird ein abgekoppelter 'jaira fetch' gestartet, wenn der letzte laenger her ist als das eingestellte Intervall (Default zehn Minuten, 'fetch-every-minutes' in settings.json), und kein Kommando wartet darauf; das Muster aus Stempel-vor-dem-Lauf, Rekursionsschutz, kein Wait, DevNull und Testbinary-Schutz steht an EINER Stelle, die der Snapshot-Lauf und dieser Fetch beide benutzen, statt ein drittes Mal abgeschrieben zu werden; JAIRA_NO_FETCH=1 schaltet es ab und ist zugleich der Rekursionsschutz des Kindes; Tests decken ab, dass faellig/nicht faellig richtig entschieden wird und dass ein abgeschaltetes Board nie faellig ist

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

