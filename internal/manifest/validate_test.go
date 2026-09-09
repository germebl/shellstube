package manifest

import (
	"strings"
	"testing"
)

const schemaPath = "../../schema/homelab.schema.json"

func TestValidateExampleManifests(t *testing.T) {
	paths := []string{
		"../../testdata/manifests/t0-minimal.yaml",
		"../../testdata/manifests/t2-wan-ueber-switch.yaml",
	}
	for _, path := range paths {
		m, err := Load(path)
		if err != nil {
			t.Fatalf("Load(%s): %v", path, err)
		}
		if err := Validate(m, schemaPath); err != nil {
			t.Errorf("Validate(%s): %v", path, err)
		}
	}
}

func TestValidateBrokenManifest(t *testing.T) {
	m, err := Load("testdata/broken.yaml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	err = Validate(m, schemaPath)
	if err == nil {
		t.Fatal("Validate: want error for broken manifest, got nil")
	}

	issues, ok := err.(Issues)
	if !ok {
		t.Fatalf("Validate: err ist %T, keine Issues", err)
	}
	if len(issues) == 0 {
		t.Fatal("Validate: keine Issues gemeldet")
	}

	want := map[string]string{
		"/tier":              "one of",
		"/nodes/cmp-01/role": "one of",
	}
	for path, wantSubstr := range want {
		found := false
		for _, issue := range issues {
			if issue.Path != path {
				continue
			}
			found = true
			if !strings.Contains(issue.Message, wantSubstr) {
				t.Errorf("issue %s: message %q enthält nicht %q", path, issue.Message, wantSubstr)
			}
		}
		if !found {
			t.Errorf("kein Issue mit Pfad %s gefunden, hatte: %v", path, issues)
		}
	}

	// Jedes Issue muss den Pfad zur fehlerhaften Stelle nennen, nicht nur
	// "ungültig" - genau das ist der Zweck von Issues.
	for _, issue := range issues {
		if issue.String() == issue.Message {
			continue // Wurzel-Verstoss ohne Pfad ist zulaessig (additionalProperties auf oberster Ebene)
		}
		if !strings.HasPrefix(issue.String(), issue.Path) {
			t.Errorf("issue.String() = %q, enthält nicht den Pfad %q", issue.String(), issue.Path)
		}
	}
}

func TestValidateMissingSchema(t *testing.T) {
	m, err := Load("../../testdata/manifests/t0-minimal.yaml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := Validate(m, "testdata/does-not-exist.json"); err == nil {
		t.Fatal("Validate: want error for missing schema, got nil")
	}
}
