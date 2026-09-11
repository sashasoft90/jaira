---
id: 01M28KJ9PZ9GA3HGBJAY17TJ1H
title: "Das Board sendet, was es einreiht, und ein abgelegtes Ticket kommt nicht zurueck"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Wer im Board arbeitet, sieht dieselbe Tafel wie alle anderen: was das Board schreibt geht raus, und was hier abgelegt wurde taucht nicht als Ref-Karte wieder auf"
context: |-
  Zwei Fehler, heute im Betrieb aufgefallen, beide aus der Ref-Arbeit von 8566KF/RFC7GA.

  Beobachtet: Alexander hat fuenf Tickets im Board abgenommen. Die Dateien wanderten korrekt ins Logbuch - aber alle fuenf standen danach WEITER in der Spalte Human Review, mit dem Marker 'pull'. Also als Ref-Karten.

  Ursache eins: das Board flusht nie. refsync.Flush laeuft nur in cli.Execute nach einem Kommando. Das TUI reiht jede Schreibung in die Outbox ein und sendet sie nie - gemessen: fuenf Eintraege in der Outbox, und die Refs auf dem Remote trugen weiter 'status: signoff', obwohl lokal schon done und abgelegt. Wer den ganzen Tag im Board arbeitet, laesst die Tafel fuer alle anderen stehenbleiben.

  Ursache zwei: refsync.Extra liefert jedes Ref, fuer das es hier keine Datei unter tickets/ gibt. Ein ins Logbuch gelegtes Ticket hat genau das - keine Datei unter tickets/ - also kommt es als Ref-Karte zurueck auf die Tafel. Das Logbuch ist aber das Gegenteil davon, auf der Tafel zu sein.

  Die beiden zusammen ergeben das Bild aus dem Screenshot: abgenommen, abgelegt, und trotzdem noch da.
definition-of-done: "das Board sendet seine Warteschlange selbst: der Hintergrundlauf flusht vor dem Fetch, sodass eine im Board gemachte Aenderung ohne weiteres Kommando beim Team ankommt; refsync.Extra ueberspringt Ids, die hier im Logbuch oder Archiv liegen, sodass ein abgelegtes Ticket nicht als Ref-Karte zurueckkommt; Tests decken beides ab - eine Board-Schreibung erreicht das Remote ohne CLI-Aufruf, und ein abgelegtes Ticket verschwindet von der Tafel obwohl sein Ref noch steht"
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T16:07:16Z
updated-at: 2026-09-11T16:09:41Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Das Board sendet, was es einreiht, und ein abgelegtes Ticket kommt nicht zurueck

## Definition of Done

- [ ] das Board sendet seine Warteschlange selbst: der Hintergrundlauf flusht vor dem Fetch, sodass eine im Board gemachte Aenderung ohne weiteres Kommando beim Team ankommt; refsync.Extra ueberspringt Ids, die hier im Logbuch oder Archiv liegen, sodass ein abgelegtes Ticket nicht als Ref-Karte zurueckkommt; Tests decken beides ab - eine Board-Schreibung erreicht das Remote ohne CLI-Aufruf, und ein abgelegtes Ticket verschwindet von der Tafel obwohl sein Ref noch steht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

