---
id: 01M291C5YMJWJBKDWVRAXK4124
title: "Die ausgewaehlte Karte leuchtet in ihrer Tag-Farbe, c schaltet um"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Fuellung der ausgewaehlten Karte traegt die Farbe ihres Tags statt eines neutralen Graus, und c schaltet zwischen beiden um"
context: "Alex am 11.09., nach ANX7Y5: die Fuellung ist da, aber sie sagt nur wo der Cursor steht. Sie koennte gleichzeitig sagen, worum das Ticket geht - die Farbe traegt die Karte ohnehin schon im Rahmen.\n\nDas Entscheidende, und der Grund warum hier zwei Zahlen stehen statt einer: in der 256-Farben-Palette gibt es keinen leisen Farbton. Gemessen ueber die acht Tag-Farben, die .jaira/tags vergibt:\n\n  tag   15%   25%   35%   45%\n  170   239.  240.  241.   96\n  45    238.   23    24    30\n  111   239.  240.   60    60\n  78    238.  239.  240.   65\n  (. = auf einem Grau gelandet, die Tag-Farbe ist weg)\n\nUnter etwa 45% Mischung faellt JEDE Tag-Farbe auf die Graustufenleiter 238/239/240, weil die Palette zwischen ihrer Grauleiter und ihrem Farbwuerfel keine schwach gesaettigten dunklen Farben hat. Ein 15%-Ton ist dort also nicht eine leisere Version der Tag-Farbe, sondern schlicht Grau.\n\nDarum zwei Mischwerte: 15% wo das Terminal 24-Bit-Farbe kann, 45% wo nicht. Wuerde man den 15%-Wert von lipgloss herunterrechnen lassen, landet er auf Grau und die Tag-Farbe ist auf jedem aelteren Terminal still verschwunden.\n\nZweite Falle, schon beim Bauen getroffen: gemischt wird AUFWAERTS von einem Grau (Index 237, RGB 58), nicht abwaerts von der Tag-Farbe Richtung Schwarz. Abwaerts landete Cyan bei 18% exakt auf dem Fuellton der Lane, und die Auswahl war unsichtbar.\n\nTaste: c. g war vergeben (erste/letzte Karte der Lane).\n\nOffen und bewusst nicht hier entschieden: Alex hat parallel rahmenlose Karten mit Zebra-Hintergrund angesehen (Mockup unter scratchpad/edgedemo). Das wuerde 81XRXX und VS5DFW beide kippen und ist eine eigene Entscheidung - dieses Ticket aendert nur die Farbe der Fuellung, nicht den Aufbau der Karte."
definition-of-done: "Die ausgewaehlte Karte mit farbigem Tag ist in dessen Ton gefuellt; c schaltet auf neutral und zurueck; die Fusszeile nennt was der naechste Druck tut; auf einem 24-Bit-Terminal wird 15% gemischt, sonst 45%; keine Tag-Farbe landet auf einem Grau; eine Karte ohne farbiges Tag bleibt neutral gefuellt; go test ./... -race gruen"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-11T20:08:36Z
updated-at: 2026-09-11T20:08:51Z
updated-by: Alexander Sacharov
---

# Die ausgewaehlte Karte leuchtet in ihrer Tag-Farbe, c schaltet um

## Definition of Done

- [ ] Die ausgewaehlte Karte mit farbigem Tag ist in dessen Ton gefuellt; c schaltet auf neutral und zurueck; die Fusszeile nennt was der naechste Druck tut; auf einem 24-Bit-Terminal wird 15% gemischt, sonst 45%; keine Tag-Farbe landet auf einem Grau; eine Karte ohne farbiges Tag bleibt neutral gefuellt; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

