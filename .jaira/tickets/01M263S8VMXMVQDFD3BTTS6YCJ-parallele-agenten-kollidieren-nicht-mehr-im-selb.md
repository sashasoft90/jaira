---
id: 01M263S8VMXMVQDFD3BTTS6YCJ
title: Parallele Agenten kollidieren nicht mehr im selben Arbeitsbaum
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Mehrere Agenten, die gleichzeitig an einem Board arbeiten, treten sich nicht gegenseitig auf die Dateien - entweder weil jaira es merkt, oder weil der Block Worktrees vorschreibt"
context: |-
  Aus einer Multi-Agent-Session am 10.09. mit bis zu vier parallelen Subagenten auf EINEM git-Index. Zwei Zwischenfaelle:
  1. Ein Subagent fuehrte zur Verifikation 'git stash' ueber das GANZE Repo aus. Der komplette uncommittete Baum aller vier Agenten verschwand kurzzeitig. Wiederhergestellt per 'git stash pop', byte-identisch geprueft.
  2. Ein 'git commit' OHNE Pathspec nahm fremde, bereits gestagte Dateien mit. Korrigiert per 'git reset --soft HEAD~1' plus Pathspec-Commit.
  Der Kern: 'jaira claim' markiert die Zustaendigkeit fuer ein TICKET. Ueber DATEIEN sagt es nichts. Zwei Agenten mit sauber geclaimten Tickets im selben Verzeichnis kollidieren trotzdem.
  Im Schema gibt es kein Feld fuer Dateibereiche - geprueft in core/ticket/schema.go.
  Dasselbe ist in dieser Session nochmal passiert: ein Agent nahm beim ersten Commit fuenf fremde ungetrackte Tickets mit und musste sie per reset --soft zuruecknehmen.
  Zwei Wege, vor dem Bauen abzuwaegen: ein Feld am Ticket ('--owns src/connectors/'), das 'jaira claim' gegen andere aktive Claims prueft - oder gar kein Code, sondern eine ausdrueckliche Worktree-Empfehlung im generierten Block. Letzteres ist billiger und loest Fall 1 mit.
definition-of-done: "Entweder prueft 'jaira claim' einen deklarierten Dateibereich gegen andere aktive Claims, oder der generierte Block schreibt fuer parallele Agenten Worktrees vor und begruendet es; die Entscheidung steht mit Begruendung im Ticket; falls Code: Tests und go test ./... -race gruen"
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:59Z
updated-at: 2026-09-10T16:52:59Z
---

# Parallele Agenten kollidieren nicht mehr im selben Arbeitsbaum

## Definition of Done

- [ ] Entweder prueft 'jaira claim' einen deklarierten Dateibereich gegen andere aktive Claims, oder der generierte Block schreibt fuer parallele Agenten Worktrees vor und begruendet es; die Entscheidung steht mit Begruendung im Ticket; falls Code: Tests und go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

