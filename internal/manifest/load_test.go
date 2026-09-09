package manifest

import "testing"

func TestLoadMinimal(t *testing.T) {
	m, err := Load("../../testdata/manifests/t0-minimal.yaml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	data, ok := m.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data ist %T, kein map[string]interface{}", m.Data)
	}
	if got := data["tier"]; got != "t0" {
		t.Errorf("tier = %v, want t0", got)
	}
	// yaml.v3 liest "version: 1" als int, jsonRoundtrip muss daraus float64
	// machen - das Format, das jsonschema erwartet.
	if got, ok := data["version"].(float64); !ok || got != 1 {
		t.Errorf("version = %v (%T), want float64(1)", data["version"], data["version"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load("testdata/does-not-exist.yaml"); err == nil {
		t.Fatal("Load: want error for missing file, got nil")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	if _, err := Load("testdata/invalid-syntax.yaml"); err == nil {
		t.Fatal("Load: want error for invalid yaml, got nil")
	}
}
