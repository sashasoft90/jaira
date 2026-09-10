---
id: 01M26J3M1RN445R66PM0RA7PFE
title: "Ein Ref verschwindet erst, wenn das Ticket im Hauptbranch angekommen ist"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Zwischen 'ins Logbuch gelegt' und 'im Hauptbranch angekommen' bleibt das Ticket fuer alle sichtbar, damit niemand dasselbe Problem ein zweites Mal aufschreibt"
context: |-
  Heute loescht logbook und archive das Ref sofort. Die Ticketdatei liegt zu diesem Zeitpunkt aber nur im Branch der Person, die es abgeraeumt hat. In dem Fenster dazwischen - Branch fertig, noch nicht gemergt - sieht das Ticket NIEMAND sonst: das Ref ist weg, der Branch ist fremd. Wer dasselbe Problem bemerkt, schreibt es ein zweites Mal auf, mit derselben Loesung.

  Das Fenster ist nicht klein. Es ist so lang wie ein Review dauert.

  Die Regel soll sein: das Ref ueberlebt bis das Ticket im Hauptbranch liegt, und erst dann wird es endgueltig entfernt. Solange traegt es den Endzustand ('done') und ist damit fuer alle sichtbar als 'fertig, wartet aufs Landen' - genau die Information, die einen doppelten Ticket verhindert.

  Gemessen, nicht vermutet: 'git rev-list -1 <hauptbranch> -- ".jaira/logbook/*/<id>*" ".jaira/archive/<id>*"' antwortet leer, solange das Ticket nur im Branch liegt, und mit einem Sha, sobald der Branch gemergt ist. Kosten: nicht messbar (0,00 s), ein Aufruf, nur fuer Tickets im Endzustand.

  Auch geprueft und ein Stolperstein: 'git symbolic-ref refs/remotes/origin/HEAD' ist in einem frischen Klon nicht zwingend gesetzt (in meinem Testrepo fehlte es). Es braucht also einen Fallback auf origin/main und origin/master und einen Eintrag in ~/.jaira/settings.json.

  Der zweite Fall, der dadurch entsteht und mitbedacht werden muss: ein Branch, der nie gemergt wird. Dann bleibt das Ref fuer immer. Es braucht eine Meldung 'fertig, aber seit N Tagen nicht angekommen' und einen ausdruecklichen Weg, das Ref trotzdem loszuwerden - stilles Aufraeumen ist es nicht, denn das waere wieder das Fenster, nur unsichtbar.
definition-of-done: "logbook und archive loeschen das Ref nicht mehr sofort, sondern schreiben den Endzustand darauf, sodass andere Klone das Ticket als 'fertig, wartet aufs Landen' sehen; das Ref wird entfernt, sobald 'git rev-list -1 <hauptbranch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>*' einen Commit findet, und das laeuft im Hintergrund nach dem Muster der Update-Pruefung statt auf dem Kommandopfad; der Hauptbranch wird aus origin/HEAD bestimmt, mit Fallback auf origin/main und origin/master und einem Eintrag in settings.json; ein Ticket, das im Endzustand steht und seit einer konfigurierbaren Frist nicht im Hauptbranch angekommen ist, wird von fetch und validate gemeldet; es gibt einen ausdruecklichen Befehl, ein solches Ref trotzdem zu entfernen, und er sagt was er tut; Tests mit zwei Klonen belegen, dass das Ticket im Fenster fuer den anderen Klon sichtbar bleibt und nach dem Merge verschwindet"
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-10T21:03:18Z
updated-at: 2026-09-10T21:05:35Z
updated-by: Alexander Sacharov
---

# Ein Ref verschwindet erst, wenn das Ticket im Hauptbranch angekommen ist

## Definition of Done

- [ ] logbook und archive loeschen das Ref nicht mehr sofort, sondern schreiben den Endzustand darauf, sodass andere Klone das Ticket als 'fertig, wartet aufs Landen' sehen; das Ref wird entfernt, sobald 'git rev-list -1 <hauptbranch> -- .jaira/logbook/*/<id>* .jaira/archive/<id>*' einen Commit findet, und das laeuft im Hintergrund nach dem Muster der Update-Pruefung statt auf dem Kommandopfad; der Hauptbranch wird aus origin/HEAD bestimmt, mit Fallback auf origin/main und origin/master und einem Eintrag in settings.json; ein Ticket, das im Endzustand steht und seit einer konfigurierbaren Frist nicht im Hauptbranch angekommen ist, wird von fetch und validate gemeldet; es gibt einen ausdruecklichen Befehl, ein solches Ref trotzdem zu entfernen, und er sagt was er tut; Tests mit zwei Klonen belegen, dass das Ticket im Fenster fuer den anderen Klon sichtbar bleibt und nach dem Merge verschwindet
- [ ] die Landebranches stehen als LISTE in settings.json ('landing-branches', z.B. main, master, develop, release/*, Glob erlaubt) - damit erledigt sich die Frage, wie die wichtige Branch bei wem heisst; ohne Eintrag gilt origin/HEAD als einziger Eintrag, und loest sich auch der nicht auf, wird KEIN Ref entfernt und einmal gesagt, was einzutragen ist
- [ ] das Entfernen selbst gehoert in den Snapshot-Lauf (PTQ3XT), nicht in logbook/archive und nicht in einen eigenen Timer: der Snapshot geht periodisch ohnehin ueber alle Refs, und er hat das Ticket unmittelbar davor in den Snapshot-Branch geschrieben - im Moment des Loeschens liegt es also in zwei Ablagen. Reihenfolge im Lauf: erst schreiben, dann jaeten

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-10 21:03 · Alexander Sacharov** — Verhaeltnis zu den anderen Tickets, damit die Reihenfolge klar ist:

- RFC7GA (human) hat das Gegenteil eingebaut: logbook/archive loescht das Ref sofort. Dieses Ticket dreht das zurueck, und zwar aus einem Grund, der dort nicht bedacht war - das Fenster zwischen 'abgeraeumt' und 'gemergt' ist so lang wie ein Review, und darin ist das Ticket fuer alle anderen unsichtbar.
- refsync.Departed ('hat die Tafel woanders verlassen, liegt hier noch') wird dadurch seltener und richtiger: es feuert erst, wenn das Ref endgueltig weg ist, also nach dem Landen - genau dann, wenn der andere Klon wirklich aufraeumen soll, und nicht schon wenn jemand nur seinen Branch fertig hat.
- PTQ3XT (Snapshot alle drei Tage) ist unabhaengig, aber die Erkennung 'im Hauptbranch angekommen' und die Snapshot-Ausloesung teilen dasselbe Muster: im Hintergrund, nie auf dem Kommandopfad. Wer zuerst gebaut wird, legt das Muster fuer den anderen.
- **2026-09-10 21:04 · Alexander Sacharov** — Gemessen, welche Antwort auf 'welcher Branch ist der wichtige' ueberhaupt verfuegbar ist:

- Klon eines nicht leeren Repos: origin/HEAD ist gesetzt (refs/remotes/origin/master), lokal und kostenlos. Das ist der Normalfall.
- Klon eines LEEREN Repos (jaira init auf einem frischen Board): nicht gesetzt, und 'git remote set-head origin -a' scheitert mit 'Cannot determine remote HEAD' - auf dem Remote gibt es noch keinen Branch.
- Nach dem ersten Push: lokal weiter nicht gesetzt, aber 'git remote set-head origin -a' funktioniert dann und setzt es dauerhaft. Ein Netzaufruf, danach nie wieder.

Wichtiger als die Reihenfolge ist zweierlei:

1. settings.json gewinnt vor origin/HEAD, nicht umgekehrt. Der HEAD des Remotes ist der Default des Hosters, nicht zwingend der Branch, auf den es der Mannschaft ankommt: develop-Flows und Release-Branches sind normal. Das ist nicht zu erraten, aber zu erfragen.

2. Loest sich nichts auf, wird kein Ref entfernt. Ein stehengelassenes Ref kostet nichts; ein faelschlich entferntes nimmt genau die Sichtbarkeit weg, fuer die dieses Ticket existiert. Der unsichere Fall muss also auf die konservative Seite fallen, und einmal sagen, was einzutragen waere.
- **2026-09-10 21:05 · Alexander Sacharov** — Entschieden (Alexander), zwei Vereinfachungen, die beide besser sind als mein Entwurf:

1. Keine 'Hauptbranch', sondern eine Liste. 'landing-branches' in settings.json mit Globs. Damit ist die Frage 'wie heisst die wichtige Branch bei euch' nicht mehr zu beantworten, sondern zu konfigurieren, und develop- oder Release-Flows sind kein Sonderfall mehr. Fehlt der Eintrag, ist origin/HEAD der einzige Eintrag der Liste.

2. Gejaetet wird im Snapshot-Lauf, nicht in logbook/archive und nicht in einem eigenen Timer. Drei Gruende, der zweite ist der eigentliche:

   - Der Snapshot geht periodisch ohnehin ueber alle Refs. Ein zweiter solcher Durchgang haette keinen eigenen Anlass.
   - Er hat das Ticket unmittelbar davor in den Snapshot-Branch geschrieben. Im Moment des Loeschens liegt es also in ZWEI Ablagen: im Snapshot und in der Branch, in der es gelandet ist. Es gibt damit buchstaeblich nichts zu verlieren - das ist der Unterschied zwischen 'aufraeumen' und 'wegwerfen'.
   - Eine Stelle entscheidet 'dieses Ref hat seinen Zweck erfuellt', statt derselben Pruefung verteilt auf logbook, archive und einen Hintergrundtimer.

   Reihenfolge im Lauf ist damit festgelegt und nicht beliebig: erst den Snapshot schreiben, dann jaeten.
