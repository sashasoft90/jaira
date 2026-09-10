---
id: 01M21PK85SPRS8BHF7MGZP3R48
title: "Es steht fest und steht geschrieben, wann jaira ein Release schneidet"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine Regel, wann getaggt wird, steht an einer Stelle im Repo, und der naechste Release-Schnitt folgt ihr"
context: |-
  Berk am 08.09.: 'wann wird released, haben wir da eine Regel?'. Antwort heute: nein.
  Gesucht und nicht gefunden: keine CONTRIBUTING.md im Repo; kein Abschnitt in README.md oder docs/ nennt Kadenz, Rhythmus oder ein Kriterium; core/release/NOTES.md regelt nur das FORMAT der Notizen, nicht den Zeitpunkt.
  Was es gibt: .github/workflows/release.yaml loest auf 'push tags v*' aus und laesst goreleaser die sechs Binaries bauen. Der Ausloeser ist also ein Tag von Hand, und wann der gesetzt wird, entscheidet bisher der Bauch.
  Bisherige Praxis als einzige Datenlage: v0.1.0 am 17.08.2026, v0.1.1 am 07.09.2026 - drei Wochen Abstand, 130+ Commits dazwischen.
  Zu entscheiden ist genau eines: nach Zeit (z.B. alle zwei Wochen), nach Inhalt (z.B. sobald NOTES.md einen ungeschnittenen Block hat), oder auf Zuruf.
  Danach hinschreiben, wo ein Fremder es findet - README oder eine neue CONTRIBUTING.md, nicht in einem Ticket.
definition-of-done: Die Release-Regel steht in einer versionierten Datei im Repo (README oder CONTRIBUTING); sie nennt den Ausloeser und wer taggt; das Ticket verlinkt die Datei
tags:
  - release
blocked-by: []
commits: []
created-at: 2026-09-08T23:45:32Z
updated-at: 2026-09-08T23:45:32Z
---

# Es steht fest und steht geschrieben, wann jaira ein Release schneidet

## Definition of Done

- [ ] Die Release-Regel steht in einer versionierten Datei im Repo (README oder CONTRIBUTING); sie nennt den Ausloeser und wer taggt; das Ticket verlinkt die Datei

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

