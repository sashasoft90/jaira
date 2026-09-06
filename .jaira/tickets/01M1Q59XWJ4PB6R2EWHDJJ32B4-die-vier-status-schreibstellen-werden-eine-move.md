---
id: 01M1Q59XWJ4PB6R2EWHDJJ32B4
title: Die vier Status-Schreibstellen werden eine Move-Funktion
status: in-progress
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Gate-Check, Claim-Regel, Mutation und Settle leben in EINER core-Funktion; CLI move, TUI applyMove/forceMove und accept rufen sie und behalten nur UI-Belange (Formatierung, pending/confirm, --from-lane-Validierung)"
context: "Berk am 04.09. (Feedback L122, zweiter Teil): der bekannte 'eigene Schnitt' aus den Learnings - vier Schreibstellen (cli/flow.go, tui applyMove, forceMove, signoff accept) teilen heute moveMutation nur paarweise; jede neue Nach-Move-Regel braucht alle vier, und grep findet accept nicht (schreibt Status direkt). Der Settle-Teil wurde bereits zentralisiert (Vorgaenger-Ticket); dieser Schnitt vereinheitlicht den Rest. Gross und invariantennah (CLI und TUI muessen identisch gaten, Interactive-Flag, Claim-beim-Pull) - bewusst NICHT im Schnellverfahren gebaut."
definition-of-done: "Eine core-Move-Funktion traegt Gate+Mutate+Settle; alle vier Aufrufer delegieren; CLI- und TUI-Verhalten byte-gleich (bestehende Tests gruen); ein Test beweist, dass eine neue Nach-Move-Regel nur noch eine Stelle braucht"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:57Z
updated-at: 2026-09-06T12:03:11Z
claimed-by: EE-3NX6GL3-1083590
claimed-at: 2026-09-06T11:51:38Z
updated-by: BeMuCa
---

# Die vier Status-Schreibstellen werden eine Move-Funktion

## Definition of Done

- [x] Eine core-Move-Funktion traegt Gate+Mutate+Settle; alle vier Aufrufer delegieren; CLI- und TUI-Verhalten byte-gleich (bestehende Tests gruen); ein Test beweist, dass eine neue Nach-Move-Regel nur noch eine Stelle braucht
  proof: core/move.Move traegt Gate+Claim+Mutate+Settle (move.go, 5 Unit-Tests); alle vier Aufrufer delegieren: flow.go (Stage persistiert, Reload, dry-run bleibt gate-only), tui applyMove/forceMove (ClaimOnPull, Force), signoff accept (Interactive); moveMutation + settleLane-Switch geloescht; Byte-Gleichheit: volle Suite -race RC=0 (test21, 16 Pakete) - einzige Testanpassung trimholds_test (pinnte den geloeschten Helper beim Namen); der 'neue Regel = eine Stelle'-Beweis sind die move_test-Faelle: Settle/Gate/Claim feuern aus EINEM Aufruf

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-06 11:54 · BeMuCa** — Bauweise: Koordinator schrieb core/move selbst (Semantik frisch: Stage persistiert VOR dem Gate wie der CLI-Move - Felder ueberleben Verweigerung -, ClaimOnPull staged im Speicher wie die TUI, Force liefert Overrode zurueck, Settle laeuft im selben Zug); Unit-Tests pinnen Verweigerung-schreibt-nichts, Force, Claim-nie-ueberschreiben, Stage-ueberlebt-Refusal, Doorway-Settle (Commit 6144c2e). Migration der vier Aufrufer parallel durch zwei Sonnet-Agenten (CLI: flow.go; TUI: applyMove/forceMove/accept), beide UNCOMMITTED - Integration und serielle Commits beim Koordinator (Index-Swallow-Lektion). dry-run bleibt bewusst beim direkten gate.CheckAdvance: kein Write-Site.
