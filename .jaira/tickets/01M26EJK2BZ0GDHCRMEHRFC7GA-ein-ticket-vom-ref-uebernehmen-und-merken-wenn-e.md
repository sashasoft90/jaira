---
id: 01M26EJK2BZ0GDHCRMEHRFC7GA
title: "Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Wer ein Ticket uebernimmt, holt es mit einem Befehl vom Ref auf die eigene Platte, und genau ein Klon tut das - und ein Ticket, das jemand anders ins Logbuch gelegt hat, wird hier gemeldet statt doppelt zu existieren"
context: |-
  Zwei Luecken, die beim Durchsprechen von 8566KF (Frage 1) gefunden wurden.

  1. Ein Ticket, das nur auf einem Ref liegt, ist heute in 'jaira list' und auf dem Board keine Karte - nur ein Zaehler in der Hinweiszeile und eine Zeile in 'jaira fetch'. Wer es bearbeiten will, hat keinen Weg, es auf die eigene Platte zu holen: er muesste die Datei per 'git show refs/jaira/tickets/<id>:<id>.md' von Hand hinschreiben.

  Der Trigger dafuer soll NICHT der Fetch sein, sondern das Uebernehmen. Grund: eine Uebernahme ist ein Schreibvorgang aufs Ref und damit ein Compare-and-Swap - genau einer gewinnt, und genau der eine Klon legt die Datei hin. Materialisiert der Fetch, legen alle Klone dieselbe Datei an, und der add/add-Fall beim Merge wird vom Ausnahme- zum Normalfall.

  2. Gemessen und belegt: legt ein Klon das Ticket ins Logbuch (oder Archiv) und ein anderer hat es noch unter tickets/, existiert es nach dem Merge zweimal - fuer git ist das rename/modify, und die Standardaufloesung behaelt beide Pfade. Das wird heute NICHT erkannt: Store.idIndex in core/ticket/store.go vergleicht Ids nur innerhalb von tickets/, weil Paths() nur dieses Verzeichnis liest. logbook/ und archive/ sieht es nicht. 'zwei Dateien, eine Id' wird also nur auf der Tafel gefunden, nicht ueber die Tafelgrenze hinweg.

  Das Signal dafuer liegt schon da und wird nur nicht benutzt: logbook und archive loeschen das Ref, und 'fetch --prune' sieht dieses Verschwinden. 'mein Ref ist weg, die Datei liegt noch auf der Tafel' heisst: jemand hat das Ticket abgeraeumt.

  Nicht gewollt: still verschieben. Ein Ticket, das unter jemandem wegwandert, muss eine Meldung sein, keine Ueberraschung.
definition-of-done: "'jaira pull <id>' holt ein Ticket, das nur auf seinem Ref liegt, als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird nichts hingelegt und die Meldung nennt, wer schneller war; 'jaira fetch' bleibt reines Lesen und legt nie eine Datei an; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von 'jaira fetch' und 'jaira validate' als 'hat die Tafel verlassen' gemeldet, mit dem Befehl der es hier nachzieht, und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg und meldet sie wie die bestehende Dublette"
tags:
  - concurrency
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T20:01:34Z
updated-at: 2026-09-10T20:01:34Z
---

# Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat

## Definition of Done

- [ ] 'jaira pull <id>' holt ein Ticket, das nur auf seinem Ref liegt, als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird nichts hingelegt und die Meldung nennt, wer schneller war; 'jaira fetch' bleibt reines Lesen und legt nie eine Datei an; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von 'jaira fetch' und 'jaira validate' als 'hat die Tafel verlassen' gemeldet, mit dem Befehl der es hier nachzieht, und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg und meldet sie wie die bestehende Dublette

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

