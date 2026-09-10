---
id: 01M266HAGY954EW9K6T08566KF
title: "Tickets reisen in eigenen Git-Refs, damit Zuweisungen ohne gemeinsamen Branch ankommen"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket und sein Besitzer erreichen den Kollegen ueber refs/jaira/tickets/<id>, ohne dass ein Branch geteilt oder gemergt werden muss, und eine Zuweisung loest bei ihm eine Benachrichtigung aus"
context: "Berk will laut Issue #7 benachrichtigt werden, wenn ihm jemand ein Ticket zuweist. Heute geht das nicht: die Ticketdatei liegt im Branch des Schreibers. Ist der Branch ungepusht, veraltet oder force-gepusht, sieht der Empfaenger nichts - weder Titel noch Id. Alle Branches abzuscannen ist teuer und trotzdem falsch. Geprueft und belegt: ein Push auf refs/jaira/... nimmt GitHub an (exit 0), ein zweiter Push aufs selbe Ref wird als non-fast-forward abgelehnt (exit 1) - das ist ein Compare-and-Swap ohne Server. Der Tree des Refs kann die Ticketdatei selbst tragen, ein Kollege liest sie per git show ohne Branch und ohne Checkout. Der Default-Refspec zieht diese Refs nicht mit, wer jaira nicht nutzt merkt nichts. Ausgeschlossen: getrennte claims-/log-Refs (ein Ref pro Ticket reicht, --force-with-lease auf den gelesenen Sha ist das CAS) und ein Release per Commit mit Parent (laesst das Ref bestehen, das Ticket klebt dann fuer immer am ersten Besitzer). Die Ticketdatei bleibt zusaetzlich im Code-Commit, damit der Reviewer Diff und Grund an einer Stelle sieht. Details, Schritte und alle Messergebnisse stehen im Body."
definition-of-done: "core/gitref liest, schreibt und loescht refs/jaira/tickets/<id> mit CAS und unterscheidet 'Rennen verloren' von 'Netz weg'; create/claim/move/note/dod pushen das Ticket zusaetzlich ins Ref und melden bei Ablehnung wer schneller war; ein Fetch holt die Refs und das Board zeigt Ref-Tickets samt Besitzer, auch ohne den Branch des Schreibers; eine Zuweisung an mich erzeugt eine Desktop-Benachrichtigung, abschaltbar; ein optionaler Hook bei move/claim wird aufgerufen; ein ins Logbuch gelegtes Ticket loescht sein Ref; Tests mit zwei Klonen eines Bare-Repos belegen Rennen, Uebernahme und den Merge zweier Branches am selben Ticket"
tags:
  - concurrency
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T17:41:04Z
updated-at: 2026-09-10T18:51:10Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-244146
claimed-at: 2026-09-10T18:17:46Z
---

## Warum Refs und nicht Branches

Ein Alert braucht einen Kanal. jaira hat genau einen: das Git-Remote. Sitzt jeder
in seinem eigenen Branch, hilft das nicht — den Branch des anderen sieht man
nicht, er kann ungepusht, veraltet oder force-gepusht sein. Alle Branches
abzuscannen ist teuer und trotzdem falsch.

Ein Ref ausserhalb von `refs/heads` loest das: es gehoert zu keinem Branch, der
Default-Refspec zieht es nicht mit (wer jaira nicht nutzt, merkt nichts), und
der Push darauf ist ein Compare-and-Swap, das die Ref-Sperre des Hosters nutzt.

## Das Ref-Layout

    refs/jaira/tickets/<id>   # ein Ref pro Ticket, Tree enthaelt die Ticketdatei

Ein Ref pro Ticket, sein Tree enthaelt `<id>.md` — nicht nur die Id, sondern den
vollen Inhalt. Damit liest ein Kollege das Ticket ohne Branch und ohne Checkout:

    git show refs/jaira/tickets/<id>:<id>.md

Lane und Besitzer stehen im Frontmatter, die Historie ist `git log` auf dem Ref.
Getrennte claims-/log-Refs sind dadurch ueberfluessig.

## Compare-and-Swap

Jeder Schreibvorgang pusht mit `--force-with-lease=refs/jaira/tickets/<id>:<sha
den man gelesen hat>`. Hat jemand dazwischen geschrieben, lehnt Git ab (exit 1),
und jaira sagt: neu einlesen, das Ticket hat sich bewegt. Ein Rennen wird sofort
erkannt, nicht drei Tage spaeter beim Merge.

## Verhaeltnis zur Ticketdatei im Commit

Beides bleibt, mit klarer Zustaendigkeit:

| | Datei im Commit | Ref |
|---|---|---|
| Antwortet auf | was war das fuer eine Aenderung, warum | wo steht das Ticket jetzt, bei wem |
| Liest | Reviewer des Diffs | Board, Alerts, andere Maschine |
| Aendert sich | einmal, mit dem Code | bei jedem move/claim/note |
| Lebt | dauerhaft in der Master-Historie | solange das Ticket auf dem Board ist |

Der Ablauf aus CLAUDE.md bleibt unveraendert: Ticket bewegen, Ticketdatei
zusammen mit dem Code committen. Zusaetzlich pusht dieselbe Operation das Ref
mit denselben Bytes. Keine zweite Wahrheit, derselbe Blob auf einem Kanal, der
ohne den Branch sichtbar ist.

Solange der Branch nicht gemergt ist, ist das Ref massgeblich fuer Lane und
Besitzer. Die Datei im Commit ist eine Momentaufnahme: auf diesem Commit stand
das Ticket hier. Genau das braucht der Reviewer.

## Zwei Branches am selben Ticket

Auf Ref-Ebene entsteht der Konflikt gar nicht erst — CAS lehnt den Zweiten ab.
Auf Dateiebene loest der bestehende Merge-Driver (`core/merge`): `status` nach
Fortschritt in der Lane-Kette, Listen (`commits`, `tags`, `blocked-by`) als
Union, Prosafelder als einziger echter Konflikt. Gegen die reale Registrierung
geprueft, siehe unten.

## Offen: Claim ohne Netz

Ohne Verbindung gibt es kein CAS. Vorschlag: `claim` verlangt Netz, alles andere
(`move`, `note`, `dod`) bleibt vollstaendig offline. Eine einzige
Online-Operation ist leicht zu erklaeren und unmoeglich misszuverstehen. Die
Alternative — lokal provisorisch claimen und beim ersten Push nach Timestamp
abgleichen — ist gebaut komplizierter und im Konfliktfall unangenehm.

## Was schon geprueft ist (gegen echtes GitHub, nicht theoretisch)

- Push eines neuen `refs/jaira/...` auf GitHub: angenommen, exit 0.
- Zweiter Push ohne Parent auf dasselbe Ref: abgelehnt, non-fast-forward, exit 1.
- `--force-with-lease` auf den gelesenen Sha: der Zweite wird abgelehnt, der
  Erste gewinnt, der Serverzustand bleibt konsistent.
- `git ls-remote origin 'refs/jaira/*'`: alle Refs in einem Roundtrip sichtbar.
- Loeschen per `:refs/jaira/...`: funktioniert.
- Branches, PRs und die GitHub-Weboberflaeche bleiben unberuehrt — die Refs sind
  dort nirgends sichtbar. Deshalb muss `jaira board` die Besitzer selbst zeigen.
- Ticketdatei im Ref-Tree, ohne jeden Branch beim Empfaenger lesbar.
- Merge-Driver: `in-progress` vs `review` ergibt `review`, obwohl `in-progress`
  den neueren Timestamp trug; `done` vs `backlog` ergibt `done`; `tags` und
  `commits` beider Seiten bleiben erhalten.

## Umsetzungsschritte

1. `core/gitref`: Refs lesen/schreiben/loeschen ueber `os/exec`, CAS mit
   `--force-with-lease`, Retry bei Netzfehler, klarer Fehlertyp fuer "Rennen
   verloren" (unterscheidbar von "Netz weg").
2. Schreibpfad: `create`, `claim`, `move`, `note`, `dod` pushen dasselbe Ticket
   zusaetzlich ins Ref. Bei Ablehnung: neu einlesen und dem Aufrufer sagen, wer
   schneller war.
3. Lesepfad: `jaira fetch` (bzw. Hintergrund-Fetch im TUI) holt
   `+refs/jaira/tickets/*:refs/jaira/tickets/*` und mischt Ref-Tickets in die
   Boardansicht — mit Kennzeichnung, dass sie noch in keinem Branch liegen.
4. Alert: Aenderung an einem Ticket mit einem selbst als Assignee loest eine
   Desktop-Benachrichtigung aus (`notify-send` / `osascript`), Intervall
   30–60 s, per Konfiguration abschaltbar.
5. Hook-Punkt fuer sofortige Zustellung: ein optionales Skript, das bei
   `move`/`claim` aufgerufen wird. jaira ruft nur auf und bringt keine
   Abhaengigkeit mit (Slack-Webhook, ntfy.sh, Telegram sind Sache des Nutzers).
6. Aufraeumen: gemergte bzw. ins Logbuch gelegte Tickets loeschen ihr Ref.
7. Tests: zwei Klone eines Bare-Repos in `t.TempDir()`, echtes Rennen um ein
   Ticket, Release und Uebernahme, Merge zweier Branches am selben Ticket.

## Was ausdruecklich nicht dazugehoert

Kein Server, kein Daemon, keine zentrale Koordination. Alles laeuft ueber das
Remote, das ohnehin da ist. Wer jaira nicht benutzt, sieht von alldem nichts.

## Options

- [x] planning

## Plan

- [x] core/gitref: Ref-Inhalt ohne Checkout bauen (hash-object -w, mktree, commit-tree mit dem geleasten Sha als Parent) und lesen (git show <sha>:<id>.md)
  proof: core/gitref/gitref.go:Write/Read, TestTicketArrivesWithoutASharedBranch
- [x] core/gitref: Push mit --force-with-lease=<ref>:<gelesener sha>, Fehler typisieren: ErrRaceLost (non-fast-forward/stale info) vs ErrOffline (Netz/Remote weg)
  proof: core/gitref/gitref.go:classify + push, TestSecondWriterLosesTheRace / TestUnreachableRemoteIsOfflineNotARace
- [ ] Entschieden: asynchron. Mutate/Create legen den Schreibvorgang in eine Outbox, gepusht wird ausserhalb - Voraussetzung fuer offline claim, loest zugleich TUI und Merge-Driver
- [ ] Schreibpfad anbinden: der Push haengt an Store.Mutate + Store.Create, nicht an create/claim/move/note/dod einzeln (alle 20+ Aufrufstellen laufen durch diese zwei)
- [ ] Bei Ablehnung: Ticket aus dem Ref neu einlesen und dem Aufrufer melden, wer schneller war (Assignee + updated-at aus dem fremden Ref)
- [ ] Lesepfad: fetch von +refs/jaira/tickets/*:refs/jaira/tickets/* plus git ls-remote als billiger Ueberblick
- [ ] Board: Ref-Tickets, die in keinem Branch liegen, in die Ansicht mischen und als solche kennzeichnen
- [ ] Settings-Datei in ~/.jaira/ anlegen (existiert heute nicht, nur projects.json und state/) - traegt den Abschalter fuer die Benachrichtigung
- [ ] Desktop-Benachrichtigung bei Zuweisung an mich: notify-send / osascript per os/exec, kein neues Modul, faellt still aus wenn nichts da ist
- [ ] Optionaler Hook bei move/claim: jaira ruft ein Skript auf und bringt keine Abhaengigkeit mit
- [ ] logbook und archive loeschen das Ref (git push origin :refs/jaira/tickets/<id>)
- [ ] Tests: Bare-Repo plus zwei Klone in t.TempDir(), echtes Rennen, Uebernahme, Merge zweier Branches am selben Ticket
- [ ] Outbox unter ~/.jaira/state/<worktree>/: ein noch nicht gepushter Schreibvorgang pro Ticket, mit dem gelesenen Sha als Lease
- [ ] Outbox abarbeiten beim naechsten Kommando, das Netz hat; abgelehnter nachgeholter Push meldet dem Benutzer, wer schneller war
- [ ] Board-Remote konfigurierbar (Default origin) - die Settings-Datei traegt es neben dem Abschalter der Benachrichtigung
- [ ] README: wozu die Refs da sind, welche Befehle sie schreiben und lesen, wie das Board-Remote konfiguriert wird, und die Fork-Grenze ausdruecklich (kein Push-Recht = Ticketdatei im Branch plus PR)

## Progress

- **2026-09-10 18:19 · Alexander Sacharov** — Warum der Plan so aussieht, drei Funde im Code (gemessen, nicht vermutet):

1. Es gibt genau zwei Schreibtore: core/ticket/store.go:644 Mutate und :620 Create. Alle Aufrufstellen (create, claim, note, dod, tags, move, TUI, Merge-Driver - 20+ per grep) laufen durch sie. Der Ref-Push gehoert deshalb dorthin, nicht in die fuenf CLI-Kommandos einzeln - genau das Muster von JJ32B4 (eine Move-Funktion, vier delegierende Aufrufer).

2. Genau daraus folgt aber ein Problem, das die Umsetzungsschritte im Ticket nicht nennen: Mutate wird auch aus dem TUI (internal/tui/edit.go, external.go, followup.go) und aus dem Merge-Driver (internal/cli/mergedriver.go:221) gerufen. Ein synchroner Push in Mutate haengt das Board am Netz und laesst den Merge-Driver mitten im git-merge pushen. Deshalb ist Schritt 3 eine Entscheidung und keine Implementierung.

3. core/gitrepo ist heute reines Lesen (git.go, 141 Zeilen: Commits, Diff, Stat, HeadSHA) und kennt kein Remote. Ein Schreib-Paket core/gitref daneben ist additiv, nichts Bestehendes muss umgebaut werden.

4. Fuer den Abschalter der Benachrichtigung gibt es noch keinen Ort: in ~/.jaira/ liegen nur projects.json und state/<worktree>/. Die Settings-Datei ist damit Teil dieses Tickets, nicht Voraussetzung.

Nicht geprueft, absichtlich: ob GitHub die Refs annimmt - das steht im Ticketbody als bereits gemessen (exit 0 / exit 1 / force-with-lease / ls-remote / Loeschen), ich habe es nicht wiederholt.
- **2026-09-10 18:44 · Alexander Sacharov** — Geprueft, ob der Plan ueber einen Fork funktioniert. Er tut es nicht.

Test: ein Ref nach origin (mein Fork sashasoft90/jaira) gepusht, dann beide Remotes abgefragt.
  git ls-remote origin   'refs/jaira/*'  -> refs/jaira/probe sichtbar
  git ls-remote upstream 'refs/jaira/*'  -> leer

Forks teilen auf GitHub den Objektspeicher, aber jedes Repository hat eigene Refs. Branches sieht man beim anderen nur, weil ein PR ein eigener Kopiermechanismus ist. Custom Refs fahren da nicht mit.

Folge fuer die Umsetzung: es braucht ein Board-Remote - ein Repository, in das alle Beteiligten die jaira-Refs pushen. Nicht 'jeder in seinen Fork'. Default origin, per Konfiguration ueberschreibbar. Ohne das schreibe ich in meinen Fork und Berk schaut in upstream, und beide sehen nichts.

Auch geprueft: ich habe push auf BeMuCa/jaira (permissions.push = true), master ist nicht geschuetzt. Upstream als Board-Remote ist also ohne neues Repository moeglich.

Nicht loesbar und deshalb festhalten: GitHub kann Rechte nicht auf ein Ref-Namespace einschraenken. push auf refs/jaira/* ohne push auf refs/heads/* gibt es nicht - Schreibrecht gilt immer fuers ganze Repository. Teilnehmer am Board = Person mit Push-Recht auf das Board-Repo. Fuer ein Team in Ordnung, fuer einen fremden Contributor nicht: der bleibt beim heutigen Weg, Ticketdatei im Branch plus PR.

Naechster Schritt: Konfigurationsfeld fuers Board-Remote in Schritt 1 der Umsetzung aufnehmen, bevor core/gitref geschrieben wird.
- **2026-09-10 18:47 · Alexander Sacharov** — Entschieden (Alexander, 10.09.2026): claim bleibt offline moeglich. Offline ist selten genug, dass es die Bedienung nicht diktieren darf.

Der Ablauf: claim schreibt lokal mit Timestamp, der Push wird nachgeholt, sobald Netz da ist, und wird er dann abgelehnt, sagt jaira es dem Benutzer (wer schneller war), statt still zu verlieren.

Damit ist die Alternative aus dem Ticketbody ('claim verlangt Netz') verworfen - der Vorschlag dort ist nicht mehr gueltig.

Folge fuer den Plan: Schritt 3 ist damit entschieden, und zwar auf asynchron. Ein synchroner Push in Mutate kann es gar nicht sein, denn ein lokal moeglicher claim braucht ohnehin einen Ort, an dem ein noch nicht gepushter Schreibvorgang liegt (Outbox unter ~/.jaira/state/<worktree>/). Das loest gleichzeitig das TUI- und Merge-Driver-Problem aus Note 1: Mutate legt nur ab, gepusht wird ausserhalb.

Damit gehoert zum Ticket zusaetzlich: die Outbox selbst, ihr Abarbeiten beim naechsten Kommando mit Netz, und die Meldung bei Ablehnung eines nachgeholten Pushs.
- **2026-09-10 18:51 · Alexander Sacharov** — Abweichung vom Plan, Schritt 1: der Ref-Commit bekommt den geleasten Sha als Parent, nicht 'ohne Parent' wie im Ticketbody vorgesehen.

Grund: mit Parent ist ein veralteter Schreibvorgang von sich aus ein non-fast-forward, den git ohnehin ablehnt - das CAS haengt dann nicht allein an --force-with-lease. Ausserdem wird 'git log refs/jaira/tickets/<id>' damit zur Bewegungshistorie des Tickets, was der Ticketbody selbst voraussetzt ('die Historie ist git log auf dem Ref') - ohne Parent gaebe es genau einen Commit und keine Historie.

--force-with-lease bleibt trotzdem gesetzt: beim allerersten Schreiben hat kein Commit einen Parent, mit dem git vergleichen koennte, und nur der leere Lease sagt 'ich erwarte, dass es das Ref noch nicht gibt'. Belegt in TestFirstWriteRefusesAnExistingRef.

Zwei Dinge, die beim Bauen aufgefallen sind und nicht im Plan standen:

- fetch braucht --prune. Ohne das bleibt ein geloeschtes Ref (Logbuch, Archiv) auf jedem anderen Klon fuer immer stehen: fuer den Besitzer weg, fuer alle anderen unsterblich. Steht jetzt als Vertrag im Doc-Kommentar von Fetch, nicht als Aufraeumen.
- Der Push setzt GIT_TERMINAL_PROMPT=0 und gitref schreibt Author/Committer selbst in die Env. Sonst haengt ein Push, der im Hintergrund eines fremden Kommandos laeuft, an einem Credential-Prompt, und auf einer Maschine ohne user.name (CI, frischer Container) scheitert commit-tree mit 'please tell me who you are'.

Lokales Ref wird erst nach dem akzeptierten Push bewegt (update-ref danach), damit der naechste Lease immer ein Sha ist, das das Remote wirklich gesehen hat. Belegt am Ende von TestUnreachableRemoteIsOfflineNotARace.

## Definition of Done

- [x] core/gitref liest, schreibt und loescht refs/jaira/tickets/<id> mit CAS und unterscheidet 'Rennen verloren' von 'Netz weg'
  proof: core/gitref/gitref.go (Write/Read/Delete/SHA, ErrRaceLost vs ErrOffline), 7 Tests in core/gitref/gitref_test.go gegen ein Bare-Repo mit zwei Klonen
- [ ] Ein konfigurierbares Board-Remote entscheidet, wohin die Refs gehen; Default origin
- [ ] create, claim, move, note und dod pushen das Ticket zusaetzlich ins Ref und melden bei Ablehnung, wer schneller war
- [ ] Ein Fetch holt die Refs und das Board zeigt Ref-Tickets samt Besitzer, auch ohne den Branch des Schreibers
- [ ] Eine Zuweisung an mich erzeugt eine Desktop-Benachrichtigung, per Konfiguration abschaltbar
- [ ] Ein optionaler Hook bei move und claim wird aufgerufen, ohne dass jaira eine Abhaengigkeit mitbringt
- [ ] Ein ins Logbuch gelegtes Ticket loescht sein Ref
- [ ] Tests mit zwei Klonen eines Bare-Repos belegen Rennen, Uebernahme und den Merge zweier Branches am selben Ticket
- [ ] Die README erklaert die Funktion: wozu die Refs da sind, welche Befehle sie schreiben und lesen, und wie man das Board-Remote konfiguriert
- [ ] Die README nennt die Fork-Grenze ausdruecklich: Refs reisen nicht ueber Forks, alle Beteiligten pushen in dasselbe Board-Repo, und Teilnahme setzt Push-Recht darauf voraus - ein Contributor ohne Push-Recht bleibt bei Ticketdatei im Branch plus PR
