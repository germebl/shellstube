# web

Website und Dokumentations-Site. Liegt bewusst im selben Repo wie die Daten, die sie
darstellt — Hardwareliste, Dienstkatalog und `docs/` werden beim Bauen direkt von hier
aus gelesen, nicht kopiert und nicht als Submodul eingebunden.

Ein Pull Request kann damit Daten, Doku und Darstellung in einem Zug ändern, und die CI
kann durchsetzen, dass keins ohne das andere kommt.

Der Auslieferungs-Workflow filtert nach Pfaden, damit eine Änderung am Go-Teil nicht die
Website neu baut und umgekehrt.

Noch leer — kommt mit M1.
