# shellstube

Ein Homelab-Baukasten, der aus einer Absicht eine Anlage macht — vom Anschluss über
Router, Switch und Access Points bis zu Compute, Speicher und Diensten.

Wir empfehlen **einen** Stack und pflegen **eine** geprüfte Hardwareliste. Nur weil beides
feststeht, können wir am Ende echte Konfigurationsdateien liefern statt guter Ratschläge.

## Was das hier ist

`stube` liest ein `homelab.yaml` — deine Absicht — und schreibt daraus lesbare Dateien auf
deine Geräte: Quadlet-Units für rootless Podman, `/etc/config/*` für OpenWrt, nftables-Regeln,
Sicherungsaufträge, dazu ein Runbook, das jede erzeugte Datei erklärt.

Kein Daemon. Kein Agent auf dem Host. Kein Zustand, den nur wir kennen.

## Die fünf Grundsätze

1. **Erklärbarkeit vor Automatik.** Alles landet als lesbare Datei auf der Platte. Lösch uns
   morgen — dein System läuft weiter.
2. **Rootless als Standard.** Wo es nicht geht, steht der Grund im Katalog, nicht in einer Fußnote.
3. **Kein verstecktes Wissen.** Der einzige Zustand ist `homelab.yaml` plus die erzeugten
   Dateien. Beides gehört in dein eigenes Git.
4. **Rückwärts vor vorwärts.** Ein Dienst kommt erst in den Katalog, wenn Sicherung *und*
   Wiederherstellung getestet sind.
5. **Deinstallation ist ein Feature.** Jeder Dienst geht vollständig wieder weg. Kein toter Code,
   keine toten Dienste.

## Stand

Konzeptphase abgeschlossen, Bau beginnt. Es gibt noch nichts zu installieren.
Der Weg steht in [docs/de/bauplan.md](docs/de/bauplan.md).

## Lizenz

Code AGPL-3.0-or-later · Dokumentation CC BY-SA 4.0 · Beispielkonfigurationen MIT.
Name und Logo sind ausgenommen, siehe [TRADEMARK.md](TRADEMARK.md).
