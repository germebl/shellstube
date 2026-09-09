package manifest

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Manifest ist der Inhalt einer homelab.yaml, geladen und in JSON-kompatible
// Go-Werte umgewandelt (map[string]interface{}, []interface{}, float64,
// string, bool, nil) - das Format, das Validate und spaeter internal/plan
// erwarten.
type Manifest struct {
	Path string
	Data interface{}
}

// Load liest eine YAML-Datei und wandelt sie in JSON-kompatible Werte um.
// Es wird nicht gegen das Schema geprueft - das macht Validate.
func Load(path string) (*Manifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("manifest %s: %w", path, err)
	}

	var doc interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("manifest %s: yaml ungültig: %w", path, err)
	}

	// yaml.v3 liefert Mappings bereits als map[string]interface{}, aber Zahlen
	// je nach Schreibweise als int, int64 oder float64. Der Umweg über JSON
	// normalisiert das auf die Typen, die das Schema-Paket erwartet.
	data, err := jsonRoundtrip(doc)
	if err != nil {
		return nil, fmt.Errorf("manifest %s: %w", path, err)
	}

	return &Manifest{Path: path, Data: data}, nil
}

func jsonRoundtrip(doc interface{}) (interface{}, error) {
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("nicht als JSON darstellbar: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return data, nil
}
