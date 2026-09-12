package hardware

import "testing"

func TestLoadRealList(t *testing.T) {
	devices, err := Load("../../hardware")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(devices.List) == 0 {
		t.Fatal("Load: keine Geraete gefunden")
	}

	d, ok := devices.ByID("zyxel-xgs1210-12")
	if !ok {
		t.Fatal("ByID(zyxel-xgs1210-12): nicht gefunden")
	}
	if got := d.Data["vendor"]; got != "Zyxel" {
		t.Errorf("vendor = %v, want Zyxel", got)
	}
}

func TestLoadMissingDir(t *testing.T) {
	devices, err := Load("testdata/does-not-exist")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(devices.List) != 0 {
		t.Fatalf("Load: erwartete leere Liste, got %d", len(devices.List))
	}
}
