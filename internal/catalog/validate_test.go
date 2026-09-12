package catalog

import "testing"

func TestValidateRealCatalog(t *testing.T) {
	services, err := Load("../../catalog/services")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if issues := Validate(services, "../.."); len(issues) != 0 {
		t.Errorf("Validate(catalog/services/): erwartete keine Befunde, got:\n%s", issues.Error())
	}
}

// TestValidateBroken vergleicht die Go-Befunde mit dem, was
// tools/check_catalog.py fuer dieselben Eingaben meldet (von Hand
// gegengeprueft). 17 Verstoesse, verteilt auf 8 Dateien.
func TestValidateBroken(t *testing.T) {
	services, err := Load("testdata/broken")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	issues := Validate(services, "testdata/broken")
	if len(issues) != 17 {
		t.Fatalf("Validate: %d Befund(e), want 17:\n%s", len(issues), issues.Error())
	}

	want := []Issue{
		{File: "testdata/broken/host-mode.yaml", Message: "network.mode host ohne host_reason"},
		{File: "testdata/broken/latest-tag.yaml", Message: "tag_policy latest ohne Begruendung in caveats"},
		{File: "testdata/broken/memory.yaml", Message: "resources.memory_source fehlt oder unbekannt"},
		{File: "testdata/broken/no-strategy.yaml", Message: "backup.strategy fehlt"},
		{File: "testdata/broken/ports.yaml", Message: "Port 'notaport' nicht lesbar"},
		{File: "testdata/broken/ports.yaml", Message: "Port 80 unter 1024 — rootless bindet das nicht, nftables leitet weiter"},
		{File: "testdata/broken/rootless-false.yaml", Message: "rootless: false ohne rootless_exception"},
		{File: "testdata/broken/supported-gate.yaml", Message: "supported ohne getestete Wiederherstellung — das entscheidet integrations-tester"},
		{File: "testdata/broken/supported-gate.yaml", Message: "supported ohne vorhandene Doku-Seite (None)"},
		{File: "testdata/broken/unknown-status.yaml", Message: "status 'beta' unbekannt"},
		{File: "testdata/broken/wrong-filename.yaml", Message: "Dateiname passt nicht zu id 'actual-id'"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'id' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'category' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'summary' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'status' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "status '' unbekannt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "backup.strategy fehlt"},
	}
	if len(want) != 17 {
		t.Fatalf("Testfehler: want hat %d Eintraege, sollten 17 sein", len(want))
	}
	for _, w := range want {
		found := false
		for _, got := range issues {
			if got == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("erwarteter Befund fehlt: %s", w)
		}
	}
}
