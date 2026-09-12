package plan

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Render stellt den Plan als Tabellen dar, ein Abschnitt je Ableitung.
func Render(p *Plan) string {
	var b strings.Builder
	renderRoles(&b, p.Roles)
	renderPorts(&b, p.Ports)
	renderVLANs(&b, p.VLANs)
	renderLinks(&b, p.Links)
	renderMemory(&b, p.Memory)
	return b.String()
}

func newSection(b *strings.Builder, title string) *tabwriter.Writer {
	fmt.Fprintf(b, "== %s ==\n", title)
	return tabwriter.NewWriter(b, 0, 2, 2, ' ', 0)
}

func renderRoles(b *strings.Builder, roles []NodeRole) {
	w := newSection(b, "Rollen")
	fmt.Fprintln(w, "KNOTEN\tROLLE\tHARDWARE\tVENDOR\tMODELL\tSTATUS\tHINWEIS")
	for _, r := range roles {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", r.Node, r.Role, dash(r.Hardware), dash(r.Vendor), dash(r.Model), dash(r.Status), r.Note)
	}
	w.Flush()
	b.WriteString("\n")
}

func renderPorts(b *strings.Builder, ports []Port) {
	w := newSection(b, "Ports")
	if len(ports) == 0 {
		fmt.Fprintln(w, "(kein Switch mit Portbelegung im Manifest)")
	} else {
		fmt.Fprintln(w, "SWITCH\tPORT\tART\tMODUS\tVLAN\tPEER\tPOE\tLABEL")
		for _, p := range ports {
			vlan := p.VLAN
			if p.Mode == "trunk" {
				vlan = strings.Join(p.TrunkVLANs, ",")
			}
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", p.Switch, p.N, p.Kind, p.Mode, dash(vlan), dash(p.Peer), yesNo(p.PoE), p.Label)
		}
	}
	w.Flush()
	b.WriteString("\n")
}

func renderVLANs(b *strings.Builder, vlans []VLAN) {
	w := newSection(b, "VLANs")
	if len(vlans) == 0 {
		fmt.Fprintln(w, "(keine VLANs im Manifest)")
	} else {
		fmt.Fprintln(w, "NAME\tID\tSUBNET\tEGRESS\tBENUTZT VON\tKOMMENTAR")
		for _, v := range vlans {
			used := strings.Join(v.UsedBy, ", ")
			if used == "" {
				used = "(nichts, vermutlich vergessen)"
			}
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\n", v.Name, v.ID, v.Subnet, v.Egress, used, v.Comment)
		}
	}
	w.Flush()
	b.WriteString("\n")
}

func renderLinks(b *strings.Builder, links []Link) {
	w := newSection(b, "Strecke")
	if len(links) == 0 {
		fmt.Fprintln(w, "(keine Links im Manifest)")
	} else {
		fmt.Fprintln(w, "A\tB\tGESCHW. GBE\tNOETIG GBE\tBEFUND\tHINWEIS")
		for _, l := range links {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", l.A, l.B, gbe(l.SpeedGbe), gbe(l.RequiredGbe), linkVerdict(l), l.Note)
		}
	}
	w.Flush()
	b.WriteString("\n")
}

func linkVerdict(l Link) string {
	switch {
	case l.Unknown:
		return "unklar"
	case l.RequiredGbe == 0:
		return "kein target_speed_gbe gesetzt"
	case l.Bottleneck:
		return "ENGSTELLE"
	default:
		return "ok"
	}
}

func renderMemory(b *strings.Builder, mem []Memory) {
	w := newSection(b, "Speicher")
	if len(mem) == 0 {
		fmt.Fprintln(w, "(kein Compute-Knoten im Manifest)")
	} else {
		fmt.Fprintln(w, "KNOTEN\tKAPAZITAET MB\tBEDARF MB\tBEFUND\tHINWEIS")
		for _, m := range mem {
			verdict := "ok"
			if m.Bottleneck {
				verdict = "ENGSTELLE"
			}
			fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\n", m.Node, m.CapacityMB, m.RequiredMB, verdict, m.Note)
		}
	}
	w.Flush()
	b.WriteString("\n")
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func yesNo(b bool) string {
	if b {
		return "ja"
	}
	return "nein"
}

func gbe(v float64) string {
	if v == 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f", v)
}
