#!/usr/bin/env python3
"""Prueft hardware/*.yaml gegen die Regeln aus CONTRIBUTING.md.

Was hier geprueft wird, muss kein Mensch und kein Modell mehr pruefen.
"""
import glob, sys, yaml

CLASSES = {"ont","modem","router","switch","ap","compute","storage","ups","accessory"}
STATUS  = {"recommended","tolerated","experimental","unsupported"}
POE     = {"full","manual","broken","none"}
PORTS   = {"rj45_1g","rj45_2g5","rj45_10g","sfp","sfp_plus","qsfp","combo"}

errors = []

def err(f, msg):
    errors.append(f"{f}: {msg}")

for f in sorted(glob.glob("hardware/*.yaml")):
    d = yaml.safe_load(open(f, encoding="utf-8")) or {}

    for key in ("id", "class", "vendor", "model", "status"):
        if not d.get(key):
            err(f, f"Pflichtfeld '{key}' fehlt")

    if d.get("id") and not f.endswith(f"/{d['id']}.yaml"):
        err(f, f"Dateiname passt nicht zu id '{d['id']}'")
    if d.get("class") not in CLASSES:
        err(f, f"unbekannte class '{d.get('class')}'")
    if d.get("status") not in STATUS:
        err(f, f"unbekannter status '{d.get('status')}'")

    # Regel 1: Revision ist Pflicht.
    if not d.get("revision"):
        err(f, "revision fehlt — gleiche Modellnummer, andere Platine, anderes Verhalten")

    # Regel 2: keine Herstellerzahlen abschreiben.
    power = d.get("power") or {}
    if "power" in d and "idle_w" not in power:
        err(f, "power ohne idle_w — lieber null als geraten")
    if power.get("idle_w") is not None and power.get("source") in (None, "unknown"):
        err(f, "idle_w gesetzt, aber source fehlt — woher stammt die Zahl?")
    if power.get("source") == "vendor":
        err(f, "source 'vendor' ist nicht zulaessig, Herstellerangaben stimmen im Leerlauf nie")

    # Regel 3: ehrlich sein, was geprueft wurde.
    ver = d.get("verified") or {}
    if "verified" not in d:
        err(f, "verified fehlt — 'by: none' ist eine gueltige und ehrliche Angabe")
    elif ver.get("by") not in (None, "none", "ci", "maintainer", "community"):
        err(f, f"verified.by '{ver.get('by')}' unbekannt")

    # Empfohlen heisst: wir liefern etwas.
    if d.get("status") == "recommended" and not d.get("generates"):
        err(f, "status recommended ohne generates — dann ist es hoechstens tolerated")
    if d.get("status") == "unsupported" and not d.get("status_reason"):
        err(f, "unsupported ohne status_reason — die Frage 'warum nicht' muss beantwortet sein")

    # Ports und PoE.
    for k in (d.get("ports") or {}):
        if k not in PORTS:
            err(f, f"unbekannter Porttyp '{k}'")
    poe = d.get("poe") or {}
    if poe.get("ports") and poe.get("openwrt") not in POE:
        err(f, f"poe.ports gesetzt, aber poe.openwrt '{poe.get('openwrt')}' unbekannt")
    if poe.get("openwrt") in ("manual", "broken") and not poe.get("note"):
        err(f, "poe.openwrt manual/broken ohne note — was genau geht nicht?")

    # Bewegliche Ziele brauchen ein Datum.
    for c in (d.get("caveats") or []):
        low = str(c).lower()
        if any(w in low for w in ("inzwischen", "mittlerweile", "neuerdings")) and not any(
                m in low for m in ("januar","februar","märz","april","mai","juni","juli",
                                   "august","september","oktober","november","dezember","20")):
            err(f, f"caveat ohne Datum: {str(c)[:60]}…")

    if d.get("price_eur") and not d.get("price_asof"):
        err(f, "price_eur ohne price_asof — ein Preis ohne Datum ist keine Angabe")

# Querverweise
ids = {(yaml.safe_load(open(f, encoding="utf-8")) or {}).get("id") for f in glob.glob("hardware/*.yaml")}
for f in sorted(glob.glob("hardware/*.yaml")):
    d = yaml.safe_load(open(f, encoding="utf-8")) or {}
    for a in (d.get("alternatives") or []):
        if a not in ids:
            err(f, f"alternatives verweist auf unbekannte id '{a}'")

if errors:
    print(f"{len(errors)} Befund(e):\n")
    for e in errors:
        print("  " + e)
    sys.exit(1)
print(f"hardware/: {len(ids)} Eintraege, alles sauber")
