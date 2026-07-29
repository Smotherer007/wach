/*
Wach — macOS Mouse Mover für Apple Silicon

Wach (deutsch: "wach" / "awake") ist eine minimalistische Menüleisten-App,
die den Mauszeiger in regelmässigen Abständen bewegt, wenn das System
inaktiv ist. Dadurch bleibt der Mac "wach" und Messaging-Apps (Slack,
Teams, etc.) setzen den Status nicht auf "Abwesend".

Im Gegensatz zu Caffeine-Apps, die nur den Bildschirmschoner verhindern,
simuliert Wach echte Benutzeraktivität durch Mausbewegungen.

Funktionsweise:
  - Alle 10 Sekunden wird die Leerlaufzeit via CoreGraphics geprüft
  - Bei >60s Inaktivität + wachem Display: Maus wird um 10px bewegt
  - Die Richtung wechselt bei jeder Bewegung (kein Abdriften)
  - Bei Benutzeraktivität oder schlafendem Display: keine Aktion

Fehlerbehandlung:
  - Fehlt die Berechtigung, erscheint nach 10 Fehlversuchen ein Hinweis
  - Maximal eine Fehlermeldung pro 24 Stunden

Lizenz: MIT (siehe LICENSE)
*/
package wach
