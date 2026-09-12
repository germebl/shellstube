package plan

import (
	"fmt"

	"github.com/shellstube/shellstube/internal/hardware"
)

// deriveRoles ordnet jedem Knoten aus dem Manifest die Hardware zu, auf die
// er sich beruft - inklusive Status, wenn die Hardwareliste ihn kennt.
func deriveRoles(nodes map[string]interface{}, devices *hardware.Devices) []NodeRole {
	names := sortedKeys(nodes)
	roles := make([]NodeRole, 0, len(names))
	for _, name := range names {
		n := asMap(nodes[name])
		r := NodeRole{Node: name, Role: asString(n["role"]), Hardware: asString(n["hardware"])}

		switch {
		case r.Hardware == "":
			r.Note = "keine Hardware angegeben"
		default:
			if dev, ok := devices.ByID(r.Hardware); ok {
				r.Vendor = asString(dev.Data["vendor"])
				r.Model = asString(dev.Data["model"])
				r.Status = asString(dev.Data["status"])
			} else {
				r.Note = fmt.Sprintf("hardware '%s' nicht in der Hardwareliste gefunden", r.Hardware)
			}
		}
		roles = append(roles, r)
	}
	return roles
}
