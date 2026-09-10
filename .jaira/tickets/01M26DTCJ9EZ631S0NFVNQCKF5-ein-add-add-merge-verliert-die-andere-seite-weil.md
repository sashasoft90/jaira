---
id: 01M26DTCJ9EZ631S0NFVNQCKF5
title: "Ein add/add-Merge verliert die andere Seite, weil die leere Basis abgelehnt wird"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Zwei Branches, die dieselbe Ticketdatei anlegen, mergen feldweise wie zwei Branches, die sie aendern - ohne dass eine Seite still verschwindet"
context: |-
  Belegt, nicht vermutet: git ruft den Merge-Driver bei einem add/add-Konflikt mit einer LEEREN Basisdatei auf. core/merge.Merge parst die Basis mit ticket.ParseDoc, das lehnt sie ab ('file does not begin with ---'), der Driver beendet mit 1, git markiert den Pfad als konfliktbehaftet und laesst 'ours' stehen.

  Der Schaden ist still: in der Datei stehen KEINE Konfliktmarker. Es sieht aus wie ein erfolgreicher Merge, aber Lane und updated-by der anderen Seite sind weg.

  Reproduktion (gemessen mit dem gebauten Binary): geteiltes Board, Ticket auf Branch ada anlegen und committen, dieselbe Datei auf Branch berk anlegen mit status: todo und updated-by: berk, dann 'git merge berk'. Ergebnis: 'CONFLICT (add/add)', 0 Konfliktmarker, status: backlog - berks Seite verloren. Direkt nachgestellt mit 'jaira merge-driver <leere-datei> ours theirs' - exit 1 mit genau dieser Meldung.

  Heute ist der Fall schwer zu treffen, weil zwei Tickets verschiedene Ids und damit verschiedene Dateinamen haben. Er wird zum Normalfall, sobald ein per Ref empfangenes Ticket lokal als Datei landet (Frage 1 an 8566KF): dann existiert dieselbe Datei auf zwei Branches. Der Fix gehoert trotzdem hierher und nicht dorthin - der Bug besteht auch ohne das.
definition-of-done: "core/merge.Merge behandelt eine leere oder nur aus Leerraum bestehende Basis als 'kein gemeinsamer Vorfahre' statt sie abzulehnen; ein add/add-Merge zweier Branches an derselben Ticketdatei laeuft ohne Konfliktmarker durch, behaelt die weitere Lane und vereinigt Listen beider Seiten; Prosafelder, die sich auf beiden Seiten unterscheiden, bleiben ein echter Konflikt und landen in conflict-theirs-<field>; ein Test gegen echtes git deckt genau den add/add-Fall ab"
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-10T19:48:21Z
updated-at: 2026-09-10T19:51:22Z
claimed-by: DESKTOP-RFTCH11-353175
claimed-at: 2026-09-10T19:48:30Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
---

# Ein add/add-Merge verliert die andere Seite, weil die leere Basis abgelehnt wird

## Definition of Done

- [x] core/merge.Merge behandelt eine leere oder nur aus Leerraum bestehende Basis als 'kein gemeinsamer Vorfahre' statt sie abzulehnen; ein add/add-Merge zweier Branches an derselben Ticketdatei laeuft ohne Konfliktmarker durch, behaelt die weitere Lane und vereinigt Listen beider Seiten; Prosafelder, die sich auf beiden Seiten unterscheiden, bleiben ein echter Konflikt und landen in conflict-theirs-<field>; ein Test gegen echtes git deckt genau den add/add-Fall ab
  proof: core/merge/merge.go: leere oder nur aus Leerraum bestehende Basis wird zu emptyDoc ('---\n---\n') statt abgelehnt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-10 19:51 · Alexander Sacharov** — Ursache genau: git ruft den Merge-Driver beim add/add-Konflikt mit einer leeren %O-Datei auf. ParseDoc lehnt sie ab, der Driver endet mit 1, git laesst 'ours' liegen und markiert den Pfad - ohne Konfliktmarker in der Datei. Das ist der stille Teil: es liest sich wie ein sauberer Merge.

Fix: eine leere Basis ist ein ABWESENDER Vorfahre, kein kaputtes File. Sie wird durch ein leeres Ticket ('---\n---\n') ersetzt, das kleinste, das ParseDoc annimmt. Damit greifen alle bestehenden Regeln so, wie sie es fuer ein Feld tun, das keine Seite geerbt hat: weitere Lane gewinnt, Listen vereinigen, Prosa mit zwei verschiedenen Saetzen bleibt ein echter Konflikt. Keine neue Regel, keine Sonderbehandlung im Driver.

Warum nicht im Driver abfangen: der Driver ist nur der Aufrufer. Die Frage 'was heisst keine gemeinsame Basis' ist eine Merge-Regel und gehoert zu den anderen, sonst hat der zweite Aufrufer (das Board via refsync.Reconcile, das denselben Fall hat, wenn ein Ref noch keinen Parent-Commit hat) dieselbe Luecke wieder.

Handprobe nach dem Fix, gemessen: 'git merge berk' laeuft durch (1 file changed), 0 Konfliktmarker, status: todo obwohl adas Zeitstempel neuer war, tags: cli + concurrency. updated-by bleibt ada - das ist die dokumentierte Regel fuer uebrige Skalare (neueres updated-at gewinnt), nicht ein Verlust.

Damit ist Option A aus Frage 1 an 8566KF (empfangenes Ticket landet als Datei) nicht mehr durch diesen Bug blockiert.
