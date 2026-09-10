---
id: 01M26F9333EWZS2TN8EHPTQ3XT
title: "Ein Snapshot-Branch traegt die Tafel als Dateien, ohne dass jemand ihn auscheckt"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Ein elternloser Branch jaira/board haelt zu jedem Zeitpunkt genau die Tickets, die gerade auf Refs liegen, wird per Plumbing ohne Checkout geschrieben und kann nie mit den Arbeitsdateien kollidieren"
context: |-
  Sobald ein Ticket nur noch auf seinem Ref lebt und erst beim Uebernehmen als Datei landet (siehe RFC7GA), liegt der unbearbeitete Backlog nur auf dem Remote. Geht das Remote verloren, ist er weg - heute steht er in der git-Historie.

  Der Snapshot-Branch schliesst das, und zwar als Opt-in: er ist ein Backup und nicht der Speicherort.

  Gemessen, nicht vermutet - der ganze Ablauf wurde von Hand mit Plumbing durchgespielt:

  - Das Baumobjekt wird mit hash-object/mktree/commit-tree gebaut, genau wie core/gitref es fuer ein Ticket-Ref schon tut. Kein Checkout, kein Anfassen des Arbeitsbaums: der Test lief auf master, git status blieb leer.
  - Der erste Snapshot-Commit hat keinen Parent (elternlos), jeder weitere den vorigen. Damit ist 'git log jaira/board' die Geschichte der Tafel und 'git diff jaira/board~1 jaira/board' sagt, was sich geaendert hat. Belegt: 'board/01AAA.md | 6 ------ / board/01CCC.md | 6 ++++++'.
  - Hinzufuegen und Entfernen ergeben sich von selbst, ohne Vergleich: der Baum wird jedes Mal aus dem AKTUELLEN Satz Refs neu gebaut. Ein Ticket, dessen Ref beim Logbuch geloescht wurde, ist im naechsten Snapshot einfach nicht mehr drin - und bleibt in den vorigen Commits auffindbar.

  Zwei Entscheidungen, die nicht offensichtlich sind und beide Gruende haben:

  1. Die Dateien liegen unter board/<id>.md und NICHT unter .jaira/tickets/. Wuerden sie am selben Pfad liegen, bekaeme der erste, der diesen Branch versehentlich merged, einen add/add-Konflikt ueber alle Tickets auf einmal - die schlimmere Version des Problems, von dem dieser Entwurf weggeht. Ein anderer Pfad macht die Kollision konstruktiv unmoeglich.
  2. Angehaengt statt force-gepusht. Ein Snapshot ist ein Commit mit dem vorigen als Parent, damit es Historie gibt und zwei gleichzeitige Snapshots mit demselben --force-with-lease aufloesen, das core/gitref schon benutzt.

  Bekannte Schwaeche, ausdruecklich: ein Snapshot ist ein Stand von damals. Er braucht einen Befehl und optional einen Timer, sonst veraltet er. Das Arbeitsstand liegt immer auf den Refs; der Snapshot ist das Backup.
definition-of-done: "'jaira snapshot' baut den Branch aus dem aktuellen Satz refs/jaira/tickets/* per hash-object, mktree und commit-tree, ohne Checkout und ohne den Arbeitsbaum zu beruehren, und laesst sich in jedem Repo-Zustand aufrufen; der erste Commit ist elternlos, jeder weitere haengt am vorigen, sodass git log die Geschichte der Tafel ist; ein Ticket, dessen Ref verschwunden ist, fehlt im neuen Snapshot und bleibt in den vorigen Commits lesbar; die Dateien liegen unter board/<id>.md, damit ein versehentlicher Merge dieses Branches nie mit .jaira/tickets/ kollidiert; zwei gleichzeitige Snapshots loesen per --force-with-lease auf statt sich zu ueberschreiben; der Branchname ist konfigurierbar und der Befehl sagt, was er hinzugefuegt und entfernt hat"
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-10T20:13:52Z
updated-at: 2026-09-10T20:13:52Z
---

# Ein Snapshot-Branch traegt die Tafel als Dateien, ohne dass jemand ihn auscheckt

## Definition of Done

- [ ] 'jaira snapshot' baut den Branch aus dem aktuellen Satz refs/jaira/tickets/* per hash-object, mktree und commit-tree, ohne Checkout und ohne den Arbeitsbaum zu beruehren, und laesst sich in jedem Repo-Zustand aufrufen; der erste Commit ist elternlos, jeder weitere haengt am vorigen, sodass git log die Geschichte der Tafel ist; ein Ticket, dessen Ref verschwunden ist, fehlt im neuen Snapshot und bleibt in den vorigen Commits lesbar; die Dateien liegen unter board/<id>.md, damit ein versehentlicher Merge dieses Branches nie mit .jaira/tickets/ kollidiert; zwei gleichzeitige Snapshots loesen per --force-with-lease auf statt sich zu ueberschreiben; der Branchname ist konfigurierbar und der Befehl sagt, was er hinzugefuegt und entfernt hat

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

