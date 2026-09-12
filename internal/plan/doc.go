// Package plan leitet aus einem geprueften Manifest und der Hardwareliste ab,
// was sich ableiten laesst: Rollenverteilung, Portbelegung, VLAN-Plan,
// Streckenpruefung gegen target_speed_gbe und Arbeitsspeicherbedarf.
//
// Wofür dieses Paket zuständig ist und wie es in die Kette passt, steht in
// docs/de/architektur.md.
package plan
