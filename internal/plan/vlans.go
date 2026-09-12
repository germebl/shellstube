package plan

import "fmt"

// deriveVLANs liest den vlans-Block und ergaenzt fuer jeden Eintrag, welche
// Knoten und Ports ihn tatsaechlich benutzen - ein VLAN ohne Benutzer ist
// meist ein Hinweis auf ein vergessenes Kabel oder einen Tippfehler.
func deriveVLANs(vlans map[string]interface{}, nodes map[string]interface{}, ports []Port) []VLAN {
	result := make([]VLAN, 0, len(vlans))
	for _, name := range sortedKeys(vlans) {
		v := asMap(vlans[name])

		egress := asString(v["egress"])
		if egress == "" {
			egress = "filtered"
		}

		result = append(result, VLAN{
			Name:    name,
			ID:      asInt(v["id"]),
			Subnet:  asString(v["subnet"]),
			Egress:  egress,
			Comment: asString(v["comment"]),
			UsedBy:  usersOfVLAN(name, nodes, ports),
		})
	}
	return result
}

func usersOfVLAN(name string, nodes map[string]interface{}, ports []Port) []string {
	seen := map[string]bool{}
	var users []string
	add := func(s string) {
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		users = append(users, s)
	}

	for _, nodeName := range sortedKeys(nodes) {
		if asString(asMap(nodes[nodeName])["vlan"]) == name {
			add(nodeName)
		}
	}

	for _, p := range ports {
		if p.VLAN == name {
			add(fmt.Sprintf("%s:%d", p.Switch, p.N))
			continue
		}
		for _, t := range p.TrunkVLANs {
			if t == name {
				add(fmt.Sprintf("%s:%d", p.Switch, p.N))
				break
			}
		}
	}
	return users
}
