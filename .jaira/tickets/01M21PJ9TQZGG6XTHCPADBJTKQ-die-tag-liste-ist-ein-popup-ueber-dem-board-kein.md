---
id: 01M21PJ9TQZGG6XTHCPADBJTKQ
title: "Die Tag-Liste ist ein Popup ueber dem Board, keine eigene Seite"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "'t' oeffnet die Tag-Liste als kleines Fenster mittig ueber dem Board; das Board bleibt ringsum sichtbar"
context: |-
  Berk am 08.09. mit zwei Screenshots.
  Heute ersetzt 't' den ganzen Bildschirm: eine Seite mit Ueberschrift 'Tags', der Liste und der Fusszeile 'esc close / t close'. Das Board ist weg, solange sie offen ist.
  Gewuenscht: dasselbe als Popup. Berk hat den Rahmen auf dem Board-Screenshot weiss eingezeichnet - ein hochkant stehendes Kaestchen, mittig, ueber Brainstorm/Todo, das Board bleibt links und rechts stehen.
  Inhalt der Liste bleibt wie er ist: Farbquadrat plus Tag-Name, eine Zeile je Tag, alphabetisch.
  Die Tasten bleiben: esc schliesst, t schliesst.
  Ort suchen: das Tag-Screen-Rendering in internal/tui - es ist heute ein eigener View-Zweig, kein Overlay.
  Vorbild im Code pruefen, bevor etwas Neues erfunden wird: das Board hat schon andere Zustaende, die ueber allem liegen; wenn es dafuer bereits ein Overlay-Muster gibt, dieses benutzen.
definition-of-done: "'t' zeichnet die Tag-Liste als Kasten ueber dem Board, das Board ist ringsum weiter sichtbar; esc und t schliessen wie bisher; ein Test pinnt, dass Board-Inhalt neben dem Popup stehen bleibt; go test ./... -race gruen"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-08T23:45:01Z
updated-at: 2026-09-08T23:45:01Z
---

# Die Tag-Liste ist ein Popup ueber dem Board, keine eigene Seite

## Definition of Done

- [ ] 't' zeichnet die Tag-Liste als Kasten ueber dem Board, das Board ist ringsum weiter sichtbar; esc und t schliessen wie bisher; ein Test pinnt, dass Board-Inhalt neben dem Popup stehen bleibt; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

