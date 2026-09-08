---
title: Architektur
description: Wie aus einem Manifest Dateien auf Geräten werden.
sidebar_order: 20
updated: 2026-09-08
---

# Architektur

    homelab.yaml ──▶ manifest ──▶ plan ──▶ render ──▶ Dateien
      (Absicht)      (prüfen)   (ableiten) (schreiben)  (auf deinen Geräten)

## Die vier Schritte

**manifest** liest und prüft. Ein Manifest, das durchkommt, ist widerspruchsfrei — aber noch
nicht vollständig: Vieles leiten wir ab, statt es abzufragen.

**plan** leitet ab, was sich ableiten lässt: welches Gerät welche Rolle bekommt, welcher Port
in welchem VLAN liegt, ob die Strecke die Zielgeschwindigkeit trägt, wie viel Arbeitsspeicher
die gewählten Dienste brauchen. Hier entstehen auch die Lücken und Empfehlungen.

**render** schreibt. Quadlet-Units, `/etc/config/*` für OpenWrt, nftables-Regeln,
Sicherungsaufträge, Kabelliste — und das Runbook, das jede erzeugte Datei erklärt.

## Was hier bewusst fehlt

Es gibt keinen Schritt, der auf ein Netzwerkgerät zugreift. `render` schreibt Dateien; wie sie
auf den Switch kommen, entscheidest du. Ab Stufe T2 erzeugen wir für Firewall und Switch
ohnehin nur Plan und Prüfliste — eine Firewall fernzukonfigurieren, während man über sie
zugreift, ist die zuverlässigste Art, sich auszusperren.

## Geheimnisse

Das Manifest enthält keine. Wo eines gebraucht wird, steht ein Verweis (`credentials_ref`),
und `stube` liest den Wert aus einer lokalen Datei, die nie in ein Repo gehört. Deshalb kann
ein Manifest bedenkenlos geteilt, versioniert und in der Stube gespeichert werden.
