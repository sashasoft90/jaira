---
id: 01M1KCKH0MCQ5P1BKTSHP1AE82
title: Die jaira-Version steht links oben im Projektfenster
status: critique
ready: true
creator: BeMuCa
goal: "Das TUI zeigt die laufende Binary-Version sichtbar im Board, links oben in der Ecke des Projektfensters"
context: "Berk am 03.09.: beim Arbeiten mit mehreren Binary-Staenden (self upgrade, go build, ~/.local/bin) ist nicht sichtbar, welche Version laeuft. Wunsch: Version ganz links oben in der Ecke des Projektfensters. ACHTUNG Kollision: DNAEPN (done, 03.09.) hat gerade entschieden, dass die bestehende Versionszeile (internal/tui/updatecheck.go, gezeichnet in home.go und view.go/Fusszeile) auf dev-Builds SCHWEIGT, weil 'jaira dev' nichts beantwortet; auf Releases zeigt sie 'jaira 0.1.1 - up to date'. Vor dem Bauen klaeren: Platzierung links oben zusaetzlich oder statt der Fusszeile, und ob ein dev-Build 'dev' zeigen soll (widersprueche DNAEPN) oder weiter nichts."
definition-of-done: "Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T10:21:34Z
updated-at: 2026-09-09T07:53:16Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4099823
claimed-at: 2026-09-08T23:46:46Z
outcome-what: "Die Versionszeile ist aus beiden Fusszeilen nach links oben gewandert, auf eine eigene erste Zeile: im Board vor der Kopfzeile mit dem Ticket-Zaehler (internal/tui/view.go renderBoard, neues versionHead), im Launcher ueber dem Wordmark (internal/tui/home.go render). statusBar zeichnet sie nicht mehr. versionLine() gibt auf einem dev-Build 'jaira dev' zurueck statt \"\", weiter vor selfupdate.PollCache. Ihr Doc-Kommentar ist nachgezogen. Die Hoehenrechnung in renderBoard misst den Kopfblock statt eine Zeile anzunehmen."
outcome-why: "Bei mehreren Binary-Staenden (self upgrade, go build, ~/.local/bin) war nicht sichtbar, welche Version laeuft: die Zeile stand unter einer umbrechenden Hinweisleiste, und auf einem dev-Build - also bei jedem Contributor - stand sie gar nicht da. Links oben ist sie eine Identitaetsangabe und keine Upgrade-Empfehlung mehr, deshalb darf ein dev-Build sich dort nennen."
outcome-resolves: "Das Board zeigt links oben die Version, ein dev-Build zeigt 'dev', und die Platzierung ist als Geometrie getestet (erste Zeile, linksbuendig, Kopfzeile mit Ticket-Zaehler darunter) - nicht nur als 'irgendwo in der Ausgabe'."
---

# Die jaira-Version steht links oben im Projektfenster

## Definition of Done

- [x] Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab
  proof: internal/tui/view.go:183 (eigene erste Zeile) + :771 versionHead; internal/tui/updatecheck.go:59 dev zeigt 'jaira dev'; internal/tui/home.go:342; Tests TestBoardVersionSitsInTheTopLeftCorner, TestBoardHeadNamesADevBuild, TestBoardHeadCarriesTheVersionIndicator, TestHomeHeadNamesADevBuild; go test ./... -race gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] versionLine(): dev-Zweig gibt 'jaira dev' statt "" zurueck, weiter vor PollCache
- [x] Doc-Kommentar von versionLine() nachziehen: Fusszeile -> links oben, dev-Absatz umschreiben, DNAEPN-Teilumkehr benennen
- [x] Board: eigene erste Zeile links oben in renderBoard, Kopfblock-Hoehe in bodyHeight einrechnen
- [x] Board: versionLine aus statusBar entfernen
- [x] Launcher: versionLine aus der Fusszeile in Home.render entfernen und links oben setzen
- [x] Tests nachziehen: die vier Fusszeilen-Zusicherungen in devfooter_test.go und updatecheck_test.go
- [x] Neuer Test deckt die Platzierung ab: Version in der ersten Zeile links, Fusszeile ohne Version, dev zeigt 'dev'
- [x] go build -o /tmp/jaira-p1ae82 ./cmd/jaira und go test ./... -race gruen

## Progress
- **2026-09-08 23:45 · BeMuCa** — Berk hat am 08.09. die beiden offenen Fragen aus dem Kontext entschieden.
1. Platzierung: ERSETZEN, nicht zusaetzlich. Die Versionszeile verschwindet aus der Fusszeile und steht nur noch links oben im Projektkopf. Damit loest sich die notierte Kollision mit DNAEPN auf - es gibt danach nur noch eine Stelle, die eine Version behauptet.
2. Dev-Build: zeigt 'dev'. Das steht so schon im DoD und geht bewusst gegen DNAEPN, das die Fusszeile auf dev-Builds schweigen liess. Begruendung: links oben ist die Zeile eine Identitaetsangabe ('welche Binary laeuft hier'), nicht mehr eine Upgrade-Empfehlung - und genau die Frage hatte Berk beim Wechsel zwischen self upgrade, go build und ~/.local/bin.
Folgeticket T8A8KM haengt daran und setzt spaeter eine Pille '^ <version>' unter diese Zeile. Die Platzierung hier muss also Platz fuer eine zweite Zeile darunter lassen.
- **2026-09-09 07:36 · BeMuCa** — pre-process, was ich verifiziert und entschieden habe.

Zwei Zeichenstellen, nicht eine. versionLine() (internal/tui/updatecheck.go:35) hat heute zwei Aufrufer: Model.versionLine, gezeichnet in statusBar (internal/tui/view.go:829) = Board-Fusszeile, und Home.versionLine, gezeichnet in Home.render (internal/tui/home.go:394) = Launcher-Fusszeile. 'Projektkopf' meint nur das Board. Die Regel 'in der Fusszeile steht keine Version mehr' trifft aber beide Fusszeilen. Entscheidung: die Zeile wandert in BEIDEN Screens von unten nach links oben, damit der Launcher die Version nicht verliert und beide Screens dieselbe Stelle benutzen. Wenn Berk den Launcher lieber ganz ohne Version will, ist das eine Zeile weniger - bewusst so gebaut, nicht uebersehen.

Platzierung: eigene Zeile GANZ oben, nicht links in die Kopfzeile mit '13 tickets'. Grund ist T8A8KM: die Pille steht 'unter der Versionszeile'. Saesse die Version links in der Kopfzeile, waere die Zeile darunter schon die Board-Leiste ('1 requirementsgenie | 2 jAIra') - kein Platz. Als eigene erste Zeile kann T8A8KM eine zweite Zeile einfach anhaengen. Deshalb zaehle ich die Hoehe des Kopfblocks (strings.Count+1) statt eine 1 zu addieren: eine zweite Zeile rechnet sich dann von selbst mit, ohne dass T8A8KM die Hoehenrechnung anfassen muss.

Der dev-Fall bricht DNAEPN nur halb. DNAEPN hatte zwei Gruende: (a) 'jaira dev' sagt nichts Neues, (b) release.Current='dev' mit einer veroeffentlichten Version verglichen ergibt immer 'unterschiedlich', die Zeile wuerde also ein Upgrade auf AELTEREN Code bewerben und auf ein Kommando zeigen, das mit dev_build ablehnt. P1AE82 kippt nur (a). (b) bleibt: der dev-Zweig kehrt weiter VOR selfupdate.PollCache um und gibt nur 'jaira dev' zurueck - kein 'up to date', keine Pille, und kein abgekoppelter Check-Kindprozess fuer Contributor.

Keine Golden-/Snapshot-Tests in internal/tui - kein testdata-Verzeichnis, kein Treffer auf 'golden'. Was bricht, sind vier gezielte Zusicherungen in internal/tui/devfooter_test.go und internal/tui/updatecheck_test.go, die die Version genau in der Fusszeile erwarten.
- **2026-09-09 07:47 · BeMuCa** — in-progress, umgesetzt und am gerenderten Bild geprueft.

So sieht der Board-Kopf jetzt aus (stripANSI-Dump aus einem Wegwerf-Test, 120x24):
  jaira 0.1.1 - 0.1.2 available - run: jaira self upgrade
  001                                                        6 tickets
  +--- Backlog ---+ +--- Brainstorm ---+ ...
Fusszeile: '! 1 enter open . v compact . ... . q quit' - keine Version mehr.
Der Launcher hat die Zeile ueber dem Wordmark, ebenfalls links oben, Fusszeile dort ebenfalls ohne Version.

Was ich pruefen musste, statt es anzunehmen: internal/tui hat KEINE Golden-Tests (kein testdata/, kein 'golden'-Treffer). Was brach, waren vier handgeschriebene Zusicherungen, die die Version genau in der Fusszeile erwarteten. Statt sie stumpf umzustellen habe ich sie umgedreht: jeder Test behauptet jetzt beides - oben steht sie, unten nicht - damit ein Rueckfall nach unten auffliegt und nicht nur ein Verschwinden.

devfooter_test.go heisst jetzt devbuild_test.go (git mv). Der alte Name behauptete eine Fusszeile, die es fuer diese Zeile nicht mehr gibt.

Die Hoehenrechnung in renderBoard zaehlt den Kopfblock (strings.Count(head,'\\n')+1) statt eine 1 abzuziehen. Absicht: T8A8KM haengt eine zweite Zeile an, und dann stimmt bodyHeight von selbst. Wer dort '- 1' hinschreibt, schiebt die Spalten unter den Bildschirmrand.

Sackgasse, die ich verworfen habe: die Version links in die Kopfzeile selbst (dort, wo rechts '6 tickets' steht). Kostet keine Zeile und ist auch 'links oben' - aber die Zeile DARUNTER ist dann schon die Board-Leiste, und T8A8KM will die Pille darunter. Haette T8A8KM gezwungen, erst eine Zeile einzuschieben.

Ausserhalb des Scope aufgefallen, NICHT angefasst: internal/cli/update.go nudgeIfStale haelt die CLI bewusst still ueber Releases - unberuehrt, die Aenderung betrifft nur das TUI.
