package catalog

import "testing"

func TestLoadRealList(t *testing.T) {
	services, err := Load("../../catalog/services")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(services.List) == 0 {
		t.Fatal("Load: keine Dienste gefunden")
	}

	s, ok := services.ByID("adguard-home")
	if !ok {
		t.Fatal("ByID(adguard-home): nicht gefunden")
	}
	if got := s.Data["category"]; got != "dns" {
		t.Errorf("category = %v, want dns", got)
	}
}

func TestLoadSkipsTemplate(t *testing.T) {
	services, err := Load("testdata/broken")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, s := range services.List {
		if s.Path == "testdata/broken/_template.yaml" {
			t.Fatal("Load: _template.yaml wurde nicht uebersprungen")
		}
	}
}
