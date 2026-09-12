package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Service ist ein Eintrag aus catalog/services/*.yaml, so wie er auf der
// Platte steht. Data traegt das gesamte Dokument als generische Werte - das
// Schema ist noch im Fluss, ein fester Go-Typ wuerde jede neue Eigenschaft zu
// einer Code-Aenderung machen.
type Service struct {
	Path string
	ID   string
	Data map[string]interface{}
}

// Services ist der Dienstkatalog, geladen und nach id indiziert.
type Services struct {
	List []*Service
	byID map[string]*Service
}

// ByID sucht einen Dienst ueber seine id.
func (ss *Services) ByID(id string) (*Service, bool) {
	s, ok := ss.byID[id]
	return s, ok
}

// Load liest jede Datei dir/*.yaml ein, ausser solchen, deren Name mit "_"
// beginnt (die Vorlage) - in derselben sortierten Reihenfolge wie
// tools/check_catalog.py (sorted(glob.glob(...))).
func Load(dir string) (*Services, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("catalog %s: %w", dir, err)
	}
	sort.Strings(paths)

	ss := &Services{byID: map[string]*Service{}}
	for _, path := range paths {
		if strings.HasPrefix(filepath.Base(path), "_") {
			continue
		}
		svc, err := loadService(path)
		if err != nil {
			return nil, err
		}
		ss.List = append(ss.List, svc)
		if svc.ID != "" {
			ss.byID[svc.ID] = svc
		}
	}
	return ss, nil
}

func loadService(path string) (*Service, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("catalog %s: %w", path, err)
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("catalog %s: yaml ungültig: %w", path, err)
	}
	if doc == nil {
		doc = map[string]interface{}{}
	}

	id, _ := doc["id"].(string)
	return &Service{Path: path, ID: id, Data: doc}, nil
}
