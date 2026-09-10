---
id: 01M263PXVTAVWMADABE8PAAAJH
title: "Aus einer Menschen-Lane gibt es einen Rueckweg, wenn die Arbeit unvollstaendig ist"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Findet ein Nachpruefer in review, dass ein Ticket unvollstaendig ist, gibt es dafuer einen Weg zurueck statt nur einer Notiz"
context: |-
  Aus einer Multi-Agent-Session am 10.09.: NDM8SF liegt in review, ein benannter Teil des Auftrags wurde nie angefasst. Der Agent konnte nur eine Notiz anhaengen. Das Ticket steht weiter als abnahmebereit da, und Berk sieht die Notiz erst, wenn er es oeffnet.
  Nachgemessen: aus review kommt ein Agent in KEINE Richtung. Zurueck nach in-progress, weiter nach done, seitwaerts nach blocked - alle drei EXIT=3, immer dieselbe Meldung 'lane "review" is a human checkpoint ... An agent cannot move a ticket out of it'.
  review.md hat kein 'rejects-to'. Die anderen Lanes haben es: critique.md und optimize.md tragen 'rejects-to: in-progress'.
  Das ist kein Fehler in der Sperre - die soll genau so sein, und sie hat sich bewaehrt (ein Agent mit Freigabe hat trotzdem kein --force benutzt). Es fehlt der Kanal daneben.
  Zu klaeren beim Bauen: reicht ein sichtbares Kennzeichen am Ticket ('ein Nachpruefer haelt das fuer unvollstaendig'), das auf dem Board auffaellt, ohne die Lane zu wechseln? Oder soll ein rejects-to auch fuer eine Menschen-Lane gelten, das nur ein Mensch ausloest? Ersteres bricht die Sperre nicht auf und ist vermutlich der kleinere Eingriff.
definition-of-done: "Ein Agent kann an einem Ticket in einer Menschen-Lane vermerken, dass es unvollstaendig ist, ohne die Lane zu wechseln; das Kennzeichen ist auf dem Board sichtbar, nicht erst im geoeffneten Ticket; die Sperre der Menschen-Lane bleibt unangetastet; Tests decken es ab; go test ./... -race gruen"
tags:
  - gates
blocked-by: []
commits: []
created-at: 2026-09-10T16:51:42Z
updated-at: 2026-09-10T16:51:42Z
---

# Aus einer Menschen-Lane gibt es einen Rueckweg, wenn die Arbeit unvollstaendig ist

## Definition of Done

- [ ] Ein Agent kann an einem Ticket in einer Menschen-Lane vermerken, dass es unvollstaendig ist, ohne die Lane zu wechseln; das Kennzeichen ist auf dem Board sichtbar, nicht erst im geoeffneten Ticket; die Sperre der Menschen-Lane bleibt unangetastet; Tests decken es ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

