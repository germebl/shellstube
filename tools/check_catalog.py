#!/usr/bin/env python3
"""Prueft catalog/services/*.yaml. Das Qualitaetstor als ausfuehrbarer Code."""
import glob, os, sys, yaml

errors = []
def err(f, m): errors.append(f"{f}: {m}")

for f in sorted(glob.glob("catalog/services/*.yaml")):
    if os.path.basename(f).startswith("_"):
        continue
    d = yaml.safe_load(open(f, encoding="utf-8")) or {}

    for key in ("id", "name", "category", "summary", "status"):
        if not d.get(key):
            err(f, f"Pflichtfeld '{key}' fehlt")
    if d.get("id") and not f.endswith(f"/{d['id']}.yaml"):
        err(f, f"Dateiname passt nicht zu id '{d['id']}'")
    if d.get("status") not in ("supported", "experimental"):
        err(f, f"status '{d.get('status')}' unbekannt")

    # Rootless ist der Standard, Abweichung braucht einen Grund.
    if d.get("rootless") is False and not d.get("rootless_exception"):
        err(f, "rootless: false ohne rootless_exception")

    # Keine privilegierten Ports.
    net = d.get("network") or {}
    for p in (net.get("ports") or []):
        try:
            num = int(str(p).split("/")[0])
        except ValueError:
            err(f, f"Port '{p}' nicht lesbar"); continue
        if num < 1024:
            err(f, f"Port {num} unter 1024 — rootless bindet das nicht, nftables leitet weiter")
    if net.get("mode") == "host" and not net.get("host_reason"):
        err(f, "network.mode host ohne host_reason")

    # Das Qualitaetstor.
    bk = d.get("backup") or {}
    if d.get("status") == "supported":
        if not bk.get("restore_tested"):
            err(f, "supported ohne getestete Wiederherstellung — das entscheidet integrations-tester")
        if not d.get("docs") or not os.path.exists(d["docs"]):
            err(f, f"supported ohne vorhandene Doku-Seite ({d.get('docs')})")
    if not bk.get("strategy"):
        err(f, "backup.strategy fehlt")

    # Ressourcen ehrlich kennzeichnen.
    res = d.get("resources") or {}
    if res.get("memory_mb") and res.get("memory_source") not in ("measured", "upstream", "estimate"):
        err(f, "resources.memory_source fehlt oder unbekannt")

    if "latest" == (d.get("image") or {}).get("tag_policy") and not d.get("caveats"):
        err(f, "tag_policy latest ohne Begruendung in caveats")

if errors:
    print(f"{len(errors)} Befund(e):\n")
    for e in errors: print("  " + e)
    sys.exit(1)
n = len([f for f in glob.glob("catalog/services/*.yaml") if not os.path.basename(f).startswith("_")])
print(f"catalog/: {n} Eintraege, alles sauber")
