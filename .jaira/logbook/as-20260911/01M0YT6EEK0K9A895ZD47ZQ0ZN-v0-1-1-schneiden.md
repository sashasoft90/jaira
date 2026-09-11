---
id: 01M0YT6EEK0K9A895ZD47ZQ0ZN
title: v0.1.1 schneiden
status: done
ready: true
creator: BeMuCa
goal: Alles nach 62989f1 ist als Release veroeffentlicht
context: "Seit dem letzten Release ist sehr viel gelandet: die Liste von sieben, drei fremde PRs, der board-bewusste Block. Schneiden heisst: ein Block '## 0.1.1' in core/release/NOTES.md, taggen, Tag pushen. In die Notizen gehoert die Verhaltensaenderung, nicht nur die Funktionsliste: ein handgeschriebenes '[-]' in einer Checkliste galt frueher als offen und gilt jetzt als zurueckgezogen, blockiert also den Abschluss nicht mehr."
definition-of-done: "core/release/NOTES.md hat einen 0.1.1-Block, der die [-]-Aenderung nennt; der Tag ist gepusht; jaira self upgrade --check findet das Release"
blocked-by: []
commits:
  - 81f277099dabe6a8f7299ae150ea5d09f3182002
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-26T10:35:02Z
updated-at: 2026-09-11T14:50:16Z
claimed-by: EE-3NX6GL3-2976569
claimed-at: 2026-08-31T19:05:43Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Tag setzen? Zwei Befehle: git tag v0.1.1 && git push origin v0.1.1 - goreleaser schneidet dann das Release, danach findet jaira self upgrade --check es und der letzte DoD-Punkt ist erfuellt. Der NOTES-Block ist schon auf master."
outcome-what: "core/release/NOTES.md hat den 0.1.1-Block: 13 Einzeiler in Anweisungs-Stimme, neueste zuerst, die [-]-Verhaltensaenderung als erste Zeile (DoD-Pflicht); committet und gepusht"
outcome-why: "Release-Schnitt war faellig - 130+ Commits seit v0.1.0 inkl. zweier Renames unreleased Surface; der Block ist die Haelfte des DoD, die ein Agent liefern kann, der Tag-Push veroeffentlicht und gehoert dir"
outcome-resolves: go test ./core/release -count=1 gruen (der NOTES-Parser liest den Block); Zeilenregel eingehalten (jede Aenderung genau eine Zeile)
---

# v0.1.1 schneiden

## Definition of Done

- [x] core/release/NOTES.md hat einen 0.1.1-Block, der die [-]-Aenderung nennt; der Tag ist gepusht; jaira self upgrade --check findet das Release
  proof: geprueft am 11.09.2026: NOTES.md hat den 0.1.1-Block mit 15 Zeilen, die [-]-Aenderung steht darin; die Tags v0.1.0, v0.1.1 und v0.1.2 liegen auf dem Remote des Originals; 'jaira self upgrade --check' antwortet '0.1.2 is already the latest release'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 14:50 · Alexander Sacharov** — Uebernommen, weil BeMuCa im Urlaub ist und das Ticket laengst erledigt war - nur nicht abgehakt. Alle drei Teile der DoD am 11.09.2026 gegen den Ist-Zustand geprueft, nicht aus dem Gedaechtnis: der Block steht in NOTES.md, die Tags liegen auf dem Remote, und self upgrade --check findet das Release (inzwischen sogar 0.1.2).

Ein Ticket, das in der HITL-Lane steht, obwohl die Arbeit draussen ist, sagt allen Falsches: es sieht aus, als waere das Release nicht erschienen.
