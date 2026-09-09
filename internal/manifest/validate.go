package manifest

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Issue ist ein einzelner Schemaverstoss mit dem Pfad zur fehlerhaften Stelle
// im Manifest, z.B. "/nodes/cmp-01/role".
type Issue struct {
	Path    string
	Message string
}

func (i Issue) String() string {
	if i.Path == "" {
		return i.Message
	}
	return fmt.Sprintf("%s: %s", i.Path, i.Message)
}

// Issues sammelt alle Verstoesse eines Validate-Laufs. Ein leeres Issues
// bedeutet nicht "gueltig" - Validate gibt dafuer nil zurueck.
type Issues []Issue

func (is Issues) Error() string {
	lines := make([]string, len(is))
	for i, issue := range is {
		lines[i] = issue.String()
	}
	return strings.Join(lines, "\n")
}

// Validate prueft m gegen das JSON-Schema unter schemaPath. Bei Verstoessen
// ist der zurueckgegebene Fehler vom Typ Issues, jeder Eintrag mit dem Pfad
// zur betroffenen Stelle im Manifest.
func Validate(m *Manifest, schemaPath string) error {
	schema, err := compileSchema(schemaPath)
	if err != nil {
		return err
	}

	err = schema.Validate(m.Data)
	if err == nil {
		return nil
	}

	ve, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return fmt.Errorf("manifest %s: %w", m.Path, err)
	}
	return issuesFrom(ve)
}

func compileSchema(schemaPath string) (*jsonschema.Schema, error) {
	f, err := os.Open(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("schema %s: %w", schemaPath, err)
	}
	defer f.Close()

	c := jsonschema.NewCompiler()
	if err := c.AddResource(schemaPath, f); err != nil {
		return nil, fmt.Errorf("schema %s: %w", schemaPath, err)
	}
	schema, err := c.Compile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("schema %s: %w", schemaPath, err)
	}
	return schema, nil
}

// issuesFrom flacht den Ursachenbaum von jsonschema ab. ve selbst trägt nur
// die generische Meldung "doesn't validate with <schema>" ohne eigenen Pfad -
// die eigentlichen Meldungen stehen in seinen Causes.
func issuesFrom(ve *jsonschema.ValidationError) Issues {
	var issues Issues
	var walk func(*jsonschema.ValidationError)
	walk = func(ve *jsonschema.ValidationError) {
		if ve.Message != "" {
			issues = append(issues, Issue{Path: ve.InstanceLocation, Message: ve.Message})
		}
		for _, cause := range ve.Causes {
			walk(cause)
		}
	}
	for _, cause := range ve.Causes {
		walk(cause)
	}
	sort.SliceStable(issues, func(i, j int) bool { return issues[i].Path < issues[j].Path })
	return issues
}
