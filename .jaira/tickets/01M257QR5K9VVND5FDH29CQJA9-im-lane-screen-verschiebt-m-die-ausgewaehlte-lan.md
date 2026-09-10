---
id: 01M257QR5K9VVND5FDH29CQJA9
title: "Im Lane-Screen verschiebt m die ausgewaehlte Lane, und die Taste steht in der Fusszeile"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine ausgewaehlte Lane laesst sich mit m im Board verschieben, und wer den Screen oeffnet, sieht ohne Suchen, dass das geht"
context: |-
  Berk am 09.09.: 'man kann unter settings lanes adden, aber man soll ein ausgewaehltes mit m moven koennen, damit man entscheiden kann wo im board es sitzt'.
  WICHTIG - Verschieben gibt es schon. internal/tui/lanes.go:277-283 belegt H und L mit moveLane(-1) und moveLane(+1), das ruft lane.MoveLane und schreibt die Reihenfolge in .jaira/lanes/order.
  Das eigentliche Problem ist Sichtbarkeit: die Fusszeile des Lane-Screens (internal/tui/lanes.go:838) listet 'E edit, x remove, p publish, n new, R refresh, esc back'. H und L stehen NICHT drin. Die Funktion existiert und ist unauffindbar.
  Zwei Dinge sind also zu tun: die Taste auf m umlegen, weil das Board selbst m schon fuer 'move' benutzt (view.go:824, 'm move') und zwei Vokabeln fuer dieselbe Handlung eine zu viel sind - und sie in die Fusszeile aufnehmen.
  Zu klaeren beim Bauen: m allein kann nicht in zwei Richtungen schieben. Entweder m schaltet in einen Verschiebe-Modus, in dem h/l die Lane bewegen und enter/esc ihn beendet, oder m verschiebt nach rechts und M nach links. Das Board hat fuer 'm move' schon ein Muster - dort nachsehen und es uebernehmen statt ein zweites zu erfinden.
  H und L danach nicht stillschweigend loeschen: wer sie kennt, verliert sie sonst wortlos. Entweder beibehalten oder in der Meldung nennen.
definition-of-done: Der Lane-Screen verschiebt die ausgewaehlte Lane per m in beide Richtungen; die Fusszeile nennt die Taste; die neue Reihenfolge steht in .jaira/lanes/order; ein Test deckt Verschieben ueber m und die Fusszeile ab; go test ./... -race gruen
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-10T08:42:49Z
updated-at: 2026-09-10T08:42:49Z
---

# Im Lane-Screen verschiebt m die ausgewaehlte Lane, und die Taste steht in der Fusszeile

## Definition of Done

- [ ] Der Lane-Screen verschiebt die ausgewaehlte Lane per m in beide Richtungen; die Fusszeile nennt die Taste; die neue Reihenfolge steht in .jaira/lanes/order; ein Test deckt Verschieben ueber m und die Fusszeile ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

