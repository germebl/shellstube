package hardware

import "testing"

func TestValidateRealHardware(t *testing.T) {
	devices, err := Load("../../hardware")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if issues := Validate(devices); len(issues) != 0 {
		t.Errorf("Validate(hardware/): erwartete keine Befunde, got:\n%s", issues.Error())
	}
}

// TestValidateBroken vergleicht die Go-Befunde mit dem, was
// tools/check_hardware.py fuer dieselben Eingaben meldet (von Hand
// gegengeprueft). 23 Verstoesse, verteilt auf 14 Dateien.
func TestValidateBroken(t *testing.T) {
	devices, err := Load("testdata/broken")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	issues := Validate(devices)
	if len(issues) != 23 {
		t.Fatalf("Validate: %d Befund(e), want 23:\n%s", len(issues), issues.Error())
	}

	want := []Issue{
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'id' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'class' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'vendor' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "Pflichtfeld 'status' fehlt"},
		{File: "testdata/broken/zz-missing.yaml", Message: "unbekannte class ''"},
		{File: "testdata/broken/zz-missing.yaml", Message: "unbekannter status ''"},
		{File: "testdata/broken/zz-missing.yaml", Message: "revision fehlt — gleiche Modellnummer, andere Platine, anderes Verhalten"},
		{File: "testdata/broken/zz-missing.yaml", Message: "verified fehlt — 'by: none' ist eine gueltige und ehrliche Angabe"},
		{File: "testdata/broken/wrong-filename.yaml", Message: "Dateiname passt nicht zu id 'actual-id'"},
		{File: "testdata/broken/power-issues.yaml", Message: "source 'vendor' ist nicht zulaessig, Herstellerangaben stimmen im Leerlauf nie"},
		{File: "testdata/broken/power-no-idle.yaml", Message: "power ohne idle_w — lieber null als geraten"},
		{File: "testdata/broken/power-idle-no-source.yaml", Message: "idle_w gesetzt, aber source fehlt — woher stammt die Zahl?"},
		{File: "testdata/broken/verified-bad.yaml", Message: "verified.by 'vendor' unbekannt"},
		{File: "testdata/broken/status-generates.yaml", Message: "status recommended ohne generates — dann ist es hoechstens tolerated"},
		{File: "testdata/broken/status-reason-missing.yaml", Message: "unsupported ohne status_reason — die Frage 'warum' muss beantwortet sein"},
		{File: "testdata/broken/bestand-both.yaml", Message: "status bestand ohne generates — dann ist es hoechstens tolerated"},
		{File: "testdata/broken/bestand-both.yaml", Message: "bestand ohne status_reason — die Frage 'warum' muss beantwortet sein"},
		{File: "testdata/broken/ports-poe.yaml", Message: "unbekannter Porttyp 'weird_port'"},
		{File: "testdata/broken/ports-poe.yaml", Message: "poe.ports gesetzt, aber poe.openwrt 'crazy' unbekannt"},
		{File: "testdata/broken/poe-manual-no-note.yaml", Message: "poe.openwrt manual/broken ohne note — was genau geht nicht?"},
		{File: "testdata/broken/caveat-no-date.yaml", Message: "caveat ohne Datum: Das Problem tritt mittlerweile seltener auf, aber ohne Zeita…"},
		{File: "testdata/broken/price-no-asof.yaml", Message: "price_eur ohne price_asof — ein Preis ohne Datum ist keine Angabe"},
		{File: "testdata/broken/alt-unknown.yaml", Message: "alternatives verweist auf unbekannte id 'does-not-exist'"},
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
