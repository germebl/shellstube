// Command stube liest ein homelab.yaml und schreibt daraus lesbare Konfigurationsdateien.
//
// Es gibt keinen Daemon und keinen Fernzugriff. Alles, was stube tut, landet als Datei
// auf der Platte — und ist damit auch ohne stube nachvollziehbar.
package main

import (
	"flag"
	"fmt"
	"os"
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
  stube version

Der Trockenlauf ist die Vorgabe. Wer wirklich schreiben will, sagt --apply.
`)
}

// cmdValidate prüft ein Manifest, ohne irgendetwas zu verändern.
func cmdValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	schema := fs.String("schema", "schema/homelab.schema.json", "Pfad zum Schema")
	_ = fs.Parse(args)
	_ = schema
	// TODO(M2): manifest.Load + manifest.Validate
	fmt.Fprintln(os.Stderr, "noch nicht gebaut")
	return 1
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
