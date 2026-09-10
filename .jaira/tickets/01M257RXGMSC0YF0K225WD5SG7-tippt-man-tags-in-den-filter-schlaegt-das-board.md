---
id: 01M257RXGMSC0YF0K225WD5SG7
title: "Tippt man tags: in den Filter, schlaegt das Board die Tags des Boards vor"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Nach 'tags:' im Filterfeld erscheint eine Liste aller Tags des Boards; hoch/runter waehlt einen aus, enter uebernimmt ihn, und das Board zeigt nur noch dessen Tickets"
context: |-
  Berk am 09.09.: '/tags: soll alle tags als vorschlag anzeigen, hoch runter klickend eins auswaehlen, dann werden nur die Tickets mit dem Tag angezeigt'.
  Was schon geht: das Filtern selbst. internal/tui/model.go:532 matches() versteht key:value, und 'tag:' wie 'tags:' sind beide belegt (model.go:555). Der Abgleich ist bewusst EXAKT, nicht Teilstring - 'tag:cur' soll nicht alles mit 'security' finden, und der Board-Filter soll dasselbe sagen wie 'jaira list --tag'. Diese Exaktheit ist der Grund, warum ein Vorschlag hier ueberhaupt noetig ist: wer den Namen nicht exakt trifft, bekommt eine leere Spalte statt einer Naeherung.
  Was fehlt: jede Form von Vorschlag. Eine Suche nach suggest, completion und autocomplete in internal/tui findet nichts - das Filterfeld ist reine Texteingabe.
  Die Tag-Liste selbst existiert schon: das Board kennt sie hinter 't', und die Farben stehen in .jaira/tags. Dieselbe Quelle benutzen, keine zweite Leseroutine bauen.
  Haengt an DBJTKQ: dort wird die Tag-Liste von einer eigenen Seite zu einem Popup. Wenn das zuerst gebaut ist, gibt es ein Overlay-Muster, an das sich der Vorschlag anhaengen kann, statt eines zweiten zu erfinden. Erst DBJTKQ, dann dieses.
  Zu klaeren beim Bauen: was passiert bei 'tag:' (Einzahl) - derselbe Vorschlag oder nur bei 'tags:'? Und was, wenn schon Zeichen hinter dem Doppelpunkt stehen: Liste filtern oder ausblenden?
definition-of-done: "Nach 'tags:' im Filterfeld steht eine Liste der Board-Tags; hoch/runter bewegt die Auswahl, enter setzt den Namen ein und das Board zeigt nur dessen Tickets; die Liste kommt aus derselben Quelle wie die Tag-Ansicht; ein Test deckt Anzeigen, Auswaehlen und das gefilterte Board ab; go test ./... -race gruen"
tags:
  - tui
blocked-by:
  - 01M21PJ9TQZGG6XTHCPADBJTKQ
commits: []
created-at: 2026-09-10T08:43:27Z
updated-at: 2026-09-10T08:43:28Z
updated-by: BeMuCa
---

# Tippt man tags: in den Filter, schlaegt das Board die Tags des Boards vor

## Definition of Done

- [ ] Nach 'tags:' im Filterfeld steht eine Liste der Board-Tags; hoch/runter bewegt die Auswahl, enter setzt den Namen ein und das Board zeigt nur dessen Tickets; die Liste kommt aus derselben Quelle wie die Tag-Ansicht; ein Test deckt Anzeigen, Auswaehlen und das gefilterte Board ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

