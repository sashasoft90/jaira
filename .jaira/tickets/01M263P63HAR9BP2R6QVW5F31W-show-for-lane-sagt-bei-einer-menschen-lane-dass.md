---
id: 01M263P63HAR9BP2R6QVW5F31W
title: "show --for-lane sagt bei einer Menschen-Lane, dass ein Agent dort nicht herauskommt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Agent erfaehrt VOR der Arbeit, dass die Ziel-Lane ein menschlicher Checkpoint ist - nicht erst, wenn der Move verweigert wird"
context: |-
  Nachgemessen am 10.09.: 'jaira show <id> --for-lane review --json' liefert praktisch ein leeres Objekt - prompt leer, produces null, model_tier leer, missing null, und complete: true.
  Das liest sich als 'hier ist nichts zu tun', nicht als 'hier darfst du hinein, aber nicht heraus'.
  Erst beim Move kommt die Wahrheit, dann aber klar: 'lane "review" is a human checkpoint: open the board and sign this off, or accept it there. An agent cannot move a ticket out of it' - EXIT=3.
  Folge, aus einer Multi-Agent-Session am 10.09.: ein Agent, der ein Ticket 'bis done' durchtreiben soll, plant gegen eine Wand und merkt es erst nach getaner Arbeit.
  Die Information liegt vor: review.md traegt 'agentic: false' und 'requires-human-exit: true'. Sie wird nur nicht ausgeliefert.
definition-of-done: "'show --for-lane' meldet fuer eine Lane mit requires-human-exit ausdruecklich, dass ein Agent sie nicht verlassen kann, im JSON als eigenes Feld und im Text; complete:true steht nicht mehr allein fuer eine Lane ohne Vertrag; ein Test deckt es ab; go test ./... -race gruen"
tags:
  - gates
blocked-by: []
commits: []
created-at: 2026-09-10T16:51:18Z
updated-at: 2026-09-10T16:51:18Z
---

# show --for-lane sagt bei einer Menschen-Lane, dass ein Agent dort nicht herauskommt

## Definition of Done

- [ ] 'show --for-lane' meldet fuer eine Lane mit requires-human-exit ausdruecklich, dass ein Agent sie nicht verlassen kann, im JSON als eigenes Feld und im Text; complete:true steht nicht mehr allein fuer eine Lane ohne Vertrag; ein Test deckt es ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

