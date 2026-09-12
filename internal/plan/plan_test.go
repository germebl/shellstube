package plan

import (
	"strings"
	"testing"

	"github.com/shellstube/shellstube/internal/catalog"
	"github.com/shellstube/shellstube/internal/hardware"
	"github.com/shellstube/shellstube/internal/manifest"
)

func derive(t *testing.T, manifestPath string) *Plan {
	t.Helper()
	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatalf("manifest.Load: %v", err)
	}
	devices, err := hardware.Load("../../hardware")
	if err != nil {
		t.Fatalf("hardware.Load: %v", err)
	}
	services, err := catalog.Load("../../catalog/services")
	if err != nil {
		t.Fatalf("catalog.Load: %v", err)
	}
	p, err := Derive(m, devices, services)
	if err != nil {
		t.Fatalf("Derive: %v", err)
	}
	return p
}

func TestDeriveWanOverSwitchIsBottleneck(t *testing.T) {
	p := derive(t, "../../testdata/manifests/t2-wan-ueber-switch.yaml")

	var link *Link
	for i, l := range p.Links {
		if (l.A == "sw-01" && l.B == "rtr-01") || (l.A == "rtr-01" && l.B == "sw-01") {
			link = &p.Links[i]
		}
	}
	if link == nil {
		t.Fatal("Links: keine Verbindung sw-01/rtr-01 gefunden")
	}
	if !link.Doubled {
		t.Error("Doubled = false, want true - wan.path steht auf switch")
	}
	if !link.Bottleneck {
		t.Errorf("Bottleneck = false, want true (2.5 Gbit/s Leitung, 2x2.5 = 5 Gbit/s noetig)")
	}
	if link.RequiredGbe != 5 {
		t.Errorf("RequiredGbe = %v, want 5 (2x target_speed_gbe)", link.RequiredGbe)
	}
}

func TestRenderMentionsEngstelle(t *testing.T) {
	p := derive(t, "../../testdata/manifests/t2-wan-ueber-switch.yaml")
	out := Render(p)
	if !strings.Contains(strings.ToLower(out), "engstelle") {
		t.Errorf("Render: erwartete 'Engstelle' im Ausgabetext, got:\n%s", out)
	}
}

func TestDeriveMinimalHasNoLinks(t *testing.T) {
	p := derive(t, "../../testdata/manifests/t0-minimal.yaml")
	if len(p.Links) != 0 {
		t.Errorf("Links = %d Eintraege, want 0 (kein links-Block im Manifest)", len(p.Links))
	}
	if len(p.Roles) != 1 || p.Roles[0].Node != "cmp-01" {
		t.Fatalf("Roles = %+v, want ein Eintrag cmp-01", p.Roles)
	}
	if p.Roles[0].Status != "recommended" {
		t.Errorf("Status = %s, want recommended (aus hardware/x86-mini-n150.yaml)", p.Roles[0].Status)
	}
}

func TestDeriveMemoryFitsMinimal(t *testing.T) {
	p := derive(t, "../../testdata/manifests/t0-minimal.yaml")
	if len(p.Memory) != 1 {
		t.Fatalf("Memory = %d Eintraege, want 1", len(p.Memory))
	}
	mem := p.Memory[0]
	if mem.Bottleneck {
		t.Errorf("Bottleneck = true, want false (16 GB gegen adguard-home mit 150 MB)")
	}
	if mem.RequiredMB != 150 {
		t.Errorf("RequiredMB = %d, want 150 (adguard-home resources.memory_mb)", mem.RequiredMB)
	}
}

func TestDeriveVLANsTrackUsers(t *testing.T) {
	p := derive(t, "../../testdata/manifests/t2-wan-ueber-switch.yaml")
	for _, v := range p.VLANs {
		if v.Name == "srv" {
			if len(v.UsedBy) == 0 {
				t.Error("srv: UsedBy ist leer, erwarte cmp-01 und Switch-Ports")
			}
			return
		}
	}
	t.Fatal("VLANs: srv nicht gefunden")
}
