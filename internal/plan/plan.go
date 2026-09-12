package plan

import (
	"fmt"

	"github.com/shellstube/shellstube/internal/catalog"
	"github.com/shellstube/shellstube/internal/hardware"
	"github.com/shellstube/shellstube/internal/manifest"
)

// Derive leitet aus einem geladenen Manifest, der Hardwareliste und dem
// Dienstkatalog alles ab, was sich ableiten laesst. Das Manifest sollte
// vorher gegen das Schema geprueft sein - Derive verlaesst sich auf dessen
// Form, prueft sie aber nicht erneut.
func Derive(m *manifest.Manifest, devices *hardware.Devices, services *catalog.Services) (*Plan, error) {
	data, ok := m.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("plan %s: unerwartetes Format, ist das Manifest geprueft?", m.Path)
	}
	nodes := asMap(data["nodes"])

	p := &Plan{}
	p.Roles = deriveRoles(nodes, devices)
	p.Ports = derivePorts(nodes)
	p.VLANs = deriveVLANs(asMap(data["vlans"]), nodes, p.Ports)
	p.Links = deriveLinks(data, nodes, p.Ports, devices)
	p.Memory = deriveMemory(data, nodes, services)
	return p, nil
}
