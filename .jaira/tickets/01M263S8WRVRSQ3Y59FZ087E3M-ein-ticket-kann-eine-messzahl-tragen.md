---
id: 01M263S8WRVRSQ3Y59FZ087E3M
title: Ein Ticket kann eine Messzahl tragen
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Der zahlenmaessige Ertrag eines Tickets steht in einem eigenen Feld statt in Freitext, ist in 'jaira list' sichtbar und ueber Tickets hinweg vergleichbar"
context: |-
  Aus einer Multi-Agent-Session am 10.09.: der eigentliche Ertrag eines Tickets ist oft eine Zahl - Ausschussrate 36 Prozent auf 0, Benchmark-Note 0.985, Trefferquote 16 von 17.
  Heute leben die alle in Freitext-Notizen und in outcome-resolves. Damit sind sie nicht vergleichbar, nicht filterbar und ueber Tickets hinweg nicht auswertbar. Wer in drei Monaten fragt 'wurde die Trefferquote besser', muss Prosa lesen.
  Im Schema gibt es kein solches Feld - geprueft in core/ticket/schema.go: es gibt goal, context, definition-of-done, outcome-what/why/resolves, review-*, aber nichts Numerisches.
  Vorgeschlagen wurde 'jaira metric <id> <name> <wert>'.
  Vor dem Bauen gegen das Mass des Projekts halten: 'is this smaller than paca?'. Eine Messzahl ist ein neues Feld im Format, und das Format ist die API. Zu klaeren: ein Wert oder mehrere pro Ticket, Einheit mitfuehren oder nicht, und ob 'jaira list' es zeigt oder nur 'show' und --json.
definition-of-done: "Ein Ticket kann eine benannte Messzahl tragen; sie steht im Frontmatter, ueberlebt einen Round-Trip ohne andere Felder anzufassen, und ist in --json auslesbar; die Entscheidung ueber Anzahl und Einheit steht begruendet im Ticket; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:59Z
updated-at: 2026-09-10T16:52:59Z
---

# Ein Ticket kann eine Messzahl tragen

## Definition of Done

- [ ] Ein Ticket kann eine benannte Messzahl tragen; sie steht im Frontmatter, ueberlebt einen Round-Trip ohne andere Felder anzufassen, und ist in --json auslesbar; die Entscheidung ueber Anzahl und Einheit steht begruendet im Ticket; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

