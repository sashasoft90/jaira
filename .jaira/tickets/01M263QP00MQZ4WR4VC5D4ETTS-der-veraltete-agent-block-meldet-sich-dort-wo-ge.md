---
id: 01M263QP00MQZ4WR4VC5D4ETTS
title: "Der veraltete Agent-Block meldet sich dort, wo gearbeitet wird"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Aendert jemand die Lanes, erfaehrt die naechste Session davon, ohne dass jemand 'jaira validate' aufruft"
context: |-
  Aus einer Multi-Agent-Session am 10.09.: der jaira-Block in CLAUDE.md nannte eine Lane-Reihenfolge mit 'signoff', die es auf dem Board nicht gab; critique und optimize fehlten. Ein Subagent lief in 'no lane "signoff" is installed', drei Tickets wurden an critique und optimize vorbeigetrieben.
  jaira ERKENNT das bereits. 'jaira validate' meldet woertlich: 'the jaira block in AGENTS.md and CLAUDE.md no longer matches this board's lanes; run jaira update'.
  Der Zustand ist inzwischen behoben - ich habe am 09.09. die Lanes angeglichen und 'jaira update' laufen lassen; 'grep -c signoff' ist in CLAUDE.md und AGENTS.md jetzt 0 und beide zeigen dieselbe Order-Zeile.
  Es fehlt also keine Erkennung. Es fehlt, dass jemand sie sieht: 'jaira validate' ruft im Alltag niemand auf, und der Block ist genau die Datei, die eine Session als Erstes liest - veraltet, ohne es zu sagen.
  Naheliegende Stellen: die erste Board-Ausgabe einer Session, oder 'jaira next' - beides Befehle, die eine Session ohnehin absetzt. Der Stop-Hook aus 'jaira hook print' ist das Vorbild fuer 'jaira sagt selbst Bescheid'.
definition-of-done: "Eine Session, die auf einem Board mit veraltetem Block arbeitet, wird darauf hingewiesen, ohne 'jaira validate' aufzurufen; der Hinweis nennt den Befehl; Skripte und --json bleiben unverschmutzt; ein Test deckt es ab; go test ./... -race gruen"
tags:
  - docs
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:07Z
updated-at: 2026-09-10T16:52:07Z
---

# Der veraltete Agent-Block meldet sich dort, wo gearbeitet wird

## Definition of Done

- [ ] Eine Session, die auf einem Board mit veraltetem Block arbeitet, wird darauf hingewiesen, ohne 'jaira validate' aufzurufen; der Hinweis nennt den Befehl; Skripte und --json bleiben unverschmutzt; ein Test deckt es ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

