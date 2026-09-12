package hardware

import (
	"fmt"
	"sort"
	"strings"
)

var (
	classes    = toSet("ont", "modem", "router", "switch", "ap", "compute", "storage", "ups", "accessory")
	statuses   = toSet("recommended", "bestand", "tolerated", "experimental", "unsupported")
	poeOpenwrt = toSet("full", "manual", "broken", "none")
	portTypes  = toSet("rj45_1g", "rj45_2g5", "rj45_10g", "sfp", "sfp_plus", "qsfp", "combo")
	verifiedBy = toSet("none", "ci", "maintainer", "community")
	months     = []string{
		"januar", "februar", "märz", "april", "mai", "juni", "juli",
		"august", "september", "oktober", "november", "dezember", "20",
	}
	movingTargetWords = []string{"inzwischen", "mittlerweile", "neuerdings"}
)

func toSet(items ...string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, it := range items {
		m[it] = true
	}
	return m
}

// Issue ist ein einzelner Befund zu einer Datei aus hardware/*.yaml.
type Issue struct {
	File    string
	Message string
}

func (i Issue) String() string { return fmt.Sprintf("%s: %s", i.File, i.Message) }

// Issues sammelt alle Befunde eines Validate-Laufs.
type Issues []Issue

func (is Issues) Error() string {
	lines := make([]string, len(is))
	for i, issue := range is {
		lines[i] = issue.String()
	}
	return strings.Join(lines, "\n")
}

// Validate bildet die Regeln aus tools/check_hardware.py nach: dieselben
// Eingaben muessen dieselben Befunde liefern. Die Python-Fassung bleibt die
// zweite Meinung, solange beide im Repo stehen.
func Validate(devices *Devices) Issues {
	var issues Issues
	add := func(path, format string, args ...interface{}) {
		issues = append(issues, Issue{File: path, Message: fmt.Sprintf(format, args...)})
	}

	ids := map[string]bool{}
	for _, d := range devices.List {
		if d.ID != "" {
			ids[d.ID] = true
		}
	}

	for _, d := range devices.List {
		data := d.Data

		for _, key := range []string{"id", "class", "vendor", "model", "status"} {
			if !truthy(data[key]) {
				add(d.Path, "Pflichtfeld '%s' fehlt", key)
			}
		}

		if truthy(data["id"]) && !strings.HasSuffix(d.Path, "/"+d.ID+".yaml") {
			add(d.Path, "Dateiname passt nicht zu id '%s'", d.ID)
		}

		class, _ := data["class"].(string)
		if !classes[class] {
			add(d.Path, "unbekannte class '%s'", class)
		}

		status, _ := data["status"].(string)
		if !statuses[status] {
			add(d.Path, "unbekannter status '%s'", status)
		}

		// Regel 1: Revision ist Pflicht.
		if !truthy(data["revision"]) {
			add(d.Path, "revision fehlt — gleiche Modellnummer, andere Platine, anderes Verhalten")
		}

		// Regel 2: keine Herstellerzahlen abschreiben.
		power := subMap(data, "power")
		if _, ok := data["power"]; ok {
			if _, ok := power["idle_w"]; !ok {
				add(d.Path, "power ohne idle_w — lieber null als geraten")
			}
		}
		if power["idle_w"] != nil {
			source, isStr := power["source"].(string)
			if power["source"] == nil || (isStr && source == "unknown") {
				add(d.Path, "idle_w gesetzt, aber source fehlt — woher stammt die Zahl?")
			}
		}
		if source, ok := power["source"].(string); ok && source == "vendor" {
			add(d.Path, "source 'vendor' ist nicht zulaessig, Herstellerangaben stimmen im Leerlauf nie")
		}

		// Regel 3: ehrlich sein, was geprueft wurde.
		if _, ok := data["verified"]; !ok {
			add(d.Path, "verified fehlt — 'by: none' ist eine gueltige und ehrliche Angabe")
		} else if by := subMap(data, "verified")["by"]; !verifiedByOK(by) {
			add(d.Path, "verified.by '%v' unbekannt", by)
		}

		// Empfohlen heisst: wir liefern etwas.
		if (status == "recommended" || status == "bestand") && !truthy(data["generates"]) {
			add(d.Path, "status %s ohne generates — dann ist es hoechstens tolerated", status)
		}
		if (status == "unsupported" || status == "bestand") && !truthy(data["status_reason"]) {
			add(d.Path, "%s ohne status_reason — die Frage 'warum' muss beantwortet sein", status)
		}

		// Ports und PoE.
		ports := subMap(data, "ports")
		portKeys := make([]string, 0, len(ports))
		for k := range ports {
			portKeys = append(portKeys, k)
		}
		sort.Strings(portKeys)
		for _, k := range portKeys {
			if !portTypes[k] {
				add(d.Path, "unbekannter Porttyp '%s'", k)
			}
		}

		poe := subMap(data, "poe")
		if truthy(poe["ports"]) {
			openwrt, _ := poe["openwrt"].(string)
			if !poeOpenwrt[openwrt] {
				add(d.Path, "poe.ports gesetzt, aber poe.openwrt '%s' unbekannt", openwrt)
			}
		}
		if openwrt, _ := poe["openwrt"].(string); openwrt == "manual" || openwrt == "broken" {
			if !truthy(poe["note"]) {
				add(d.Path, "poe.openwrt manual/broken ohne note — was genau geht nicht?")
			}
		}

		// Bewegliche Ziele brauchen ein Datum.
		for _, raw := range subList(data, "caveats") {
			c := fmt.Sprintf("%v", raw)
			low := strings.ToLower(c)
			if !containsAny(low, movingTargetWords) {
				continue
			}
			if containsAny(low, months) {
				continue
			}
			add(d.Path, "caveat ohne Datum: %s…", truncateRunes(c, 60))
		}

		if truthy(data["price_eur"]) && !truthy(data["price_asof"]) {
			add(d.Path, "price_eur ohne price_asof — ein Preis ohne Datum ist keine Angabe")
		}
	}

	// Querverweise
	for _, d := range devices.List {
		for _, raw := range subList(d.Data, "alternatives") {
			a, _ := raw.(string)
			if !ids[a] {
				add(d.Path, "alternatives verweist auf unbekannte id '%s'", a)
			}
		}
	}

	return issues
}

func verifiedByOK(by interface{}) bool {
	if by == nil {
		return true
	}
	s, ok := by.(string)
	return ok && verifiedBy[s]
}

func subMap(d map[string]interface{}, key string) map[string]interface{} {
	m, _ := d[key].(map[string]interface{})
	if m == nil {
		return map[string]interface{}{}
	}
	return m
}

func subList(d map[string]interface{}, key string) []interface{} {
	l, _ := d[key].([]interface{})
	return l
}

// truthy bildet Pythons Wahrheitswert fuer "not d.get(key)" nach: nil, "",
// 0, leere Listen und leere Maps sind falsch, alles andere ist wahr.
func truthy(v interface{}) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case string:
		return x != ""
	case int:
		return x != 0
	case int64:
		return x != 0
	case float64:
		return x != 0
	case []interface{}:
		return len(x) > 0
	case map[string]interface{}:
		return len(x) > 0
	default:
		return true
	}
}

func containsAny(haystack string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
