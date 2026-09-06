---
id: 01M1Q59XWJ4PB6R2EWHDJJ32B4
title: Die vier Status-Schreibstellen werden eine Move-Funktion
status: done
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Gate-Check, Claim-Regel, Mutation und Settle leben in EINER core-Funktion; CLI move, TUI applyMove/forceMove und accept rufen sie und behalten nur UI-Belange (Formatierung, pending/confirm, --from-lane-Validierung)"
context: "Berk am 04.09. (Feedback L122, zweiter Teil): der bekannte 'eigene Schnitt' aus den Learnings - vier Schreibstellen (cli/flow.go, tui applyMove, forceMove, signoff accept) teilen heute moveMutation nur paarweise; jede neue Nach-Move-Regel braucht alle vier, und grep findet accept nicht (schreibt Status direkt). Der Settle-Teil wurde bereits zentralisiert (Vorgaenger-Ticket); dieser Schnitt vereinheitlicht den Rest. Gross und invariantennah (CLI und TUI muessen identisch gaten, Interactive-Flag, Claim-beim-Pull) - bewusst NICHT im Schnellverfahren gebaut."
definition-of-done: "Eine core-Move-Funktion traegt Gate+Mutate+Settle; alle vier Aufrufer delegieren; CLI- und TUI-Verhalten byte-gleich (bestehende Tests gruen); ein Test beweist, dass eine neue Nach-Move-Regel nur noch eine Stelle braucht"
tags: []
blocked-by: []
commits:
  - a6aede828b3407655da43a48e4fd3d8a34fc6707
  - 6144c2ed64fb274ea9b7d4939c2a6e7b05c0a7e7
  - 47b17371e9eaa82b633949c62b20db3148e0a600
created-at: 2026-09-04T21:30:57Z
updated-at: 2026-09-06T12:05:00Z
claimed-by: EE-3NX6GL3-1083590
claimed-at: 2026-09-06T11:51:38Z
updated-by: BeMuCa
outcome-what: "core/move.Move buendelt Gate-Check, Claim-on-Pull, Status-Schreiben und Settle; CLI-move, TUI applyMove/forceMove und accept delegieren und behalten nur Staging/Flags/Wortlaut/Pending-UX; moveMutation und der settleLane-Switch sind geloescht"
outcome-why: "Feedback L122 zweiter Teil: vier Schreibstellen hiessen vierfache Verdrahtung jeder Nach-Move-Regel - das Doorway-Review fand accept() nur durch Suchen; ab jetzt braucht eine neue Regel eine Stelle"
outcome-resolves: "5 core-Unit-Tests (Verweigerung schreibt nichts, Force+Overrode, Claim nie ueberschreiben, Stage ueberlebt Refusal, Doorway-Settle); alle vier Aufrufer migriert von zwei parallelen Agenten mit eigenen Gates; Byte-Gleichheit: volle Suite -race RC=0 mit EINER Namens-Testanpassung"
executed-by: fable+sonnet
review-summary: "Kritik: Schnittlinie richtig gezogen - ins core wandert, was identisch sein MUSS (Gate-Reihenfolge, Claim-Regel, Write, Settle), draussen bleibt, was legitim verschieden ist (CLI-Staging/Flags, TUI-Pending, accept-Wortlaut); dry-run bewusst beim direkten Gate-Check belassen (kein Write-Site). Kein Verhalten erfunden: Stage-persistiert-vor-Gate und Claim-im-Speicher sind die dokumentierten Alt-Semantiken beider Seiten, jetzt als benannte Request-Modi."
review-gaps: "Nichts entfernt daruber hinaus. Gelassen: forceMove baut seine Meldung weiter aus den beim Pending gemerkten Refusals (p.refusals) statt aus res.Overrode - gleicher Text, und ein zwischenzeitlich veraendertes Board wuerde res ohnehin frisch liefern; gate.Request wird an einer Stelle konstruiert, Question/Reason reisen als Felder durch."
review-check: |-
  1. Ein normaler Move, ein --force-Move, ein a-Accept im TUI: identisches Verhalten wie gestern (Meldungen, Doorway-Filing, Overrode-Reihenfolge).
  2. grep -rn "SetScalar(ticket.FieldStatus" internal/ -> nur noch core/move schreibt Status (plus merge, das kein Move ist).
  3. go test ./core/move -v zeigt die fuenf Vertrags-Tests.
  4. jaira move --dry-run schreibt weiter nichts.
test-verdict: "pass: volle Suite -race nach Cache-Loeschung RC=0 (test21, 16 Pakete inkl. core/move); DoD am Baum: grep zeigt null Status-SetScalar ausserhalb core/move (nur merge, kein Move); Verhalten exerziert: alle Moves DIESES Tickets (todo->in-progress->critique->optimize->testing) liefen bereits durch move.Move am neuen Binary"
review-verdict: "accept (Struktur offengelegt: core von Fable-Koordinator mit 5 Vertrags-Tests, Migrationen von zwei parallelen Sonnet-Agenten mit je eigenen verifizierten Gates, Integration + volle Suite -race RC=0 beim Koordinator; damit haben beide Migrationsdiffs ein vom Implementierer verschiedenes Modell gesehen. Staerkster Beweis ist der Betrieb: jeder Board-Move seit dem Commit - inklusive der Zuege dieses Tickets und des Doorway-Filings am Ende - laeuft durch die eine Funktion.)"
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
