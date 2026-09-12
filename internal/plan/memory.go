package plan

import (
	"strings"

	"github.com/shellstube/shellstube/internal/catalog"
)

// deriveMemory rechnet den Speicherbedarf der gewaehlten Dienste gegen die
// Compute-Knoten. Das Manifest kennt noch keine Zuordnung Dienst-zu-Knoten -
// solange es nur einen Compute-Knoten gibt, ist das eindeutig; bei mehreren
// nimmt der erste (nach Namen) alle Dienste, mit einem Hinweis darauf.
func deriveMemory(data map[string]interface{}, nodes map[string]interface{}, services *catalog.Services) []Memory {
	var computeNodes []string
	for _, name := range sortedKeys(nodes) {
		if asString(asMap(nodes[name])["runtime"]) == "podman-rootless" {
			computeNodes = append(computeNodes, name)
		}
	}
	if len(computeNodes) == 0 {
		return nil
	}

	requiredMB, unresolved := serviceMemoryMB(asList(data["services"]), services)

	result := make([]Memory, 0, len(computeNodes))
	for i, name := range computeNodes {
		n := asMap(nodes[name])
		m := Memory{Node: name, CapacityMB: asInt(asFloat(n["memory_gb"]) * 1024)}

		if i == 0 {
			m.RequiredMB = requiredMB
			if len(computeNodes) > 1 {
				m.Note = "mehrere Compute-Knoten, Zuordnung fehlt im Manifest - Dienste hier angenommen"
			}
			if len(unresolved) > 0 {
				m.Note = joinNotes(m.Note, "unbekannt im Katalog: "+strings.Join(unresolved, ", "))
			}
		} else {
			m.Note = "keine Dienste zugeordnet, Manifest kennt keine Knotenzuordnung"
		}

		m.Bottleneck = m.CapacityMB > 0 && m.RequiredMB > m.CapacityMB
		result = append(result, m)
	}
	return result
}

func serviceMemoryMB(ids []interface{}, services *catalog.Services) (int, []string) {
	total := 0
	var unresolved []string
	for _, raw := range ids {
		id := asString(raw)
		svc, ok := services.ByID(id)
		if !ok {
			unresolved = append(unresolved, id)
			continue
		}
		total += asInt(asMap(svc.Data["resources"])["memory_mb"])
	}
	return total, unresolved
}
