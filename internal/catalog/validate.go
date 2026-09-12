package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Issue ist ein einzelner Befund zu einer Datei aus catalog/services/*.yaml.
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

// Validate bildet die Regeln aus tools/check_catalog.py nach: dieselben
// Eingaben muessen dieselben Befunde liefern. root ist das Verzeichnis, gegen
// das relative docs-Pfade aufgeloest werden (bei einem Lauf aus dem
// Projektwurzelverzeichnis ".").
func Validate(services *Services, root string) Issues {
	var issues Issues
	add := func(path, format string, args ...interface{}) {
		issues = append(issues, Issue{File: path, Message: fmt.Sprintf(format, args...)})
	}

	for _, s := range services.List {
		data := s.Data

		for _, key := range []string{"id", "name", "category", "summary", "status"} {
			if !truthy(data[key]) {
				add(s.Path, "Pflichtfeld '%s' fehlt", key)
			}
		}
		if truthy(data["id"]) && !strings.HasSuffix(s.Path, "/"+s.ID+".yaml") {
			add(s.Path, "Dateiname passt nicht zu id '%s'", s.ID)
		}
		status, _ := data["status"].(string)
		if status != "supported" && status != "experimental" {
			add(s.Path, "status '%s' unbekannt", status)
		}

		// Rootless ist der Standard, Abweichung braucht einen Grund.
		if rootless, ok := data["rootless"].(bool); ok && !rootless && !truthy(data["rootless_exception"]) {
			add(s.Path, "rootless: false ohne rootless_exception")
		}

		// Keine privilegierten Ports.
		net := subMap(data, "network")
		for _, raw := range subList(net, "ports") {
			p := fmt.Sprintf("%v", raw)
			num, err := strconv.Atoi(strings.SplitN(p, "/", 2)[0])
			if err != nil {
				add(s.Path, "Port '%s' nicht lesbar", p)
				continue
			}
			if num < 1024 {
				add(s.Path, "Port %d unter 1024 — rootless bindet das nicht, nftables leitet weiter", num)
			}
		}
		if mode, _ := net["mode"].(string); mode == "host" && !truthy(net["host_reason"]) {
			add(s.Path, "network.mode host ohne host_reason")
		}

		// Das Qualitaetstor.
		bk := subMap(data, "backup")
		if status == "supported" {
			if !truthy(bk["restore_tested"]) {
				add(s.Path, "supported ohne getestete Wiederherstellung — das entscheidet integrations-tester")
			}
			docs, _ := data["docs"].(string)
			if !truthy(data["docs"]) || !fileExists(filepath.Join(root, docs)) {
				add(s.Path, "supported ohne vorhandene Doku-Seite (%s)", pyStr(data["docs"]))
			}
		}
		if !truthy(bk["strategy"]) {
			add(s.Path, "backup.strategy fehlt")
		}

		// Ressourcen ehrlich kennzeichnen.
		res := subMap(data, "resources")
		if truthy(res["memory_mb"]) {
			src, ok := res["memory_source"].(string)
			if !ok || (src != "measured" && src != "upstream" && src != "estimate") {
				add(s.Path, "resources.memory_source fehlt oder unbekannt")
			}
		}

		if tagPolicy, _ := subMap(data, "image")["tag_policy"].(string); tagPolicy == "latest" && !truthy(data["caveats"]) {
			add(s.Path, "tag_policy latest ohne Begruendung in caveats")
		}
	}

	return issues
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

func pyStr(v interface{}) string {
	if v == nil {
		return "None"
	}
	return fmt.Sprintf("%v", v)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
