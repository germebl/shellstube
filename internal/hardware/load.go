package hardware

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Device ist ein Eintrag aus hardware/*.yaml, so wie er auf der Platte steht.
// Data traegt das gesamte Dokument als generische Werte - das Schema ist noch
// im Fluss, ein fester Go-Typ wuerde jede neue Eigenschaft zu einer
// Code-Aenderung machen.
type Device struct {
	Path string
	ID   string
	Data map[string]interface{}
}

// Devices ist die Hardwareliste, geladen und nach id indiziert.
type Devices struct {
	List []*Device
	byID map[string]*Device
}

// ByID sucht ein Geraet ueber seine id.
func (ds *Devices) ByID(id string) (*Device, bool) {
	d, ok := ds.byID[id]
	return d, ok
}

// Load liest jede Datei dir/*.yaml ein, in derselben sortierten Reihenfolge
// wie tools/check_hardware.py (sorted(glob.glob(...))).
func Load(dir string) (*Devices, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("hardware %s: %w", dir, err)
	}
	sort.Strings(paths)

	ds := &Devices{byID: map[string]*Device{}}
	for _, path := range paths {
		dev, err := loadDevice(path)
		if err != nil {
			return nil, err
		}
		ds.List = append(ds.List, dev)
		if dev.ID != "" {
			ds.byID[dev.ID] = dev
		}
	}
	return ds, nil
}

func loadDevice(path string) (*Device, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("hardware %s: %w", path, err)
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("hardware %s: yaml ungültig: %w", path, err)
	}
	if doc == nil {
		doc = map[string]interface{}{}
	}

	id, _ := doc["id"].(string)
	return &Device{Path: path, ID: id, Data: doc}, nil
}
