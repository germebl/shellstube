package plan

import (
	"fmt"
	"math"

	"github.com/shellstube/shellstube/internal/hardware"
)

// portSpeedGbe uebersetzt einen Porttyp in seine Geschwindigkeit. sfp ohne
// Zusatz wird als 1G gefuehrt - das haeufigste Modul, alles schnellere traegt
// im Manifest sfp_plus oder qsfp.
var portSpeedGbe = map[string]float64{
	"rj45_1g":  1,
	"rj45_2g5": 2.5,
	"rj45_10g": 10,
	"sfp":      1,
	"sfp_plus": 10,
	"qsfp":     40,
}

// deriveLinks prueft jede Verbindung aus manifest.links gegen
// target_speed_gbe. Steht wan.path auf "switch", traegt die Leitung
// zwischen Switch und Router den Internetverkehr doppelt - einmal roh
// hinein, einmal gefiltert zurueck zu den Geraeten hinter dem Switch. Diese
// Leitung braucht deshalb die doppelte Zielgeschwindigkeit.
func deriveLinks(data map[string]interface{}, nodes map[string]interface{}, ports []Port, devices *hardware.Devices) []Link {
	raw := asList(data["links"])
	if len(raw) == 0 {
		return nil
	}

	target := asFloat(data["target_speed_gbe"])
	wanPath := asString(asMap(data["wan"])["path"])
	if wanPath == "" {
		wanPath = "router"
	}

	links := make([]Link, 0, len(raw))
	for _, item := range raw {
		pair := asList(item)
		if len(pair) != 2 {
			continue
		}
		a, b := asString(pair[0]), asString(pair[1])
		na, nb := asMap(nodes[a]), asMap(nodes[b])

		speedA, noteA := nodeCapacity(a, na, b, ports, devices)
		speedB, noteB := nodeCapacity(b, nb, a, ports, devices)
		speed := math.Min(speedA, speedB)
		unknown := speedA == 0 || speedB == 0

		doubled := wanPath == "switch" && isSwitchRouterPair(asString(na["role"]), asString(nb["role"]))
		required := target
		if doubled {
			required = target * 2
		}

		link := Link{
			A: a, B: b,
			SpeedGbe:    speed,
			RequiredGbe: required,
			Doubled:     doubled,
			Unknown:     unknown,
			Note:        joinNotes(noteA, noteB),
		}
		if doubled {
			link.Note = joinNotes(link.Note, "WAN-Verkehr laeuft hier doppelt: roh hinein, gefiltert zurueck")
		}
		if !unknown && target > 0 && speed < required {
			link.Bottleneck = true
		}
		links = append(links, link)
	}
	return links
}

func isSwitchRouterPair(roleA, roleB string) bool {
	return (roleA == "switch" && roleB == "router") || (roleA == "router" && roleB == "switch")
}

// nodeCapacity ermittelt, mit welcher Geschwindigkeit name an peer
// angebunden ist. Bei einem Switch zaehlt der Port, der peer als peer
// eintraegt; sonst link_gbe, und wenn das fehlt, der schnellste Porttyp aus
// dem Hardware-Datenblatt.
func nodeCapacity(name string, n map[string]interface{}, peer string, ports []Port, devices *hardware.Devices) (float64, string) {
	if asString(n["role"]) == "switch" {
		for _, p := range ports {
			if p.Switch != name || p.Peer != peer {
				continue
			}
			if speed, ok := portSpeedGbe[p.Kind]; ok {
				return speed, ""
			}
			return 0, fmt.Sprintf("%s: unbekannter Porttyp '%s'", name, p.Kind)
		}
		return 0, fmt.Sprintf("%s: keine Portzuordnung zu '%s' im Manifest", name, peer)
	}

	if lg, ok := n["link_gbe"]; ok {
		return asFloat(lg), ""
	}

	if dev, ok := devices.ByID(asString(n["hardware"])); ok {
		if speed, ok := maxPortSpeed(dev.Data); ok {
			return speed, fmt.Sprintf("%s: link_gbe fehlt, aus Hardwareliste geschaetzt", name)
		}
	}
	return 0, fmt.Sprintf("%s: keine Geschwindigkeit ermittelbar", name)
}

func maxPortSpeed(data map[string]interface{}) (float64, bool) {
	var max float64
	found := false
	for kind, count := range asMap(data["ports"]) {
		if asFloat(count) <= 0 {
			continue
		}
		if speed, ok := portSpeedGbe[kind]; ok {
			found = true
			if speed > max {
				max = speed
			}
		}
	}
	return max, found
}

func joinNotes(notes ...string) string {
	out := ""
	for _, n := range notes {
		if n == "" {
			continue
		}
		if out != "" {
			out += "; "
		}
		out += n
	}
	return out
}
