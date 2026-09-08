# Sicherheit

## Etwas gefunden?

Schreib an **security@shellstube.de**. Bitte nicht als öffentliches Issue.

Wir melden uns innerhalb von 72 Stunden. Wenn du eine Frist setzen willst: 90 Tage sind für uns
in Ordnung, kürzer bei aktiver Ausnutzung. Wir nennen dich in der Veröffentlichung, wenn du das
möchtest.

## Was besonders zählt

Dieses Projekt erzeugt Firewall-Regeln, Netztrennung und Zugriffswege. Ein Fehler darin ist
schwerwiegender als ein Absturz. Vorrangig sind deshalb:

- Erzeugte Regelwerke, die mehr erlauben als das Manifest beschreibt
- VLAN-Zuordnungen, die Trennung aufheben
- Rootless-Ausnahmen, die weiter reichen als dokumentiert
- Geheimnisse, die versehentlich in erzeugte Dateien oder Protokolle geraten

## Was kein Sicherheitsproblem ist

- Dass `homelab.yaml` deinen Netzaufbau beschreibt. Das ist Absicht — es enthält keine Geheimnisse.
- Dass `--dry-run` zeigt, was passieren würde. Das ist der Sinn der Sache.
