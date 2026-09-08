# Beitragen

## Kurzfassung

- **Ein Gerät beitragen:** eine YAML-Datei unter `hardware/`. Kein Go nötig.
- **Einen Dienst beitragen:** eine YAML-Datei unter `catalog/services/`. Kein Go nötig.
- **Am Werkzeug arbeiten:** Go, siehe unten.

Sprache: Code, Bezeichner, Commit-Nachrichten und Issue-Titel auf **Englisch**.
Dokumentation und alles, was Nutzer lesen, auf **Deutsch** (Englisch folgt ab M3).

## Herkunftsnachweis statt Lizenzvertrag

Wir verlangen **keinen CLA**. Stattdessen den Developer Certificate of Origin: Du bestätigst,
dass du das Recht hast, deinen Beitrag unter unserer Lizenz einzubringen.

    git commit -s -m "add zyxel gs1900-24hp"

Das hängt eine `Signed-off-by:`-Zeile an. Ohne sie schließt die CI den Pull Request nicht.

Was das bedeutet: Eine spätere Doppellizenzierung ist damit dauerhaft ausgeschlossen. Das ist
Absicht — niedrige Hürde für dich ist uns wichtiger als eine Option, die wir nicht ziehen wollen.

## Ein Gerät auf die Hardwareliste

Kopiere `schema/hardware.schema.yaml` als Vorlage. Drei Regeln, die keine Verhandlungssache sind:

1. **Revision angeben.** Dieselbe Modellnummer verhält sich auf einer neueren Platine anders.
   `revision: "*"` nur, wenn du nachweisen kannst, dass alle gleich laufen.
2. **`power.idle_w: null`, solange du nicht selbst gemessen hast.** Schreib keine
   Herstellerangaben ab — die stimmen im Leerlauf nie.
3. **`verified.by: none`, wenn du es nicht getestet hast.** Eine ehrliche Lücke ist brauchbarer
   als eine geratene Zahl.

Ein Eintrag mit Status `unsupported` ist genauso wertvoll wie eine Empfehlung. Wer wissen will,
warum sein Gerät nicht dabei ist, soll die Antwort finden statt sie zu suchen.

## Einen Dienst in den Katalog

Vorlage: `catalog/services/_template.yaml`. Ein Dienst gilt erst als `supported`, wenn alles
davon zutrifft:

- läuft rootless (oder die Ausnahme ist im Feld `rootless_exception` begründet)
- Sicherung **und** Wiederherstellung sind getestet
- eine Doku-Seite existiert
- ein CI-Test läuft gegen eine echte VM
- `stube remove` hinterlässt nichts

Alles andere trägt sichtbar `experimental`. Wir nehmen unfertige Dienste auf — aber sie sind
als unfertig markiert, nicht versteckt.

## Commits

Betreff im Imperativ, klein geschrieben, unter 60 Zeichen, kein Punkt am Ende. Rumpf nur,
wenn das Warum nicht aus dem Diff hervorgeht.

    poe-budget je port rendern
    revision beim dgs-1210 nachtragen

    Auf F3 kann realtek-poe den controller nicht ansprechen. Der eintrag
    stand vorher ohne revision da und war damit irreführend.

Keine Emoji, keine Präfixe wie `feat(...)`, keine Fußzeilen außer `Signed-off-by`.
Ein Commit ist ein Gedanke — nicht vierzig Änderungen in einem, aber auch nicht jede Zeile
einzeln.

## Am Werkzeug arbeiten

    make lint      # golangci-lint, yamllint, markdownlint, vale, shellcheck
    make test      # Unit-Tests
    make test-vm   # Integrationstests gegen eine frische VM (langsam)
    make schema    # homelab.yaml gegen das Schema prüfen

Bevor du einen Pull Request öffnest: `make lint test`. Was ein Werkzeug prüfen kann, prüft ein
Werkzeug — dafür brauchst du keine Rückmeldung von uns und wir keine Diskussion.

## Was wir nicht annehmen

- Geräte, die ohne Herstellerkonto oder Cloud nicht laufen
- Dienste ohne getestete Wiederherstellung
- Automatik, die Netzwerkgeräte ohne Rückfrage umkonfiguriert
- Telemetrie in irgendeiner Form
