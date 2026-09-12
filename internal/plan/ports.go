package plan

// derivePorts flacht die ports-Liste jedes Switch-Knotens zu einer
// gemeinsamen Tabelle - die Quelle fuer Kabelliste und VLAN-Zuordnung.
func derivePorts(nodes map[string]interface{}) []Port {
	var ports []Port
	for _, name := range sortedKeys(nodes) {
		n := asMap(nodes[name])
		if asString(n["role"]) != "switch" {
			continue
		}
		for _, raw := range asList(n["ports"]) {
			p := asMap(raw)

			var trunk []string
			for _, v := range asList(p["trunk_vlans"]) {
				trunk = append(trunk, asString(v))
			}

			mode := asString(p["mode"])
			if mode == "" {
				mode = "access"
			}

			ports = append(ports, Port{
				Switch:     name,
				N:          asInt(p["n"]),
				Kind:       asString(p["kind"]),
				Mode:       mode,
				VLAN:       asString(p["vlan"]),
				TrunkVLANs: trunk,
				Peer:       asString(p["peer"]),
				PoE:        asBool(p["poe"]),
				Label:      asString(p["label"]),
			})
		}
	}
	return ports
}
