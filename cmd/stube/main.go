// Command stube liest ein homelab.yaml und schreibt daraus lesbare Konfigurationsdateien.
//
// Es gibt keinen Daemon und keinen Fernzugriff. Alles, was stube tut, landet als Datei
// auf der Platte — und ist damit auch ohne stube nachvollziehbar.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shellstube/shellstube/internal/catalog"
	"github.com/shellstube/shellstube/internal/hardware"
	"github.com/shellstube/shellstube/internal/manifest"
)

var version = "dev" // wird beim Bauen gesetzt

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "validate":
		os.Exit(cmdValidate(os.Args[2:]))
	case "plan":
		os.Exit(cmdPlan(os.Args[2:]))
	case "apply":
		os.Exit(cmdApply(os.Args[2:]))
	case "remove":
		os.Exit(cmdRemove(os.Args[2:]))
	case "hardware":
		os.Exit(cmdHardware(os.Args[2:]))
	case "catalog":
		os.Exit(cmdCatalog(os.Args[2:]))
	case "version":
		fmt.Println(version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `stube — Homelab aus einem Manifest

  stube validate <manifest>   Manifest gegen das Schema prüfen
  stube plan     <manifest>   Rollen, Ports, VLANs und Strecke ableiten und zeigen
  stube apply    <manifest>   Artefakte schreiben. Ohne --apply nur ein Trockenlauf.
  stube remove   <manifest>   Erzeugtes wieder entfernen, vollständig
  stube hardware list         id je Zeile, aus hardware/*.yaml
  stube catalog  list         id je Zeile, aus catalog/services/*.yaml
  stube version

Der Trockenlauf ist die Vorgabe. Wer wirklich schreiben will, sagt --apply.
`)
}

// cmdValidate prüft ein oder mehrere Manifeste, ohne irgendetwas zu verändern.
func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	schema := fs.String("schema", "schema/homelab.schema.json", "Pfad zum Schema")
	_ = fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "validate: mindestens ein Manifest angeben")
		return 2
	}

	ok := true
	for _, path := range paths {
		if err := validateOne(path, *schema); err != nil {
			fmt.Fprintln(os.Stderr, err)
			ok = false
			continue
		}
		fmt.Printf("%s: gültig\n", path)
	}
	if !ok {
		return 1
	}
	return 0
}

// validateOne lädt und prüft ein einzelnes Manifest. Bei Verstößen listet die
// Fehlermeldung jeden davon mit dem Pfad zur betroffenen Stelle auf.
func validateOne(path, schemaPath string) error {
	m, err := manifest.Load(path)
	if err != nil {
		return err
	}
	if err := manifest.Validate(m, schemaPath); err != nil {
		issues, ok := err.(manifest.Issues)
		if !ok {
			return fmt.Errorf("%s: %w", path, err)
		}
		lines := make([]string, len(issues))
		for i, issue := range issues {
			lines[i] = "  " + issue.String()
		}
		return fmt.Errorf("%s: ungültig\n%s", path, strings.Join(lines, "\n"))
	}
	return nil
}

func cmdPlan(args []string) int {
	// TODO(M2): plan.Derive und Ausgabe als Tabelle
	fmt.Fprintln(os.Stderr, "noch nicht gebaut")
	return 1
}

// cmdApply schreibt Artefakte. Ohne --apply passiert nichts außer Ausgabe.
func cmdApply(args []string) int {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	really := fs.Bool("apply", false, "wirklich schreiben statt nur zeigen")
	_ = fs.Parse(args)
	if !*really {
		fmt.Fprintln(os.Stderr, "Trockenlauf. Mit --apply wird geschrieben.")
	}
	// TODO(M2): render.All
	fmt.Fprintln(os.Stderr, "noch nicht gebaut")
	return 1
}

func cmdRemove(args []string) int {
	// TODO(M2): render.Remove — Units, Volumes, Proxy-Einträge, Sicherungsaufträge
	fmt.Fprintln(os.Stderr, "noch nicht gebaut")
	return 1
}

// cmdHardware liest hardware/*.yaml und listet, was sich einlesen lässt.
func cmdHardware(args []string) int {
	fs := flag.NewFlagSet("hardware", flag.ExitOnError)
	dir := fs.String("dir", "hardware", "Verzeichnis mit hardware/*.yaml")
	_ = fs.Parse(args)

	if fs.NArg() != 1 || fs.Arg(0) != "list" {
		fmt.Fprintln(os.Stderr, "hardware: stube hardware list")
		return 2
	}

	devices, err := hardware.Load(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for _, d := range devices.List {
		fmt.Println(d.ID)
	}
	return 0
}

// cmdCatalog liest catalog/services/*.yaml und listet, was sich einlesen lässt.
func cmdCatalog(args []string) int {
	fs := flag.NewFlagSet("catalog", flag.ExitOnError)
	dir := fs.String("dir", "catalog/services", "Verzeichnis mit catalog/services/*.yaml")
	_ = fs.Parse(args)

	if fs.NArg() != 1 || fs.Arg(0) != "list" {
		fmt.Fprintln(os.Stderr, "catalog: stube catalog list")
		return 2
	}

	services, err := catalog.Load(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for _, s := range services.List {
		fmt.Println(s.ID)
	}
	return 0
}
