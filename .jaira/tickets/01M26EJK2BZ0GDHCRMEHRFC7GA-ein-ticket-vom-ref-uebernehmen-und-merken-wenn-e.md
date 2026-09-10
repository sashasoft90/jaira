---
id: 01M26EJK2BZ0GDHCRMEHRFC7GA
title: "Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Ein Ticket lebt bis zur Uebernahme nur auf seinem Ref, und erst 'jaira pull' legt es hier als Datei hin - damit existiert es zu jeder Zeit in genau einem Klon und ein Merge kann es nicht doppeln"
context: |-
  Die Konflikte, die beim Durchsprechen von 8566KF gefunden wurden, sollen nicht aufgeloest, sondern unmoeglich gemacht werden. Beides ist gemessen, nicht vermutet.

  Was heute passiert:

  1. Ein Ticket, das nur auf einem Ref liegt, ist in 'jaira list' und auf dem Board keine Karte, und es gibt keinen Befehl, es auf die eigene Platte zu holen - man muesste 'git show refs/jaira/tickets/<id>:<id>.md' von Hand hinschreiben.
  2. Legt ein Klon das Ticket ins Logbuch, waehrend ein anderer es noch unter tickets/ hat, existiert es nach dem Merge zweimal (fuer git ist das rename/modify). Erkannt wird das nicht: Store.idIndex in core/ticket/store.go vergleicht Ids nur innerhalb von tickets/, weil Paths() nur dieses Verzeichnis liest.

  Der Entwurf, auf den wir gekommen sind: solange niemand das Ticket bearbeitet, ist das Ref der Speicherort und die Datei existiert nirgends. Das Uebernehmen legt sie hin - und weil ein Uebernehmen ein Schreibvorgang aufs Ref und damit ein Compare-and-Swap ist, gewinnt genau einer. Es bleibt kein Ehrenwort ('jeder zieht es nur einmal'), sondern die Mechanik erzwingt es. Ein zweiter Puller wird abgelehnt und erfaehrt, wer schneller war.

  Damit gibt es zu jeder Zeit genau eine Ticketdatei. add/add und die Logbuch-Dublette treten nicht mehr auf, statt behandelt zu werden.

  Drei Preise, ausdruecklich abgewogen und nicht uebersehen:

  1. 'Klonen und die Tafel sehen' gilt nicht mehr wortwoertlich. Der Default-Refspec zieht die Refs nicht, also wird aus 'clone, dann jaira' ein 'clone, dann jaira fetch'. Die README-Zeile muss mitgeaendert werden statt still falsch zu werden.
  2. Der unbearbeitete Backlog liegt nur auf dem Remote. Abgefedert durch die Regel 'die Datei wird bei pull und beim Verlassen der Tafel committet' - alles, woran gearbeitet wurde, und alles Abgeschlossene steht damit in der Historie. Fuer den unberuehrten Backlog ist der Snapshot-Branch aus PTQ3XT das Backup.
  3. Zwei Modi: ein Board mit Remote hat das Ref als Speicherort, ein Board ohne Remote (jaira init gitignored .jaira/) die Datei. Diese Verzweigung muss an einer Stelle stehen und benannt sein, sonst wird sie ueberall halb nachgebaut.
definition-of-done: "'jaira pull <id>' holt ein Ticket, das nur auf seinem Ref liegt, als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird nichts hingelegt und die Meldung nennt, wer schneller war; 'jaira fetch' bleibt reines Lesen und legt nie eine Datei an; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von 'jaira fetch' und 'jaira validate' als 'hat die Tafel verlassen' gemeldet, mit dem Befehl der es hier nachzieht, und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg und meldet sie wie die bestehende Dublette"
tags:
  - concurrency
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T20:01:34Z
updated-at: 2026-09-10T20:14:42Z
updated-by: Alexander Sacharov
---

# Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat

## Definition of Done

- [ ] 'jaira pull <id>' holt ein Ticket vom Ref als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird keine Datei angelegt und die Meldung nennt, wer schneller war; 'jaira create' legt auf einem Board mit Remote keine lokale Datei mehr an, sondern nur das Ref, und sagt das; 'jaira fetch' bleibt reines Lesen; ein Board ohne Remote verhaelt sich unveraendert wie heute, und diese Verzweigung steht an genau einer Stelle im Code; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von fetch und validate als 'hat die Tafel verlassen' gemeldet und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg; die README-Zeile ueber 'clone and see the same board' ist auf den neuen Ablauf korrigiert

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-10 20:14 · Alexander Sacharov** — Reihenfolge, in der das gebaut werden muss: erst pull (dieses Ticket), dann darf create aufhoeren, lokal zu schreiben. Umgekehrt gaebe es einen Zustand, in dem ein Ticket auf einem Ref liegt und niemand es holen kann.

PTQ3XT (Snapshot-Branch) ist die Abfederung von Preis 2 und kann parallel laufen - es haengt nicht an diesem Ticket, aber dieses Ticket sollte nicht ohne es in Produktion gehen, sonst liegt der Backlog eine Zeit lang ohne Backup nur auf dem Remote.
